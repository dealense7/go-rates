package utils

import (
	"fmt"
	_ "github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", buildDSN())
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(400 * time.Millisecond)

	return db
}

func NewCacheClient() *redis.Client {
	addr, password := buildRedisPort()
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           0,
		Password:     password,
		PoolSize:     10,
		MinIdleConns: 10,
	})
}

func NewCache() (*redis_store.RedisStore, error) {
	cache := redis_store.NewRedis(NewCacheClient())

	return cache, nil
}

func buildDSN() string {
	var config = LoadDBEnv()

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}

func buildRedisPort() (string, string) {
	var config = loadRedisEnv()

	return fmt.Sprintf("%s:%s",
		config.Host,
		config.Port,
	), config.Password
}
