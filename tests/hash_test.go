package tests

import (
	"github.com/0xIsRookie/Helper/hash"
	"testing"
)

func TestStringMd5(test *testing.T) {
	testString := map[string]string{
		"admin":    "21232f297a57a5a743894a0e4a801fc3",
		"password": "5f4dcc3b5aa765d61d8327deb882cf99",
		"1@2,.":    "dfdb015cc9c20558f363f2c3d386bb64",
	}

	for k, v := range testString {
		if t := hash.StringMd5(k); t != v {
			test.Fatalf("[x] StringMd5 Test results do not match, encrypted string:%s result: %s != %s", k, t, v)
		}
	}

}
