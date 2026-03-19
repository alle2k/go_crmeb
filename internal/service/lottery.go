package service

import (
	"boilerplate/internal/dto"
	"boilerplate/pkg/constants"
	"boilerplate/pkg/redis"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LotteryService struct{}

func (l *LotteryService) Add(c *gin.Context, req *dto.LotteryAddReq) (int64, error) {
	userId := c.MustGet("userId")
	log.Printf("用户参与活动, userId: %v, ruleId: %v\n", userId, req.RuleId)
	if err := loadLotteryRule(c, req.RuleId); err != nil {
		log.Printf("加载奖励规则失败，ruleId: %v, err: %v\n", req.RuleId, err)
		return 0, err
	}
	return 1, nil
}

func loadLotteryRule(c *gin.Context, ruleId int64) error {
	rows, err := redis.RDB.Exists(c, constants.LotteryKey+fmt.Sprintf("%d", ruleId)).Result()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	key := constants.LotteryRuleEditKey + fmt.Sprintf("%d", ruleId)
	identifier, err := redis.Acquire(c, key, 3*time.Second, 5000*time.Millisecond)
	if err != nil {
		return err
	}
	defer redis.Release(c, key, identifier)
	// double check
	rows, err = redis.RDB.Exists(c, constants.LotteryKey+fmt.Sprintf("%d", ruleId)).Result()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	return nil
}
