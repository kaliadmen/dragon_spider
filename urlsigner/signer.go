package urlsigner

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"net/url"
	"strings"
	"time"
)

// Signature is the type for the package. Secret is the urlsigner secret,
// a hard to guess string used to sign things. The secret must not exceed 64 characters.
type Signature struct {
	Secret string
}

// SignUrl generates a signed url and returns it, stripping off http:// and https://
func (s *Signature) SignUrl(data string) (string, error) {
	_, err := url.ParseRequestURI(data)
	if err != nil {
		return "", err
	}

	var (
		stringToSign string
		urlToSign    string
	)

	ex := strings.Split(data, "//")
	exploded := strings.Split(ex[1], "/")
	domain := exploded[0]
	exploded[0] = ""

	if strings.Contains(ex[1], "-") {
		stringToSign = strings.Join(exploded, "/")
	} else {
		stringToSign = strings.Join(exploded, "/")
	}

	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)

	if strings.Contains(stringToSign, "?") {
		// handle case where URL contains query parameters
		urlToSign = fmt.Sprintf("%s&hash=", stringToSign)
	} else {
		// no query parameters
		urlToSign = fmt.Sprintf("%s?hash=", stringToSign)
	}

	tokenBytes := pen.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return fmt.Sprintf("%s//%s%s", ex[0], domain, token), nil
}

// VerifyUrl verifies a signed url and returns true if it is valid,
// false if it is not. http:// and https:// are stripped off
// before verification
func (s *Signature) VerifyUrl(data string) (bool, error) {
	_, err := url.ParseRequestURI(data)
	if err != nil {
		return false, err
	}

	ex := strings.Split(data, "//")

	var exploded []string
	var stringToVerify string

	exploded = strings.Split(ex[1], "/")
	exploded[0] = ""

	if strings.Contains(ex[1], "-") {
		stringToVerify = strings.Join(exploded, "/")

	} else {
		stringToVerify = strings.Join(exploded, "/")
	}

	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)

	_, err = pen.Unsign([]byte(stringToVerify))
	if err != nil {
		// signature is not valid.
		return false, err
	}

	// valid signature
	return true, nil

}

// IsExpired checks to see if a token has expired. It returns true if
// the token was created within minutesUntilExpire, and false otherwise.
func (s *Signature) IsExpired(data string, minutesUntilExpire int) bool {
	exploded := strings.Split(data, "//")

	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)
	ts := pen.Parse([]byte(exploded[1]))

	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
