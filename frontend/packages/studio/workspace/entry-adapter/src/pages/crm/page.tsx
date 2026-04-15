/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { type FC, type FormEvent, type ReactNode, useState } from 'react';

import { useRequest } from 'ahooks';
import {
  Button,
  Empty,
  Layout,
  Spin,
  Typography,
} from '@coze-arch/coze-design';
import { crm as crmApi } from '@coze-studio/api-schema';
import type {
  CustomerData,
  DashboardOverviewData,
  OpportunityData,
  SalesOrderData,
} from '@coze-studio/api-schema/crm';

import { SimpleLineChart } from './components/simple-line-chart';
import { CRMStatusTag } from './components/status-tag';

interface CRMManagePageProps {
  spaceId: string;
}

const CARD_CLASS_NAME =
  'rounded-[20px] border border-solid coz-stroke-primary coz-bg-max p-[20px] shadow-[0_16px_40px_rgba(15,23,42,0.04)]';
const INPUT_CLASS_NAME =
  'w-full rounded-[12px] border border-solid coz-stroke-primary bg-white px-[12px] py-[10px] text-[14px] outline-none transition focus:border-[#2563eb]';

interface CreateCustomerFormState {
  customerName: string;
  ownerUserName: string;
  industry: string;
  remark: string;
}

const METRIC_CARDS = [
  {
    key: 'customer_total',
    title: 'Customers',
    background:
      'bg-[linear-gradient(135deg,#f8fafc_0%,#ffffff_52%,#dbeafe_100%)]',
  },
  {
    key: 'new_customers_this_month',
    title: 'New Customers This Month',
    background:
      'bg-[linear-gradient(135deg,#fff7ed_0%,#ffffff_52%,#fde68a_100%)]',
  },
  {
    key: 'opportunity_total_amount',
    title: 'Opportunity Amount',
    background:
      'bg-[linear-gradient(135deg,#f0fdf4_0%,#ffffff_52%,#bbf7d0_100%)]',
  },
  {
    key: 'new_opportunities_this_month',
    title: 'New Opportunities This Month',
    background:
      'bg-[linear-gradient(135deg,#eff6ff_0%,#ffffff_52%,#bfdbfe_100%)]',
  },
  {
    key: 'sales_order_total_amount',
    title: 'Sales Order Amount',
    background:
      'bg-[linear-gradient(135deg,#fefce8_0%,#ffffff_52%,#fde68a_100%)]',
  },
] as const;

