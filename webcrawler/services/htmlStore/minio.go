package htmlStore

import (
	"book-search/webcrawler/config"
	"context"
	"io"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetMinioClient() (*minio.Client, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, err
	}

	useSSL := false
	instance, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func StoreHTML(ctx context.Context, client *minio.Client, html string, htmlHash string) error {
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	bucketExist, err := client.BucketExists(ctx, env.MinioBucket)
	if err != nil {
		return err
	}

	if !bucketExist {
		err = client.MakeBucket(ctx, env.MinioBucket, minio.MakeBucketOptions{})
		log.Println("Created bucket:", env.MinioBucket)
		if err != nil {
			return err
		}
	}

	_, err = client.PutObject(
		ctx,
		env.MinioBucket,
		htmlHash,
		strings.NewReader(html),
		int64(len(html)),
		minio.PutObjectOptions{
			ContentType: "text/html",
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func GetHTML(ctx context.Context, client *minio.Client, htmlHash string) (string, error) {
	env, err := config.GetEnv()
	if err != nil {
		return "", err
	}

	object, err := client.GetObject(ctx, env.MinioBucket, htmlHash, minio.GetObjectOptions{})
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
