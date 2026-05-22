package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/firstlfq/rec/internal/storage"
)

// HelpCmd 显示帮助
type HelpCmd struct{}

func (c *HelpCmd) Name() string          { return "help" }
func (c *HelpCmd) Aliases() []string     { return []string{"--help", "-h"} }
func (c *HelpCmd) Execute(args []string, store storage.Storage) error {
	_ = args
	_ = store
	fmt.Println(`rec — 个人日志记忆工具

用法:
  rec <消息>            记录一条日志
  rec today             查看今天的日志
  rec yesterday         查看昨天的日志
  rec 搜 <关键词>        搜索所有日志
  rec plan <任务>        添加工作计划
  rec plan              查看工作计划
  rec done <编号>        标记任务完成
  rec help              显示帮助

数据保存在 ~/.rec/ 目录下，纯文本格式。`)
	os.Exit(0)
	return nil
}

// LogCmd 记录日志：rec <消息>
type LogCmd struct{}

func (c *LogCmd) Name() string      { return "log" }
func (c *LogCmd) Aliases() []string { return []string{} }
func (c *LogCmd) Execute(args []string, store storage.Storage) error {
	if len(args) < 1 {
		return fmt.Errorf("请输入要记录的内容")
	}
	content := strings.Join(args, " ")
	if err := store.AppendEntry(time.Now(), content); err != nil {
		return fmt.Errorf("记录失败: %w", err)
	}
	fmt.Printf("✅ 已记录 (%s)\n", time.Now().Format("15:04"))
	return nil
}

// TodayCmd 查看今天日志：rec today
type TodayCmd struct{}

func (c *TodayCmd) Name() string      { return "today" }
func (c *TodayCmd) Aliases() []string { return []string{"t"} }
func (c *TodayCmd) Execute(args []string, store storage.Storage) error {
	_ = args
	entries, err := store.GetEntries(time.Now())
	if err != nil {
		return err
	}
	date := time.Now().Format("2006-01-02")
	fmt.Printf("📋 %s\n", date)
	fmt.Println(strings.Repeat("─", 40))
	if len(entries) == 0 {
		fmt.Println("  (暂无记录)")
		return nil
	}
	for _, e := range entries {
		fmt.Printf("  %s  %s\n", e.Time.Format("15:04"), e.Content)
	}
	return nil
}

// YesterdayCmd 查看昨天日志：rec yesterday
type YesterdayCmd struct{}

func (c *YesterdayCmd) Name() string      { return "yesterday" }
func (c *YesterdayCmd) Aliases() []string { return []string{"y"} }
func (c *YesterdayCmd) Execute(args []string, store storage.Storage) error {
	_ = args
	date := time.Now().AddDate(0, 0, -1)
	entries, err := store.GetEntries(date)
	if err != nil {
		return err
	}
	fmt.Printf("📋 %s\n", date.Format("2006-01-02"))
	fmt.Println(strings.Repeat("─", 40))
	if len(entries) == 0 {
		fmt.Println("  (暂无记录)")
		return nil
	}
	for _, e := range entries {
		fmt.Printf("  %s  %s\n", e.Time.Format("15:04"), e.Content)
	}
	return nil
}
