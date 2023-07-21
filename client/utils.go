package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetJson(resp *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return err
}