export const CRMManagePage: FC<CRMManagePageProps> = ({ spaceId }) => {
  const [showCreateCustomerForm, setShowCreateCustomerForm] = useState(false);
  const [createCustomerForm, setCreateCustomerForm] = useState<CreateCustomerFormState>(getInitialCreateCustomerForm());
  const [createCustomerLoading, setCreateCustomerLoading] = useState(false);
  const [createCustomerError, setCreateCustomerError] = useState('');

  const overviewRequest = useRequest(
    async () => {
      try {
        const resp = await crmApi.GetDashboardOverview(
          { space_id: spaceId },
          { __disableErrorToast: true },
        );
        return resp.data;
      } catch {
        return undefined;
      }
    },
    {
      refreshDeps: [spaceId],
    },
  );

  const tablesRequest = useRequest(
    async () => {
      const [customers, opportunities, salesOrders] = await Promise.all([
        loadCustomers(spaceId),
        loadOpportunities(spaceId),
        loadSalesOrders(spaceId),
      ]);

      return {
        customers,
        opportunities,
        salesOrders,
      };
    },
    {
      refreshDeps: [spaceId],
    },
  );

  const overview = overviewRequest.data;
  const tables = tablesRequest.data;
  const loading = overviewRequest.loading || tablesRequest.loading;

  const onCreateCustomerSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!createCustomerForm.customerName.trim()) {
      setCreateCustomerError('Customer name is required.');
      return;
    }

    setCreateCustomerLoading(true);
    setCreateCustomerError('');
    try {
      await crmApi.CreateCustomer(
        {
          space_id: spaceId,
          customer_name: createCustomerForm.customerName.trim(),
          owner_user_name: createCustomerForm.ownerUserName.trim(),
          industry: createCustomerForm.industry.trim(),
          remark: createCustomerForm.remark.trim(),
          status: 'active',
        },
        { __disableErrorToast: true },
      );
      setCreateCustomerForm(getInitialCreateCustomerForm());
      setShowCreateCustomerForm(false);
      overviewRequest.refresh();
      tablesRequest.refresh();
    } catch {
      setCreateCustomerError('Create customer failed. Please try again.');
    } finally {
      setCreateCustomerLoading(false);
    }
  };

  return (
    <Layout className="h-full overflow-hidden">
      <Layout.Header className="pb-0">
        <div className="flex flex-wrap items-start justify-between gap-[16px]">
          <div>
            <Typography.Title heading={4} className="!mb-[8px]">
              CRM Dashboard
            </Typography.Title>
            <Typography.Paragraph className="!mb-0 max-w-[720px] coz-fg-secondary">
              Monitor baseline customer, opportunity and sales order metrics,
              all queried directly from MySQL CRM tables.
            </Typography.Paragraph>
          </div>
          <div className="flex flex-wrap gap-[12px]">
            <Button
              color="secondary"
              onClick={() => {
                overviewRequest.refresh();
                tablesRequest.refresh();
              }}
            >
              Refresh
            </Button>
            <Button
              onClick={() => {
                setCreateCustomerError('');
                setShowCreateCustomerForm(value => !value);
              }}
            >
              {showCreateCustomerForm ? 'Hide Quick Create' : 'Quick Create Customer'}
            </Button>
          </div>
        </div>
      </Layout.Header>
      <Layout.Content className="!h-auto !min-h-0 !flex-1 overflow-auto pb-[24px]">
        <Spin spinning={loading}>
          <div className="mt-[12px] flex flex-col gap-[16px]">
            {showCreateCustomerForm ? (
              <section className={`${CARD_CLASS_NAME} bg-[linear-gradient(135deg,#eff6ff_0%,#ffffff_56%,#dbeafe_100%)]`}>
                <div className="mb-[16px]">
                  <Typography.Title heading={6} className="!mb-[4px]">
                    Quick Create Customer
                  </Typography.Title>
                  <Typography.Paragraph className="!mb-0 coz-fg-secondary">
                    Create a minimal customer record without leaving the dashboard.
                  </Typography.Paragraph>
                </div>
                <form className="grid grid-cols-1 gap-[16px] md:grid-cols-2" onSubmit={onCreateCustomerSubmit}>
                  <FormField label="Customer Name" required>
                    <input
                      aria-label="Customer Name"
                      className={INPUT_CLASS_NAME}
                      value={createCustomerForm.customerName}
                      onChange={event => {
                        setCreateCustomerForm(value => ({
                          ...value,
                          customerName: event.target.value,
                        }));
                      }}
                    />
                  </FormField>
                  <FormField label="Owner">
                    <input
                      aria-label="Owner"
                      className={INPUT_CLASS_NAME}
                      value={createCustomerForm.ownerUserName}
                      onChange={event => {
                        setCreateCustomerForm(value => ({
                          ...value,
                          ownerUserName: event.target.value,
                        }));
                      }}
                    />
                  </FormField>
                  <FormField label="Industry">
                    <input
                      aria-label="Industry"
                      className={INPUT_CLASS_NAME}
                      value={createCustomerForm.industry}
                      onChange={event => {
                        setCreateCustomerForm(value => ({
                          ...value,
                          industry: event.target.value,
                        }));
                      }}
                    />
                  </FormField>
                  <FormField label="Remark" className="md:col-span-2">
                    <textarea
                      aria-label="Remark"
                      className={`${INPUT_CLASS_NAME} min-h-[96px] resize-y`}
                      value={createCustomerForm.remark}
                      onChange={event => {
                        setCreateCustomerForm(value => ({
                          ...value,
                          remark: event.target.value,
                        }));
                      }}
                    />
                  </FormField>
                  <div className="md:col-span-2 flex flex-wrap items-center justify-between gap-[12px]">
                    <Typography.Text className="coz-fg-secondary text-[12px]">
                      {createCustomerError || 'Only customer name is required for the minimal demo flow.'}
                    </Typography.Text>
                    <div className="flex gap-[12px]">
                      <Button
                        color="secondary"
                        htmlType="button"
                        onClick={() => {
                          setCreateCustomerError('');
                          setCreateCustomerForm(getInitialCreateCustomerForm());
                          setShowCreateCustomerForm(false);
                        }}
                      >
                        Cancel
                      </Button>
                      <Button htmlType="submit" loading={createCustomerLoading}>
                        Create Customer
                      </Button>
                    </div>
                  </div>
                </form>
              </section>
            ) : null}

            <section className="grid grid-cols-1 gap-[16px] md:grid-cols-2 xl:grid-cols-5">
              {METRIC_CARDS.map(item => (
                <MetricCard
                  key={item.key}
                  title={item.title}
                  background={item.background}
                  value={getMetricValue(overview, item.key)}
                />
              ))}
            </section>

            <section className={`${CARD_CLASS_NAME} min-h-[360px]`}>
              <div className="mb-[16px] flex flex-wrap items-end justify-between gap-[12px]">
                <div>
                  <Typography.Title heading={6} className="!mb-[4px]">
                    Recent 30-Day Order Trend
                  </Typography.Title>
                  <Typography.Paragraph className="!mb-0 coz-fg-secondary">
                    The line is drawn by order amount and scoped to the current
                    tenant workspace.
                  </Typography.Paragraph>
                </div>
                <Typography.Text className="coz-fg-secondary text-[12px]">
                  {overview?.recent_order_trend?.length ?? 0} daily points
                </Typography.Text>
              </div>
              <SimpleLineChart data={overview?.recent_order_trend ?? []} />
            </section>

            <section className="grid grid-cols-1 gap-[16px] 2xl:grid-cols-3">
              <DashboardTableCard
                title="Customer Snapshot"
                description="Most recently updated customers"
                headers={['Customer', 'Owner', 'Status', 'Updated']}
                rows={(tables?.customers ?? []).map(item => [
                  item.customer_name || '-',
                  item.owner_user_name || '-',
                  (
                    <CRMStatusTag
                      key={`customer-status-${item.customer_id}`}
                      value={item.status}
                    />
                  ),
                  formatDateTime(item.updated_at),
                ])}
              />
              <DashboardTableCard
                title="Opportunity Snapshot"
                description="Most recently updated opportunities"
                headers={['Opportunity', 'Stage', 'Amount', 'Status']}
                rows={(tables?.opportunities ?? []).map(item => [
                  item.opportunity_name || '-',
                  item.stage || '-',
                  formatAmount(item.amount),
                  (
                    <CRMStatusTag
                      key={`opportunity-status-${item.opportunity_id}`}
                      value={item.status}
                    />
                  ),
                ])}
              />
              <DashboardTableCard
                title="Order Snapshot"
                description="Most recently updated sales orders"
                headers={['Product', 'Amount', 'Order Date', 'Status']}
                rows={(tables?.salesOrders ?? []).map(item => [
                  item.product_name || '-',
                  formatAmount(item.amount),
                  item.order_date || '-',
                  (
                    <CRMStatusTag
                      key={`order-status-${item.sales_order_id}`}
                      value={item.status}
                    />
                  ),
                ])}
              />
            </section>
          </div>
        </Spin>
      </Layout.Content>
    </Layout>
  );
};

