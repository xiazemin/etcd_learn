package ratelimit

import "github.com/micro/go-micro/client/selector"

func TestRateClientLimit(t *testing.T) {
	// setup
	r := mock.NewRegistry()
	s := selector.NewSelector(selector.Registry(r))

	testRates := []int{1, 10, 20, 100}

	for _, limit := range testRates {
		b := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

		c := client.NewClient(
			// set the selector
			client.Selector(s),
			// add the breaker wrapper
			client.Wrap(NewClientWrapper(b, false)),//fasle=快速失败?
		)

		req := c.NewRequest(
			"test.service",
			"Test.Method",
			&TestRequest{},
			client.WithContentType("application/json"),
		)
		rsp := TestResponse{}

		for j := 0; j < limit; j++ {
			err := c.Call(context.TODO(), req, &rsp)
			e := errors.Parse(err.Error())
			if e.Code == 429 {
				t.Errorf("Unexpected rate limit error: %v", err)
			}
		}

		err := c.Call(context.TODO(), req, rsp)
		e := errors.Parse(err.Error())
		if e.Code != 429 {
			t.Errorf("Expected rate limit error, got: %v", err)
		}
	}
}
