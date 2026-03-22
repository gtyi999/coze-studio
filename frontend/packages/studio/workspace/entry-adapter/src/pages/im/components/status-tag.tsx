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

import { Tag } from '@coze-arch/coze-design';

import { CHANNEL_STATUS_META_MAP, TASK_STATUS_META_MAP } from '../constants';

interface IMStatusTagProps {
  value?: string | null;
  type: 'channel' | 'task';
}

export const IMStatusTag: FC<IMStatusTagProps> = ({ value, type }) => {
  const meta =
    type === 'channel'
      ? CHANNEL_STATUS_META_MAP[value ?? '']
      : TASK_STATUS_META_MAP[value ?? ''];

  if (!meta) {
    return (
      <Tag size="mini" color="grey">
        {value || '-'}
      </Tag>
    );
  }

  return (
    <Tag size="mini" color={meta.color} className="font-[500]">
      {meta.text}
    </Tag>
  );
};
