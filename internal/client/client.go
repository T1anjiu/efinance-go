package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
)

var DefaultClient *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
	}
	DefaultClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

var defaultHeaders = http.Header{
	"User-Agent":      {"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko"},
	"Accept":          {"*/*"},
	"Accept-Language":  {"zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2"},
}

func Get(rawURL string, params url.Values, headers http.Header) ([]byte, error) {
	var body []byte
	err := retry.Do(func() error {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		if params != nil {
			q := u.Query()
			for k, vs := range params {
				for _, v := range vs {
					q.Set(k, v)
				}
			}
			u.RawQuery = q.Encode()
		}
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return err
		}
		for k, vs := range defaultHeaders {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
		if headers != nil {
			for k, vs := range headers {
				for _, v := range vs {
					req.Header.Set(k, v)
				}
			}
		}
		resp, err := DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(time.Second))
	return body, err
}

func PostJSON(rawURL string, body io.Reader, headers http.Header) ([]byte, error) {
	var result []byte
	err := retry.Do(func() error {
		req, err := http.NewRequest("POST", rawURL, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		for k, vs := range defaultHeaders {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
		if headers != nil {
			for k, vs := range headers {
				for _, v := range vs {
					req.Header.Set(k, v)
				}
			}
		}
		resp, err := DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		result, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(time.Second))
	return result, err
}

func PostForm(rawURL string, data url.Values, headers http.Header) ([]byte, error) {
	var result []byte
	err := retry.Do(func() error {
		req, err := http.NewRequest("POST", rawURL, strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for k, vs := range defaultHeaders {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
		if headers != nil {
			for k, vs := range headers {
				for _, v := range vs {
					req.Header.Set(k, v)
				}
			}
		}
		resp, err := DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		result, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(time.Second))
	return result, err
}
