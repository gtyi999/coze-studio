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

import type { CRMNLQueryResult, CRMQueryExample } from './types';

export const CRM_QUERY_EXAMPLES: CRMQueryExample[] = [
  {
    key: 'customer-count',
    title: 'Customer Count',
    description: 'Check the current active customer volume.',
    question: 'How many customers do I have now?',
    intent_type: 'customer_count',
  },
  {
    key: 'top-sales',
    title: 'Top Sales',
    description: 'Find the best sales performer in the current quarter.',
    question: 'Which sales rep has the best performance?',
    intent_type: 'top_sales_current_quarter',
  },
  {
    key: 'forecast-hot-product',
    title: 'Hot Product Forecast',
    description: 'Predict the best selling product for next quarter.',
    question: 'Which product will sell best next quarter?',
    intent_type: 'forecast_hot_product',
  },
];

export const CRM_QUERY_MOCK_RESULT_MAP: Record<string, CRMNLQueryResult> = {
  customer_count: {
    question: 'How many customers do I have now?',
    answer: 'There are 1,277 active customers in the current workspace.',
    intent_type: 'customer_count',
    data: [
      {
        customer_count: 1277,
      },
    ],
    table_columns: [
      {
        key: 'customer_count',
        title: 'Active Customers',
        dataIndex: 'customer_count',
        valueType: 'number',
      },
    ],
    chart: {
      type: 'stat',
      title: 'Active Customer Count',
      metric_label: 'Customers',
      series: [
        {
          name: 'Customers',
          data: [{ label: 'Current', value: 1277 }],
        },
      ],
    },
    meta: {
      definition: 'Counts customer records where status equals active.',
      updated_at: '2026-04-19 10:32',
      scope_note: 'All authorized customers in the current workspace',
      applied_filters: ['customer.status = active'],
    },
  },
  top_sales_current_quarter: {
    question: 'Which sales rep has the best performance?',
    answer: 'Zhang San leads this quarter with 947,818 yuan in sales.',
    intent_type: 'top_sales_current_quarter',
    data: [
      { sales_user_name: 'Zhang San', sales_amount: 947818, rank_no: 1 },
      { sales_user_name: 'Li Si', sales_amount: 882540, rank_no: 2 },
      { sales_user_name: 'Wang Wu', sales_amount: 798230, rank_no: 3 },
      { sales_user_name: 'Zhao Liu', sales_amount: 733410, rank_no: 4 },
      { sales_user_name: 'Sun Qi', sales_amount: 695880, rank_no: 5 },
    ],
    table_columns: [
      {
        key: 'rank_no',
        title: 'Rank',
        dataIndex: 'rank_no',
        valueType: 'number',
      },
      {
        key: 'sales_user_name',
        title: 'Sales Rep',
        dataIndex: 'sales_user_name',
        valueType: 'string',
      },
      {
        key: 'sales_amount',
        title: 'Sales Amount',
        dataIndex: 'sales_amount',
        valueType: 'currency',
      },
    ],
    chart: {
      type: 'bar',
      title: 'Quarterly Sales Top 5',
      x_field: 'sales_user_name',
      y_field: 'sales_amount',
      metric_label: 'Sales Amount',
      series: [
        {
          name: 'Sales Amount',
          data: [
            { label: 'Zhang San', value: 947818 },
            { label: 'Li Si', value: 882540 },
            { label: 'Wang Wu', value: 798230 },
            { label: 'Zhao Liu', value: 733410 },
            { label: 'Sun Qi', value: 695880 },
          ],
        },
      ],
    },
    meta: {
      definition: 'Aggregates crm_sales_order.amount by sales rep.',
      updated_at: '2026-04-19 10:32',
      scope_note: 'Authorized sales orders in the current workspace',
      applied_filters: ['time_range = current_quarter', 'status != draft'],
    },
  },
  forecast_hot_product: {
    question: 'Which product will sell best next quarter?',
    answer:
      'AI Smart Seat Pack is the leading candidate for next quarter because its recent growth is steady and its volatility is low.',
    intent_type: 'forecast_hot_product',
    data: [
      {
        product_name: 'AI Smart Seat Pack',
        trend_score: 92.4,
        last_3m_avg: 1860,
        growth_rate: 0.24,
      },
      {
        product_name: 'Multi-channel Leads Suite',
        trend_score: 87.8,
        last_3m_avg: 1642,
        growth_rate: 0.19,
      },
      {
        product_name: 'Sales Automation Flow',
        trend_score: 83.6,
        last_3m_avg: 1510,
        growth_rate: 0.12,
      },
    ],
    table_columns: [
      {
        key: 'product_name',
        title: 'Product',
        dataIndex: 'product_name',
        valueType: 'string',
      },
      {
        key: 'trend_score',
        title: 'Trend Score',
        dataIndex: 'trend_score',
        valueType: 'number',
      },
      {
        key: 'last_3m_avg',
        title: 'Last 3-Month Avg',
        dataIndex: 'last_3m_avg',
        valueType: 'number',
      },
      {
        key: 'growth_rate',
        title: 'Growth Rate',
        dataIndex: 'growth_rate',
        valueType: 'percent',
      },
    ],
    chart: {
      type: 'line',
      title: 'Recent 6-Month Sales Trend',
      x_field: 'month',
      y_field: 'sales_qty',
      metric_label: 'Sales Qty',
      series: [
        {
          name: 'AI Smart Seat Pack',
          data: [
            { label: 'Nov', value: 1220 },
            { label: 'Dec', value: 1380 },
            { label: 'Jan', value: 1490 },
            { label: 'Feb', value: 1680 },
            { label: 'Mar', value: 1930 },
            { label: 'Apr', value: 1970 },
          ],
        },
        {
          name: 'Multi-channel Leads Suite',
          data: [
            { label: 'Nov', value: 1130 },
            { label: 'Dec', value: 1240 },
            { label: 'Jan', value: 1360 },
            { label: 'Feb', value: 1490 },
            { label: 'Mar', value: 1650 },
            { label: 'Apr', value: 1785 },
          ],
        },
      ],
    },
    disclaimer:
      'Forecasts are based on the last 6 months of historical trends and do not guarantee business outcomes.',
    meta: {
      definition:
        'Combines 6-month sales history, recent weighted average, growth rate and volatility into a trend score.',
      updated_at: '2026-04-19 10:32',
      scope_note: 'Authorized product sales in the current workspace',
      applied_filters: ['time_range = last_6_months', 'status != draft'],
    },
  },
};
