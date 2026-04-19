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

package service

import (
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func validateQueryRequest(req *entity.QueryRequest) error {
	if req == nil {
		return errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query request is required"))
	}
	if err := validateScope(&req.Scope); err != nil {
		return err
	}

	req.Question = strings.TrimSpace(req.Question)
	req.RequestID = strings.TrimSpace(req.RequestID)
	req.ConversationID = strings.TrimSpace(req.ConversationID)
	if req.Question == "" {
		return errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "question is required"))
	}

	return nil
}

func validateSemanticCatalogRequest(req *entity.SemanticCatalogRequest) error {
	if req == nil {
		return errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "semantic catalog request is required"))
	}
	if err := validateScope(&req.Scope); err != nil {
		return err
	}

	req.Keyword = strings.TrimSpace(req.Keyword)
	return nil
}

func validateQueryLogFilter(filter *entity.QueryLogFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query log filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}

	normalizePage(&filter.PageOption)
	filter.Question = strings.TrimSpace(filter.Question)
	return nil
}

func validateForecastRequest(req *entity.ForecastRequest) error {
	if req == nil {
		return errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "forecast request is required"))
	}
	if err := validateScope(&req.Scope); err != nil {
		return err
	}

	req.MetricKey = strings.TrimSpace(req.MetricKey)
	if req.MetricKey == "" {
		req.MetricKey = "product_sales_qty"
	}
	if req.Months <= 0 {
		req.Months = 6
	}
	if req.Limit <= 0 {
		req.Limit = 5
	}

	return nil
}
