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

type SensitivityLevel string
type MaskPolicy string

const (
	SensitivityLevelPublic    SensitivityLevel = "public"
	SensitivityLevelInternal  SensitivityLevel = "internal"
	SensitivityLevelSensitive SensitivityLevel = "sensitive"
)

const (
	MaskPolicyNone MaskPolicy = "none"
	MaskPolicyMask MaskPolicy = "mask"
	MaskPolicyDeny MaskPolicy = "deny"
	MaskPolicyRole MaskPolicy = "role_based"
)

type SemanticCatalogRequest struct {
	Scope

	Keyword         string
	TableKeys       []string
	IncludeInactive bool
}

type SemanticCatalog struct {
	Tables    []*SemanticTable
	Columns   []*SemanticColumn
	Metrics   []*SemanticMetric
	Relations []*SemanticRelation
}

type SemanticTable struct {
	ID int64
	Scope

	TableKey             string
	TableName            string
	TableDesc            string
	PhysicalTableName    string
	PrimaryTimeColumnKey string
	Status               string
	DefaultScopeJSON     string
	OwnerDomain          string
	VersionNo            int32
}

type SemanticColumn struct {
	ID int64
	Scope

	TableKey         string
	ColumnKey        string
	ColumnName       string
	ColumnDesc       string
	SourceColumnName string
	DataType         string
	ColumnRole       string
	DefaultAggFunc   string
	SensitivityLevel SensitivityLevel
	MaskPolicy       MaskPolicy
	AllowRoles       []string
	IsPrimaryKey     bool
	IsTimeKey        bool
	IsFilterable     bool
	IsGroupable      bool
	AllowFilterOnly  bool
}

type SemanticMetric struct {
	ID int64
	Scope

	MetricKey            string
	MetricName           string
	MetricDesc           string
	TableKey             string
	MetricType           string
	AggFunc              string
	MeasureColumnKey     string
	FormulaExpr          string
	DefaultTimeColumnKey string
	DefaultScopeJSON     string
	DefaultGroupBy       []string
	Unit                 string
	DisclaimerText       string
}

type SemanticRelation struct {
	ID int64
	Scope

	RelationKey    string
	LeftTableKey   string
	RightTableKey  string
	RelationType   string
	JoinType       string
	LeftColumnKey  string
	RightColumnKey string
}
