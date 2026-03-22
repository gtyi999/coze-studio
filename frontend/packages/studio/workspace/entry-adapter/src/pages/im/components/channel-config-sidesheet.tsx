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

import { type FC, useEffect, useMemo, useRef, useState } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@coze-arch/i18n';
import {
  Button,
  Form,
  FormInput,
  FormSelect,
  FormTextArea,
  SideSheet,
  Spin,
  Toast,
  Typography,
  type FormApi,
} from '@coze-arch/coze-design';
import { im as imApi } from '@coze-studio/api-schema';
import { type CreateIMChannelConfigRequest } from '@coze-studio/api-schema/im';

import {
  CHANNEL_STATUS_OPTIONS,
  SESSION_SCOPE_OPTIONS,
} from '../constants';
import { parseExtJSON, stringifyExtJSON } from '../utils';
import type {
  ChannelFormValues,
  IMChannelConfig,
  IMPlatformRecord,
} from '../types';

interface ChannelConfigSideSheetProps {
  visible: boolean;
  mode: 'create' | 'edit';
  spaceId: string;
  value?: IMChannelConfig;
  platforms: IMPlatformRecord[];
  onCancel: () => void;
  onSuccess: () => void;
}

function toFormValues(
  value: IMChannelConfig | undefined,
): ChannelFormValues {
  return {
    name: value?.name ?? '',
    platform: value?.platform ?? '',
    bot_id: value?.bot_id ?? '',
    connector_id: value?.connector_id ?? '',
    tenant_key: value?.tenant_key ?? '',
    app_id: value?.app_id ?? '',
    bot_code: value?.bot_code ?? '',
    session_scope: value?.session_scope ?? 'chat',
    status: value?.status ?? 'enabled',
    platform_config:
      value?.platform_config ?? value?.masked_platform_config ?? '',
    ext_json: stringifyExtJSON(value?.ext),
  };
}

