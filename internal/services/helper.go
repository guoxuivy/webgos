package services

import (
	"context"
	"webgos/internal/xdb"

	"gorm.io/gorm"
)

// ctxDB 获取主库连接（写操作）
func ctxDB(ctx context.Context) *gorm.DB {
	return xdb.GetDB().WithContext(ctx)
}

// ctxSDB 获取从库连接（读操作）
// 注意：从库存在复制延迟，写后紧跟的读、强一致性读请使用 ctxDB 走主库
func ctxSDB(ctx context.Context) *gorm.DB {
	if xdb.GetSlaveDB() == nil {
		return xdb.GetDB().WithContext(ctx)
	}
	return xdb.GetSlaveDB().WithContext(ctx)
}
