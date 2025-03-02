package minio

import (
	"book-search/webcrawler/config"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var instance *minio.Client
var once sync.Once

func generateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getMinioClient() (*minio.Client, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		log.Println("Initializing Minio storage")
		useSSL := false

		instance, err = minio.New(env.MinioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatal("Minio client error: ", err)
		}
	})

	return instance, nil
}

func StoreHTML(ctx context.Context, html string) (string, error) {
	env, err := config.GetEnv()
	if err != nil {
		return "", err
	}

	client, err := getMinioClient()
	if err != nil {
		return "", err
	}

	bucketExist, err := client.BucketExists(ctx, env.MinioBucket)
	if err != nil {
		return "", err
	}

	if !bucketExist {
		err = client.MakeBucket(ctx, env.MinioBucket, minio.MakeBucketOptions{})
		log.Println("Created bucket:", env.MinioBucket)
		if err != nil {
			return "", err
		}
	}

	contentHash := generateContentHash(html)

	_, err = client.PutObject(
		ctx,
		env.MinioBucket,
		contentHash,
		strings.NewReader(html),
		int64(len(html)),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)

	if err != nil {
		return "", err
	}

	return contentHash, nil
}

func GetHTML(ctx context.Context, contentHash string) (string, error) {
	env, err := config.GetEnv()
	if err != nil {
		return "", err
	}

	client, err := getMinioClient()
	if err != nil {
		return "", err
	}

	object, err := client.GetObject(ctx, env.MinioBucket, contentHash, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer object.Close()

	// Read the object contents
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, object); err != nil {
		return "", err
	}

	return buf.String(), nil
}
