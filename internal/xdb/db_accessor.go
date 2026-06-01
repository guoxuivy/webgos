package xdb

import (
	"math/rand"
	"sync"
	"webgos/internal/config"

	"gorm.io/gorm"
)

var (
	slaveIndex int
	slaveLock  sync.Mutex
)

// GetDB 获取主库连接（写操作）
func GetDB() *gorm.DB {
	return masterDB
}

// GetSlaveDB 获取备库连接（读操作）
func GetSlaveDB() *gorm.DB {
	if !config.GlobalConfig.Database.ReadWriteSeparation || len(slaveDBs) == 0 {
		return nil
	}
	switch config.GlobalConfig.Database.SlaveLoadBalance {
	case "random":
		return getRandomSlave()
	case "round_robin":
		return getRoundRobinSlave()
	default:
		return getRandomSlave()
	}
}

func getRandomSlave() *gorm.DB {
	slaveLock.Lock()
	defer slaveLock.Unlock()

	index := rand.Intn(len(slaveDBs))
	return slaveDBs[index]
}

func getRoundRobinSlave() *gorm.DB {
	slaveLock.Lock()
	defer slaveLock.Unlock()

	slave := slaveDBs[slaveIndex]
	slaveIndex = (slaveIndex + 1) % len(slaveDBs)
	return slave
}
