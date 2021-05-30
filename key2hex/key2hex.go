package key2hex

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
)

func FromFile(fname string) (string, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}

	return Convert(string(data))
}

func Convert(data string) (string, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data))-1)
	_, err := base64.StdEncoding.Decode(buf, []byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
