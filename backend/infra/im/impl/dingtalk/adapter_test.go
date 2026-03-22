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

package dingtalk

import (
	"context"
	"testing"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	"github.com/stretchr/testify/require"
)

func TestAdapterParseCallback(t *testing.T) {
	t.Parallel()

	adapter := &Adapter{
		cfg: Config{
			AppKey: "ding_test_app",
			AgentMap: map[string]int64{
				"ding_test_bot": 1002,
			},
		},
		client: newClient(),
	}

	msg, resp, err := adapter.ParseCallback(context.Background(), &imEntity.CallbackRequest{
		Platform: imEntity.PlatformDingTalk,
		Body: []byte(`{
			"msgtype":"text",
			"msgId":"msg_1",
			"conversationId":"cid_group_1",
			"conversationType":"2",
			"conversationTitle":"Test Group",
			"sessionWebhook":"https://example.com/session",
			"senderStaffId":"staff_1",
			"senderNick":"Alice",
			"robotCode":"ding_test_bot",
			"text":{"content":"@机器人 hello dingtalk"},
			"atUsers":[{"dingtalkId":"bot_1","staffId":"bot_staff"}]
		}`),
		Headers: map[string]string{},
		Query:   map[string]string{},
	})

	require.NoError(t, err)
	require.Nil(t, resp)
	require.NotNil(t, msg)
	require.Equal(t, int64(1002), msg.AgentID)
	require.Equal(t, "msg_1", msg.MessageID)
	require.Equal(t, "staff_1", msg.UserID)
	require.Equal(t, "Alice", msg.UserName)
	require.Equal(t, "hello dingtalk", msg.Text)
	require.True(t, msg.IsAtBot)
	require.Equal(t, "https://example.com/session", msg.Metadata["session_webhook"])
}
