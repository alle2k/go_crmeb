package dto

type LotteryAddReq struct {
	RuleId int64 `json:"ruleId" binding:"required"`
}