const MetricCard: FC<{
  title: string;
  value: string;
  background: string;
}> = ({ title, value, background }) => (
  <div className={`${CARD_CLASS_NAME} ${background}`}>
    <Typography.Text className="coz-fg-secondary text-[12px]">
      {title}
    </Typography.Text>
    <div className="mt-[18px] flex items-end justify-between gap-[12px]">
      <Typography.Title heading={3} className="!mb-0">
        {value}
      </Typography.Title>
      <div className="h-[48px] w-[48px] rounded-full bg-[rgba(255,255,255,0.72)]" />
    </div>
  </div>
);

const DashboardTableCard: FC<{
  title: string;
  description: string;
  headers: string[];
  rows: Array<Array<ReactNode>>;
}> = ({ title, description, headers, rows }) => (
  <div className={`${CARD_CLASS_NAME} min-h-[320px]`}>
    <div className="mb-[16px]">
      <Typography.Title heading={6} className="!mb-[4px]">
        {title}
      </Typography.Title>
      <Typography.Paragraph className="!mb-0 coz-fg-secondary">
        {description}
      </Typography.Paragraph>
    </div>
    {rows.length ? (
      <div className="overflow-hidden rounded-[14px] border border-solid coz-stroke-primary">
        <table className="w-full border-collapse text-left text-[12px]">
          <thead className="bg-[rgba(15,23,42,0.03)]">
            <tr>
              {headers.map(header => (
                <th key={header} className="px-[12px] py-[10px] font-[600]">
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, rowIndex) => (
              <tr
                key={`${title}-${rowIndex}`}
                className="border-t border-solid coz-stroke-primary"
              >
                {row.map((cell, cellIndex) => (
                  <td
                    key={`${title}-${rowIndex}-${cellIndex}`}
                    className="px-[12px] py-[12px]"
                  >
                    {cell}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    ) : (
      <div className="flex min-h-[220px] items-center justify-center">
        <Empty title={`No data for ${title}`} />
      </div>
    )}
  </div>
);

const FormField: FC<{
  label: string;
  required?: boolean;
  className?: string;
  children: ReactNode;
}> = ({ label, required, className, children }) => (
  <label className={`block ${className ?? ''}`}>
    <div className="mb-[8px] text-[12px] font-[600] coz-fg-secondary">
      {label}
      {required ? ' *' : ''}
    </div>
    {children}
  </label>
);

function getMetricValue(
  overview: DashboardOverviewData | undefined,
  key: (typeof METRIC_CARDS)[number]['key'],
): string {
  switch (key) {
    case 'customer_total':
      return formatCount(overview?.customer_total);
    case 'new_customers_this_month':
      return formatCount(overview?.new_customers_this_month);
    case 'opportunity_total_amount':
      return formatAmount(overview?.opportunity_total_amount);
    case 'new_opportunities_this_month':
      return formatCount(overview?.new_opportunities_this_month);
    case 'sales_order_total_amount':
      return formatAmount(overview?.sales_order_total_amount);
    default:
      return '0';
  }
}

async function loadCustomers(spaceId: string): Promise<CustomerData[]> {
  try {
    const resp = await crmApi.ListCustomers(
      {
        space_id: spaceId,
        page: 1,
        page_size: 5,
      },
      { __disableErrorToast: true },
    );
    return resp.data?.list ?? [];
  } catch {
    return [];
  }
}

async function loadOpportunities(spaceId: string): Promise<OpportunityData[]> {
  try {
    const resp = await crmApi.ListOpportunities(
      {
        space_id: spaceId,
        page: 1,
        page_size: 5,
      },
      { __disableErrorToast: true },
    );
    return resp.data?.list ?? [];
  } catch {
    return [];
  }
}

async function loadSalesOrders(spaceId: string): Promise<SalesOrderData[]> {
  try {
    const resp = await crmApi.ListSalesOrders(
      {
        space_id: spaceId,
        page: 1,
        page_size: 5,
      },
      { __disableErrorToast: true },
    );
    return resp.data?.list ?? [];
  } catch {
    return [];
  }
}

function formatCount(value?: number): string {
  return new Intl.NumberFormat('zh-CN', {
    maximumFractionDigits: 0,
  }).format(value || 0);
}

function formatAmount(value?: string): string {
  return new Intl.NumberFormat('zh-CN', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(Number(value || 0));
}

function formatDateTime(value?: string): string {
  const timestamp = Number(value || 0);
  if (!timestamp) {
    return '-';
  }

  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) {
    return '-';
  }

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

function pad(value: number): string {
  return String(value).padStart(2, '0');
}

function getInitialCreateCustomerForm(): CreateCustomerFormState {
  return {
    customerName: '',
    ownerUserName: '',
    industry: '',
    remark: '',
  };
}
