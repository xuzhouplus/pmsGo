package helper

import (
	"log"
	"math/rand"
	"reflect"
	"time"
)

func FirstToUpper(input string) string {
	if input == "" {
		return ""
	}
	tmp := []byte(input)
	first := tmp[0]
	if first > 96 && first < 123 {
		tmp[0] = first - 32
		return string(tmp)
	}
	return input
}

func FirstToLower(input string) string {
	if input == "" {
		return ""
	}
	tmp := []byte(input)
	first := tmp[0]
	if first > 64 && first < 91 {
		tmp[0] = first + 32
		return string(tmp)
	}
	return input
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // 伪随机种子
	baseStr := "abcdefghigklmnopqistuvwxyzABCDEFGHIGKLMNOPQISTUVWXYZ0123456789"
	salt := make([]byte, length)
	for n := 0; n < length; n++ {
		salt[n] = baseStr[rand.Int31n(int32(len(baseStr)))]
	}
	return string(salt)
}

func IsInSlice(slice []interface{}, val string) (int, bool) {
	log.Println(slice)
	if len(slice) == 0 {
		return -1, false
	}
	for i, item := range slice {
		log.Println(item)
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func DynamicInvoke(object interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	//动态调用方法
	reflect.ValueOf(object).MethodByName(methodName).Call(inputs)

	//动态访问属性
	reflect.ValueOf(object).Elem().FieldByName("Name")
}

func CallStructMethod()  {
	
}
