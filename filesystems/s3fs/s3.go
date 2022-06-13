package s3fs

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kaliadmen/dragon_spider/filesystems"
	"log"
	"net/http"
	"os"
	"path"
)

type S3 struct {
	Key      string
	Secret   string
	Region   string
	Endpoint string
	Bucket   string
}

func (s *S3) getCredentials() *credentials.Credentials {
	c := credentials.NewStaticCredentials(s.Key, s.Secret, "")
	return c
}

func (s *S3) Put(fileName, directory string) error {
	c := s.getCredentials()

	amSess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &s.Endpoint,
		Region:      &s.Region,
		Credentials: c,
	}))

	uploader := s3manager.NewUploader(amSess)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("s3: ", err)
			return
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	var size = fileInfo.Size()

	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	_, err = uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Body:        fileBytes,
		Bucket:      aws.String(s.Bucket),
		ContentType: aws.String(fileType),
		Key:         aws.String(fmt.Sprintf("%s/%s", directory, path.Base(fileName))),
		Metadata: map[string]*string{
			"Key": aws.String("MetadataValue"),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *S3) Get(destination string, items ...string) error {
	c := s.getCredentials()

	amSess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &s.Endpoint,
		Region:      &s.Region,
		Credentials: c,
	}))

	for _, item := range items {
		err := func() error {
			file, err := os.Create(fmt.Sprintf("%s/%s", destination, item))
			if err != nil {
				return err
			}

			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Println("s3: ", err)
					return
				}
			}(file)

			downloader := s3manager.NewDownloader(amSess)

			_, err = downloader.Download(file, &s3.GetObjectInput{
				Bucket: aws.String(s.Bucket),
				Key:    aws.String(item),
			})

			if err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *S3) List(prefix string) ([]filesystems.Listing, error) {
	if prefix == "/" {
		prefix = ""
	}

	c := s.getCredentials()
	// create AWS session
	amSess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &s.Endpoint,
		Region:      &s.Region,
		Credentials: c,
	}))

	//create AWS service
	svc := s3.New(amSess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(prefix),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		return nil, err
	}

	var listing []filesystems.Listing

	for _, item := range result.Contents {
		fileSize := filesystems.ToMegabytes(float64(*item.Size))
		currItem := filesystems.Listing{
			Etag:         *item.ETag,
			LastModified: *item.LastModified,
			Key:          *item.Key,
			Size:         fileSize,
		}
		listing = append(listing, currItem)
	}

	return listing, nil
}

func (s *S3) Delete(itemsToDelete []string) (bool, error) {
	c := s.getCredentials()

	amSess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &s.Endpoint,
		Region:      &s.Region,
		Credentials: c,
	}))

	svc := s3.New(amSess)

	for _, item := range itemsToDelete {
		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.Bucket),
			Delete: &s3.Delete{
				Objects: []*s3.ObjectIdentifier{
					{
						Key: aws.String(item),
					},
				},
				Quiet: aws.Bool(false),
			},
		}

		_, err := svc.DeleteObjects(input)
		if err != nil {
			if amErr, ok := err.(awserr.Error); ok {
				switch amErr.Code() {
				default:
					log.Println("Amazon S3 Error:", amErr.Error())
					return false, err
				}
			} else {
				log.Println("s3:", err)
				return false, err
			}
		}
	}
	return true, nil
}
