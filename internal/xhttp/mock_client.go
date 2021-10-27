package xhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type MockClient struct {
	ResponseData interface{}
}

func (m *MockClient) Do(*http.Request) (*http.Response, error) {
	response := &http.Response{
		StatusCode: http.StatusOK,
	}

	if m.ResponseData == nil {
		return response, nil
	}

	data, err := json.Marshal(m.ResponseData)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewReader(data))

	return response, nil
}
