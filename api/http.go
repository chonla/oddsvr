package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type client struct {
	AccessToken string
}

func (c *client) Get(url string, output interface{}) error {
	httpClient := &http.Client{}
	req, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return e
	}
	if c.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	}

	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		e = json.Unmarshal(bodyBytes, output)
		if e != nil {
			return e
		}
		return nil
	}
	return fmt.Errorf("error: %s", resp.Status)
}

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
