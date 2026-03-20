package model

import "time"

type EbLotteryUserRecord struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	LotteryRuleId int64     `gorm:"column:lottery_rule_id"`
	UserId        int       `gorm:"column:user_id"`
	Reward        int64     `gorm:"column:reward"`
	CreateUser    int       `gorm:"column:create_user"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateUser    int       `gorm:"column:update_user"`
	UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime"`
}
