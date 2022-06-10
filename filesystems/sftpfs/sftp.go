package sftpfs

import (
	"fmt"
	"github.com/kaliadmen/dragon_spider/filesystems"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type SFTP struct {
	Host     string
	User     string
	Password string
	Port     string
}

func (s *SFTP) getCredentials() (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	config := &ssh.ClientConfig{
		Config: ssh.Config{},
		User:   s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *SFTP) Put(fileName, directory string) error {
	client, err := s.getCredentials()
	if err != nil {
		log.Println("sftp: ", err)
		return err
	}

	defer func(client *sftp.Client) {
		err := client.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(client)

	uploadFile, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(uploadFile)

	fileOnServer, err := client.Create(fmt.Sprintf("%s/%s", directory, path.Base(fileName)))
	if err != nil {
		return err
	}

	defer func(f2 *sftp.File) {
		err := f2.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(fileOnServer)

	if _, err := io.Copy(fileOnServer, uploadFile); err != nil {
		return err
	}

	return nil
}

func (s *SFTP) Get(destination string, items ...string) error {
	client, err := s.getCredentials()
	if err != nil {
		log.Println("sftp: ", err)
		return err
	}

	defer func(client *sftp.Client) {
		err := client.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(client)

	for _, item := range items {
		err := func() error {
			destFile, err := os.Create(fmt.Sprintf("%s/%s", destination, path.Base(item)))
			if err != nil {
				return err
			}

			defer func(destFile *os.File) {
				err := destFile.Close()
				if err != nil {
					log.Println("sftp: ", err)
					return
				}
			}(destFile)

			srcFile, err := client.Open(item)
			if err != nil {
				return err
			}

			defer func(srcFile *sftp.File) {
				err := srcFile.Close()
				if err != nil {
					log.Println("sftp: ", err)
					return
				}
			}(srcFile)

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}

			//flush in memory copy
			err = destFile.Sync()
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

func (s *SFTP) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	client, err := s.getCredentials()
	if err != nil {
		log.Println("sftp: ", err)
		return listing, err
	}

	defer func(client *sftp.Client) {
		err := client.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(client)

	files, err := client.ReadDir(prefix)
	if err != nil {
		return listing, err
	}

	for _, file := range files {
		var item filesystems.Listing

		if !strings.HasPrefix(file.Name(), ".") {
			fileSize := filesystems.ToMegabytes(float64(file.Size()))

			item.Key = file.Name()
			item.Size = fileSize
			item.LastModified = file.ModTime()
			item.IsDir = file.IsDir()
			listing = append(listing, item)
		}
	}

	return listing, nil
}

func (s *SFTP) Delete(itemsToDelete []string) (bool, error) {
	client, err := s.getCredentials()
	if err != nil {
		log.Println("sftp: ", err)
		return false, err
	}

	defer func(client *sftp.Client) {
		err := client.Close()
		if err != nil {
			log.Println("sftp: ", err)
			return
		}
	}(client)

	for _, item := range itemsToDelete {
		deleteErr := client.Remove(item)
		if deleteErr != nil {
			return false, err
		}
	}

	return true, nil
}
