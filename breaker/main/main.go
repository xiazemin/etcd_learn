package main

import (
	"fmt"
	"github.com/xiazemin/etcd_learn/breaker/goBreaker"
	"log"
	"time"
)


func main() {
	var cmds = []int32{1, 2, 289, 55}
	var options = goBreaker.Options{
		BucketTime:        150 * time.Millisecond,
		BucketNums:        200,
		BreakerRate:       0.6,
		BreakerMinSamples: 300,
		CoolingTimeout:    3 * time.Second,
		DetectTimeout:     150 * time.Millisecond,
		HalfOpenSuccess:   3,
	}
	cb := goBreaker.InitCircuitBreakers(cmds, options)
	fmt.Println(cb,cb.IsTriggerBreaker(289))

		// downStream service is broken
		for i:=0;i<310;i++ {
			body, err := cb.Get("http://www.baidu1.com", cb.GetBreaker(289))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(len(string(body)))
			if cb.IsTriggerBreaker(289) {
				fmt.Println("triger")
			}
		}

	}

/*
beforeRequest()
fail: 299 CLOSED
afterRequest false {0xc42006e1e0 {{0 0} 0 0 0 0} 2 {0 0 <nil>} {0 0 <nil>} 0 {150000000 200 0.6 0 300 3000000000 150000000 3 0x11f8070 0x11f7d00 0x107b490} 0x107b490}
0
beforeRequest()
2019/10/09 20:07:25 breaker state change, command 289: CLOSED -> OPEN, (succ: 0, err: 300, timeout: 0, rate: 1.00)
fail: 300 OPEN
afterRequest false {0xc42006e1e0 {{0 0} 0 0 0 0} 0 {13789900285145578098 539982538 0x13d1fe0} {0 0 <nil>} 0 {150000000 200 0.6 0 300 3000000000 150000000 3 0x11f8070 0x11f7d00 0x107b490} 0x107b490}
0
triger
 */
