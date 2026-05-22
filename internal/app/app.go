package app

import (
	"fmt"
	"os"

	"github.com/firstlfq/rec/internal/storage"
)

// Command 命令接口 —— 每个子命令实现此接口
type Command interface {
	// Name 返回命令名称
	Name() string
	// Aliases 返回别名列表
	Aliases() []string
	// Execute 执行命令
	Execute(args []string, store storage.Storage) error
}

// App 是命令路由的主入口
type App struct {
	commands []Command
	store    storage.Storage
	name     string // 程序名，如 "rec"
}

// New 创建应用实例
func New(name string, store storage.Storage) *App {
	return &App{
		name: name,
		store: store,
	}
}

// Register 注册一个命令
func (a *App) Register(cmd Command) {
	a.commands = append(a.commands, cmd)
}

// Run 根据 os.Args 执行对应的命令
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.showHelp()
		return nil
	}

	cmdName := args[1]

	// 匹配命令
	for _, cmd := range a.commands {
		if cmd.Name() == cmdName {
			return cmd.Execute(args[2:], a.store)
		}
		for _, alias := range cmd.Aliases() {
			if alias == cmdName {
				return cmd.Execute(args[2:], a.store)
			}
		}
	}

	// 没有匹配的子命令：看作日志内容
	// 例如 "rec 今天天气不错" → 直接记录
	for _, cmd := range a.commands {
		if cmd.Name() == "log" {
			return cmd.Execute(args[1:], a.store)
		}
	}

	return fmt.Errorf("未知命令: %s", cmdName)
}

func (a *App) showHelp() {
	fmt.Printf(`%s — 个人日志记忆工具

用法:
  %s <消息>            记录一条日志
  %s today             查看今天的日志
  %s yesterday         查看昨天的日志
  %s 搜 <关键词>        搜索所有日志
  %s plan <任务>        添加工作计划
  %s plan              查看工作计划
  %s done <编号>        标记任务完成
  %s help              显示帮助

数据保存在 ~/.%s/ 目录下，纯文本格式。
`,
		a.name, a.name, a.name, a.name, a.name, a.name, a.name, a.name, a.name)
	os.Exit(0)
}
