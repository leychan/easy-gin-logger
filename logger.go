package easyginlogger

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

type CustomLogger struct {
    Logger *slog.Logger
}

// 自定义等级名称
var LevelNames = map[slog.Leveler]string{
    LevelTrace:      "TRACE",
    LevelNotice:     "NOTICE",
    LevelFatal:      "FATAL",
}

const (
    LevelTrace  = slog.Level(-8)
    LevelNotice = slog.Level(2)
    LevelFatal  = slog.Level(12)
)

// NewCustomLogger 构造函数
func NewCustomLogger(hdl slog.Handler) *CustomLogger {
    return &CustomLogger{
        Logger: slog.New(hdl),
    }
}

// Log 通用记录日志
func (cl *CustomLogger) Log(ctx *gin.Context, msg string, level slog.Level, subType string, attrs []slog.Attr) {
    appendAttrs := cl.attachAttributes(ctx)
    attrs = append(attrs, appendAttrs...)
    var anyAttrs []any
    for _, attr := range attrs {
        anyAttrs = append(anyAttrs, attr)
    }

    if subType != "" {
        anyAttrs = append(anyAttrs, slog.String("sub_type", subType))
    }
    cl.Logger.Log(context.Background(), level, msg, anyAttrs...)
}

// Trace
func (cl *CustomLogger) Trace(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, LevelTrace, subType, attrs)
}

// Info
func (cl *CustomLogger) Info(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, slog.LevelInfo, subType, attrs)
}

// Debug
func (cl *CustomLogger) Debug(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, slog.LevelDebug, subType, attrs)
}

// Warn
func (cl *CustomLogger) Warn(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, slog.LevelWarn, subType, attrs)
}

// Error
func (cl *CustomLogger) Error(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, slog.LevelError, subType, attrs)
}

// Fatal
func (cl *CustomLogger) Fatal(ctx *gin.Context, msg string, subType string, attrs []slog.Attr) {
    cl.Log(ctx, msg, LevelFatal, subType, attrs)
}

// attachAttributes 添加属性
func (cl *CustomLogger) attachAttributes(ctx *gin.Context) []slog.Attr {
    commonAttrs := getCommonAttributes()
    if  ctx == nil {
        return commonAttrs
    }

    return append(commonAttrs, getGinContextAttributes(ctx)...)
}

// getCommonAttributes 获取公共属性
func getCommonAttributes() []slog.Attr {
    hostname, _ := os.Hostname()
    attrs := []slog.Attr{
        {Key: "host", Value: slog.StringValue(hostname)},
    }
    return attrs
}

// getGinContextAttributes 获取gin上下文属性
func getGinContextAttributes(ctx *gin.Context) []slog.Attr {
    attrs := []slog.Attr{
        {Key: "method", Value: slog.StringValue(ctx.Request.Method)},
        {Key: "host", Value: slog.StringValue(ctx.Request.Host)},
        {Key: "UA", Value: slog.StringValue(ctx.Request.UserAgent())},
        {Key: "client_ip", Value: slog.StringValue(ctx.ClientIP())},
        {Key: "query", Value: slog.StringValue(ctx.Request.URL.RawQuery)},
        {Key: "url", Value: slog.StringValue(ctx.Request.URL.Path)},
        {Key: "env", Value: slog.StringValue(ctx.Request.Header.Get("env"))},
        {Key: "request_id", Value: slog.StringValue(ctx.Request.Header.Get("X-Request-ID"))},
    }
    return attrs
}
