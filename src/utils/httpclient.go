package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

func fetch(
	method, url, querystr string,
	body []byte,
	headers *map[string]string,
) (
	res *http.Response,
	data []byte,
	err error,
) {
	if querystr != "" {
		url = url + "?" + querystr
	}
	bodyReader := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	if headers != nil {
		for key, val := range *headers {
			req.Header.Add(key, val)
		}
	}

	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return res, data, nil
}

func Get(url string, querys interface{}, headers *map[string]string) (res *http.Response, data []byte, err error) {
	q, err := query.Values(querys)
	if err != nil {
		return nil, nil, err
	}

	return fetch("GET", url, q.Encode(), nil, headers)
}

func Post(url string, body []byte, headers *map[string]string) (res *http.Response, data []byte, err error) {
	return fetch("POST", url, "", body, headers)
}
