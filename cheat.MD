# Go Syntax Cheatsheet

## 1 Package & Imports
```
package main       // main package for executable

import (
    "fmt"         // standard lib
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib" // side-effect import
)
```

## 2 Functions
```
func add(a int, b int) int {
    return a + b
}

// multiple return values
func divide(a, b int) (int, error) {
    if b == 0 { return 0, fmt.Errorf("divide by zero") }
    return a / b, nil
}
```

## 3 Variables & Constants
```
var x int = 10      // explicit
y := 20             // type inferred
const pi = 3.14     // constant
```

## 4 Structs
```
type User struct {
    ID       int
    Username string
    Email    string
}

u := User{ID:1, Username:"alice", Email:"a@b.com"}
fmt.Println(u.Username)
```

## 5 Pointers
```
x := 10
px := &x   // pointer
*px = 20   // dereference
```

## 6 Control Flow
### If/Else
```
if x > 10 {
    fmt.Println("large")
} else {
    fmt.Println("small")
}
```

### For loops
```
for i := 0; i < 5; i++ { fmt.Println(i) }
sum := 0
for sum < 10 { sum++ }  // while-style
```

### Switch
```
switch day := "Mon"; day {
case "Mon": fmt.Println("Start")
default: fmt.Println("Other")
}
```

## 7 Error Handling
```
result, err := divide(10,0)
if err != nil {
    fmt.Println(err)
    return
}
```

## 8 Defer
```
f, _ := os.Open("file.txt")
defer f.Close()  // runs at function exit
```

## 9 Concurrency (optional)
```
go func() {
    fmt.Println("in goroutine")
}()
```

go → starts a lightweight thread (goroutine)