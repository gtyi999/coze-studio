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

import { type FC, useEffect } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@coze-arch/i18n';
import {
  SideSheet,
  Spin,
  Typography,
} from '@coze-arch/coze-design';
import { im as imApi } from '@coze-studio/api-schema';

import { formatDateTime, stringifyJSONBlock } from '../utils';
import { IMStatusTag } from './status-tag';
import type { IMTaskRecord } from '../types';

interface TaskDetailSideSheetProps {
  visible: boolean;
  task?: IMTaskRecord;
  onCancel: () => void;
}

const DetailItem = ({
  label,
  value,
}: {
  label: string;
  value?: string | number | null;
}) => (
  <div>
    <div className="text-[12px] coz-fg-secondary">{label}</div>
    <Typography.Text>{value || '-'}</Typography.Text>
  </div>
);

const JSONBlock = ({ title, content }: { title: string; content?: unknown }) => {
  const text = stringifyJSONBlock(content);
  if (!text) {
    return null;
  }

  return (
    <div className="mt-[16px]">
      <Typography.Text strong>{title}</Typography.Text>
      <pre className="mt-[8px] max-h-[220px] overflow-auto rounded-[12px] coz-bg-primary p-[12px] whitespace-pre-wrap break-words text-[12px] leading-[20px]">
        {text}
      </pre>
    </div>
  );
};

export const TaskDetailSideSheet: FC<TaskDetailSideSheetProps> = ({
  visible,
  task,
  onCancel,
}) => {
  const detailRequest = useRequest(
    async (taskId: string) => {
      const resp = await imApi.GetIMTaskDetail(
        { task_id: taskId },
        { __disableErrorToast: true },
      );

      return resp.data;
    },
    {
      manual: true,
    },
  );

  useEffect(() => {
    if (!visible || !task?.task_id) {
      return;
    }

    detailRequest.run(task.task_id);
  }, [visible, task?.task_id]);

  const currentTask = detailRequest.data ?? task;

  return (
    <SideSheet
      visible={visible}
      width={720}
      title={I18n.t('im_task_detail', {}, 'Task detail')}
      onCancel={onCancel}
    >
      {detailRequest.loading ? (
        <div className="flex h-full items-center justify-center">
          <Spin spinning={true} />
        </div>
      ) : currentTask ? (
        <div className="flex flex-col gap-[16px]">
          <div className="grid gap-[16px] sm:grid-cols-2">
            <DetailItem
              label={I18n.t('im_task_id', {}, 'Task ID')}
              value={currentTask.task_id}
            />
            <div>
              <div className="text-[12px] coz-fg-secondary">
                {I18n.t('api_status_1', {}, 'Status')}
              </div>
              <div className="mt-[4px]">
                <IMStatusTag value={currentTask.status} type="task" />
              </div>
            </div>
            <DetailItem
              label={I18n.t('im_platform', {}, 'Platform')}
              value={currentTask.platform}
            />
            <DetailItem
              label={I18n.t('trace_id', {}, 'Trace ID')}
              value={currentTask.trace_id}
            />
            <DetailItem
              label={I18n.t('im_task_retry', {}, 'Retry count')}
              value={
                typeof currentTask.retry_count === 'number'
                  ? `${currentTask.retry_count}/${currentTask.max_retry_count ?? 0}`
                  : '-'
              }
            />
            <DetailItem
              label={I18n.t('im_session_id', {}, 'Session ID')}
              value={currentTask.session_id}
            />
            <DetailItem
              label={I18n.t('im_event_id', {}, 'Event ID')}
              value={currentTask.event_id}
            />
            <DetailItem
              label={I18n.t('updated_time', {}, 'Updated time')}
              value={formatDateTime(currentTask.updated_at)}
            />
            <DetailItem
              label={I18n.t('created_time', {}, 'Created time')}
              value={formatDateTime(currentTask.created_at)}
            />
            <DetailItem
              label={I18n.t('im_next_retry_time', {}, 'Next retry time')}
              value={formatDateTime(currentTask.next_retry_at)}
            />
          </div>
          <div>
            <div className="text-[12px] coz-fg-secondary">
              {I18n.t('Error_message', {}, 'Error message')}
            </div>
            <Typography.Text>{currentTask.error_msg || '-'}</Typography.Text>
          </div>
          <JSONBlock
            title={I18n.t('im_message_snapshot', {}, 'Message snapshot')}
            content={currentTask.message_snapshot}
          />
          <JSONBlock
            title={I18n.t('im_request_snapshot', {}, 'Request snapshot')}
            content={currentTask.request_snapshot}
          />
          <JSONBlock
            title={I18n.t('im_gateway_ticket', {}, 'Gateway ticket')}
            content={currentTask.gateway_ticket}
          />
          <JSONBlock
            title={I18n.t('im_result_snapshot', {}, 'Result snapshot')}
            content={currentTask.result_snapshot}
          />
        </div>
      ) : (
        <div className="flex h-full items-center justify-center coz-fg-secondary">
          {I18n.t('im_task_detail_empty', {}, 'No task detail available')}
        </div>
      )}
    </SideSheet>
  );
};
