# CRM API Design

## 概述

本次 CRM 一期 IDL 设计面向 `customer`、`contact`、`opportunity`、`follow_record`、`product`、`sales_order` 六个核心实体，目标是先补齐基础 CRUD 和列表查询能力，并与当前 coze-studio 仓库的 Thrift / Hertz 生成链路保持兼容。

本次实际落点：

- `idl/crm/common.thrift`
- `idl/crm/crm.thrift`
- `idl/api.thrift`

其中：

- `idl/crm/common.thrift` 负责承载 CRM 共享结构，包括分页、过滤、排序、删除结果
- `idl/crm/crm.thrift` 作为 CRM 的统一对外入口，集中声明六类实体的 request / response / service
- `idl/api.thrift` 增加 CRM service 聚合，方便后续继续沿用仓库当前的 master idl 生成方式

## 接口分组

CRM 一期按 bounded context 下的实体进行 service 拆分：

- `CrmCustomerService`
- `CrmContactService`
- `CrmOpportunityService`
- `CrmFollowRecordService`
- `CrmProductService`
- `CrmSalesOrderService`

每个 service 都覆盖 5 类基础动作：

- `CreateXxx`
- `UpdateXxx`
- `DeleteXxx`
- `GetXxxDetail`
- `ListXxx`

对应 API path 统一为：

- `/api/crm/customer/*`
- `/api/crm/contact/*`
- `/api/crm/opportunity/*`
- `/api/crm/follow_record/*`
- `/api/crm/product/*`
- `/api/crm/sales_order/*`

## 命名理由

本次命名遵循两个优先级：

1. 先兼容仓库现有 Thrift / Hertz 代码生成习惯
2. 再保持 CRM 语义清晰、前后端字段稳定

具体取舍如下：

- service 名称采用 `CrmXxxService`
原因：仓库现有 thrift service 基本都使用 `PascalCase`，直接沿用更稳，后续生成 Go 接口和路由代码时不会出现导出名不一致的问题。

- method 名称采用 `CreateCustomer`、`UpdateCustomer`、`GetCustomerDetail` 这类实体前缀形式
原因：虽然从业务语义看，单个 service 内部可以只叫 `create / update / delete / getDetail / list`，但 master idl 聚合到 `idl/api.thrift` 后，若多个 service 都只有通用方法名，后续生成的中间件、handler 入口很容易发生同名冲突。实体前缀可以降低代码生成冲突风险。

- request / response 使用独立 struct
原因：与仓库现有 IDL 组织方式一致，便于后端 handler/model 生成、前端 schema 生成和后续接口演进。

- 字段命名尽量使用 `snake_case`
原因：与仓库当前 thrift 字段风格、数据库字段命名和 API 入参风格更一致。

- 新增业务字段统一使用 `optional`
原因：满足一期阶段的渐进式开发需求，给后端校验和后续字段扩展保留空间。

## 通用约定

### Request / Response 基类

所有 request 统一带：

- `255: optional base.Base Base`

所有 response 统一带：

- `253: optional i64 code`
- `254: optional string msg`
- `255: optional base.BaseResp BaseResp`

### 分页查询

所有 list request 统一包含以下字段：

- `page_no`
- `page_size`
- `keyword`
- `filters`
- `sorts`

其中共享结构定义在 `idl/crm/common.thrift`：

- `CrmFilter`
- `CrmSort`
- `PageInfo`

### 数值字段约定

以下字段在 IDL 中采用 `string`：

- `amount`
- `unit_price`
- `quantity`

原因：

- 对应 MySQL `decimal` 字段
- 避免后续生成到前端 TS 后发生精度丢失

以下字段采用 `i64 + js_conv`：

- 主键 ID
- `tenant_id`
- `space_id`
- 审计字段
- 毫秒时间戳字段

这样可以继续贴合仓库已有的 JS 大整数转换方式。

## 与当前仓库的关系

本次设计没有引入新的 IDL 框架或新的接口组织方式，而是沿用当前仓库已有约定：

- `idl/api.thrift` 作为 master idl
- `api.category` / `api.gen_path` / `agw.preserve_base` 继续沿用现有注解方式
- CRM 继续放在 `idl/crm` 目录下，不单独新建额外协议层

这意味着下一阶段可以直接在现有生成链路上继续推进：

- 后端 thrift model 生成
- Hertz router/handler 生成或增量更新
- 前端 `frontend/packages/arch/api-schema` 重新生成 CRM schema

## 后续建议

建议按下面顺序推进：

1. 先基于这份 IDL 重新生成后端 model / router 骨架，确认生成结果没有符号冲突。
2. 再同步生成前端 `api-schema/crm`，替换当前仓库内较早期的 CRM 手写 schema。
3. 最后让 `backend/api/handler/coze/crm_service.go`、`backend/api/router/coze/crm_manual.go` 与新 IDL 对齐，收敛掉旧版 `get` / `description` / `customer_type` 等试验字段。
