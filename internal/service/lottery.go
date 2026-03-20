package service

import (
	"boilerplate/internal/dto"
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	"boilerplate/pkg/constants"
	"boilerplate/pkg/redis"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LotteryService struct {
	RuleRepo       *repository.LotteryRuleRepository
	UserRecordRepo *repository.LotteryUserRecordRepository
}

func NewLotteryService() *LotteryService {
	return &LotteryService{RuleRepo: new(repository.LotteryRuleRepository), UserRecordRepo: new(repository.LotteryUserRecordRepository)}
}

func (l *LotteryService) Add(c *gin.Context, req *dto.LotteryAddReq) (int64, error) {
	userId := c.MustGet("userId").(int)
	log.Printf("用户参与活动, userId: %v, ruleId: %v\n", userId, req.RuleId)
	if err := l.loadLotteryRule(c, req.RuleId); err != nil {
		log.Printf("加载奖励规则失败，ruleId: %v, err: %v\n", req.RuleId, err)
		return 0, err
	}
	userIdStr := strconv.Itoa(userId)
	if err := l.loadLotteryMembers(c, req.RuleId, userIdStr); err != nil {
		return 0, err
	}
	ruleIdStr := fmt.Sprintf("%d", req.RuleId)
	res, err := redis.RDB.HGet(c, constants.LotteryMembersKey+ruleIdStr, userIdStr).Result()
	if err != nil {
		return 0, err
	}
	rank, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	size, err := redis.RDB.LLen(c, constants.LotteryKey+ruleIdStr).Result()
	if err != nil {
		return 0, err
	}
	if int64(rank) > size {
		log.Printf("参与名额已用完，用户ID: %v, 活动ID: %v", userId, req.RuleId)
		redis.RDB.HDel(c, constants.LotteryMembersKey+ruleIdStr, userIdStr)
		redis.RDB.Decr(c, constants.LotteryUserOrderKey+ruleIdStr)
		return 0, errors.New("来晚啦，活动参与名额已用完～")
	}
	res, err = redis.RDB.LIndex(c, constants.LotteryKey+ruleIdStr, int64(rank)-1).Result()
	if err != nil {
		return 0, err
	}
	reward, _ := strconv.Atoi(res)
	if err := l.UserRecordRepo.Save(c, req.RuleId, userId, int64(reward)); err != nil {
		log.Printf("保存数据库失败，用户ID: %d, 活动ID: %v", userId, req.RuleId)
		redis.RDB.HDel(c, constants.LotteryMembersKey+ruleIdStr, userIdStr)
		redis.RDB.Decr(c, constants.LotteryUserOrderKey+ruleIdStr)
		return 0, err
	}
	return int64(reward), nil
}

func (l *LotteryService) loadLotteryRule(c *gin.Context, ruleId int64) error {
	lotteryKey := constants.LotteryKey + fmt.Sprintf("%d", ruleId)
	rows, err := redis.RDB.Exists(c, lotteryKey).Result()
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
	rows, err = redis.RDB.Exists(c, lotteryKey).Result()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}
	rule, err := l.RuleRepo.GetById(c, ruleId)
	if err != nil {
		return err
	}
	if rule == nil || rule.Status != model.LotteryRuleStatusEnabled {
		log.Printf("规则不存在或被禁用，ruleId: %d", ruleId)
		return errors.New("规则不存在或被禁用")
	}
	redis.RDB.RPush(c, lotteryKey, strings.Split(rule.Extra, ","))
	redis.RDB.Expire(c, lotteryKey, 7*24*time.Hour)
	return nil
}

func (l *LotteryService) loadLotteryMembers(c *gin.Context, ruleId int64, userId string) error {
	ruleIdStr := fmt.Sprintf("%d", ruleId)
	lotteryMembersKey := constants.LotteryMembersKey + ruleIdStr
	lotteryUserOrderKey := constants.LotteryUserOrderKey + ruleIdStr
	rows, err := redis.RDB.Exists(c, lotteryMembersKey).Result()
	if err != nil {
		return err
	}
	if rows > 0 {
		if err = addUserCache(c, ruleId, userId); err != nil {
			return err
		}
		return nil
	}
	membersLock := constants.LotteryMembersLock + ruleIdStr
	identifier, err := redis.Acquire(c, membersLock, 3*time.Second, 5*time.Second)
	if err != nil {
		return err
	}
	defer redis.Release(c, membersLock, identifier)
	// double check
	rows, err = redis.RDB.Exists(c, lotteryMembersKey).Result()
	if err != nil {
		return err
	}
	if rows > 0 {
		if err = addUserCache(c, ruleId, userId); err != nil {
			return err
		}
		return nil
	}
	redis.RDB.Set(c, lotteryUserOrderKey, 0, 8*24*time.Hour)
	list, err := l.UserRecordRepo.ListByLotteryRuleId(c, ruleId)
	if err != nil {
		return err
	}
	if nil == list || len(*list) < 1 {
		val, err := redis.RDB.Incr(c, lotteryUserOrderKey).Result()
		if err != nil {
			return err
		}
		_, err = redis.RDB.HSet(c, lotteryMembersKey, userId, val).Result()
		if err != nil {
			redis.RDB.Del(c, lotteryUserOrderKey)
			return err
		}
		redis.RDB.Expire(c, lotteryMembersKey, 7*24*time.Hour)
	} else {
		m := make(map[int]int64)
		for _, item := range *list {
			val, err := redis.RDB.Incr(c, lotteryUserOrderKey).Result()
			if err != nil {
				redis.RDB.Del(c, lotteryUserOrderKey)
				return err
			}
			m[item.UserId] = val
		}
		_, err = redis.RDB.HSet(c, lotteryMembersKey, m).Result()
		if err != nil {
			redis.RDB.Del(c, lotteryUserOrderKey)
			return err
		}
		redis.RDB.Expire(c, lotteryMembersKey, 7*24*time.Hour)
		err = addUserCache(c, ruleId, userId)
		if err != nil {
			redis.RDB.Del(c, lotteryUserOrderKey, lotteryMembersKey)
			return err
		}
	}
	return nil
}

func addUserCache(c *gin.Context, ruleId int64, userId string) error {
	exists, err := redis.RDB.HExists(c, constants.LotteryMembersKey+fmt.Sprintf("%d", ruleId), userId).Result()
	if err != nil {
		return err
	}
	if exists {
		log.Printf("用户已经参与过活动，用户ID: %d，活动ID: %d", userId, ruleId)
		return errors.New("您已参与过本次活动，感谢参与～")
	}
	val, err := redis.RDB.Incr(c, constants.LotteryUserOrderKey+fmt.Sprintf("%d", ruleId)).Result()
	if err != nil {
		return err
	}
	_, err = redis.RDB.HSet(c, constants.LotteryMembersKey+fmt.Sprintf("%d", ruleId), userId, val).Result()
	if err != nil {
		return err
	}
	return nil
}
