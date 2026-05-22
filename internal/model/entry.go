package model

import "time"

// Entry 一条日志记录
type Entry struct {
	Time    time.Time `json:"time"`
	Content string    `json:"content"`
}

// Plan 一条工作计划
type Plan struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// SearchResult 搜索结果
type SearchResult struct {
	Date    string // 文件名，如 "2026-05-22"
	Time    string // 时间，如 "14:09"
	Content string // 匹配的行内容
}