// eslint-disable-next-line @coze-arch/max-line-per-function
export const ChannelConfigSideSheet: FC<ChannelConfigSideSheetProps> = ({
  visible,
  mode,
  spaceId,
  value,
  platforms,
  onCancel,
  onSuccess,
}) => {
  const formApiRef = useRef<FormApi<ChannelFormValues>>();
  const [currentPlatform, setCurrentPlatform] = useState<string>(
    value?.platform ?? platforms[0]?.platform ?? '',
  );

  const platformOptions = useMemo(
    () =>
      platforms.map(item => ({
        label: item.name,
        value: item.platform,
      })),
    [platforms],
  );

  const platformInfo = useMemo(
    () => platforms.find(item => item.platform === currentPlatform),
    [currentPlatform, platforms],
  );

  const detailRequest = useRequest(
    async (configId: string) => {
      const resp = await imApi.GetIMChannelConfig(
        { config_id: configId },
        { __disableErrorToast: true },
      );

      return resp.data;
    },
    {
      manual: true,
    },
  );

  const saveRequest = useRequest(
    async (payload: CreateIMChannelConfigRequest) => {
      if (mode === 'edit' && value?.config_id) {
        await imApi.UpdateIMChannelConfig(
          {
            config_id: value.config_id,
            name: payload.name,
            bot_id: payload.bot_id,
            connector_id: payload.connector_id,
            tenant_key: payload.tenant_key,
            app_id: payload.app_id,
            bot_code: payload.bot_code,
            session_scope: payload.session_scope,
            status: payload.status,
            platform_config: payload.platform_config,
            ext: payload.ext,
          },
          { __disableErrorToast: true },
        );
        return;
      }

      await imApi.CreateIMChannelConfig(payload, {
        __disableErrorToast: true,
      });
    },
    {
      manual: true,
      onSuccess: () => {
        Toast.success(
          I18n.t(
            mode === 'edit' ? 'Update_success' : 'Create_success',
            {},
            mode === 'edit' ? 'Updated successfully' : 'Created successfully',
          ),
        );
        onSuccess();
      },
      onError: error => {
        Toast.error(
          error instanceof Error
            ? error.message
            : I18n.t('im_channel_save_failed', {}, 'Failed to save channel'),
        );
      },
    },
  );

  useEffect(() => {
    if (!visible) {
      return;
    }

    if (mode === 'edit' && value?.config_id) {
      detailRequest.run(value.config_id);
      return;
    }

    const initialValue = toFormValues(value);
    const initialPlatform = initialValue.platform || platforms[0]?.platform || '';
    initialValue.platform = initialPlatform;
    formApiRef.current?.setValues(initialValue, { isOverride: true });
    setCurrentPlatform(initialPlatform);
  }, [visible, mode, value, platforms]);

  useEffect(() => {
    if (!visible || !formApiRef.current) {
      return;
    }

    const detail = detailRequest.data ?? value;
    if (!detail) {
      return;
    }

    formApiRef.current.setValues(toFormValues(detail), { isOverride: true });
    setCurrentPlatform(detail.platform ?? platforms[0]?.platform ?? '');
  }, [detailRequest.data, visible, value, platforms]);

  const handleSubmit = async () => {
    try {
      const values = await formApiRef.current?.validate();
      if (!values) {
        return;
      }

      const payload: CreateIMChannelConfigRequest = {
        platform: values.platform,
        name: values.name,
        space_id: spaceId,
        bot_id: values.bot_id,
        connector_id: values.connector_id || platformInfo?.connector_id,
        tenant_key: values.tenant_key,
        app_id: values.app_id,
        bot_code: values.bot_code,
        session_scope: values.session_scope,
        status: values.status,
        platform_config: values.platform_config,
        ext: parseExtJSON(values.ext_json),
      };

      await saveRequest.runAsync(payload);
    } catch (error) {
      if (error instanceof Error && error.message) {
        Toast.error(error.message);
      }
    }
  };

  return (
    <SideSheet
      visible={visible}
      width={560}
      title={I18n.t(
        mode === 'edit' ? 'im_channel_edit' : 'im_channel_create',
        {},
        mode === 'edit' ? 'Edit channel' : 'Create channel',
      )}
      onCancel={onCancel}
    >
      {detailRequest.loading ? (
        <div className="flex h-full items-center justify-center">
          <Spin spinning={true} />
        </div>
      ) : (
        <div className="flex h-full flex-col gap-[16px]">
          <Form<ChannelFormValues>
            getFormApi={api => {
              formApiRef.current = api;
            }}
            onValueChange={values => setCurrentPlatform(values.platform ?? '')}
          >
            <FormInput
              field="name"
              label={I18n.t('im_channel_name', {}, 'Channel name')}
              rules={[{ required: true }]}
              maxLength={64}
            />
            <FormSelect
              field="platform"
              label={I18n.t('im_platform', {}, 'Platform')}
              optionList={platformOptions}
              rules={[{ required: true }]}
              disabled={mode === 'edit'}
            />
            <div className="mb-[16px] rounded-[16px] coz-bg-primary p-[12px]">
              <Typography.Text strong>
                {I18n.t('callback_channel_config', {}, 'Callback')}
              </Typography.Text>
              <div className="mt-[8px] text-[12px] coz-fg-secondary">
                {I18n.t('im_callback_url', {}, 'Callback URL')}
              </div>
              <Typography.Text ellipsis={{ showTooltip: true }}>
                {platformInfo?.callback_url ?? '-'}
              </Typography.Text>
              <div className="mt-[8px] text-[12px] coz-fg-secondary">
                {I18n.t('im_connector_id', {}, 'Connector ID')}
              </div>
              <Typography.Text>{platformInfo?.connector_id ?? '-'}</Typography.Text>
            </div>
            <FormInput
              field="bot_id"
              label={I18n.t('im_bot_id', {}, 'Bot ID')}
              maxLength={32}
            />
            <FormInput
              field="tenant_key"
              label={I18n.t('im_tenant_key', {}, 'Tenant key')}
              maxLength={128}
            />
            <FormInput
              field="app_id"
              label={I18n.t('im_app_id', {}, 'App ID')}
              maxLength={128}
            />
            <FormInput
              field="bot_code"
              label={I18n.t('im_bot_code', {}, 'Bot code')}
              maxLength={128}
            />
            <FormSelect
              field="session_scope"
              label={I18n.t('im_session_scope', {}, 'Session scope')}
              optionList={SESSION_SCOPE_OPTIONS}
            />
            <FormSelect
              field="status"
              label={I18n.t('api_status_1', {}, 'Status')}
              optionList={CHANNEL_STATUS_OPTIONS}
            />
            <FormTextArea
              field="platform_config"
              label={I18n.t('im_platform_config', {}, 'Platform config JSON')}
              autosize={{ minRows: 5, maxRows: 12 }}
            />
            <FormTextArea
              field="ext_json"
              label={I18n.t('im_extra_meta', {}, 'Extra metadata JSON')}
              autosize={{ minRows: 3, maxRows: 8 }}
            />
          </Form>
          <div className="mt-auto flex justify-end gap-[8px] border-t coz-stroke-primary pt-[16px]">
            <Button color="secondary" onClick={onCancel}>
              {I18n.t('Cancel', {}, 'Cancel')}
            </Button>
            <Button loading={saveRequest.loading} onClick={handleSubmit}>
              {I18n.t('Confirm', {}, 'Confirm')}
            </Button>
          </div>
        </div>
      )}
    </SideSheet>
  );
};
