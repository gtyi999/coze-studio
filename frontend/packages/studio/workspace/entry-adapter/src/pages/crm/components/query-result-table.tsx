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

import { type FC, useMemo } from 'react';

import {
  Empty,
  Table,
  Typography,
  type ColumnProps,
} from '@coze-arch/coze-design';

import type { CRMNLQueryResult, CRMQueryTableColumn } from '../types';

interface QueryResultTableProps {
  result?: CRMNLQueryResult;
}

type TableRow = Record<string, unknown> & {
  __rowKey: string;
};

export const QueryResultTable: FC<QueryResultTableProps> = ({ result }) => {
  const rows = useMemo<TableRow[]>(
    () =>
      (result?.data ?? []).map((item, index) => ({
        __rowKey: `${result?.intent_type ?? 'query'}-${index}`,
        ...item,
      })),
    [result],
  );

  const columns = useMemo<ColumnProps<TableRow>[]>(
    () =>
      (result?.table_columns ?? []).map(column => ({
        title: column.title,
        dataIndex: column.dataIndex,
        render: value => renderCellValue(value, column),
      })),
    [result?.table_columns],
  );

  if (!result) {
    return (
      <div className="flex min-h-[180px] items-center justify-center rounded-[16px] border border-dashed coz-stroke-primary bg-[rgba(248,250,252,0.72)]">
        <Empty title="Run a query to render the structured table here." />
      </div>
    );
  }

  if (!rows.length || !columns.length) {
    return (
      <div className="flex min-h-[180px] items-center justify-center rounded-[16px] border border-dashed coz-stroke-primary bg-[rgba(248,250,252,0.72)]">
        <Empty title="The current query did not return a tabular result." />
      </div>
    );
  }

  return (
    <div className="rounded-[16px] border border-solid coz-stroke-primary bg-white">
      <div className="border-b border-solid coz-stroke-primary px-[16px] py-[12px]">
        <Typography.Title heading={6} className="!mb-0">
          Structured Result
        </Typography.Title>
      </div>
      <Table
        offsetY={320}
        tableProps={{
          rowKey: '__rowKey',
          columns,
          dataSource: rows,
          pagination: false,
          scroll: {
            x: Math.max(columns.length * 180, 640),
          },
        }}
      />
    </div>
  );
};

function renderCellValue(value: unknown, column: CRMQueryTableColumn) {
  switch (column.valueType) {
    case 'currency':
      return formatCurrency(value);
    case 'percent':
      return formatPercent(value);
    case 'number':
      return formatNumber(value);
    default:
      return (
        <Typography.Text ellipsis={{ showTooltip: true }}>
          {String(value ?? '-')}
        </Typography.Text>
      );
  }
}

function formatNumber(value: unknown): string {
  return new Intl.NumberFormat('zh-CN', {
    maximumFractionDigits: 2,
  }).format(Number(value ?? 0));
}

function formatCurrency(value: unknown): string {
  return `${new Intl.NumberFormat('zh-CN', {
    maximumFractionDigits: 2,
  }).format(Number(value ?? 0))} yuan`;
}

function formatPercent(value: unknown): string {
  const numericValue = Number(value ?? 0);
  return new Intl.NumberFormat('zh-CN', {
    style: 'percent',
    maximumFractionDigits: 1,
  }).format(numericValue);
}
