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

import { type FC, useMemo, useState } from 'react';

import { useRequest } from 'ahooks';
import { I18n } from '@coze-arch/i18n';
import {
  Button,
  Layout,
  Search,
  Select,
  Space,
  Table,
  Tabs,
  TabPane,
  Toast,
  Typography,
  type ColumnProps,
} from '@coze-arch/coze-design';
import { im as imApi } from '@coze-studio/api-schema';

import {
  CHANNEL_FILTER_DEFAULT,
  CHANNEL_STATUS_OPTIONS,
  TASK_FILTER_DEFAULT,
  TASK_STATUS_OPTIONS,
} from './constants';
import { ChannelConfigSideSheet } from './components/channel-config-sidesheet';
import { ConnectivityTestPanel } from './components/connectivity-test-panel';
import { TaskDetailSideSheet } from './components/task-detail-sidesheet';
import { IMStatusTag } from './components/status-tag';
import { formatDateTime, getPlatformLabel } from './utils';
import type {
  ChannelFilterState,
  IMChannelConfig,
  IMPageTab,
  IMPlatformRecord,
  IMTaskRecord,
  TaskFilterState,
} from './types';

export const IMManagePage: FC<{ spaceId: string }> = ({ spaceId }) => {
  const [activeTab, setActiveTab] = useState<IMPageTab>('channels');
  const [channelFilters, setChannelFilters] =
    useState<ChannelFilterState>(CHANNEL_FILTER_DEFAULT);
  const [taskFilters, setTaskFilters] =
    useState<TaskFilterState>(TASK_FILTER_DEFAULT);
  const [channelSheetVisible, setChannelSheetVisible] = useState(false);
  const [channelSheetMode, setChannelSheetMode] = useState<'create' | 'edit'>(
    'create',
  );
  const [editingChannel, setEditingChannel] = useState<IMChannelConfig>();
  const [selectedTestConfigId, setSelectedTestConfigId] = useState<string>();
  const [selectedTask, setSelectedTask] = useState<IMTaskRecord>();
  const [taskDetailVisible, setTaskDetailVisible] = useState(false);

  const platformRequest = useRequest(
    async () => {
      const resp = await imApi.ListIMPlatforms(
        {},
        { __disableErrorToast: true },
      );
      return resp.data ?? [];
    },
    {
      refreshDeps: [spaceId],
    },
  );

  const channelListRequest = useRequest(
    async () => {
      try {
        const resp = await imApi.ListIMChannelConfigs(
          {
            space_id: spaceId,
            platform: channelFilters.platform || undefined,
            status: channelFilters.status || undefined,
            keyword: channelFilters.keyword || undefined,
          },
          { __disableErrorToast: true },
        );
        return resp.data?.list ?? [];
      } catch {
        return [] as IMChannelConfig[];
      }
    },
    {
      refreshDeps: [
        spaceId,
        channelFilters.platform,
        channelFilters.status,
        channelFilters.keyword,
      ],
    },
  );

  const taskListRequest = useRequest(
    async () => {
      try {
        const resp = await imApi.ListIMTaskRecords(
          {
            space_id: spaceId,
            platform: taskFilters.platform || undefined,
            status: taskFilters.status || undefined,
            config_id: taskFilters.configId || undefined,
            task_id: taskFilters.taskId || undefined,
          },
          { __disableErrorToast: true },
        );
        return resp.data?.list ?? [];
      } catch {
        return [] as IMTaskRecord[];
      }
    },
    {
      refreshDeps: [
        spaceId,
        taskFilters.platform,
        taskFilters.status,
        taskFilters.configId,
        taskFilters.taskId,
      ],
    },
  );

  const retryTaskRequest = useRequest(
    async (taskId: string) =>
      imApi.RetryIMTask({ task_id: taskId }, { __disableErrorToast: true }),
    {
      manual: true,
      onSuccess: () => {
        Toast.success(I18n.t('retry_success', {}, 'Retry submitted'));
        taskListRequest.refresh();
      },
      onError: error => {
        Toast.error(
          error instanceof Error
            ? error.message
            : I18n.t('retry_failed', {}, 'Retry failed'),
        );
      },
    },
  );

  const platformOptions = useMemo(
    () => [
      {
        label: I18n.t('library_filter_tags_all_types', {}, 'All platforms'),
        value: '',
      },
      ...(platformRequest.data ?? []).map((item: IMPlatformRecord) => ({
        label: item.name,
        value: item.platform,
      })),
    ],
    [platformRequest.data],
  );

  const channelOptions = useMemo(
    () =>
      (channelListRequest.data ?? []).map(item => ({
        label: item.name || item.config_id || '-',
        value: item.config_id || '',
      })),
    [channelListRequest.data],
  );

  const openCreateChannel = () => {
    setEditingChannel(undefined);
    setChannelSheetMode('create');
    setChannelSheetVisible(true);
  };

  const openEditChannel = (record: IMChannelConfig) => {
    setEditingChannel(record);
    setChannelSheetMode('edit');
    setChannelSheetVisible(true);
  };

  const openConnectivityTest = (configId?: string) => {
    setSelectedTestConfigId(configId);
    setActiveTab('connectivity');
  };

  const openTaskList = (configId?: string) => {
    setTaskFilters(prev => ({
      ...prev,
      configId: configId ?? '',
    }));
    setActiveTab('tasks');
  };

  const openTaskDetail = (task: IMTaskRecord) => {
    setSelectedTask(task);
    setTaskDetailVisible(true);
  };

  const channelColumns: ColumnProps<IMChannelConfig>[] = [
    {
      title: I18n.t('im_channel_name', {}, 'Channel'),
      dataIndex: 'name',
      width: 220,
      render: (_value, record) => (
        <div className="flex flex-col gap-[4px]">
          <Typography.Text strong>{record.name || '-'}</Typography.Text>
          <Typography.Text size="small" className="coz-fg-secondary">
            {record.callback_path || '-'}
          </Typography.Text>
        </div>
      ),
    },
    {
      title: I18n.t('im_platform', {}, 'Platform'),
      dataIndex: 'platform',
      width: 120,
      render: value => getPlatformLabel(String(value)),
    },
    {
      title: I18n.t('callback_channel_config', {}, 'Binding'),
      dataIndex: 'bot_id',
      width: 220,
      render: (_value, record) => (
        <div className="flex flex-col gap-[4px]">
          <Typography.Text>
            {I18n.t('im_bot_id', {}, 'Bot ID')}: {record.bot_id || '-'}
          </Typography.Text>
          <Typography.Text size="small" className="coz-fg-secondary">
            {I18n.t('im_tenant_key', {}, 'Tenant key')}: {record.tenant_key || '-'}
          </Typography.Text>
        </div>
      ),
    },
    {
      title: I18n.t('im_callback_url', {}, 'Callback URL'),
      dataIndex: 'callback_url',
      width: 260,
      render: value => (
        <Typography.Text ellipsis={{ showTooltip: true }}>
          {String(value || '-')}
        </Typography.Text>
      ),
    },
    {
      title: I18n.t('im_session_scope', {}, 'Session scope'),
      dataIndex: 'session_scope',
      width: 140,
      render: value => String(value || '-'),
    },
    {
      title: I18n.t('api_status_1', {}, 'Status'),
      dataIndex: 'status',
      width: 120,
      render: value => <IMStatusTag value={String(value)} type="channel" />,
    },
    {
      title: I18n.t('library_edited_time', {}, 'Updated time'),
      dataIndex: 'updated_at',
      width: 180,
      render: value => formatDateTime(String(value)),
    },
    {
      title: I18n.t('library_actions', {}, 'Actions'),
      dataIndex: 'config_id',
      width: 220,
      fixed: 'right',
      render: (_value, record) => (
        <Space spacing={8}>
          <Typography.Text link onClick={() => openEditChannel(record)}>
            {I18n.t('Edit', {}, 'Edit')}
          </Typography.Text>
          <Typography.Text
            link
            onClick={() => openConnectivityTest(record.config_id)}
          >
            {I18n.t('im_test_connectivity', {}, 'Test')}
          </Typography.Text>
          <Typography.Text link onClick={() => openTaskList(record.config_id)}>
            {I18n.t('im_task_records', {}, 'Tasks')}
          </Typography.Text>
        </Space>
      ),
    },
  ];

  const taskColumns: ColumnProps<IMTaskRecord>[] = [
    {
      title: I18n.t('im_task_id', {}, 'Task ID'),
      dataIndex: 'task_id',
      width: 220,
      render: value => (
        <Typography.Text ellipsis={{ showTooltip: true }}>
          {String(value || '-')}
        </Typography.Text>
      ),
    },
    {
      title: I18n.t('im_platform', {}, 'Platform'),
      dataIndex: 'platform',
      width: 120,
      render: value => getPlatformLabel(String(value)),
    },
    {
      title: I18n.t('api_status_1', {}, 'Status'),
      dataIndex: 'status',
      width: 120,
      render: value => <IMStatusTag value={String(value)} type="task" />,
    },
    {
      title: I18n.t('im_task_retry', {}, 'Retry'),
      dataIndex: 'retry_count',
      width: 120,
      render: (_value, record) =>
        `${record.retry_count ?? 0}/${record.max_retry_count ?? 0}`,
    },
    {
      title: I18n.t('trace_id', {}, 'Trace ID'),
      dataIndex: 'trace_id',
      width: 180,
      render: value => (
        <Typography.Text ellipsis={{ showTooltip: true }}>
          {String(value || '-')}
        </Typography.Text>
      ),
    },
    {
      title: I18n.t('updated_time', {}, 'Updated time'),
      dataIndex: 'updated_at',
      width: 180,
      render: value => formatDateTime(String(value)),
    },
    {
      title: I18n.t('library_actions', {}, 'Actions'),
      dataIndex: 'task_id',
      width: 180,
      fixed: 'right',
      render: (_value, record) => (
        <Space spacing={8}>
          <Typography.Text link onClick={() => openTaskDetail(record)}>
            {I18n.t('workflow_trigger_user_create_list_read', {}, 'View')}
          </Typography.Text>
          <Typography.Text
            link
            onClick={() => {
              if (!record.task_id) {
                return;
              }

              retryTaskRequest.run(record.task_id);
            }}
          >
            {I18n.t('retry', {}, 'Retry')}
          </Typography.Text>
        </Space>
      ),
    },
  ];

  return (
    <>
      <Layout className="h-full overflow-hidden">
        <Layout.Header className="pb-0">
          <div className="flex flex-wrap items-start justify-between gap-[16px]">
            <div>
              <Typography.Title heading={4} className="!mb-[8px]">
                {I18n.t('navigation_workspace_im', {}, 'IM')}
              </Typography.Title>
              <Typography.Paragraph className="!mb-0 coz-fg-secondary">
                {I18n.t(
                  'im_workspace_desc',
                  {},
                  'Manage IM channel bindings, connectivity checks and async task records.',
                )}
              </Typography.Paragraph>
            </div>
            {activeTab === 'channels' ? (
              <Button onClick={openCreateChannel}>
                {I18n.t('im_channel_create', {}, 'Create channel')}
              </Button>
            ) : null}
          </div>
        </Layout.Header>
        <Layout.Content className="!h-auto !min-h-0 !flex-1 overflow-auto pb-[24px]">
          <Tabs
            className="flex min-h-0 flex-1 flex-col [&_.semi-tabs-content]:min-h-0 [&_.semi-tabs-content]:flex-1"
            activeKey={activeTab}
            onChange={key => setActiveTab(key as IMPageTab)}
          >
            <TabPane
              itemKey="channels"
              tab={I18n.t('im_channel_configs', {}, 'Channel configs')}
            >
              <div className="mb-[16px] mt-[12px] flex flex-wrap items-center justify-between gap-[12px]">
                <Space spacing={8}>
                  <Select
                    className="min-w-[168px]"
                    optionList={platformOptions}
                    value={channelFilters.platform}
                    onChange={value =>
                      setChannelFilters(prev => ({
                        ...prev,
                        platform: String(value || ''),
                      }))
                    }
                  />
                  <Select
                    className="min-w-[168px]"
                    optionList={[
                      {
                        label: I18n.t('library_all_status', {}, 'All status'),
                        value: '',
                      },
                      ...CHANNEL_STATUS_OPTIONS,
                    ]}
                    value={channelFilters.status}
                    onChange={value =>
                      setChannelFilters(prev => ({
                        ...prev,
                        status: String(value || ''),
                      }))
                    }
                  />
                </Space>
                <Space spacing={8}>
                  <Search
                    width={260}
                    showClear={true}
                    placeholder={I18n.t(
                      'im_search_channel',
                      {},
                      'Search channel name',
                    )}
                    onSearch={value =>
                      setChannelFilters(prev => ({
                        ...prev,
                        keyword: value,
                      }))
                    }
                  />
                  <Button
                    color="secondary"
                    onClick={() => channelListRequest.refresh()}
                  >
                    {I18n.t('workflow_trigger_user_create_refresh', {}, 'Refresh')}
                  </Button>
                </Space>
              </div>
              <Table
                offsetY={248}
                tableProps={{
                  rowKey: 'config_id',
                  loading: channelListRequest.loading,
                  columns: channelColumns,
                  dataSource: channelListRequest.data,
                  pagination: false,
                  scroll: { x: 1460 },
                }}
              />
            </TabPane>
            <TabPane
              itemKey="connectivity"
              tab={I18n.t('im_connectivity_test', {}, 'Connectivity')}
            >
              <div className="mt-[12px]">
                <ConnectivityTestPanel
                  spaceId={spaceId}
                  channels={channelListRequest.data ?? []}
                  defaultConfigId={selectedTestConfigId}
                />
              </div>
            </TabPane>
            <TabPane
              itemKey="tasks"
              tab={I18n.t('im_task_records', {}, 'Task records')}
            >
              <div className="mb-[16px] mt-[12px] flex flex-wrap items-center justify-between gap-[12px]">
                <Space spacing={8}>
                  <Select
                    className="min-w-[168px]"
                    optionList={platformOptions}
                    value={taskFilters.platform}
                    onChange={value =>
                      setTaskFilters(prev => ({
                        ...prev,
                        platform: String(value || ''),
                      }))
                    }
                  />
                  <Select
                    className="min-w-[168px]"
                    optionList={[
                      {
                        label: I18n.t('library_all_status', {}, 'All status'),
                        value: '',
                      },
                      ...TASK_STATUS_OPTIONS,
                    ]}
                    value={taskFilters.status}
                    onChange={value =>
                      setTaskFilters(prev => ({
                        ...prev,
                        status: String(value || ''),
                      }))
                    }
                  />
                  <Select
                    className="min-w-[220px]"
                    optionList={[
                      {
                        label: I18n.t(
                          'im_all_channel_configs',
                          {},
                          'All channel configs',
                        ),
                        value: '',
                      },
                      ...channelOptions,
                    ]}
                    value={taskFilters.configId}
                    onChange={value =>
                      setTaskFilters(prev => ({
                        ...prev,
                        configId: String(value || ''),
                      }))
                    }
                  />
                </Space>
                <Space spacing={8}>
                  <Search
                    width={260}
                    showClear={true}
                    placeholder={I18n.t('im_search_task', {}, 'Search task ID')}
                    onSearch={value =>
                      setTaskFilters(prev => ({
                        ...prev,
                        taskId: value,
                      }))
                    }
                  />
                  <Button color="secondary" onClick={() => taskListRequest.refresh()}>
                    {I18n.t('workflow_trigger_user_create_refresh', {}, 'Refresh')}
                  </Button>
                </Space>
              </div>
              <Table
                offsetY={248}
                tableProps={{
                  rowKey: 'task_id',
                  loading: taskListRequest.loading,
                  columns: taskColumns,
                  dataSource: taskListRequest.data,
                  pagination: false,
                  scroll: { x: 1120 },
                }}
              />
            </TabPane>
          </Tabs>
        </Layout.Content>
      </Layout>
      <ChannelConfigSideSheet
        visible={channelSheetVisible}
        mode={channelSheetMode}
        spaceId={spaceId}
        value={editingChannel}
        platforms={platformRequest.data ?? []}
        onCancel={() => setChannelSheetVisible(false)}
        onSuccess={() => {
          setChannelSheetVisible(false);
          channelListRequest.refresh();
        }}
      />
      <TaskDetailSideSheet
        visible={taskDetailVisible}
        task={selectedTask}
        onCancel={() => setTaskDetailVisible(false)}
      />
    </>
  );
};
