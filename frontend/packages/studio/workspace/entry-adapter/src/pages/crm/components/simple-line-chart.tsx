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

import { Empty, Typography } from '@coze-arch/coze-design';
import type { DashboardOrderTrendData } from '@coze-studio/api-schema/crm';

interface SimpleLineChartProps {
  data: DashboardOrderTrendData[];
}

const CHART_WIDTH = 640;
const CHART_HEIGHT = 240;
const PADDING_X = 32;
const PADDING_TOP = 24;
const PADDING_BOTTOM = 40;

export const SimpleLineChart: FC<SimpleLineChartProps> = ({ data }) => {
  const { labels, maxValue, points, areaPoints } = useMemo(() => {
    if (!data.length) {
      return {
        labels: [] as { x: number; label: string }[],
        maxValue: 0,
        points: '',
        areaPoints: '',
      };
    }

    const values = data.map(item => Number(item.order_amount || 0));
    const max = Math.max(...values, 1);
    const chartInnerWidth = CHART_WIDTH - PADDING_X * 2;
    const chartInnerHeight = CHART_HEIGHT - PADDING_TOP - PADDING_BOTTOM;

    const rawPoints = data.map((item, index) => {
      const ratio = data.length === 1 ? 0 : index / (data.length - 1);
      const x = PADDING_X + chartInnerWidth * ratio;
      const y =
        PADDING_TOP +
        chartInnerHeight -
        (Number(item.order_amount || 0) / max) * chartInnerHeight;
      return { x, y, label: formatDateLabel(String(item.date || '')) };
    });

    const labelsStep = Math.max(Math.floor(data.length / 6), 1);
    const chartLabels = rawPoints.filter(
      (_item, index) =>
        index % labelsStep === 0 || index === rawPoints.length - 1,
    );

    const chartPoints = rawPoints.map(item => `${item.x},${item.y}`).join(' ');
    const chartAreaPoints = [
      `${rawPoints[0]?.x},${CHART_HEIGHT - PADDING_BOTTOM}`,
      chartPoints,
      `${rawPoints[rawPoints.length - 1]?.x},${CHART_HEIGHT - PADDING_BOTTOM}`,
    ].join(' ');

    return {
      labels: chartLabels,
      maxValue: max,
      points: chartPoints,
      areaPoints: chartAreaPoints,
    };
  }, [data]);

  if (!data.length) {
    return (
      <div className="flex min-h-[280px] items-center justify-center">
        <Empty title="No order trend data in the last 30 days" />
      </div>
    );
  }

  return (
    <div className="w-full">
      <div className="mb-[8px] flex items-center justify-between">
        <Typography.Text className="coz-fg-secondary text-[12px]">
          Amount trend
        </Typography.Text>
        <Typography.Text className="coz-fg-secondary text-[12px]">
          Peak {formatAmount(maxValue)}
        </Typography.Text>
      </div>
      <svg
        viewBox={`0 0 ${CHART_WIDTH} ${CHART_HEIGHT}`}
        className="h-[260px] w-full overflow-visible"
        preserveAspectRatio="none"
      >
        {[0, 1, 2, 3].map(index => {
          const y =
            PADDING_TOP +
            ((CHART_HEIGHT - PADDING_TOP - PADDING_BOTTOM) / 3) * index;
          return (
            <line
              key={index}
              x1={PADDING_X}
              x2={CHART_WIDTH - PADDING_X}
              y1={y}
              y2={y}
              stroke="rgba(15, 23, 42, 0.08)"
              strokeDasharray="4 6"
            />
          );
        })}
        <polyline
          points={areaPoints}
          fill="rgba(14, 165, 233, 0.12)"
          stroke="none"
        />
        <polyline
          points={points}
          fill="none"
          stroke="#0284c7"
          strokeWidth="3"
          strokeLinejoin="round"
          strokeLinecap="round"
        />
        {labels.map(item => (
          <text
            key={`${item.label}-${item.x}`}
            x={item.x}
            y={CHART_HEIGHT - 12}
            textAnchor="middle"
            fill="rgba(15, 23, 42, 0.55)"
            fontSize="11"
          >
            {item.label}
          </text>
        ))}
      </svg>
    </div>
  );
};

function formatDateLabel(value: string): string {
  if (!value) {
    return '-';
  }
  const parts = value.split('-');
  if (parts.length !== 3) {
    return value;
  }
  return `${parts[1]}-${parts[2]}`;
}

function formatAmount(value: number): string {
  return new Intl.NumberFormat('zh-CN', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(value || 0);
}
