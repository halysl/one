package routers

import (
	. "github.com/halysl/one/log"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func removeComma(s string) (number int) {
	s = strings.Replace(s, ",", "", 2)
	number, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		ErrorLogger.Println(err)
	}
	return
}

func requestCommon(method, url string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	InfoLogger.Println(req.URL.String())

	if err != nil {
		ErrorLogger.Println(err)
		return []byte{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		ErrorLogger.Println(err)
		return []byte{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ErrorLogger.Println(err)
		return []byte{}, nil
	}
	return data, nil
}
