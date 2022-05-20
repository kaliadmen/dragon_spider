package urlsigner

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte
}

func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign string

	crypt := goalone.New(s.Secret, goalone.Timestamp)

	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

func (s *Signer) IsValidToken(token string) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	if err != nil {
		return false
	}

	return true
}

func (s *Signer) IsExpired(token string, minutesUntilExpiration int) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	timeSince := crypt.Parse([]byte(token))

	return time.Since(timeSince.Timestamp) > time.Duration(minutesUntilExpiration)*time.Minute
}
