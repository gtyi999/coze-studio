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

interface CRMStatusTagProps {
  value?: string;
}

const STATUS_META: Record<string, { color: 'green' | 'blue' | 'grey' }> = {
  active: { color: 'green' },
  open: { color: 'green' },
  draft: { color: 'blue' },
  inactive: { color: 'grey' },
};

export const CRMStatusTag: FC<CRMStatusTagProps> = ({ value }) => {
  const text = String(value || '-');
  const meta = STATUS_META[text] ?? { color: 'grey' as const };

  return (
    <Tag size="mini" color={meta.color} className="font-[500]">
      {text}
    </Tag>
  );
};
