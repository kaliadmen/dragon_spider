package dragonSpider

import (
	"crypto/rand"
	"os"
)

const (
	randomString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_+"
)

//CreateDir creates a directory if it doesn't already exist'
func (ds *DragonSpider) CreateDir(path string) error {
	const mode = 0755

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *DragonSpider) CreateFile(path string) error {
	if !FileExists(path) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}

	return nil
}

//RandomString generates a random string length n from values in constant randomString
func (ds *DragonSpider) RandomString(length int) string {
	s, r := make([]rune, length), []rune(randomString)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

//FileExists checks if a file exists
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
