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

import { PLATFORM_TEXT_MAP } from './constants';

export function formatDateTime(value?: string | number | null) {
  if (value === null || typeof value === 'undefined' || value === '') {
    return '-';
  }

  const numericValue =
    typeof value === 'string' ? Number.parseInt(value, 10) : value;
  const timestamp = Number.isFinite(numericValue)
    ? numericValue
    : Number.parseInt(String(value), 10);
  const date = new Date(timestamp);

  if (Number.isNaN(date.getTime())) {
    return String(value);
  }

  return date.toLocaleString();
}

export function getPlatformLabel(platform?: string | null) {
  if (!platform) {
    return '-';
  }

  return PLATFORM_TEXT_MAP[platform] ?? platform;
}

export function parseExtJSON(text?: string) {
  if (!text?.trim()) {
    return undefined;
  }

  const parsed = JSON.parse(text) as Record<string, unknown>;
  return Object.fromEntries(
    Object.entries(parsed).map(([key, value]) => [key, String(value)]),
  );
}

export function stringifyExtJSON(value?: Record<string, string>) {
  if (!value || Object.keys(value).length === 0) {
    return '';
  }

  return JSON.stringify(value, null, 2);
}

export function stringifyJSONBlock(value?: unknown) {
  if (!value) {
    return '';
  }

  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }

  return JSON.stringify(value, null, 2);
}
