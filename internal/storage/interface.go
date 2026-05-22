package storage

import (
	"time"

	"github.com/firstlfq/rec/internal/model"
)

// Storage 存储接口 —— 所有数据操作通过此接口
// 目前有文件系统实现 (FSStorage)，以后可切换为 SQLite 等
type Storage interface {
	// === 日志 ===

	// AppendEntry 追加一条日志到指定日期
	AppendEntry(date time.Time, content string) error

	// GetEntries 获取指定日期的全部日志
	GetEntries(date time.Time) ([]model.Entry, error)

	// === 搜索 ===

	// SearchEntries 在所有日志中搜索关键词
	SearchEntries(keyword string) ([]model.SearchResult, error)

	// === 计划 ===

	// AddPlan 添加一条计划，返回计划编号
	AddPlan(content string) (int, error)

	// GetPlans 获取所有计划
	GetPlans() ([]model.Plan, error)

	// CompletePlan 标记计划为已完成
	CompletePlan(id int) error
}
