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

import { type FC } from 'react';

import { Empty, Typography } from '@coze-arch/coze-design';

import type { CRMNLQueryResult } from '../types';

interface QueryResultChartProps {
  result?: CRMNLQueryResult;
}

const CHART_WIDTH = 680;
const CHART_HEIGHT = 260;
const CHART_PADDING_X = 32;
const CHART_PADDING_TOP = 24;
const CHART_PADDING_BOTTOM = 44;
const CHART_COLORS = ['#0284c7', '#0f766e', '#d97706', '#7c3aed'];

export const QueryResultChart: FC<QueryResultChartProps> = ({ result }) => {
  const chart = result?.chart;

  if (!result || !chart || chart.type === 'table' || chart.type === 'unknown') {
    return (
      <div className="flex min-h-[220px] items-center justify-center rounded-[16px] border border-dashed coz-stroke-primary bg-[rgba(248,250,252,0.72)]">
        <Empty title="No chart is available for the current result." />
      </div>
    );
  }

  if (chart.type === 'stat') {
    const value = chart.series?.[0]?.data?.[0]?.value ?? 0;
    return (
      <div className="rounded-[16px] border border-solid coz-stroke-primary bg-[linear-gradient(135deg,#eff6ff_0%,#ffffff_55%,#dbeafe_100%)] p-[20px]">
        <Typography.Text className="coz-fg-secondary text-[12px]">
          {chart.metric_label || chart.title || 'Primary metric'}
        </Typography.Text>
        <Typography.Title heading={2} className="!mb-[8px] !mt-[12px]">
          {formatMetricNumber(value)}
        </Typography.Title>
        <Typography.Paragraph className="!mb-0 coz-fg-secondary">
          {chart.title || 'The most important KPI from the current query.'}
        </Typography.Paragraph>
      </div>
    );
  }

  const seriesList = chart.series ?? [];
  if (!seriesList.length) {
    return (
      <div className="flex min-h-[220px] items-center justify-center rounded-[16px] border border-dashed coz-stroke-primary bg-[rgba(248,250,252,0.72)]">
        <Empty title="The chart data is empty." />
      </div>
    );
  }

  const labels = seriesList[0]?.data?.map(item => item.label) ?? [];
  const values = seriesList.flatMap(item => item.data.map(point => point.value));
  const maxValue = Math.max(...values, 1);
  const innerWidth = CHART_WIDTH - CHART_PADDING_X * 2;
  const innerHeight =
    CHART_HEIGHT - CHART_PADDING_TOP - CHART_PADDING_BOTTOM;

  return (
    <div className="rounded-[16px] border border-solid coz-stroke-primary bg-white p-[16px]">
      <div className="mb-[12px] flex flex-wrap items-center justify-between gap-[12px]">
        <div>
          <Typography.Title heading={6} className="!mb-[4px]">
            {chart.title || 'Result Chart'}
          </Typography.Title>
          <Typography.Paragraph className="!mb-0 coz-fg-secondary">
            {chart.metric_label
              ? `${chart.metric_label} visualization`
              : 'The structured rows are rendered as a lightweight chart.'}
          </Typography.Paragraph>
        </div>
        <div className="flex flex-wrap gap-[8px]">
          {seriesList.map((item, index) => (
            <span
              key={`${item.name}-${index}`}
              className="inline-flex items-center gap-[6px] rounded-full bg-[rgba(248,250,252,0.9)] px-[10px] py-[4px] text-[12px] coz-fg-secondary"
            >
              <span
                className="h-[8px] w-[8px] rounded-full"
                style={{ background: CHART_COLORS[index % CHART_COLORS.length] }}
              />
              {item.name}
            </span>
          ))}
        </div>
      </div>
      <svg
        viewBox={`0 0 ${CHART_WIDTH} ${CHART_HEIGHT}`}
        className="h-[280px] w-full overflow-visible"
        preserveAspectRatio="none"
      >
        {[0, 1, 2, 3].map(index => {
          const y = CHART_PADDING_TOP + (innerHeight / 3) * index;
          return (
            <line
              key={`grid-${index}`}
              x1={CHART_PADDING_X}
              x2={CHART_WIDTH - CHART_PADDING_X}
              y1={y}
              y2={y}
              stroke="rgba(15, 23, 42, 0.08)"
              strokeDasharray="4 6"
            />
          );
        })}

        {chart.type === 'bar'
          ? renderBars(
              seriesList[0]?.data ?? [],
              maxValue,
              innerWidth,
              innerHeight,
            )
          : renderLines(seriesList, maxValue, innerWidth, innerHeight)}

        {labels.map((label, index) => {
          const x =
            CHART_PADDING_X +
            innerWidth *
              (labels.length <= 1 ? 0 : index / Math.max(labels.length - 1, 1));
          return (
            <text
              key={`${label}-${index}`}
              x={x}
              y={CHART_HEIGHT - 14}
              textAnchor="middle"
              fill="rgba(15, 23, 42, 0.55)"
              fontSize="11"
            >
              {label}
            </text>
          );
        })}
      </svg>
    </div>
  );
};

function renderBars(
  points: Array<{ label: string; value: number }>,
  maxValue: number,
  innerWidth: number,
  innerHeight: number,
) {
  const barWidth = Math.max(innerWidth / Math.max(points.length * 1.8, 1), 28);

  return [
    ...points.map((point, index) => {
      const ratio = points.length <= 1 ? 0 : index / Math.max(points.length - 1, 1);
      const x = CHART_PADDING_X + innerWidth * ratio - barWidth / 2;
      const barHeight = (point.value / maxValue) * innerHeight;
      const y = CHART_PADDING_TOP + innerHeight - barHeight;

      return (
        <g key={`${point.label}-${index}`}>
          <rect
            x={x}
            y={y}
            width={barWidth}
            height={Math.max(barHeight, 4)}
            rx={10}
            fill="url(#crm-bar-gradient)"
          />
          <text
            x={x + barWidth / 2}
            y={y - 10}
            textAnchor="middle"
            fill="rgba(15, 23, 42, 0.65)"
            fontSize="11"
          >
            {formatMetricNumber(point.value)}
          </text>
        </g>
      );
    }),
    <defs key="bar-gradient">
      <linearGradient id="crm-bar-gradient" x1="0" x2="0" y1="0" y2="1">
        <stop offset="0%" stopColor="#38bdf8" />
        <stop offset="100%" stopColor="#2563eb" />
      </linearGradient>
    </defs>,
  ];
}

function renderLines(
  seriesList: NonNullable<CRMNLQueryResult['chart']>['series'],
  maxValue: number,
  innerWidth: number,
  innerHeight: number,
) {
  return (seriesList ?? []).map((series, seriesIndex) => {
    const points = series.data
      .map((item, pointIndex) => {
        const ratio =
          series.data.length <= 1
            ? 0
            : pointIndex / Math.max(series.data.length - 1, 1);
        const x = CHART_PADDING_X + innerWidth * ratio;
        const y =
          CHART_PADDING_TOP +
          innerHeight -
          (item.value / maxValue) * innerHeight;
        return { x, y };
      })
      .map(item => `${item.x},${item.y}`)
      .join(' ');

    return (
      <polyline
        key={`${series.name}-${seriesIndex}`}
        points={points}
        fill="none"
        stroke={CHART_COLORS[seriesIndex % CHART_COLORS.length]}
        strokeWidth="3"
        strokeLinejoin="round"
        strokeLinecap="round"
      />
    );
  });
}

function formatMetricNumber(value: number): string {
  return new Intl.NumberFormat('zh-CN', {
    maximumFractionDigits: 2,
  }).format(value || 0);
}
