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

package wecom

import (
	"context"
	"testing"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	"github.com/stretchr/testify/require"
)

func TestAdapterParseVerifyRequest(t *testing.T) {
	t.Parallel()

	adapter := &Adapter{
		cfg: Config{
			CorpID:         "ww_test",
			Token:          "token",
			EncodingAESKey: "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
			AgentMap: map[string]int64{
				"1000002": 1003,
			},
		},
		client: newClient(Config{}),
	}

	echo, err := encryptMessage(adapter.cfg.EncodingAESKey, adapter.cfg.CorpID, "verify_ok")
	require.NoError(t, err)

	msg, resp, parseErr := adapter.ParseCallback(context.Background(), &imEntity.CallbackRequest{
		Platform: imEntity.PlatformWeCom,
		Method:   "GET",
		Query: map[string]string{
			"timestamp":     "1710000000",
			"nonce":         "nonce_1",
			"echostr":       echo,
			"msg_signature": signature(adapter.cfg.Token, "1710000000", "nonce_1", echo),
		},
	})

	require.NoError(t, parseErr)
	require.Nil(t, msg)
	require.NotNil(t, resp)
	require.Equal(t, "verify_ok", string(resp.Body))
}

func TestAdapterParseEncryptedMessage(t *testing.T) {
	t.Parallel()

	adapter := &Adapter{
		cfg: Config{
			CorpID:         "ww_test",
			Token:          "token",
			EncodingAESKey: "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
			AgentMap: map[string]int64{
				"1000002": 1003,
			},
		},
		client: newClient(Config{}),
	}

	plainXML := `<xml><ToUserName><![CDATA[ww_test]]></ToUserName><FromUserName><![CDATA[user_1]]></FromUserName><CreateTime>1710000000</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[@机器人 hello wecom]]></Content><MsgId>msg_1</MsgId><AgentID>1000002</AgentID><ChatId><![CDATA[chat_1]]></ChatId></xml>`
	encrypted, err := encryptMessage(adapter.cfg.EncodingAESKey, adapter.cfg.CorpID, plainXML)
	require.NoError(t, err)

	msg, resp, parseErr := adapter.ParseCallback(context.Background(), &imEntity.CallbackRequest{
		Platform: imEntity.PlatformWeCom,
		Method:   "POST",
		Body:     []byte(`<xml><ToUserName><![CDATA[ww_test]]></ToUserName><AgentID><![CDATA[1000002]]></AgentID><Encrypt><![CDATA[` + encrypted + `]]></Encrypt></xml>`),
		Query: map[string]string{
			"timestamp":     "1710000000",
			"nonce":         "nonce_1",
			"msg_signature": signature(adapter.cfg.Token, "1710000000", "nonce_1", encrypted),
		},
	})

	require.NoError(t, parseErr)
	require.Nil(t, resp)
	require.NotNil(t, msg)
	require.Equal(t, int64(1003), msg.AgentID)
	require.Equal(t, "msg_1", msg.MessageID)
	require.Equal(t, "user_1", msg.UserID)
	require.Equal(t, "chat_1", msg.ChatID)
	require.Equal(t, "@机器人 hello wecom", msg.Text)
}
