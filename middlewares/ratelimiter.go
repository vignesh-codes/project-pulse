package middlewares

import (
	"pulse-service/apps/repository/adapter"
	"pulse-service/logger"
	"pulse-service/utils/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const USER_RATE_LIMIT_SECONDS = 30

// - 1 req per 30s for a user
func CheckRateLimit(repository *adapter.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_uid, found := c.Get("user_uid")
		if !found {
			status := response.ValidationError(response.ErrEmptyParam, "Login required")
			c.JSON(status.Status(), status)
			c.Abort()
			return
		}
		from := c.Request.URL.Query().Get("from")
		some_info := c.Request.URL.Query().Get("some_info")
		key := user_uid.(string) + ":" + "user_rate_limiter" + ":" + from + ":" + some_info
		out, err := repository.RedDB.GetKey(key)
		if err == redis.Nil {
			_ = repository.RedDB.SetKey(key, []byte(strconv.Itoa(int(time.Now().UnixMilli()))), time.Second*USER_RATE_LIMIT_SECONDS)
		} else if err != nil {
			logger.Logger.Error(logger.RedisError, zap.Any(logger.KEY_ERROR, err.Error()))
			c.Next()
			return
		} else {
			lastReqTimestamp, _ := strconv.Atoi(string(out))
			if int64(lastReqTimestamp+(1000*int(USER_RATE_LIMIT_SECONDS-1))) >= time.Now().UnixMilli() {
				// 29s
				status := response.RateLimitExceedError(response.ErrRateLimitExceed, "Too many requests")
				c.JSON(status.Status(), status)
				c.Abort()
				return
			}
		}

		c.Next()

	}
}
