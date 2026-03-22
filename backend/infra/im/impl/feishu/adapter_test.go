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

package feishu

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
			AppID:             "cli_test",
			AppSecret:         "secret",
			VerificationToken: "token",
			AgentMap: map[string]int64{
				"cli_test": 1001,
			},
		},
		client: newClient(Config{}),
	}

	testCases := []struct {
		name       string
		body       string
		assertFunc func(t *testing.T, msg *imEntity.IMessage, resp *imEntity.RawResponse, err error)
	}{
		{
			name: "challenge",
			body: `{"token":"token","challenge":"challenge_value"}`,
			assertFunc: func(t *testing.T, msg *imEntity.IMessage, resp *imEntity.RawResponse, err error) {
				require.NoError(t, err)
				require.Nil(t, msg)
				require.NotNil(t, resp)
				require.JSONEq(t, `{"challenge":"challenge_value"}`, string(resp.Body))
			},
		},
		{
			name: "group at bot text message",
			body: `{
				"token":"token",
				"header":{"event_type":"im.message.receive_v1","event_id":"evt_1","app_id":"cli_test","tenant_key":"tenant_1"},
				"event":{
					"sender":{"sender_id":{"user_id":"ou_user_1"},"sender_type":"user","sender_name":"Alice"},
					"message":{
						"message_id":"om_1",
						"chat_id":"oc_group_1",
						"thread_id":"omt_1",
						"chat_type":"group",
						"message_type":"text",
						"content":"{\"text\":\"<at user_id=\\\"ou_bot\\\">Bot</at> hello feishu\"}",
						"mentions":[{"key":"@_user_1","name":"Bot","id":{"open_id":"ou_bot"},"type":"user"}]
					}
				}
			}`,
			assertFunc: func(t *testing.T, msg *imEntity.IMessage, resp *imEntity.RawResponse, err error) {
				require.NoError(t, err)
				require.Nil(t, resp)
				require.NotNil(t, msg)
				require.Equal(t, int64(1001), msg.AgentID)
				require.Equal(t, "evt_1", msg.EventID)
				require.Equal(t, "om_1", msg.MessageID)
				require.Equal(t, "Alice", msg.UserName)
				require.Equal(t, "ou_user_1", msg.UserID)
				require.Equal(t, "hello feishu", msg.Text)
				require.True(t, msg.IsAtBot)
				require.NotEmpty(t, msg.SessionID)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg, resp, err := adapter.ParseCallback(context.Background(), &imEntity.CallbackRequest{
				Platform: imEntity.PlatformFeishu,
				Body:     []byte(tc.body),
				Headers:  map[string]string{},
				Query:    map[string]string{},
			})

			tc.assertFunc(t, msg, resp, err)
		})
	}
}
