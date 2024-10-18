package snowflake

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
	once          sync.Once // 确保初始化只执行一次
)

// getMachineID 用于获取机器 ID
func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

// initSnowflake 内部初始化函数，只会被调用一次
func initSnowflake() {
	t, _ := time.Parse("2006-01-02", "2020-01-01") // 自定义开始时间
	settings := sonyflake.Settings{
		StartTime: t,
		MachineID: getMachineID,
	}

	sonyFlake = sonyflake.NewSonyflake(settings)
	if sonyFlake == nil {
		panic("failed to initialize Snowflake")
	}
}

// GetID 返回生成的唯一 ID（懒加载）
func GetID() (int64, error) {
	// 使用 sync.Once 确保 Snowflake 只初始化一次
	once.Do(initSnowflake)

	if sonyFlake == nil {
		return 0, fmt.Errorf("sonyflake not initialized")
	}

	id, err := sonyFlake.NextID()
	return int64(id), err
}

// SetMachineID 设置机器 ID，必须在第一次调用 GetID 之前设置
func SetMachineID(machineID uint16) {
	sonyMachineID = machineID
}
