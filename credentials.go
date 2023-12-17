package hugging_face_api

import (
	"os"
	"reflect"
)

type Credentials struct {
	HuggingFaceToken string `env:"HUGGING_FACE_API_TOKEN"`
}

func GetCredentials() Credentials {
	var result Credentials = Credentials{}
	nf := reflect.TypeOf(result).NumField()
	for i := 0; i < nf; i += 1 {
		envTag := reflect.TypeOf(result).Field(i).Tag.Get("env")
		envVar := os.Getenv(envTag)
		field := reflect.ValueOf(&result).Elem().Field(i)
		field.SetString(envVar)
	}

	return result
}
