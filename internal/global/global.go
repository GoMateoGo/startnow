package global

import (
	"second_hand_mall/config"

	"go.uber.org/zap"
)

var (
	GVAL_CONFIG config.Config // 配置
	GVAL_LOG    *zap.Logger   // 日志
)
