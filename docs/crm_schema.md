# CRM Schema

## 概述

本次新增的 CRM 一期表结构全部采用 MySQL 作为主存，不引入 MongoDB、Elasticsearch 或其他新的主存储。

本次覆盖六张核心业务表：

- `crm_customer`
- `crm_contact`
- `crm_opportunity`
- `crm_follow_record`
- `crm_product`
- `crm_sales_order`

统一约束：

- 主键均为 `id bigint`
- 每张表均包含 `tenant_id`
- 每张表均包含 `created_by` / `updated_by`
- 每张表均包含 `created_at` / `updated_at`
- 每张表均包含 `is_deleted`
- 表级主存均为 MySQL

补充说明：

- 为了贴近当前仓库以工作空间为主的访问模型，DDL 中额外保留了 `space_id`
- `created_at` / `updated_at` 采用毫秒时间戳，便于与仓库现有后端风格保持一致
- `is_deleted` 采用逻辑删除标记，默认 `0`

## DDL 文件位置

规范 DDL 源文件：

- [crm_phase1.sql](/C:/GoProject/src/coze-studio/backend/types/ddl/crm_phase1.sql)

初始化入口同步位置：

- `docker/volumes/mysql/schema.sql`
- `docker/atlas/opencoze_latest_schema.hcl`

初始化脚本：

- [apply_crm_schema.sh](/C:/GoProject/src/coze-studio/scripts/setup/apply_crm_schema.sh)

## 表说明

### 1. crm_customer

业务含义：

- 存储 CRM 客户主数据
- 作为联系人、商机、跟进记录、销售订单的上游主体

核心字段：

- `customer_name`: 客户名称
- `customer_code`: 客户编码
- `industry`: 所属行业
- `level`: 客户等级
- `owner_user_id` / `owner_user_name`: 客户负责人
- `status`: 客户状态
- `mobile` / `email` / `address`: 联系信息
- `remark`: 备注

核心索引：

- `tenant_id + is_deleted`
- `owner_user_id`
- `status`
- `created_at`
- `customer_code`
- `customer_name`

### 2. crm_contact

业务含义：

- 存储客户下的联系人信息
- 支持一个客户维护多个联系人

核心字段：

- `customer_id`: 所属客户
- `contact_name`: 联系人姓名
- `mobile` / `email`: 联系方式
- `title`: 职务
- `is_primary`: 是否主联系人
- `status`: 联系人状态
- `remark`: 备注

核心索引：

- `tenant_id + is_deleted`
- `customer_id`
- `status`
- `created_at`
- `contact_name`

### 3. crm_opportunity

业务含义：

- 存储销售商机
- 与客户形成一对多关系

核心字段：

- `customer_id`: 所属客户
- `opportunity_name`: 商机名称
- `stage`: 商机阶段
- `amount`: 预计金额
- `expected_close_date`: 预计成交日期
- `owner_user_id` / `owner_user_name`: 商机负责人
- `status`: 商机状态
- `remark`: 备注

核心索引：

- `tenant_id + is_deleted`
- `customer_id`
- `owner_user_id`
- `status`
- `created_at`
- `stage`

### 4. crm_follow_record

业务含义：

- 存储客户或联系人的跟进记录
- 可用于沉淀电话、拜访、微信、会议等跟进动作

核心字段：

- `customer_id`: 所属客户
- `contact_id`: 关联联系人
- `follow_type`: 跟进类型
- `content`: 跟进内容
- `next_follow_time`: 下次跟进时间
- `owner_user_id` / `owner_user_name`: 跟进负责人
- `status`: 跟进记录状态

核心索引：

- `tenant_id + is_deleted`
- `customer_id`
- `contact_id`
- `owner_user_id`
- `status`
- `created_at`
- `next_follow_time`

### 5. crm_product

业务含义：

- 存储 CRM 销售产品主数据
- 供销售订单引用

核心字段：

- `product_name`: 产品名称
- `product_code`: 产品编码
- `category`: 产品分类
- `unit_price`: 标准单价
- `status`: 产品状态
- `remark`: 备注

核心索引：

- `tenant_id + is_deleted`
- `status`
- `created_at`
- `product_code`
- `product_name`
- `category`

### 6. crm_sales_order

业务含义：

- 存储销售订单事实数据
- 与客户、商机、产品建立关联

核心字段：

- `customer_id`: 所属客户
- `opportunity_id`: 来源商机
- `product_id`: 关联产品
- `product_name`: 产品名称快照
- `sales_user_id` / `sales_user_name`: 销售人员
- `quantity`: 数量
- `amount`: 订单金额
- `order_date`: 下单日期
- `status`: 订单状态
- `remark`: 备注

核心索引：

- `tenant_id + is_deleted`
- `customer_id`
- `opportunity_id`
- `product_id`
- `sales_user_id`
- `status`
- `created_at`
- `order_date`

## 初始化执行说明

### 方式一：对已有 MySQL 容器补 CRM 表

仓库内已提供脚本：

```bash
bash scripts/setup/apply_crm_schema.sh
```

Windows 环境下建议在 Git Bash 或 WSL 中执行该脚本。

脚本默认行为：

- 读取 `backend/types/ddl/crm_phase1.sql`
- 连接运行中的 `coze-mysql`
- 将 CRM 六张表的 DDL 应用到 `opencoze` 库

可覆盖环境变量：

- `MYSQL_CONTAINER`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DATABASE`

### 方式二：通过 Docker 初始化新库

已将 CRM 表结构同步到：

- `docker/volumes/mysql/schema.sql`
- `docker/atlas/opencoze_latest_schema.hcl`

对于全新 MySQL 数据目录，可按仓库默认方式启动：

```bash
docker compose -f docker/docker-compose.yml up -d mysql
```

如果是已存在旧数据目录，建议先确认是否需要清理旧库或旧卷后再重新初始化。

### 方式三：手动执行 DDL

```bash
docker exec -i coze-mysql mysql -uroot -proot opencoze < backend/types/ddl/crm_phase1.sql
```

如果本地 `.env` 中账号密码不同，请替换为实际连接信息。
