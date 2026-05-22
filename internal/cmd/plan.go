package cmd

import (
	"fmt"
	"strings"

	"github.com/firstlfq/rec/internal/storage"
)

// PlanCmd 计划：rec plan / rec plan <任务>
type PlanCmd struct{}

func (c *PlanCmd) Name() string      { return "plan" }
func (c *PlanCmd) Aliases() []string { return []string{"p"} }
func (c *PlanCmd) Execute(args []string, store storage.Storage) error {
	if len(args) == 0 {
		// rec plan：查看计划
		return showPlans(store)
	}
	// rec plan <任务>：添加计划
	content := strings.Join(args, " ")
	id, err := store.AddPlan(content)
	if err != nil {
		return fmt.Errorf("添加计划失败: %w", err)
	}
	fmt.Printf("✅ 已添加任务 #%d: %s\n", id, content)
	return nil
}

func showPlans(store storage.Storage) error {
	plans, err := store.GetPlans()
	if err != nil {
		return err
	}

	fmt.Println("📋 工作计划")
	fmt.Println(strings.Repeat("─", 40))
	if len(plans) == 0 {
		fmt.Println("  (暂无计划)")
		return nil
	}

	hasActive := false
	for _, p := range plans {
		if !p.Done {
			hasActive = true
			fmt.Printf("  🔴 [%d] %s\n", p.ID, p.Content)
		}
	}
	if !hasActive {
		fmt.Println("  所有任务已完成 🎉")
	}
	fmt.Println()

	// 显示已完成任务
	hasDone := false
	for _, p := range plans {
		if p.Done {
			if !hasDone {
				fmt.Println("  ✅ 已完成:")
				hasDone = true
			}
			fmt.Printf("    [%d] %s\n", p.ID, p.Content)
		}
	}
	return nil
}

// DoneCmd 完成任务：rec done <编号>
type DoneCmd struct{}

func (c *DoneCmd) Name() string      { return "done" }
func (c *DoneCmd) Aliases() []string { return []string{"d", "完成"} }
func (c *DoneCmd) Execute(args []string, store storage.Storage) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: rec done <任务编号>")
	}

	var id int
	if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
		return fmt.Errorf("无效的编号: %s", args[0])
	}

	if err := store.CompletePlan(id); err != nil {
		return fmt.Errorf("操作失败: %w", err)
	}

	fmt.Printf("✅ 任务 #%d 已完成\n", id)
	return nil
}
