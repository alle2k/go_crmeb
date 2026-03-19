package middleware

import (
	"boilerplate/config"
	"boilerplate/pkg/constants"
	"boilerplate/pkg/redis"
	"boilerplate/pkg/response"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(config.AppConfig.Token.Header)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.FailWithMsg("Please login first"))
			return
		}
		if strings.HasPrefix(token, constants.UserTokenRedisKeyPrefix) {
			token = strings.TrimPrefix(token, constants.UserTokenRedisKeyPrefix)
		}
		key := constants.UserTokenRedisKeyPrefix + token
		res, err := redis.RDB.Get(c.Request.Context(), key).Result()
		if err != nil {
			log.Printf("get value from redis error: %#v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.FailWithMsg("Please login first"))
			return
		}
		id, err := strconv.Atoi(res)
		if err != nil {
			log.Printf("获取用户id异常: %#v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.FailWithMsg("Please login first"))
			return
		}
		c.Set("userId", id)
		redis.RDB.Set(c.Request.Context(), key, id, time.Duration(config.AppConfig.Token.ExpireTime)*time.Minute)
	}
}
