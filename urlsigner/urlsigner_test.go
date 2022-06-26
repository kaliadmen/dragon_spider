package urlsigner

import (
	"testing"
)

var signTests = []struct {
	name     string
	url      string
	validUrl bool
	hasError bool
}{
	{
		name:     "query params",
		url:      "https://anything.com/test?id=1",
		validUrl: true,
		hasError: false,
	},
	{
		name:     "no query params",
		url:      "https://anything.com/test",
		validUrl: true,
		hasError: false,
	},
	{
		name:     "empty url",
		url:      "",
		validUrl: false,
		hasError: true,
	},
	{
		name:     "not url",
		url:      "invalid",
		validUrl: false,
		hasError: true,
	},
}

func TestSignature_SignUrl(t *testing.T) {
	sign := Signature{Secret: "1qaz@WSX"}

	for _, e := range signTests {
		signed, err := sign.SignUrl(e.url)

		if err == nil && e.hasError {
			t.Errorf("%s: does not have error, when it should", e.name)
		}

		if err != nil && !e.hasError {
			t.Errorf("%s: has error, when it should not have one", e.name)
		}

		if len(signed) > 0 && len(e.url) != 0 && e.validUrl && e.hasError {
			t.Errorf("%s: failed to sign valid url", e.name)
		}

		if !e.validUrl && err == nil {
			t.Errorf("%s: signed invalid url", e.name)
		}
	}
}

var verifyTests = []struct {
	name       string
	url        string
	validUrl   bool
	shouldPass bool
}{
	{
		name:       "valid url and signature",
		url:        "https://anything.com/test?id=1",
		shouldPass: true,
		validUrl:   true,
	},
	{
		name:       "valid url and invalid signature",
		url:        "https://www.anthing.com/some/url",
		shouldPass: false,
		validUrl:   false,
	},
	{
		name:       "not a url",
		url:        "not a url",
		shouldPass: false,
		validUrl:   false,
	},
}

func TestSignature_VerifyToken(t *testing.T) {
	sign := Signature{Secret: "1qaz@WSX"}

	for _, e := range verifyTests {
		var signed string

		if e.validUrl {
			signed, _ = sign.SignUrl(e.url)
		} else {
			signed = e.url
		}

		valid, err := sign.VerifyUrl(signed)

		if err != nil && e.validUrl {
			t.Errorf("%s: error when validating url %s", e.name, e.url)
		}
		if !valid && e.shouldPass {
			t.Errorf("%s: valid token shows as invalid", e.name)
		}
		if valid && !e.validUrl {
			t.Errorf("%s: returned valid on invalid url %s", e.name, e.url)
		}
	}
}

func TestSignature_IsExpired(t *testing.T) {
	sign := Signature{Secret: "1qaz@WSX"}

	signed, _ := sign.SignUrl("http://anything.com/test?id=1")

	expired := sign.IsExpired(signed, 1)

	if expired {
		t.Error("token shows expired when it should not")
	}

	expired = sign.IsExpired(signed, -1)
	if !expired {
		t.Error("token shows not expired when it should be")
	}
}
