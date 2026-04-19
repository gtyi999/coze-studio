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

package entity

type ForecastRequest struct {
	Scope

	UserID     int64
	MetricKey  string
	Months     int
	Limit      int
	ProductIDs []int64
}

type ForecastFeature struct {
	ProductID     int64
	ProductName   string
	MetricKey     string
	Period        string
	MetricValue   float64
	GrowthRate    float64
	TrendSlope    float64
	WeightedAvg3M float64
	Volatility    float64
	Score         float64
}

type ForecastResult struct {
	MetricKey      string
	TopProductID   int64
	TopProductName string
	Features       []*ForecastFeature
	Reasons        []string
	Disclaimer     string
	GeneratedAt    int64
}
