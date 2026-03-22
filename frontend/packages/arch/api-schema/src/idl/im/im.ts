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

import { createAPI } from './../../api/config';

export type IMPlatform = 'feishu' | 'dingtalk' | 'wecom';
export type IMChannelStatus = 'enabled' | 'disabled';
export type IMSessionScope = 'user' | 'chat' | 'thread' | 'chat_user';
export type IMTaskStatus =
  | 'pending'
  | 'running'
  | 'success'
  | 'failed'
  | 'retrying';

export interface ListIMPlatformsRequest {}

export interface IMPlatformInfo {
  platform: IMPlatform | string;
  name: string;
  connector_id: string;
  callback_path: string;
  callback_url: string;
  enabled: boolean;
}

export interface ListIMPlatformsResponse {
  data: IMPlatformInfo[];
  code: number;
  msg: string;
}

export interface IMChannelConfigData {
  config_id?: string;
  platform?: IMPlatform | string;
  name?: string;
  space_id?: string;
  bot_id?: string;
  connector_id?: string;
  tenant_key?: string;
  app_id?: string;
  bot_code?: string;
  callback_path?: string;
  callback_url?: string;
  session_scope?: IMSessionScope | string;
  status?: IMChannelStatus | string;
  platform_config?: string;
  masked_platform_config?: string;
  ext?: Record<string, string>;
  created_at?: string;
  updated_at?: string;
}

export interface ListIMChannelConfigsRequest {
  space_id?: string;
  platform?: IMPlatform | string;
  status?: IMChannelStatus | string;
  keyword?: string;
  page?: number;
  page_size?: number;
}

export interface ListIMChannelConfigsResponseData {
  list?: IMChannelConfigData[];
  total?: number;
}

export interface ListIMChannelConfigsResponse {
  data?: ListIMChannelConfigsResponseData;
  code: number;
  msg: string;
}

export interface GetIMChannelConfigRequest {
  config_id: string;
}

export interface GetIMChannelConfigResponse {
  data?: IMChannelConfigData;
  code: number;
  msg: string;
}

export interface CreateIMChannelConfigRequest {
  platform?: IMPlatform | string;
  name?: string;
  space_id?: string;
  bot_id?: string;
  connector_id?: string;
  tenant_key?: string;
  app_id?: string;
  bot_code?: string;
  session_scope?: IMSessionScope | string;
  status?: IMChannelStatus | string;
  platform_config?: string;
  ext?: Record<string, string>;
}

export interface CreateIMChannelConfigResponse {
  data?: IMChannelConfigData;
  code: number;
  msg: string;
}

export interface UpdateIMChannelConfigRequest {
  config_id: string;
  name?: string;
  bot_id?: string;
  connector_id?: string;
  tenant_key?: string;
  app_id?: string;
  bot_code?: string;
  session_scope?: IMSessionScope | string;
  status?: IMChannelStatus | string;
  platform_config?: string;
  ext?: Record<string, string>;
}

export interface UpdateIMChannelConfigResponse {
  data?: IMChannelConfigData;
  code: number;
  msg: string;
}

export interface TestIMChannelConnectivityRequest {
  config_id?: string;
  space_id?: string;
  platform?: IMPlatform | string;
  sender_id?: string;
  sender_name?: string;
  chat_id?: string;
  thread_id?: string;
  message_text?: string;
  is_at_bot?: boolean;
  prefer_async?: boolean;
}

export interface IMConnectivityTestData {
  accepted?: boolean;
  status?: string;
  task_id?: string;
  trace_id?: string;
  message?: string;
  reply_preview?: string;
}

export interface TestIMChannelConnectivityResponse {
  data?: IMConnectivityTestData;
  code: number;
  msg: string;
}

