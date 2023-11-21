package initializer

import (
	"github.com/spf13/viper"
	"github.com/zzj-custom/pkg/pAntsPool"
	"github.com/zzj-custom/pkg/pMysql"
	"github.com/zzj-custom/pkg/pRedis"
	"log"
	"log/slog"
)

func InitResource() {
	resource := new(Resource)

	// 初始化mysql
	resource.MysqlInit()

	// 初始化redis
	resource.RedisInit()

	// 初始化ants_pool
	resource.AntsPoolInit()

}

type ResourceInterface interface {
	MysqlInit()
	RedisInit()
	AntsPoolInit()
}

type Resource struct{}

func (*Resource) MysqlInit() {
	dbs := make(map[string]*pMysql.Database)
	err := viper.UnmarshalKey("database", &dbs)
	if err != nil {
		log.Fatal("数据库配置获取失败")
	}

	if _, err = pMysql.DbInit(dbs); err != nil {
		log.Fatal("初始化mysql失败")
	}
}

func (*Resource) RedisInit() {
	dbs := make([]*pRedis.MultiDialConfig, 0)
	err := viper.UnmarshalKey("redis", &dbs)
	if err != nil {
		log.Fatal("获取redis配置失败")
	}

	if err = pRedis.InitMultiPools(dbs); err != nil {
		log.Fatal("初始化redis配置失败")
	}
}

func (*Resource) AntsPoolInit() {
	size := viper.GetInt("ants_pool.size")
	if _, err := pAntsPool.InitAsyncTaskPool(size); err != nil {
		slog.Info("初始化ant_pool失败")
	}
}

func ReleaseResource() {
	// 释放mysql
	dbs := pMysql.Dbs()
	for key, db := range dbs {
		if conn, err := db.DB(); err != nil {
			_ = conn.Close()
		} else {
			slog.With("key", key).Info("")
		}
	}
}
