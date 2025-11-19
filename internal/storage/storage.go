package storage

import (
    "context"
    "io"
    "os"
    "time"

    "github.com/BlaccStacc/blaccend/internal/config"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func Init(c *config.Config) error {
    cfg := config.Load()

    client, err := minio.New(cfg.StorageEndpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
        Secure: false,
        Region: cfg.StorageRegion,
    })
    if err != nil {
        return err
    }

    Client = client

    ctx := context.Background()
    exists, err := Client.BucketExists(ctx, cfg.StorageBucket)
    if err != nil {
        return err
    }

    if !exists {
        return Client.MakeBucket(ctx, cfg.StorageBucket, minio.MakeBucketOptions{
            Region: cfg.StorageRegion,
        })
    }

    return nil
}

func UploadFile(path string, key string, contentType string) error {
    cfg := config.Load()

    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()

    info, err := f.Stat()
    if err != nil {
        return err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err = Client.PutObject(
        ctx,
        cfg.StorageBucket,
        key,
        f,
        info.Size(),
        minio.PutObjectOptions{ContentType: contentType},
    )

    return err
}

func GetFile(key string) (io.ReadCloser, error) {
    cfg := config.Load()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    obj, err := Client.GetObject(
        ctx,
        cfg.StorageBucket,
        key,
        minio.GetObjectOptions{},
    )

    if err != nil {
        return nil, err
    }

    return obj, nil
}
