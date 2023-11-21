package middleware

import (
	"admin/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/zzj-custom/pkg/pRedis"
	"log/slog"
	"strconv"
)

type LimitConfig struct {
	// GenerationKey 根据业务生成key 下面CheckOrMark查询生成
	GenerationKey func(c *gin.Context) string
	// 检查函数,用户可修改具体逻辑,更加灵活
	CheckOrMark func(key string, expire int, limit int) error
	// Expire key 过期时间
	Expire int
	// Limit 周期时间
	Limit int
	Err   error
}

func (c LimitConfig) LimitWithTime(ctx *gin.Context) {
	if err := c.CheckOrMark(c.GenerationKey(ctx), c.Expire, c.Limit); err != nil {
		response.FailWithCode(response.RateLimitError, ctx)
		ctx.Abort()
		return
	}
	ctx.Next()
}

// DefaultGenerationKey 默认生成key
func DefaultGenerationKey(c *gin.Context) string {
	return "rate-limit:" + c.ClientIP()
}

func DefaultCheckOrMark(key string, expire int, limit int) (err error) {
	// 判断是否开启redis
	redisPool, err := pRedis.Pool()
	if err != nil {
		return errors.Wrapf(err, "获取redis连接池失败")
	}

	rds := redisPool.Get()
	defer func() { _ = rds.Close() }()

	if err = SetLimitWithTime(key, limit, expire); err != nil {
		slog.With("key", key, "limit", limit, "expire", expire).With(err).Error("setLimitWithTime失败")
	}
	return errors.Wrapf(err, "设置访问次数失败")
}

func DefaultLimit(ctx *gin.Context) {
	slog.Info("测试rate-limit")
	LimitConfig{
		GenerationKey: DefaultGenerationKey,
		CheckOrMark:   DefaultCheckOrMark,
		Expire:        viper.GetInt("rate_limit.expired"),
		Limit:         viper.GetInt("rate_limit.limit"),
	}.LimitWithTime(ctx)
}

// SetLimitWithTime 设置访问次数
func SetLimitWithTime(key string, limit int, expiration int) error {
	redisPool, err := pRedis.Pool()
	if err != nil {
		return errors.Wrapf(err, "获取redis连接池失败")
	}

	rds := redisPool.Get()
	defer func() { _ = rds.Close() }()

	count, err := redis.Int(rds.Do("EXISTS", key))
	if err != nil {
		return errors.Wrapf(err, "判断key[%s]是否存在错误", key)
	}

	t, err := redis.Int(rds.Do("TTL", key))
	if err != nil {
		return errors.Wrapf(err, "获取key[%s]的TTL失败", key)
	}

	if count == 0 || t < 0 {
		_, err = rds.Do("SET", key, 1, "EX", expiration)
		return errors.Wrapf(err,
			"自增key[%s]和过期时间expired[%d]失败,count[%d],t[%d]",
			key,
			expiration,
			count,
			t,
		)
	}

	// 获取数据
	times, err := redis.Int(rds.Do("GET", key))
	if err != nil {
		return errors.Wrapf(err, "获取redis限制key[%s]失败", key)
	}

	if times >= limit {
		return errors.New("请求太过频繁, 请 " + strconv.Itoa(t) + " 秒后尝试")
	}

	_, err = rds.Do("INCR", key)
	return errors.Wrapf(err, "自增key[%s]失败", key)
}
