package api

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GarageSpace struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Location    *string `json:"location,omitempty"` // ex: "Shelf A1"
	CreatedAt   string  `json:"created_at"`
}

type GarageItem struct {
	ID        int64   `json:"id"`
	SpaceID   int64   `json:"space_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Notes     *string `json:"notes,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// RegisterGarageRoutes attaches garage endpoints to protected group
func RegisterGarageRoutes(group fiber.Router, db *sql.DB) {
	// Spaces
	group.Get("/garage/spaces", listGarageSpaces(db))
	group.Post("/garage/spaces", createGarageSpace(db))

	// Items
	group.Get("/garage/items", listGarageItems(db))
	group.Post("/garage/items", createGarageItem(db))
	group.Put("/garage/items/:id", updateGarageItem(db))
	group.Delete("/garage/items/:id", deleteGarageItem(db))
}

// ---------- SPACES HANDLERS ----------

func listGarageSpaces(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(`SELECT id, name, description, location, created_at FROM garage_spaces ORDER BY id`)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer rows.Close()

		var spaces []GarageSpace
		for rows.Next() {
			var s GarageSpace
			var desc, loc *string
			var created time.Time

			if err := rows.Scan(&s.ID, &s.Name, &desc, &loc, &created); err != nil {
				return fiber.ErrInternalServerError
			}
			s.Description = desc
			s.Location = loc
			s.CreatedAt = created.UTC().Format(time.RFC3339)
			spaces = append(spaces, s)
		}

		return c.JSON(spaces)
	}
}

func createGarageSpace(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			Name        string  `json:"name"`
			Description *string `json:"description"`
			Location    *string `json:"location"`
		}

		if err := c.BodyParser(&body); err != nil {
			return fiber.ErrBadRequest
		}
		if body.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name is required")
		}

		var id int64
		err := db.QueryRow(`
			INSERT INTO garage_spaces (name, description, location)
			VALUES ($1, $2, $3)
			RETURNING id
		`, body.Name, body.Description, body.Location).Scan(&id)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
	}
}

// ---------- ITEMS HANDLERS ----------

func listGarageItems(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		spaceID := c.Query("space_id") // optional filter

		q := `SELECT id, space_id, name, quantity, notes, created_at, updated_at FROM garage_items`
		args := []any{}
		if spaceID != "" {
			q += ` WHERE space_id = $1`
			args = append(args, spaceID)
		}
		q += ` ORDER BY id`

		rows, err := db.Query(q, args...)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		defer rows.Close()

		var items []GarageItem
		for rows.Next() {
			var it GarageItem
			var notes *string
			var created, updated time.Time

			if err := rows.Scan(&it.ID, &it.SpaceID, &it.Name, &it.Quantity, &notes, &created, &updated); err != nil {
				return fiber.ErrInternalServerError
			}
			it.Notes = notes
			it.CreatedAt = created.UTC().Format(time.RFC3339)
			it.UpdatedAt = updated.UTC().Format(time.RFC3339)
			items = append(items, it)
		}

		return c.JSON(items)
	}
}

func createGarageItem(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			SpaceID  int64   `json:"space_id"`
			Name     string  `json:"name"`
			Quantity int     `json:"quantity"`
			Notes    *string `json:"notes"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.ErrBadRequest
		}
		if body.Name == "" || body.SpaceID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "space_id and name are required")
		}
		if body.Quantity == 0 {
			body.Quantity = 1
		}

		var id int64
		err := db.QueryRow(`
			INSERT INTO garage_items (space_id, name, quantity, notes)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, body.SpaceID, body.Name, body.Quantity, body.Notes).Scan(&id)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
	}
}

func updateGarageItem(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return fiber.ErrBadRequest
		}

		var body struct {
			SpaceID  *int64  `json:"space_id"`
			Name     *string `json:"name"`
			Quantity *int    `json:"quantity"`
			Notes    *string `json:"notes"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.ErrBadRequest
		}

		// simplu: update full row (ai putea face și patch dinamic, dar nu e obligatoriu acum)
		_, err = db.Exec(`
			UPDATE garage_items
			SET space_id = COALESCE($1, space_id),
			    name     = COALESCE($2, name),
			    quantity = COALESCE($3, quantity),
			    notes    = COALESCE($4, notes),
			    updated_at = NOW()
			WHERE id = $5
		`, body.SpaceID, body.Name, body.Quantity, body.Notes, id)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func deleteGarageItem(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return fiber.ErrBadRequest
		}

		_, err = db.Exec(`DELETE FROM garage_items WHERE id = $1`, id)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
