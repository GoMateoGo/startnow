package initialize

import (
	"fmt"
	"os"

	"second_hand_mall/internal/global"
	"second_hand_mall/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.GVAL_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.GVAL_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GVAL_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.GVAL_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger = zap.New(zapcore.NewTee(cores...))
	if global.GVAL_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
