package dragonSpider

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"time"
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

func (ds *DragonSpider) CreateDirs(path string) error {
	const mode = 0755

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
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

type Encryption struct {
	Key []byte
}

func (e *Encryption) Encrypt(text string) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *Encryption) Decrypt(encrypted string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(encrypted)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

func (ds *DragonSpider) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	//get calling function name
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)

	//regex for function name
	nameRegex := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := nameRegex.ReplaceAllString(funcObj.Name(), "$1")

	ds.InfoLog.Println(fmt.Sprintf(" \u001b[33mLoad Time\u001B[0m: \u001b[35m%s\u001B[0m took \u001b[32m%s\u001b[0m", funcName, elapsed))

}
