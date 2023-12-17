package hugging_face_api

import (
	"encoding/json"
	"io"
	http "net/http"
	"strings"
)

const API_URL = "https://api-inference.huggingface.co/models/"

type ModelContext struct {
	ModelId string
}

func (ctx *ModelContext) Request(payload interface{}) (*http.Response, error) {
	var payloadString string
	switch p := payload.(type) {
	case []byte:
		payloadString = string(p)
	default:
		json, _ := json.Marshal(p)
		payloadString = string(json)
	}
	payloadReader := strings.NewReader(string(payloadString))
	request, _ := http.NewRequest("POST", API_URL+ctx.ModelId, payloadReader)
	request.Header.Add("Authorization", "Bearer "+GetCredentials().HuggingFaceToken)
	client := &http.Client{}
	return client.Do(request)
}

func ReadAllClose(rc *io.ReadCloser) ([]byte, error) {
	defer (*rc).Close()
	return io.ReadAll(*rc)
}

type Payload interface {
	GetPayloadReader() (*strings.Reader, error)
}
