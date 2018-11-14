package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type client struct{}

func (c *client) Post(url string, data, output interface{}) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	res, e := http.Post(url, "application/json; charset=utf-8", b)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	e = json.Unmarshal(bodyBytes, output)
	if e != nil {
		return e
	}
	return nil
}
