package breaker

import (
	"sync"
	"time"
	. "github.com/xiazemin/etcd_learn/breaker/goBreaker"
)

type Breaker struct {
	Container // contains all success, error and timeout
	sync.RWMutex

	state           State
	openTime        time.Time // the time when the breaker become OPEN
	lastRetryTime   time.Time // last retry time when in HALFOPEN state
	halfopenSuccess int       // consecutive successes when HALFOPEN

	options Options

	now func() time.Time
}

// Options for Breaker
type Options struct {
					 // parameters for container
	BucketTime time.Duration // the time each bucket holds
	BucketNums int           // the number of buckets the breaker have

					 // parameters for breaker
	BreakerRate        float64
	BreakerMinQPS      int           // when instance > 1, if qps is over this value, the breaker trip will work
	BreakerMinSamples  int           // for RateTrip callback
	CoolingTimeout     time.Duration // fixed when create
	DetectTimeout      time.Duration // fixed when create
	HalfOpenSuccess    int
	ShouldTrip         TripFunc // trip callback, default is RateTrip func
	StateChangeHandler StateChangeHandler

	now func() time.Time
}


type window struct {
	sync.Mutex
	oldest  int       // oldest bucket index
	latest  int       // latest bucket index
	buckets []*bucket // buckets this window has

	bucketTime time.Duration // time each bucket holds
	bucketNums int           // the largest number of buckets of window could have
	expireTime time.Duration // expire time of this window, equals to window size
	inWindow   int           // the number of buckets in the window currently

	conseErr int64 //consecutive errors
}


