package model

import (
	"time"
)

type EbLotteryRule struct {
	ID         int64             `gorm:"column:id;primaryKey;autoIncrement"`
	UserCount  int               `gorm:"column:user_count"`
	TotalPrize int64             `gorm:"column:total_prize"`
	LowerLimit int64             `gorm:"column:lower_limit"`
	UpperLimit int64             `gorm:"column:upper_limit"`
	Extra      string            `gorm:"column:extra"`
	Status     LotteryRuleStatus `gorm:"column:status;default:1,check:status IN (0,1)"`
	CreateUser int               `gorm:"column:create_user"`
	CreateTime time.Time         `gorm:"column:create_time;autoCreateTime"`
	UpdateUser int               `gorm:"column:update_user"`
	UpdateTime time.Time         `gorm:"column:update_time;autoUpdateTime"`
}
