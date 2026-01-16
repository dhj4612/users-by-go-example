package logger

import (
	"github.com/gin-gonic/gin"
)

type Logger struct {
	ctx *gin.Context
}

func GetLogger(ctx *gin.Context) *Logger {
	return &Logger{
		ctx: ctx,
	}
}

func (l *Logger) Info() {
	//log.New()
}
