package oppopush

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout : time.Second * 60,
	}
)

const (
	RetryTimes       = 2
)

func doPost(url string, form url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tryTime := 0
tryAgain:
	resp, err := client.Do(req)
	if err != nil {
		tryTime++
		if tryTime < RetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code:"+strconv.Itoa(resp.StatusCode))
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(result)
	str, err = strconv.Unquote(str)
	if err != nil {
		str = string(result)
	}
	return []byte(str), nil
}

func doGet(url string, params string) ([]byte, error) {
	req, err := http.NewRequest("GET", url+params, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code:"+strconv.Itoa(resp.StatusCode))
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
