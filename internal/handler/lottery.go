package handler

import (
	"boilerplate/internal/dto"
	"boilerplate/internal/service"
	"boilerplate/pkg/constants"
	"boilerplate/pkg/redis"
	"boilerplate/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LotteryHandler struct{}

func (l *LotteryHandler) Add(c *gin.Context) {
	var form dto.LotteryAddReq
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, response.FailWithMsg("ruleId is required"))
		return
	}
	id := c.MustGet("userId")
	res, err := redis.Lock(c, constants.LotteryUserLock, func() (r *response.CommonResult) {
		reward, err := new(service.LotteryService).Add(c, &form)
		if err != nil {
			log.Printf("抽奖失败，用户ID: %v，错误: %v", id, err)
			c.JSON(http.StatusInternalServerError, response.Fail())
			return
		}
		return response.Success(reward)
	})
	if err != nil {
		log.Printf("抽奖失败，用户ID: %v，错误: %v", id, err)
		c.JSON(http.StatusInternalServerError, response.Fail())
		return
	}
	c.JSON(http.StatusOK, res)
}
