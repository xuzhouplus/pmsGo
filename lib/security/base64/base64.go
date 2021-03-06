package base64

import "encoding/base64"

func Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Decode(decode string) []byte {
	decodeString, err := base64.StdEncoding.DecodeString(decode)
	if err != nil {
		return nil
	}
	return decodeString
}
