package easyginlogger

import (
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCustomLogger_Log(t *testing.T) {
	fileLog, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	hdl := slog.NewJSONHandler(fileLog, &slog.HandlerOptions{AddSource: true, Level: LevelTrace, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			levelLabel, exists := LevelNames[level]
			if !exists {
				levelLabel = level.String()
			}

			a.Value = slog.StringValue(levelLabel)
		}

		return a
	},})
	type fields struct {
		Logger *slog.Logger
	}
	type args struct {
		ctx   *gin.Context
		msg   string
		level slog.Level
		attrs []slog.Attr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1", 
			fields: fields{Logger: slog.New(hdl)},
			args: args{ctx: nil, msg: "test", level: LevelTrace, attrs: []slog.Attr{
				{Key: "test", Value: slog.StringValue("trace")},
			}},
		},{
			name: "2", 
			fields: fields{Logger: slog.New(hdl)},
			args: args{ctx: nil, msg: "test", level: slog.LevelDebug, attrs: []slog.Attr{
				{Key: "test", Value: slog.StringValue("debug")},
			}},
		},{
			name: "3", 
			fields: fields{Logger: slog.New(hdl)},
			args: args{ctx: nil, msg: "test", level: slog.LevelInfo, attrs: []slog.Attr{
				{Key: "test", Value: slog.StringValue("info")},
			}},
		},{
			name: "4", 
			fields: fields{Logger: slog.New(hdl)},
			args: args{ctx: nil, msg: "test", level: slog.LevelWarn, attrs: []slog.Attr{
				{Key: "test", Value: slog.StringValue("warning")},
			}},
		},{
			name: "5", 
			fields: fields{Logger: slog.New(hdl)},
			args: args{ctx: nil, msg: "test", level: slog.LevelError, attrs: []slog.Attr{
				{Key: "test", Value: slog.StringValue("error")},
			}},
		},
	}

	i := 0
	for i < 100 {
		i++
		r := rand.New(rand.NewSource(time.Now().UnixMicro())) // 创建新的随机数生成器
		num := r.Intn(50000)
		// fmt.Println(num % 4)
		tests = append(tests, tests[num % 4])
	}
	var wg sync.WaitGroup
	wg.Add(105)

	for _, tt := range tests {
		go t.Run(tt.name, func(t *testing.T) {
			cl := &CustomLogger{
				Logger: tt.fields.Logger,
			}
			cl.Log(tt.args.ctx, tt.args.msg, tt.args.level, tt.args.attrs)
			wg.Done()
		})
	}
	wg.Wait()
}
