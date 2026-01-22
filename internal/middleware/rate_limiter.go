package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// 存储每个IP对应的令牌桶
var ipBuckets = make(map[string]*ratelimit.Bucket)
var mutex sync.Mutex

func getBucket(ip string) *ratelimit.Bucket {
	mutex.Lock()
	defer mutex.Unlock()

	bucket, exists := ipBuckets[ip]
	if !exists {
		// 创建令牌桶：每秒填充1个令牌，容量为5（允许短时突发）
		bucket = ratelimit.NewBucketWithQuantum(time.Second, 5, 1)
		ipBuckets[ip] = bucket
	}
	return bucket
}

func IpRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		bucket := getBucket(clientIP)
		// 尝试取一个令牌，如果不足则返回429
		if bucket.TakeAvailable(1) == 0 {
			c.AbortWithError(http.StatusTooManyRequests, errors.New("请求过于频繁，请稍后再试"))
			return
		}
		c.Next()
	}
}
