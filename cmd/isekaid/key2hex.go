package main

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
)

func key2hexFromFile(fname string) (string, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}

	return key2hexConvert(string(data))
}

func key2hexConvert(data string) (string, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data))-1)
	_, err := base64.StdEncoding.Decode(buf, []byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
