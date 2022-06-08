package miniofs

import (
	"context"
	"fmt"
	"github.com/kaliadmen/dragon_spider/filesystems"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"path"
	"strings"
	_ "time"
)

type Minio struct {
	Endpoint string
	Key      string
	Secret   string
	UseSSL   bool
	Region   string
	Bucket   string
}

func (m *Minio) getCredentials() *minio.Client {
	clinet, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.Key, m.Secret, ""),
		Secure: m.UseSSL,
	})

	if err != nil {
		log.Println("miniofs : ", err)
	}

	return clinet
}

func (m *Minio) Put(fileName, directory string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objName := path.Base(fileName)
	client := m.getCredentials()

	uploadInfo, err := client.FPutObject(ctx, m.Bucket, fmt.Sprintf("%s/%s", directory, objName), fileName, minio.PutObjectOptions{})
	if err != nil {
		log.Println("miniofs: Failed FPutObject")
		log.Println("miniofs:", err)
		log.Println("miniofs UploadInfo:", uploadInfo)
		return err
	}

	return nil
}

func (m *Minio) Get(destination string, items ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := m.getCredentials()

	for _, item := range items {
		err := client.FGetObject(ctx, m.Bucket, item, fmt.Sprintf("%s/%s", destination, path.Base(item)), minio.GetObjectOptions{})
		if err != nil {
			log.Println("miniofs:", err)
			return err
		}
	}

	return nil
}

func (m *Minio) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := m.getCredentials()

	ObjChan := client.ListObjects(ctx, m.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range ObjChan {
		if obj.Err != nil {
			fmt.Println("miniofs:", obj.Err)
			return listing, obj.Err
		}

		if !strings.HasPrefix(obj.Key, ".") {
			mb := filesystems.ToMegabytes(float64(obj.Size))

			item := filesystems.Listing{
				Etag:         obj.ETag,
				LastModified: obj.LastModified,
				Key:          obj.Key,
				Size:         mb,
			}

			listing = append(listing, item)

		}
	}

	return listing, nil
}

func (m *Minio) Delete(itemsToDelete []string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := m.getCredentials()

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	for _, item := range itemsToDelete {
		err := client.RemoveObject(ctx, m.Bucket, item, opts)
		if err != nil {
			fmt.Println("miniofs:", err)
			return false, err
		}
	}

	return true, nil
}
