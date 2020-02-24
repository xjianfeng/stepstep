package define

/* 全局KEY 定义 */

//REDIS KEY 定义
var (
	REDIS_PRE_FIX  = "STEP_"
	REDIS_KEY_USER = REDIS_PRE_FIX + "USER_"

	REDIS_HEKY_USER_CACHE      = "cache"
	REDIS_HEKY_USER_SESSIONKEY = "sessionkey"
	REDIS_HEKY_USER_TOKEN      = "token"

	REDIS_KEY_INCR_UID = REDIS_PRE_FIX + "UID"

	REDIS_KEY_FORMID = REDIS_PRE_FIX + "FORMID_"
)

var (
	//每两个小时领取红包
	TIMOUT_RED_REDPACK int64 = 7200
	//每两个小时领取步数
	TIMOUT_AWARD_STEP int64 = 7200
	//观看广告频率秒数
	TIMOUT_READ_VIDEO int64 = 30

	//超时键值类型
	TIMEOUT_KEY_LUCKSTEP = "luckstep"
	TIMEOUT_KEY_VIDEO    = "advideo"
	TIMEOUT_KEY_REDPACK  = "redpack"

	//不限制次数
	TIMEOUT_NOT_LIMIT_REFRESH = -1

	TIMEOUT_INIT_KEYS = []string{
		TIMEOUT_KEY_LUCKSTEP,
		TIMEOUT_KEY_VIDEO,
		TIMEOUT_KEY_REDPACK,
	}

	//观看广告次数
	AD_VIDEO_MAX_CNT = 10

	//响应格式
	RESPONSE_TYPE_JSON   = "json"
	RESPONSE_TYPE_STRING = "string"

	//奖励类型
	AWARD_TYPE_STEP  = "step"
	AWARD_TYPE_MONEY = "money"

	//定时状态
	STATUS_INIT   int8 = 0
	STATUS_ALREDY int8 = 1
	STATUS_FINISH int8 = 2

	//时间秒杀定义
	TIME_WEEK_SEC int64 = 86400 * 7

	//HTTP 头定义
	HTTP_HEDER_JSON = map[string]string{"Content-Type": "application/json"}

	//明细类型
	DESC_TYPE_TASK            = "运动奖励"
	DESC_TYPE_FIREND          = "好友红包"
	DESC_TYPE_TIMEOUT_REDPACK = "定时红包"
	DESC_TYPE_KEEP_REDPACK    = "打卡红包"

	//系统标识
	SYSTEM_TASK            = "task"
	SYSTEM_FIREND          = "friend"
	SYSTEM_TIMEOUT_REDPACK = "timeout"
	SYSTEM_KEEP_REPACK     = "keep"

	//变化类型
	MOD_TYPE_ADD int8 = 1
	MOD_TYPE_SUB int8 = 2
)

const (
	AWARD_LIST_MAX_IDX = 10

	//奖励幸运步数
	AWARD_LUCK_STEP = 1000
	//奖励视频步数
	AWARD_VIDEO_STEP = 1000

	//红包随机数
	AWARD_RED_PACK       = 4
	AWARD_FRIEND_REDPACK = 6

	//运动红包数
	SPORT_REDPACK_NUM = 2
	//保存的步数天数
	STEP_HISTORY_DAY = 3

	TIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	AWARD_STEP_INTERVAL = []int{1, 1000, 3000, 6000, 10000, 25000, 40000, 50000, 60000, 80000}
	AWARD_SPORT_REDPACK = []int{6, 20}
)
