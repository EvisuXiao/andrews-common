package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"github.com/EvisuXiao/andrews-common/exception"
)

var middleware = &Middleware{}

type Middleware struct {
	Controller
}

func (m *Middleware) RateLimiter(limit int) RouterHandler {
	bucket := ratelimit.NewBucketWithQuantum(time.Second, int64(limit), int64(limit))
	return func(c *gin.Context) bool {
		if bucket.TakeAvailable(1) < 1 {
			return m.FailureResponseWithCode(c, http.StatusTooManyRequests, exception.CustomErrWrapper("too many requests"))
		}
		return m.Next(c)
	}
}
