package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

func GetCache(redisCase usecase.RedisCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.FullPath() {
		case helper.CACHE_FETCH_EVENT_LIST:
			var events []entity.Event
			err := redisCase.Get(helper.CACHE_FETCH_EVENT_LIST, &events)
			if err == nil {
				c.AbortWithStatusJSON(http.StatusOK, helper.BuildResponse(true, "Cache!", events))
				return
			}
			c.Next()
		case helper.CACHE_FETCH_AVAILABLE_EVENT_LIST:
			var events []entity.Event
			err := redisCase.Get(helper.CACHE_FETCH_AVAILABLE_EVENT_LIST, &events)
			if err == nil {
				c.AbortWithStatusJSON(http.StatusOK, helper.BuildResponse(true, "Cache!", events))
				return
			}
			c.Next()
		default:
			c.Next()
		}
	}
}
