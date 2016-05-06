package utils

import (
	"net/http"
	"io/ioutil"
	"bytes"
)

func  PostJson(url string,body string) (int,string,error) {

	var jsonStr = []byte(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return  500,"",err
	}


	client := &http.Client{}
	response, errs := client.Do(req)
	if errs != nil {
		return  response.StatusCode,"",err
	}
	defer response.Body.Close()
	retbody, _ := ioutil.ReadAll(response.Body)
	return  response.StatusCode,string(retbody),nil
}

func  Get(url string) (int,string,error)  {
	response, err := http.Get(url)

	if err != nil {
		return response.StatusCode,"",err
	} else {

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response.StatusCode,"",err
		}
		return response.StatusCode,string(contents),nil
	}
}
