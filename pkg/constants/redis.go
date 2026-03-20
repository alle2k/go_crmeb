package constants

const (
	// UserTokenRedisKeyPrefix 用户登token redis存储前缀
	UserTokenRedisKeyPrefix = "TOKEN_USER:"
	// LotteryKey 乐透Key
	LotteryKey = "lottery:"
	// LotteryRuleEditKey 乐透配置修改锁Key
	LotteryRuleEditKey = "lotteryRuleEdit:"
	// LotteryMembersLock 乐透获取已参与用户锁Key
	LotteryMembersLock = "lotteryMemberLock:"
	// LotteryMembersKey 乐透已参与用户Key
	LotteryMembersKey = "lotteryMember:"
	// LotteryUserLock 乐透用户锁Key
	LotteryUserLock = "lotteryUser:"
	// LotteryUserOrderKey 参与乐透用户领取顺序Key
	LotteryUserOrderKey = "lotteryUserOrder:"
)
