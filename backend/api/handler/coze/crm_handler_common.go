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
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	protocolconsts "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func parseRequiredInt64Param(c *app.RequestContext, value string, field string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		invalidParamRequestResponse(c, field+" is required")
		return 0, false
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		invalidParamRequestResponse(c, field+" is invalid")
		return 0, false
	}
	return parsed, true
}

func parseOptionalInt64Param(c *app.RequestContext, value string, field string) *int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		invalidParamRequestResponse(c, field+" is invalid")
		return nil
	}
	return &parsed
}

func parseOptionalInt64Value(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func parseOptionalTimestampParam(c *app.RequestContext, value string, field string) *int64 {
	return parseOptionalInt64Param(c, value, field)
}

func parseOptionalTimestampValue(value string) int64 {
	return parseOptionalInt64Value(value)
}

func parseDecimalParam(c *app.RequestContext, value string, field string, required bool) (float64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		if required {
			invalidParamRequestResponse(c, field+" is required")
			return 0, false
		}
		return 0, true
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		invalidParamRequestResponse(c, field+" is invalid")
		return 0, false
	}
	return parsed, true
}

func parseOptionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func parseOptionalDate(value string) *string {
	return parseOptionalString(value)
}

func writeCRMSuccess(c *app.RequestContext, data any) {
	c.JSON(protocolconsts.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func formatCRMInt64(value int64) string {
	if value <= 0 {
		return ""
	}
	return strconv.FormatInt(value, 10)
}

func formatCRMFloat(value float64) string {
	if value == 0 {
		return "0"
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}
