package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedis() *redis.Client {
	once.Do(func() {
		redisClient = initRedis()
	})
	return redisClient
}

func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //No password set
		DB:       0,  // Use default DB
	})

	return client
}

const (
	lockExpire = 10 * time.Second // 锁的超时时间
	lockValue  = "1"              // 锁的内容
)

func acquireLock(lockName string, timeout time.Duration) bool {
	/*
		将lockKey作为锁的名称，lockValue作为锁的值，lockExpire作为锁的超时时间。
		如果SetNX命令返回1，说明成功地设置了键值对，并且这个锁被当前请求的客户端获取，
		如果SetNX命令返回0，说明当前锁已经被其他客户端获取，请求的这个客户端没有获取成功。
	*/
	startTime := time.Now()
	for {
		result, err := GetRedis().SetNX(lockName, lockValue, lockExpire).Result()
		if err == nil && result == true {
			return true
		}
		if timeout > 0 && time.Now().Sub(startTime) > timeout {
			return false
		}
		time.Sleep(time.Millisecond * 100)
	}

}

func releaseLock(lockName string) bool {
	result, err := GetRedis().Del(lockName).Result()
	if err != nil || result == 0 {
		return false
	}
	return true
}

// 示例代码
func main() {

	const lockKey = "mylock" // 锁的名称
	if acquireLock(lockKey, time.Second*1) {
		// 执行需要加锁的代码
		fmt.Println("DO SOMETHING")
		defer releaseLock(lockKey)
	} else {
		// 未获取到锁
		fmt.Println("Failed to acquire lock")
	}

}
