package client_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Xavier-Lam/go-wechat/internal/client"
	"github.com/Xavier-Lam/go-wechat/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestGetJsonSuccess(t *testing.T) {
	resp, _ := test.Responses.Json(`{"name":"Alice","age":25}`)

	expectedData := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Alice",
		Age:  25,
	}

	var actualData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := client.GetJson(resp, &actualData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, actualData)

	// Create a mock HTTP response with an unreadable body
	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	}

	var data interface{}
	err = client.GetJson(resp, &data)

	assert.Error(t, err)

	// Create a mock HTTP response with invalid JSON
	resp, _ = test.Responses.Json(`{"name":"Alice","age":25`)

	err = client.GetJson(resp, &data)

	assert.Error(t, err)
}
