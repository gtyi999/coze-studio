# CRM Agent Setup

本地仓库现在提供了一个可复用脚本，用来在 Coze 里编排 `CRM Agent`，并把它依赖的两张数据库表一起准备好。

运行方式:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup\init_crm_agent_bot.ps1
```

脚本会完成这些动作:

1. 复用或初始化 `crm-agent-owner@example.com` 的测试资源
2. 创建或复用名为 `CRM Agent` 的 draft bot
3. 写入 bot 的 prompt、描述和开场白
4. 创建或复用 `客户信息表`
5. 创建或复用 `销售业绩表`
6. 写入示例数据
7. 把两张表绑定到 `CRM Agent`
8. 回读 bot 配置和数据库记录
9. 将真实结果写入 `output/crm-agent-bot.json`

当前脚本写入的核心示例数据:

- `客户信息表`: `客户状态=有效`、`客户数量=1277`、`统计时间=当前`
- `销售业绩表`: `张三=947818`、`李四=825120`、`王五=801306`，统计周期均为 `本季度`

对应的目标问答:

- `我现在的客户数有多少？`
  `当前有效客户共 1277 家。`
- `哪个销售的业绩最好？`
  `销售张三业绩最高，本季度累计 947,818 元。`

## 重要说明

`CRM Agent` 的编排本身已经完成，但当前这套本地 `coze-server` 还没有可用的大模型 endpoint。没有模型 endpoint 时，bot 的真实聊天会在运行期失败，日志里会出现类似下面的错误:

```text
can't fetch endpoint sts token without endpoint
```

也就是说:

- bot 草稿、prompt、数据库绑定、示例数据都已经正确落地
- 真实对话要跑通，还需要先补一个可用的模型配置并重启服务

脚本输出的 `output/crm-agent-bot.json` 会包含:

- `bot_id`
- 绑定数据库的 `online_id` / `draft_id`
- 已回读到的示例数据
- 当前环境是否具备 live chat 条件
