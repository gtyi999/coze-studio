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

import type {
  IMChannelConfigData,
  IMConnectivityTestData,
  IMPlatformInfo,
  IMTaskData,
} from '@coze-studio/api-schema/im';

export type IMChannelConfig = IMChannelConfigData;
export type IMTaskRecord = IMTaskData;
export type IMPlatformRecord = IMPlatformInfo;
export type IMConnectivityResult = IMConnectivityTestData;

export type IMPageTab = 'channels' | 'connectivity' | 'tasks';

export interface ChannelFilterState {
  keyword: string;
  platform: string;
  status: string;
}

export interface TaskFilterState {
  platform: string;
  status: string;
  taskId: string;
  configId: string;
}

export interface ChannelFormValues {
  name?: string;
  platform?: string;
  bot_id?: string;
  connector_id?: string;
  tenant_key?: string;
  app_id?: string;
  bot_code?: string;
  session_scope?: string;
  status?: string;
  platform_config?: string;
  ext_json?: string;
}

export interface ConnectivityTestFormValues {
  config_id?: string;
  sender_id?: string;
  sender_name?: string;
  chat_id?: string;
  thread_id?: string;
  message_text?: string;
  is_at_bot?: boolean;
  prefer_async?: boolean;
}
