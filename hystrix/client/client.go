package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func run(result chan string, name string)  {
	start := time.Now().String()

	response, err := http.Get("http://127.0.0.1:8001/" + name)
	defer response.Body.Close()

	//duration := time.Now().String()

	if err != nil {
		result <- "请求开始：" + start + " url " + name + " response:" + "http service error "
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		result <- "请求开始：" + start + " url " + name + " response:" + string(body)
	}
}

func main()  {

	result := make(chan string)

	for i := 0; i < 10; i++ {
		go func() {
			run(result, "aaa")
			run(result, "bbb")
		}()
	}

	for {
		select {
		case r  := <- result:
			fmt.Print(r)
		default:
		}
	}
}