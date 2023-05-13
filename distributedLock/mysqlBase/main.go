package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		db = initDB()
	})
	return db
}

func initDB() *gorm.DB {

	var user string = "root"
	var password string = "123456"
	var host string = "localhost"
	var port string = "3306"
	var dbname string = "offermaker"
	var Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		filepath.Join(dbname),
	)

	DB, err := gorm.Open(mysql.Open(Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	fmt.Println("successfully connect database")

	return DB
}

type DistributedLock struct {
	Name       string `gorm:"type:varchar(255);primary_key"`
	ExpireTime int64  `gorm:"type:bigint(20)"`
}

func FindLock(name string) (lock *DistributedLock, err error) {
	// 查询记录
	result := GetDB().First(&lock, "name = ?", name)
	// 处理查询结果
	if result.RowsAffected == 0 {
		return nil, errors.New("lock not found")
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		return lock, nil
	}
}

func DeleteLock(name string) (bool, error) {
	db := GetDB()
	tx := db.Begin()

	lock := DistributedLock{
		Name: name,
	}
	if err := tx.Delete(&lock).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}
	return true, nil
}

const (
	LOCK_EXPIRE_TIME = 2000 // 锁的过期时间（毫秒）
)

func (lock *DistributedLock) AcquireLock(timeout time.Duration) (success bool, err error) {
	startTime := time.Now()                               //当前时间
	now := startTime.UnixNano() / int64(time.Millisecond) // 获取当前时间的毫秒数
	lock.ExpireTime = now + LOCK_EXPIRE_TIME              // 计算锁的过期时间
	//检查如果存在锁，其是否过期
	existlock, err := FindLock(lock.Name)
	if err != nil {
		fmt.Println("查找存在错误") //注：这里无论查找成功或失败与否均不要返回，因为还未尝试获取锁呢
	}
	if existlock != nil {
		fmt.Printf("查找成功,存在该锁")
		//如果过期就删除
		if now > existlock.ExpireTime {
			fmt.Println("已过期")
			isDelete, err := DeleteLock(lock.Name)
			if err != nil {
				//删除失败
				fmt.Println("删除出现错误") //注：这里无论删除成功或失败与否均不要返回，因为还未尝试获取锁呢
			}
			if isDelete != true {
				fmt.Println("未删除成功") //注：这里无论删除成功或失败与否均不要返回，因为还未尝试获取锁呢
			}

		} else {
			fmt.Println("未过期")
		}

	}
	for {
		result := GetDB().Create(lock)
		if result.Error == nil || result.RowsAffected != 0 {
			fmt.Println("AcquireLock Successfully")
			return true, nil
		}
		if timeout > 0 && time.Now().Sub(startTime) > timeout {
			return false, nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (lock *DistributedLock) ReleaseLock() (success bool, err error) {
	result := GetDB().Where("name = ?", lock.Name).Unscoped().Delete(&DistributedLock{})
	if result.Error != nil || result.RowsAffected == 0 { // 如果删除失败，则锁不存在或已过期
		fmt.Println("ReleaseLock UnSuccessfully")
		return false, nil
	}
	fmt.Println("ReleaseLock Successfully")
	return true, nil
}

func main() {
	const LOCK_NAME = "my_lock" // 定义锁名称
	// 获取锁
	lock := &DistributedLock{Name: LOCK_NAME, ExpireTime: 0}
	if ok, err := lock.AcquireLock(time.Second * 1); err == nil {
		if ok {
			// 获取锁成功后的执行业务逻辑操作
			fmt.Println("获取锁成功，我要做一些操作了")
			defer lock.ReleaseLock() //释放锁
		} else {
			// 获取锁超时
			fmt.Println("timeout: fail to get DistributedLock")
		}
	} else {
		// 获取锁失败后的执行业务逻辑操作
		fmt.Println("获取锁失败")
	}

}
