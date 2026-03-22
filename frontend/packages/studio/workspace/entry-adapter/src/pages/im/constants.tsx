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

import { I18n } from '@coze-arch/i18n';
import { type TagProps } from '@coze-arch/coze-design';

export const CHANNEL_FILTER_DEFAULT = {
  keyword: '',
  platform: '',
  status: '',
};

export const TASK_FILTER_DEFAULT = {
  platform: '',
  status: '',
  taskId: '',
  configId: '',
};

export const PLATFORM_TEXT_MAP: Record<string, string> = {
  feishu: I18n.t('im_platform_feishu', {}, 'Feishu'),
  dingtalk: I18n.t('im_platform_dingtalk', {}, 'DingTalk'),
  wecom: I18n.t('im_platform_wecom', {}, 'WeCom'),
};

export const CHANNEL_STATUS_META_MAP: Record<
  string,
  Pick<TagProps, 'color'> & { text: string }
> = {
  enabled: {
    color: 'green',
    text: I18n.t('im_status_enabled', {}, 'Enabled'),
  },
  disabled: {
    color: 'grey',
    text: I18n.t('im_status_disabled', {}, 'Disabled'),
  },
};

export const TASK_STATUS_META_MAP: Record<
  string,
  Pick<TagProps, 'color'> & { text: string }
> = {
  pending: {
    color: 'primary',
    text: I18n.t('im_task_status_pending', {}, 'Pending'),
  },
  running: {
    color: 'brand',
    text: I18n.t('im_task_status_running', {}, 'Running'),
  },
  success: {
    color: 'green',
    text: I18n.t('im_task_status_success', {}, 'Success'),
  },
  failed: {
    color: 'red',
    text: I18n.t('im_task_status_failed', {}, 'Failed'),
  },
  retrying: {
    color: 'yellow',
    text: I18n.t('im_task_status_retrying', {}, 'Retrying'),
  },
};

export const SESSION_SCOPE_OPTIONS = [
  {
    label: I18n.t('im_session_scope_chat', {}, 'Group shared'),
    value: 'chat',
  },
  {
    label: I18n.t('im_session_scope_thread', {}, 'Thread'),
    value: 'thread',
  },
  {
    label: I18n.t('im_session_scope_user', {}, '1:1 user'),
    value: 'user',
  },
  {
    label: I18n.t('im_session_scope_chat_user', {}, 'Group user isolated'),
    value: 'chat_user',
  },
];

export const CHANNEL_STATUS_OPTIONS = [
  {
    label: I18n.t('im_status_enabled', {}, 'Enabled'),
    value: 'enabled',
  },
  {
    label: I18n.t('im_status_disabled', {}, 'Disabled'),
    value: 'disabled',
  },
];

export const TASK_STATUS_OPTIONS = [
  {
    label: I18n.t('im_task_status_pending', {}, 'Pending'),
    value: 'pending',
  },
  {
    label: I18n.t('im_task_status_running', {}, 'Running'),
    value: 'running',
  },
  {
    label: I18n.t('im_task_status_success', {}, 'Success'),
    value: 'success',
  },
  {
    label: I18n.t('im_task_status_failed', {}, 'Failed'),
    value: 'failed',
  },
  {
    label: I18n.t('im_task_status_retrying', {}, 'Retrying'),
    value: 'retrying',
  },
];
