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

export type CRMQueryIntentType =
  | 'unknown'
  | 'customer_count'
  | 'top_sales_current_quarter'
  | 'sales_topn_current_quarter'
  | 'product_sales_topn'
  | 'public_pool_count'
  | 'forecast_hot_product';

export type CRMChartType = 'unknown' | 'stat' | 'bar' | 'line' | 'table';

export type CRMColumnValueType =
  | 'string'
  | 'number'
  | 'currency'
  | 'percent'
  | 'date';

export interface CRMNLQueryRequest {
  space_id: string;
  question: string;
  debug?: boolean;
}

export interface CRMQueryTableColumn {
  key: string;
  title: string;
  dataIndex: string;
  valueType?: CRMColumnValueType;
}

export interface CRMChartPoint {
  label: string;
  value: number;
}

export interface CRMChartSeries {
  name: string;
  data: CRMChartPoint[];
}

export interface CRMQueryChart {
  type: CRMChartType;
  title?: string;
  x_field?: string;
  y_field?: string;
  metric_label?: string;
  series?: CRMChartSeries[];
}

export interface CRMQueryMeta {
  definition?: string;
  updated_at?: string;
  scope_note?: string;
  applied_filters?: string[];
}

export interface CRMNLQueryResult {
  question: string;
  answer: string;
  intent_type: CRMQueryIntentType;
  data: Array<Record<string, unknown>>;
  table_columns?: CRMQueryTableColumn[];
  chart?: CRMQueryChart;
  disclaimer?: string;
  meta?: CRMQueryMeta;
}

export interface CRMQueryExample {
  key: string;
  title: string;
  description: string;
  question: string;
  intent_type: CRMQueryIntentType;
}

export interface CRMQueryHistoryItem {
  id: string;
  question: string;
  answer: string;
  intent_type: CRMQueryIntentType;
  created_at: string;
}
