package repository

import (
	"boilerplate/internal/model"
	"boilerplate/pkg/mysql"

	"github.com/gin-gonic/gin"
)

type LotteryRuleRepository struct{}

func (l *LotteryRuleRepository) GetById(c *gin.Context, id int64) (*model.EbLotteryRule, error) {
	var r model.EbLotteryRule
	if err := mysql.DB.WithContext(c).First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

type LotteryUserRecordRepository struct{}

func (l *LotteryUserRecordRepository) ListByLotteryRuleId(c *gin.Context, ruleId int64) (*[]model.EbLotteryUserRecord, error) {
	var records []model.EbLotteryUserRecord
	if err := mysql.DB.WithContext(c).Where(&model.EbLotteryUserRecord{LotteryRuleId: ruleId}).Order("id").Find(&records).Error; err != nil {
		return nil, err
	}
	return &records, nil
}

func (l *LotteryUserRecordRepository) Save(c *gin.Context, ruleId int64, userId int, reward int64) error {
	if err := mysql.DB.WithContext(c).Create(&model.EbLotteryUserRecord{LotteryRuleId: ruleId, UserId: userId, Reward: reward, CreateUser: userId, UpdateUser: userId}).Error; err != nil {
		return err
	}
	return nil
}
