package main

import (
"github.com/afex/hystrix-go/hystrix"
"net/http"
"time"
)

type response struct {
	msg string
}

func ordinaryHandler(w http.ResponseWriter, req *http.Request)  {
	w.Write([]byte("Hello\n"))
}

func AHandle(w http.ResponseWriter, r *http.Request) {
	handle(w, r, "aaa")
}

func BHandle(w http.ResponseWriter, r *http.Request) {
	handle(w, r, "bbb")
}

func handle(w http.ResponseWriter, r *http.Request, name string) {
	done := make(chan *response, 1)

	//增加fallback方法
	//fallback := func(err error) error {
	//  done <- &response{"fallback response\n"}
	//  return nil
	//}

	errChan := hystrix.Go(name, func() error {
		//请求延时
		time.Sleep(2 * time.Second)
		done <- &response{"OKK\n"}
		return nil
	}, nil)

	select {
	case err := <-errChan:
		http.Error(w, err.Error(), 500)
	case d := <-done:
		w.Write([]byte(d.msg))
	}
}

func main()  {

	hystrix.ConfigureCommand("aaa", hystrix.CommandConfig{
		Timeout: 1000,
	})

	hystrix.ConfigureCommand("bbb", hystrix.CommandConfig{
		Timeout: 5000,
	})

	http.HandleFunc("/aaa", AHandle)

	http.HandleFunc("/bbb", BHandle)

	http.HandleFunc("/hello", ordinaryHandler)

	http.ListenAndServe(":8001", nil)
}