Wrapper
Wrapper提供了一种包装机制，使得在执行某方法前先执行Wrapper，优点Filter的意思；因此可以在客户端和服务器做很多功能：熔断限流、Filter、Auth等。
client代码如下：调用greeter.Hello时先执行logWrap.Call方法，再调用RPC请求。

server代码如下:当RPC调用进来时先执行logWrapper，再执行Hello

熔断
Micro提供了两种实现，gobreaker和hystrix，熔断是在客户端实现。先看看 hystrix：

hystrix会根据这5个参数（超时时间、并发请求数、请求量、空歇床、错误率）来选择合适的服务进行调度，目前是使用的 hystrix提供的默认参数，不支持自定义参数，示例:

gobreaker方案与hystrix类似，可以自定义参数。

限流
ratelimit可以在客户端做，也可以在服务端做；micro提供了两种方案：juju/ratelimit、uber/ratelimit。

客户端实现:


NewBucketWithRate入参为速率（QPS）和容量（CAP），比如每秒5个请求，最大保持50个活动的请求
NewClientWrapper第二个参数wait，指示当受到限流时是否等待，如果是false即快速失败，返回（429,too mant request）
服务端实现(以下代码包含了客户端测试代码):


client.NewClient支持多个Wrapper，将熔断限流功能都添加上

import "github.com/micro/go-plugins/wrapper/breaker/hystrix"
import "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
c := client.NewClient(
client.Wrap(ratelimit.NewClientWrapper(b, false)),
client.Wrap(hystrixNewClientWrapper()),
)