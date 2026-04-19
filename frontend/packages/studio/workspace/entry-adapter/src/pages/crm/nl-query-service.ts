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

import { crm as crmApi } from '@coze-studio/api-schema';

import { CRM_QUERY_EXAMPLES, CRM_QUERY_MOCK_RESULT_MAP } from './mock-data';
import type {
  CRMColumnValueType,
  CRMNLQueryRequest,
  CRMNLQueryResult,
  CRMQueryChart,
  CRMQueryMeta,
  CRMQueryTableColumn,
} from './types';

type CRMNLQueryAPIFn = (
  payload: Record<string, unknown>,
  options?: Record<string, unknown>,
) => Promise<{ data?: unknown }>;

interface CRMApiWithNLQuery {
  RunCRMNLQuery?: CRMNLQueryAPIFn;
}

export async function runCRMNLQuery(
  params: CRMNLQueryRequest,
): Promise<CRMNLQueryResult> {
  const nlQueryApi = (crmApi as unknown as CRMApiWithNLQuery).RunCRMNLQuery;

  if (nlQueryApi) {
    try {
      const resp = await nlQueryApi(
        {
          space_id: params.space_id,
          question: params.question,
          debug: params.debug,
        },
        { __disableErrorToast: true },
      );

      const normalizedResult = normalizeCRMNLQueryResult(
        params.question,
        resp.data,
      );

      if (isMeaningfulCRMNLQueryResult(normalizedResult)) {
        return normalizedResult;
      }
    } catch (error) {
      console.warn(
        '[CRM NL Query] RunCRMNLQuery failed, falling back to mock result.',
        error,
      );
    }
  }

  return getMockCRMNLQueryResult(params.question);
}

export function getMockCRMNLQueryResult(question: string): CRMNLQueryResult {
  const lowerQuestion = question.toLowerCase();

  if (
    lowerQuestion.includes('\u9884\u6d4b') ||
    lowerQuestion.includes('\u70ed\u9500') ||
    lowerQuestion.includes('\u5356\u5f97\u6700\u597d') ||
    lowerQuestion.includes('forecast')
  ) {
    return {
      ...CRM_QUERY_MOCK_RESULT_MAP.forecast_hot_product,
      question,
    };
  }

  if (
    lowerQuestion.includes('\u4e1a\u7ee9') ||
    lowerQuestion.includes('\u9500\u552e') ||
    lowerQuestion.includes('\u51a0\u519b') ||
    lowerQuestion.includes('sales')
  ) {
    return {
      ...CRM_QUERY_MOCK_RESULT_MAP.top_sales_current_quarter,
      question,
    };
  }

  return {
    ...CRM_QUERY_MOCK_RESULT_MAP.customer_count,
    question,
  };
}

function normalizeCRMNLQueryResult(
  question: string,
  payload: unknown,
): CRMNLQueryResult {
  const raw = toRecord(payload);
  const rows = normalizeRows(raw.data);
  const tableColumns = normalizeColumns(raw.table_columns, rows);

  return {
    question,
    answer: String(raw.answer ?? ''),
    intent_type: String(raw.intent_type ?? 'unknown') as CRMNLQueryResult['intent_type'],
    data: rows,
    table_columns: tableColumns,
    chart: normalizeChart(raw.chart, rows),
    disclaimer: String(raw.disclaimer ?? ''),
    meta: normalizeMeta(raw.meta),
  };
}

function normalizeRows(payload: unknown): Array<Record<string, unknown>> {
  if (!Array.isArray(payload)) {
    return [];
  }

  return payload
    .map(item => toRecord(item))
    .filter(item => Object.keys(item).length > 0);
}

function normalizeColumns(
  payload: unknown,
  rows: Array<Record<string, unknown>>,
): CRMQueryTableColumn[] {
  if (Array.isArray(payload) && payload.length > 0) {
    return payload
      .map(item => toRecord(item))
      .filter(item => item.dataIndex || item.data_index)
      .map(item => ({
        key: String(item.key ?? item.dataIndex ?? item.data_index),
        title: String(item.title ?? item.key ?? item.dataIndex ?? ''),
        dataIndex: String(item.dataIndex ?? item.data_index ?? item.key ?? ''),
        valueType: String(
          item.valueType ?? item.value_type ?? 'string',
        ) as CRMColumnValueType,
      }));
  }

  const firstRow = rows[0];
  if (!firstRow) {
    return [];
  }

  return Object.keys(firstRow).map(key => ({
    key,
    title: key,
    dataIndex: key,
    valueType: inferValueType(firstRow[key]),
  }));
}

function normalizeChart(
  payload: unknown,
  rows: Array<Record<string, unknown>>,
): CRMQueryChart | undefined {
  const raw = toRecord(payload);

  if (!Object.keys(raw).length) {
    if (!rows.length) {
      return undefined;
    }

    const firstKey = Object.keys(rows[0])[0];
    const firstValue = rows[0][firstKey];
    const numericKey = Object.keys(rows[0]).find(key =>
      typeof rows[0][key] === 'number',
    );

    if (typeof firstValue === 'number' || numericKey === firstKey) {
      return {
        type: 'stat',
        title: firstKey,
        metric_label: firstKey,
        series: [
          {
            name: firstKey,
            data: [{ label: 'Current', value: Number(firstValue || 0) }],
          },
        ],
      };
    }

    return {
      type: 'table',
      title: 'Query Result',
    };
  }

  return {
    type: String(raw.type ?? 'unknown') as CRMQueryChart['type'],
    title: String(raw.title ?? ''),
    x_field: String(raw.x_field ?? raw.xField ?? ''),
    y_field: String(raw.y_field ?? raw.yField ?? ''),
    metric_label: String(raw.metric_label ?? raw.metricLabel ?? ''),
    series: Array.isArray(raw.series)
      ? raw.series.map(seriesItem => {
          const series = toRecord(seriesItem);
          return {
            name: String(series.name ?? ''),
            data: Array.isArray(series.data)
              ? series.data.map(item => {
                  const point = toRecord(item);
                  return {
                    label: String(point.label ?? ''),
                    value: Number(point.value ?? 0),
                  };
                })
              : [],
          };
        })
      : [],
  };
}

function normalizeMeta(payload: unknown): CRMQueryMeta | undefined {
  const raw = toRecord(payload);
  if (!Object.keys(raw).length) {
    return undefined;
  }

  return {
    definition: String(raw.definition ?? ''),
    updated_at: String(raw.updated_at ?? raw.updatedAt ?? ''),
    scope_note: String(raw.scope_note ?? raw.scopeNote ?? ''),
    applied_filters: Array.isArray(raw.applied_filters)
      ? raw.applied_filters.map(item => String(item))
      : [],
  };
}

function inferValueType(value: unknown): CRMColumnValueType {
  if (typeof value === 'number') {
    return 'number';
  }
  return 'string';
}

function toRecord(payload: unknown): Record<string, unknown> {
  if (!payload || typeof payload !== 'object' || Array.isArray(payload)) {
    return {};
  }
  return payload as Record<string, unknown>;
}

function isMeaningfulCRMNLQueryResult(result: CRMNLQueryResult): boolean {
  return Boolean(
    result.answer ||
      result.data.length ||
      result.table_columns?.length ||
      result.chart?.series?.length,
  );
}

export function getDefaultCRMNLQuestion(): string {
  return CRM_QUERY_EXAMPLES[0]?.question ?? '';
}
