package goBreaker

import (
	"github.com/golang/go/src/fmt"
	"net/http"
	"io/ioutil"
)

func (cb *CircuitBreaker) Execute(req func() (interface{},int, error),b *Breaker) (interface{}, error) {
	err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(false,0,b)
			//panic(e)
		}
	}()

	result,code, err := req()
	cb.afterRequest(err == nil,code,b)
	return result, err
}


func (cb *CircuitBreaker)beforeRequest()error{
	fmt.Println("beforeRequest()")
	return nil
}

func (cb *CircuitBreaker)afterRequest(sucess bool,code int,b *Breaker)  {
	if sucess{
		b.Succeed()
		fmt.Println("sucess:",b.Successes(),b.State())
	}else if code==504{
		b.Timeout()
		fmt.Println("timeout:",b.Timeouts(),b.State())
	}else{
		b.Fail()
		fmt.Println("fail:",b.Failures(),b.State())
	}
   fmt.Println("afterRequest",sucess,*b)
}



// Get wraps http.Get in CircuitBreaker.
func (cb *CircuitBreaker)Get(url string, b *Breaker) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, int, error) {
		resp, err := http.Get(url)
		code := resp.StatusCode
		if err != nil {
			return nil, code, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, code, err
		}

		return body, code, nil
	}, b)
	if body==nil || err != nil {
		return nil, err
	}

	return body.([]byte), nil
}
