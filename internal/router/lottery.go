package router

import (
	"boilerplate/internal/handler"
	"boilerplate/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerLotteryRoutes(r *gin.Engine) {
	h := new(handler.LotteryHandler)

	group := r.Group("/api/front/lottery", middleware.TokenMiddleware())
	{
		group.POST("/", h.Add)
	}
}
