package newrelic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (c *Client) doGet(path string, params fmt.Stringer, out interface{}) error {
	callPath := c.url.String() + path
	if params != nil {
		s := params.String()
		if s != "" {
			callPath += fmt.Sprintf("?%s", s)
		}
	}

	req, err := http.NewRequest("GET", callPath, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return c.doRequest(req, out)
}

func (c *Client) doPost(path string, params fmt.Stringer, body interface{}, out interface{}) error {
	callPath := c.url.String() + path
	if params != nil {
		s := params.String()
		if s != "" {
			callPath += fmt.Sprintf("?%s", s)
		}
	}

	var asJson []byte
	if body != nil {
		tmp, err := json.Marshal(body)
		if err != nil {
			return err
		}
		asJson = tmp
	}
	r := bytes.NewBuffer(asJson)
	req, err := http.NewRequest("POST", callPath, r)
	if err != nil {
		return err
	}
	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return c.doRequest(req, out)
}

func (c *Client) doUpdate(path string, params interface{}, out interface{}) error {
	var asJson []byte
	if params != nil {
		tmp, err := json.Marshal(params)
		if err != nil {
			return err
		}
		asJson = tmp
	}
	r := bytes.NewBuffer(asJson)
	req, err := http.NewRequest("PUT", c.url.String()+path, r)
	if err != nil {
		return err
	}
	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return c.doRequest(req, out)
}

func (c *Client) doDelete(path string, out interface{}) error {
	req, err := http.NewRequest("DELETE", c.url.String()+path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Api-Key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	return c.doRequest(req, out)
}

func (c *Client) doRequest(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("newrelic http error (%s): %s", resp.Status, b)
	}
	if len(b) == 0 {
		b = []byte{'{', '}'}
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}
	return nil
}

func encodeGetParams(params map[string]interface{}) string {
	s := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			val := v.(string)
			if val != "" {
				s.Add(k, val)
			}
		case int:
			val := v.(int)
			// TODO: Zero values versus not defined
			if val != 0 {
				s.Add(k, strconv.Itoa(val))
			}
		case []string:
			val := v.([]string)
			if len(val) != 0 {
				s.Add(k, strings.Join(val, ","))
			}
		case []int:
			val := v.([]int)
			arr := []string{}
			for _, v := range val {
				arr = append(arr, strconv.Itoa(v))
			}
			if len(arr) != 0 {
				s.Add(k, strings.Join(arr, ","))
			}
		case time.Time:
			val := v.(time.Time)
			if !val.IsZero() {
				s.Add(k, val.String())
			}
		case Array:
			val := v.(Array)
			for _, v := range val.arr {
				s.Add(k, v)
			}
		case bool:
			if v.(bool) {
				s.Add(k, "true")
			}
		default:
			s.Add(k, fmt.Sprintf("%v", v))
		}
	}
	return s.Encode()
}
