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

import { type FC, useEffect, useRef, useState } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@coze-arch/i18n';
import {
  Button,
  Form,
  FormInput,
  FormSelect,
  FormTextArea,
  Spin,
  Toast,
  Typography,
  type FormApi,
} from '@coze-arch/coze-design';
import { im as imApi } from '@coze-studio/api-schema';

import type {
  ConnectivityTestFormValues,
  IMChannelConfig,
  IMConnectivityResult,
} from '../types';
import { IMStatusTag } from './status-tag';

interface ConnectivityTestPanelProps {
  spaceId: string;
  channels: IMChannelConfig[];
  defaultConfigId?: string;
}

export const ConnectivityTestPanel: FC<ConnectivityTestPanelProps> = ({
  spaceId,
  channels,
  defaultConfigId,
}) => {
  const formApiRef = useRef<FormApi<ConnectivityTestFormValues>>();
  const [result, setResult] = useState<IMConnectivityResult>();

  const configOptions = channels.map(item => ({
    label: item.name || item.config_id || '-',
    value: item.config_id,
  }));

  const testRequest = useRequest(
    async (values: ConnectivityTestFormValues) => {
      const resp = await imApi.TestIMChannelConnectivity(
        {
          config_id: values.config_id,
          space_id: spaceId,
          sender_id: values.sender_id,
          sender_name: values.sender_name,
          chat_id: values.chat_id,
          thread_id: values.thread_id,
          message_text: values.message_text,
          is_at_bot: values.is_at_bot,
          prefer_async: values.prefer_async,
        },
        { __disableErrorToast: true },
      );

      return resp.data;
    },
    {
      manual: true,
      onSuccess: data => {
        setResult(data);
        Toast.success(
          I18n.t(
            'im_test_request_sent',
            {},
            'Connectivity test request submitted',
          ),
        );
      },
      onError: error => {
        Toast.error(
          error instanceof Error
            ? error.message
            : I18n.t('im_test_failed', {}, 'Connectivity test failed'),
        );
      },
    },
  );

  useEffect(() => {
    if (!defaultConfigId || !formApiRef.current) {
      return;
    }

    formApiRef.current.setValue('config_id', defaultConfigId);
  }, [defaultConfigId]);

  const handleSubmit = async () => {
    try {
      const values = await formApiRef.current?.validate();
      if (!values) {
        return;
      }

      await testRequest.runAsync(values);
    } catch (error) {
      if (error instanceof Error && error.message) {
        Toast.error(error.message);
      }
    }
  };

  return (
    <div className="grid gap-[20px] xl:grid-cols-[minmax(0,480px)_minmax(0,1fr)]">
      <div className="rounded-[20px] border coz-stroke-primary coz-bg-max p-[20px]">
        <div className="mb-[16px]">
          <Typography.Title heading={6} className="!mb-[4px]">
            {I18n.t('im_connectivity_test', {}, 'Connectivity test')}
          </Typography.Title>
          <Typography.Paragraph className="!mb-0 coz-fg-secondary">
            {I18n.t(
              'im_connectivity_test_desc',
              {},
              'Trigger one mock inbound message through the IM gateway.',
            )}
          </Typography.Paragraph>
        </div>
        <Form<ConnectivityTestFormValues>
          getFormApi={api => {
            formApiRef.current = api;
          }}
          initValues={{
            config_id: defaultConfigId,
            is_at_bot: true,
            prefer_async: true,
          }}
        >
          <FormSelect
            field="config_id"
            label={I18n.t('im_channel_config', {}, 'Channel config')}
            optionList={configOptions}
            rules={[{ required: true }]}
          />
          <FormInput
            field="sender_id"
            label={I18n.t('im_sender_id', {}, 'Sender ID')}
            rules={[{ required: true }]}
          />
          <FormInput
            field="sender_name"
            label={I18n.t('im_sender_name', {}, 'Sender name')}
          />
          <FormInput
            field="chat_id"
            label={I18n.t('im_chat_id', {}, 'Chat ID')}
          />
          <FormInput
            field="thread_id"
            label={I18n.t('im_thread_id', {}, 'Thread ID')}
          />
          <FormTextArea
            field="message_text"
            label={I18n.t('im_mock_message', {}, 'Mock message')}
            autosize={{ minRows: 4, maxRows: 8 }}
            rules={[{ required: true }]}
          />
          <Form.Checkbox field="is_at_bot" noLabel>
            {I18n.t('im_at_bot_message', {}, '@ robot message')}
          </Form.Checkbox>
          <Form.Checkbox field="prefer_async" noLabel>
            {I18n.t('im_prefer_async', {}, 'Prefer async execution')}
          </Form.Checkbox>
        </Form>
        <div className="mt-[16px] flex justify-end">
          <Button loading={testRequest.loading} onClick={handleSubmit}>
            {I18n.t('workflow_testset_edit_confirm', {}, 'Run test')}
          </Button>
        </div>
      </div>
      <div className="rounded-[20px] border coz-stroke-primary coz-bg-max p-[20px]">
        <div className="mb-[16px]">
          <Typography.Title heading={6} className="!mb-[4px]">
            {I18n.t('im_test_result', {}, 'Test result')}
          </Typography.Title>
          <Typography.Paragraph className="!mb-0 coz-fg-secondary">
            {I18n.t(
              'im_test_result_desc',
              {},
              'Accepted tasks will continue on the backend and can be tracked in the task list.',
            )}
          </Typography.Paragraph>
        </div>
        {testRequest.loading ? (
          <div className="flex h-[240px] items-center justify-center">
            <Spin spinning={true} />
          </div>
        ) : result ? (
          <div className="flex flex-col gap-[16px]">
            <div className="grid gap-[12px] sm:grid-cols-2">
              <div>
                <div className="text-[12px] coz-fg-secondary">
                  {I18n.t('api_status_1', {}, 'Status')}
                </div>
                <div className="mt-[4px]">
                  <IMStatusTag value={result.status} type="task" />
                </div>
              </div>
              <div>
                <div className="text-[12px] coz-fg-secondary">
                  {I18n.t('im_task_id', {}, 'Task ID')}
                </div>
                <Typography.Text>{result.task_id || '-'}</Typography.Text>
              </div>
              <div>
                <div className="text-[12px] coz-fg-secondary">
                  {I18n.t('trace_id', {}, 'Trace ID')}
                </div>
                <Typography.Text>{result.trace_id || '-'}</Typography.Text>
              </div>
              <div>
                <div className="text-[12px] coz-fg-secondary">
                  {I18n.t('im_result_message', {}, 'Message')}
                </div>
                <Typography.Text>{result.message || '-'}</Typography.Text>
              </div>
            </div>
            <div>
              <div className="mb-[8px] text-[12px] coz-fg-secondary">
                {I18n.t('im_reply_preview', {}, 'Reply preview')}
              </div>
              <pre className="max-h-[280px] overflow-auto rounded-[12px] coz-bg-primary p-[12px] whitespace-pre-wrap break-words text-[12px] leading-[20px]">
                {result.reply_preview || '-'}
              </pre>
            </div>
          </div>
        ) : (
          <div className="flex h-[240px] items-center justify-center coz-fg-secondary">
            {I18n.t(
              'im_no_test_result',
              {},
              'Run a connectivity test to inspect the gateway response.',
            )}
          </div>
        )}
      </div>
    </div>
  );
};
