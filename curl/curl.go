package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/EvisuXiao/andrews-common/utils"
)

const ArrayBodyKey = "rawArray"

func Request(target, method string, data map[string]interface{}, result interface{}, requestWrappers ...func(r *http.Request)) error {
	var body io.Reader
	if method == http.MethodGet {
		query := url.Values{}
		for k, v := range data {
			query.Set(k, fmt.Sprint(v))
		}
		urlObj, _ := url.Parse(target)
		urlObj.RawQuery = query.Encode()
		target = urlObj.String()
	} else {
		var j []byte
		if rawArray, ok := data[ArrayBodyKey]; ok {
			j, _ = json.Marshal(rawArray)
		} else {
			j, _ = json.Marshal(data)
		}
		body = bytes.NewReader(j)
	}
	request, err := http.NewRequest(method, target, body)
	if utils.HasErr(err) {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if !utils.IsEmpty(requestWrappers) {
		requestWrapper := requestWrappers[0]
		requestWrapper(request)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if utils.HasErr(err) {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		resBody, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(resBody, result)
		return nil
	}
	return fmt.Errorf("request error with status %d", response.StatusCode)
}
