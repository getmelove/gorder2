package logging

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// 定义了日志字段的键名常量
const (
	Method   = "method"   // 方法名
	Args     = "args"     // 参数
	Cost     = "cost_ms"  // 耗时（毫秒）
	Response = "response" // 响应结果
	Error    = "err"      // 错误信息
)

// WhenMySQL 用于记录 MySQL 操作的日志。
// 它接收上下文、方法名和参数，并返回初始日志字段和一个回调函数。
// 回调函数应在操作完成后调用，以记录响应结果、错误信息和执行耗时。
func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
	// 初始化日志字段，记录方法名和参数
	fields := logrus.Fields{
		Method: method,
		Args:   formatMySQLArgs(args),
	}
	// 记录开始时间
	start := time.Now()
	// 返回回调函数
	return fields, func(resp any, err *error) {
		level, msg := logrus.InfoLevel, "mysql_success"
		// 计算并记录耗时
		fields[Cost] = time.Since(start).Milliseconds()
		fields[Response] = resp

		// 如果发生错误，记录错误信息并将日志级别设置为 Error
		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "mysql_error"
			fields[Error] = (*err).Error()
		}

		// 使用上下文和字段记录最终日志
		logrus.WithContext(ctx).WithFields(fields).Logf(level, "%s", msg)
	}
}

// formatMySQLArgs 将多个参数格式化为字符串，参数之间用 "||" 分隔
func formatMySQLArgs(args []any) string {
	var item []string
	for _, arg := range args {
		item = append(item, formatMySQLArg(arg))
	}
	return strings.Join(item, "||")
}

// formatMySQLArg 将单个参数格式化为 JSON 字符串
func formatMySQLArg(arg any) string {
	switch v := arg.(type) {
	default:
		// 尝试将参数序列化为 JSON
		bytes, err := json.Marshal(v)
		if err != nil {
			return "unsupported type in formatMySQLArg||err=" + err.Error()
		}
		return string(bytes)
	}
}