export interface IMTaskData {
  task_id?: string;
  platform?: IMPlatform | string;
  config_id?: string;
  task_type?: string;
  status?: IMTaskStatus | string;
  conversation_id?: string;
  run_id?: string;
  bot_id?: string;
  event_id?: string;
  session_id?: string;
  retry_count?: number;
  max_retry_count?: number;
  next_retry_at?: string;
  error_code?: number;
  error_msg?: string;
  trace_id?: string;
  ext?: Record<string, string>;
  message_snapshot?: string;
  request_snapshot?: string;
  gateway_ticket?: string;
  result_snapshot?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ListIMTaskRecordsRequest {
  space_id?: string;
  platform?: IMPlatform | string;
  status?: IMTaskStatus | string;
  config_id?: string;
  task_id?: string;
  page?: number;
  page_size?: number;
}

export interface ListIMTaskRecordsResponseData {
  list?: IMTaskData[];
  total?: number;
}

export interface ListIMTaskRecordsResponse {
  data?: ListIMTaskRecordsResponseData;
  code: number;
  msg: string;
}

export interface GetIMTaskDetailRequest {
  task_id: string;
}

export interface GetIMTaskDetailResponse {
  data?: IMTaskData;
  code: number;
  msg: string;
}

export interface RetryIMTaskRequest {
  task_id: string;
}

export interface RetryIMTaskResponse {
  data?: IMTaskData;
  code: number;
  msg: string;
}

export const ListIMPlatforms = /*#__PURE__*/ createAPI<
  ListIMPlatformsRequest,
  ListIMPlatformsResponse
>({
  url: '/api/im/platforms',
  method: 'GET',
  name: 'ListIMPlatforms',
  reqType: 'ListIMPlatformsRequest',
  reqMapping: {},
  resType: 'ListIMPlatformsResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const ListIMChannelConfigs = /*#__PURE__*/ createAPI<
  ListIMChannelConfigsRequest,
  ListIMChannelConfigsResponse
>({
  url: '/api/im/channel_config/list',
  method: 'GET',
  name: 'ListIMChannelConfigs',
  reqType: 'ListIMChannelConfigsRequest',
  reqMapping: {
    query: ['space_id', 'platform', 'status', 'keyword', 'page', 'page_size'],
  },
  resType: 'ListIMChannelConfigsResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const GetIMChannelConfig = /*#__PURE__*/ createAPI<
  GetIMChannelConfigRequest,
  GetIMChannelConfigResponse
>({
  url: '/api/im/channel_config/get',
  method: 'GET',
  name: 'GetIMChannelConfig',
  reqType: 'GetIMChannelConfigRequest',
  reqMapping: {
    query: ['config_id'],
  },
  resType: 'GetIMChannelConfigResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const CreateIMChannelConfig = /*#__PURE__*/ createAPI<
  CreateIMChannelConfigRequest,
  CreateIMChannelConfigResponse
>({
  url: '/api/im/channel_config/create',
  method: 'POST',
  name: 'CreateIMChannelConfig',
  reqType: 'CreateIMChannelConfigRequest',
  reqMapping: {
    body: [
      'platform',
      'name',
      'space_id',
      'bot_id',
      'connector_id',
      'tenant_key',
      'app_id',
      'bot_code',
      'session_scope',
      'status',
      'platform_config',
      'ext',
    ],
  },
  resType: 'CreateIMChannelConfigResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const UpdateIMChannelConfig = /*#__PURE__*/ createAPI<
  UpdateIMChannelConfigRequest,
  UpdateIMChannelConfigResponse
>({
  url: '/api/im/channel_config/update',
  method: 'POST',
  name: 'UpdateIMChannelConfig',
  reqType: 'UpdateIMChannelConfigRequest',
  reqMapping: {
    body: [
      'config_id',
      'name',
      'bot_id',
      'connector_id',
      'tenant_key',
      'app_id',
      'bot_code',
      'session_scope',
      'status',
      'platform_config',
      'ext',
    ],
  },
  resType: 'UpdateIMChannelConfigResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const TestIMChannelConnectivity = /*#__PURE__*/ createAPI<
  TestIMChannelConnectivityRequest,
  TestIMChannelConnectivityResponse
>({
  url: '/api/im/channel_config/test',
  method: 'POST',
  name: 'TestIMChannelConnectivity',
  reqType: 'TestIMChannelConnectivityRequest',
  reqMapping: {
    body: [
      'config_id',
      'space_id',
      'platform',
      'sender_id',
      'sender_name',
      'chat_id',
      'thread_id',
      'message_text',
      'is_at_bot',
      'prefer_async',
    ],
  },
  resType: 'TestIMChannelConnectivityResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const ListIMTaskRecords = /*#__PURE__*/ createAPI<
  ListIMTaskRecordsRequest,
  ListIMTaskRecordsResponse
>({
  url: '/api/im/task/list',
  method: 'GET',
  name: 'ListIMTaskRecords',
  reqType: 'ListIMTaskRecordsRequest',
  reqMapping: {
    query: ['space_id', 'platform', 'status', 'config_id', 'task_id', 'page', 'page_size'],
  },
  resType: 'ListIMTaskRecordsResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const GetIMTaskDetail = /*#__PURE__*/ createAPI<
  GetIMTaskDetailRequest,
  GetIMTaskDetailResponse
>({
  url: '/api/im/task/get',
  method: 'GET',
  name: 'GetIMTaskDetail',
  reqType: 'GetIMTaskDetailRequest',
  reqMapping: {
    query: ['task_id'],
  },
  resType: 'GetIMTaskDetailResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});

export const RetryIMTask = /*#__PURE__*/ createAPI<
  RetryIMTaskRequest,
  RetryIMTaskResponse
>({
  url: '/api/im/task/retry',
  method: 'POST',
  name: 'RetryIMTask',
  reqType: 'RetryIMTaskRequest',
  reqMapping: {
    body: ['task_id'],
  },
  resType: 'RetryIMTaskResponse',
  schemaRoot: 'api://schemas/idl_im_im',
  service: 'im',
});
