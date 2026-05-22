# rec — 个人日志记忆工具

> 记性不好？让 rec 帮你记住每天的工作、命令和计划。

## 简介

`rec` 是一个极简的终端日志工具。每次在终端里做完一件事，随手记一笔。时间久了，你就有了一本可以随时搜索的「工作日记」。

## 安装

### 方式一：下载预编译二进制

```bash
# 从 GitHub Releases 下载最新版
curl -sL https://github.com/firstlfq/rec/releases/latest/download/rec-linux-amd64 -o /usr/local/bin/rec
chmod +x /usr/local/bin/rec
```

### 方式二：从源码编译

```bash
git clone https://github.com/firstlfq/rec.git
cd rec
go build -o rec .
# 移动到 PATH
sudo mv rec /usr/local/bin/
```

## 快速上手

```bash
# 记录一条日志
rec 今天在改串口驱动的 buffer 管理
rec 三点跟老王开会，协议栈 v3 接口定下来了
rec 学到的新命令：ss -tlnp 查看监听端口

# 查看今天的日志
rec today

# 搜索历史记录
rec 搜 串口

# 添加工作计划
rec plan 完成 buffer 管理改写
rec plan 给老王发协议栈文档

# 查看工作计划
rec plan

# 查看所有命令
rec help
```

## 数据存储

所有数据保存在 `~/.rec/` 目录下：

```
~/.rec/
├── 2026-05-22.md    # 每天一个 Markdown 文件
├── 2026-05-23.md
├── plans.md          # 工作计划
└── ...
```

纯文本格式，你可以用 `cat`、`grep`、`vim` 直接查看和编辑。不会被锁在任何私有格式里。

## 命令参考

| 命令 | 说明 |
|------|------|
| `rec <消息>` | 记录一条带时间戳的日志 |
| `rec today` | 查看今天的日志 |
| `rec yesterday` | 查看昨天的日志 |
| `rec 搜 <关键词>` | 搜索所有历史记录 |
| `rec plan <任务>` | 添加工作计划 |
| `rec plan` | 查看所有计划 |
| `rec done <编号>` | 标记任务完成 |
| `rec help` | 显示帮助信息 |

## 未来计划

- [ ] 模糊搜索（fuzzy find）
- [ ] 标签分类
- [ ] 按日期范围搜索
- [ ] TUI 交互界面
- [ ] AI 智能搜索（自然语言搜日志）
- [ ] AI 自动生成周报

## 许可证

MIT
