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

package coze

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	cozehandler "github.com/coze-dev/coze-studio/backend/api/handler/coze"
)

func RegisterManualIMRoutes(r *server.Hertz) {
	root := r.Group("/", rootMw()...)
	api := root.Group("/api", _apiMw()...)
	im := api.Group("/im", _imMw()...)

	channelConfig := im.Group("/channel_config")
	channelConfig.GET("/list", cozehandler.ListIMChannelConfigs)
	channelConfig.GET("/get", cozehandler.GetIMChannelConfig)
	channelConfig.POST("/create", cozehandler.CreateIMChannelConfig)
	channelConfig.POST("/update", cozehandler.UpdateIMChannelConfig)
	channelConfig.POST("/test", cozehandler.TestIMChannelConnectivity)

	task := im.Group("/task")
	task.GET("/list", cozehandler.ListIMTaskRecords)
	task.GET("/get", cozehandler.GetIMTaskDetail)
	task.POST("/retry", cozehandler.RetryIMTask)
}
