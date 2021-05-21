package json

import (
	"encoding/json"
	"fmt"
)

func Encode(data interface{}) (string, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func Decode(data string, rt interface{}) error {
	fmt.Println(data)
	err := json.Unmarshal([]byte(data), &rt)
	if err != nil {
		return err
	}
	return nil
}
