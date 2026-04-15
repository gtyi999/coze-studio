# CRM Self Review

## 已完成项

- CRM 仍然并入现有 `coze-server` 运行链路，没有新增独立容器或独立服务。
- 后端 handler 未直接写 SQL，CRM 的数据访问仍然下沉在 MySQL 仓储层。
- `backend/domain/crm` 没有反向依赖 `infra`，领域层仍然只依赖实体、仓储接口和领域服务。
- 所有 CRM 列表、详情、更新和删除路径都继续按 `tenant_id + space_id` 做隔离。
- 所有删除动作保持软删除，并补齐了删除前的依赖检查，避免删出孤儿数据。
- 联系人、商机、产品和客户的删除保护已经补齐，覆盖了关键下游引用关系。
- 前端 CRM 页面没有遗留运行时硬编码 mock 数据，mock 仅存在于测试文件中。
- 设计文档里关于软删除字段的描述已经统一为 `is_deleted`，与 DDL 和实现保持一致。
- 已补齐基础测试，包含领域规则、仓储查询和 service 删除守卫等关键路径。

## 未完成项

- 当前 CRM 仍是 phase 1 最小可用范围，没有引入字段级权限、审批流或复杂角色模型。
- 暂未补充完整的 E2E 自动化测试，现阶段以单测和页面级基础测试为主。
- 暂未增加独立报表中心或分析存储，统计仍然直接查 MySQL。

## 技术债

- 当前 `crm_opportunity` 并未持久化 `contact_id`，联系人删除目前只需要对 `follow_record` 做保护；如果二期要补联系人-商机关联，需要先统一 schema、API 和转换层。
- 当前 CRM 仍以 `space_id` 作为实际授权边界，`tenant_id` 延续现有空间映射策略。
- 审计日志只覆盖关键写操作，后续可以补审计查询页和检索能力。
- CRM 的权限控制目前仍然偏 workspace 级，后续如需更细粒度控制，需要补 RBAC 和字段级约束。

## 二期建议

- 补充按钮级和字段级权限控制。
- 增加商机阶段流转、销售过程约束和更多联动校验。
- 增加客户、商机、订单的更多统计维度和可视化报表。
- 增加审计日志查询页，方便排查关键写操作。
- 补充更完整的 E2E 覆盖，尤其是列表、详情、编辑、删除和跨租户拦截。

## 最终建议的提测清单

以下命令默认在对应 module 目录下执行：Go 命令在 `backend/`，前端命令在 `frontend/packages/studio/workspace/entry-adapter/`。

1. 领域规则测试：`go1.24.0 test ./domain/crm/service -run TestCRMService -count=1`
2. 仓储测试：`go1.24.0 test ./infra/crm/impl/mysql -run TestCRMRepository -count=1`
3. 应用层测试：`go1.24.0 test ./application/crm -count=1`
4. 前端基础测试：`corepack pnpm test -- crm-page.test.tsx`
5. 手工冒烟：访问 `/space/<space_id>/crm`，确认 Dashboard、客户、商机和订单列表可正常打开。
6. 租户隔离：用不同 `space_id` 访问同一资源，确认详情、编辑和删除都能被正确拦截。
7. 删除保护：分别验证客户、联系人、商机和产品在存在下游引用时不能删除。
8. 软删除验证：删除后记录不应出现在默认列表中，但历史审计应仍可保留。
9. 示例数据：导入 demo 数据后，CRM 页面和 Dashboard 应能稳定显示基础经营数据。
10. 部署检查：在现有 `coze-server` 启动链路下确认 CRM 路由、MySQL schema 和页面路由都已加载。
