package gormx

import (
	"gorm.io/gorm"
	"time"
)

type JsonLog struct {
	Time    string                 `json:"time"`
	Content string                 `json:"content"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func NewJsonLog(content string, data map[string]interface{}) *JsonLog {
	return &JsonLog{
		Time:    time.Now().Format(time.RFC3339),
		Content: content,
		Data:    data,
	}
}

// AddLogEntry 添加日志条目
func AddLogEntry(column string, content string) map[string]interface{} {
	logEntry := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"content": content,
	}

	return map[string]interface{}{
		column: gorm.Expr("JSON_ARRAY_APPEND(IFNULL(log, JSON_ARRAY()), '$', CAST(? AS JSON))", logEntry),
	}
}

// AddDetailedLogEntry 添加详细日志条目
func AddDetailedLogEntry(content string, data map[string]interface{}) map[string]interface{} {
	logEntry := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"content": content,
		"data":    data,
	}

	return map[string]interface{}{
		"log": gorm.Expr("JSON_ARRAY_APPEND(IFNULL(log, JSON_ARRAY()), '$', CAST(? AS JSON))", logEntry),
	}
}

// AddLogEntryWithLimit 添加日志条目并限制数量
func AddLogEntryWithLimit(column string, content string, limit int) map[string]interface{} {
	logEntry := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"content": content,
	}

	return map[string]interface{}{
		column: gorm.Expr("JSON_ARRAY_APPEND(IF(JSON_LENGTH(IFNULL("+column+", JSON_ARRAY())) >= ?, JSON_REMOVE("+column+", CONCAT('$[0]')), IFNULL("+column+", JSON_ARRAY())), '$', CAST(? AS JSON))", limit, logEntry),
	}
}

// AddLogEntryWithLimitAndData 添加带数据的日志条目并限制数量
func AddLogEntryWithLimitAndData(column string, content string, data map[string]interface{}, limit int) map[string]interface{} {
	logEntry := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"content": content,
		"data":    data,
	}

	return map[string]interface{}{
		column: gorm.Expr("JSON_ARRAY_APPEND(IF(JSON_LENGTH(IFNULL("+column+", JSON_ARRAY())) >= ?, JSON_REMOVE("+column+", CONCAT('$[0]')), IFNULL("+column+", JSON_ARRAY())), '$', CAST(? AS JSON))", limit, logEntry),
	}
}

// JsonAttachment 附件
type JsonAttachment struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
