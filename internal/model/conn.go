package model

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/seekincs/rate-limit/configs"
)

type pool struct {
	MySQL *sql.DB
	Redis *redis.Client
}

var (
	conn           pool
	once           sync.Once
	evalScriptHash string
)

func initMySQL() {
	mysql, err := sql.Open("mysql", configs.DSN)
	if err != nil {
		panic(err)
	}
	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)
	conn.MySQL = mysql
	log.Println("mysql connected.")
}

func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr: configs.RedisAddr,
	})
	if pingErr := rdb.Ping(context.Background()).Err(); pingErr != nil {
		panic(pingErr)
	}
	conn.Redis = rdb
	log.Println("redis connected.")
}

func InitDB() {
	once.Do(func() {
		initMySQL()
		initRedis()
		go func() {
			evalScriptHash = fmt.Sprintf("%x", sha1.Sum([]byte(configs.RateLimitEvalScript)))
			exists := conn.Redis.ScriptExists(context.Background(), evalScriptHash)
			if exists.Err() != nil {
				log.Panic(exists)
			}
			if ok, _ := exists.Result(); !ok[0] {
				loadResult := conn.Redis.ScriptLoad(context.Background(), evalScriptHash)
				if loadResult.Err() != nil {
					log.Panic(loadResult)
				}
				log.Printf("eval script loaded:%s|%s", evalScriptHash, configs.RateLimitEvalScript)
			}
		}()
	})
}

func EvalSha(token_bucket_key string, qps int32, n int32) bool {
	evalResult := conn.Redis.EvalSha(context.Background(), evalScriptHash,
		[]string{token_bucket_key + ":token", token_bucket_key + ":timestamp"},
		qps, qps,
		time.Now().Unix(),
		n,
	)
	ok := true
	if evalResult.Err() != nil {
		log.Println(evalResult.Err().Error())
		ok = false
	}
	return ok
}
