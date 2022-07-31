package request

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func PostData(uri string, data []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("content-type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("status is " + resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}
