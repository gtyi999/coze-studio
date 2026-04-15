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

import { createAPI } from './../../api/config';

const schemaRoot = 'api://schemas/idl_crm_crm';
const service = 'crm';

export type CRMStatus = 'active' | 'inactive' | 'open' | 'draft' | string;

export interface CustomerData {
  customer_id?: string;
  tenant_id?: string;
  space_id?: string;
  customer_name?: string;
  customer_code?: string;
  industry?: string;
  level?: string;
  customer_level?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
  mobile?: string;
  email?: string;
  address?: string;
  remark?: string;
  description?: string;
  created_by?: string;
  updated_by?: string;
  creator_id?: string;
  updater_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ContactData {
  contact_id?: string;
  tenant_id?: string;
  space_id?: string;
  customer_id?: string;
  customer_name?: string;
  contact_name?: string;
  mobile?: string;
  email?: string;
  title?: string;
  position?: string;
  is_primary?: boolean;
  status?: CRMStatus;
  remark?: string;
  description?: string;
  created_by?: string;
  updated_by?: string;
  creator_id?: string;
  updater_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface OpportunityData {
  opportunity_id?: string;
  tenant_id?: string;
  space_id?: string;
  customer_id?: string;
  opportunity_name?: string;
  stage?: string;
  amount?: string;
  expected_close_date?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
  remark?: string;
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface FollowRecordData {
  follow_record_id?: string;
  tenant_id?: string;
  space_id?: string;
  customer_id?: string;
  contact_id?: string;
  follow_type?: string;
  content?: string;
  next_follow_time?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface ProductData {
  product_id?: string;
  tenant_id?: string;
  space_id?: string;
  product_name?: string;
  product_code?: string;
  category?: string;
  unit_price?: string;
  status?: CRMStatus;
  remark?: string;
  description?: string;
  created_by?: string;
  updated_by?: string;
  creator_id?: string;
  updater_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface SalesOrderData {
  sales_order_id?: string;
  tenant_id?: string;
  space_id?: string;
  customer_id?: string;
  opportunity_id?: string;
  product_id?: string;
  product_name?: string;
  sales_user_id?: string;
  sales_user_name?: string;
  quantity?: string;
  amount?: string;
  order_date?: string;
  status?: CRMStatus;
  remark?: string;
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CustomerListData {
  list?: CustomerData[];
  total?: number;
}

export interface ContactListData {
  list?: ContactData[];
  total?: number;
}

export interface OpportunityListData {
  list?: OpportunityData[];
  total?: number;
}

export interface FollowRecordListData {
  list?: FollowRecordData[];
  total?: number;
}

export interface ProductListData {
  list?: ProductData[];
  total?: number;
}

export interface SalesOrderListData {
  list?: SalesOrderData[];
  total?: number;
}

export interface DashboardOrderTrendData {
  date?: string;
  order_count?: number;
  order_amount?: string;
}

export interface DashboardOverviewData {
  customer_total?: number;
  new_customers_this_month?: number;
  opportunity_total_amount?: string;
  new_opportunities_this_month?: number;
  sales_order_total_amount?: string;
  recent_order_trend?: DashboardOrderTrendData[];
}

export interface ListCustomersRequest {
  space_id: string;
  keyword?: string;
  status?: CRMStatus;
  owner_user_id?: string;
  created_at_start?: string;
  created_at_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListCustomersResponse {
  data?: CustomerListData;
  code: number;
  msg: string;
}

export interface GetCustomerRequest {
  space_id: string;
  customer_id: string;
}

export interface GetCustomerResponse {
  data?: CustomerData;
  code: number;
  msg: string;
}

export interface CreateCustomerRequest {
  space_id: string;
  customer_name: string;
  customer_code?: string;
  industry?: string;
  level?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
  mobile?: string;
  email?: string;
  address?: string;
  remark?: string;
}

export interface CreateCustomerResponse {
  data?: CustomerData;
  code: number;
  msg: string;
}

export interface UpdateCustomerRequest extends CreateCustomerRequest {
  customer_id: string;
}

export interface UpdateCustomerResponse {
  data?: CustomerData;
  code: number;
  msg: string;
}

export interface DeleteCustomerRequest {
  space_id: string;
  customer_id: string;
}

export interface DeleteCustomerResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface ListContactsRequest {
  space_id: string;
  customer_id?: string;
  keyword?: string;
  status?: CRMStatus;
  created_at_start?: string;
  created_at_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListContactsResponse {
  data?: ContactListData;
  code: number;
  msg: string;
}

export interface GetContactRequest {
  space_id: string;
  contact_id: string;
}

export interface GetContactResponse {
  data?: ContactData;
  code: number;
  msg: string;
}

export interface CreateContactRequest {
  space_id: string;
  customer_id: string;
  contact_name: string;
  mobile?: string;
  email?: string;
  title?: string;
  is_primary?: boolean;
  status?: CRMStatus;
  remark?: string;
}

export interface CreateContactResponse {
  data?: ContactData;
  code: number;
  msg: string;
}

export interface UpdateContactRequest extends CreateContactRequest {
  contact_id: string;
}

export interface UpdateContactResponse {
  data?: ContactData;
  code: number;
  msg: string;
}

export interface DeleteContactRequest {
  space_id: string;
  contact_id: string;
}

export interface DeleteContactResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface ListOpportunitiesRequest {
  space_id: string;
  customer_id?: string;
  owner_user_id?: string;
  keyword?: string;
  status?: CRMStatus;
  created_at_start?: string;
  created_at_end?: string;
  expected_close_date_start?: string;
  expected_close_date_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListOpportunitiesResponse {
  data?: OpportunityListData;
  code: number;
  msg: string;
}

export interface GetOpportunityRequest {
  space_id: string;
  opportunity_id: string;
}

export interface GetOpportunityResponse {
  data?: OpportunityData;
  code: number;
  msg: string;
}

export interface CreateOpportunityRequest {
  space_id: string;
  customer_id: string;
  opportunity_name: string;
  stage?: string;
  amount?: string;
  expected_close_date?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
  remark?: string;
}

export interface CreateOpportunityResponse {
  data?: OpportunityData;
  code: number;
  msg: string;
}

export interface UpdateOpportunityRequest extends CreateOpportunityRequest {
  opportunity_id: string;
}

export interface UpdateOpportunityResponse {
  data?: OpportunityData;
  code: number;
  msg: string;
}

export interface DeleteOpportunityRequest {
  space_id: string;
  opportunity_id: string;
}

export interface DeleteOpportunityResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface ListFollowRecordsRequest {
  space_id: string;
  customer_id?: string;
  contact_id?: string;
  owner_user_id?: string;
  keyword?: string;
  status?: CRMStatus;
  created_at_start?: string;
  created_at_end?: string;
  next_follow_time_start?: string;
  next_follow_time_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListFollowRecordsResponse {
  data?: FollowRecordListData;
  code: number;
  msg: string;
}

export interface GetFollowRecordRequest {
  space_id: string;
  follow_record_id: string;
}

export interface GetFollowRecordResponse {
  data?: FollowRecordData;
  code: number;
  msg: string;
}

export interface CreateFollowRecordRequest {
  space_id: string;
  customer_id: string;
  contact_id?: string;
  follow_type: string;
  content: string;
  next_follow_time?: string;
  owner_user_id?: string;
  owner_user_name?: string;
  status?: CRMStatus;
}

export interface CreateFollowRecordResponse {
  data?: FollowRecordData;
  code: number;
  msg: string;
}

export interface UpdateFollowRecordRequest extends CreateFollowRecordRequest {
  follow_record_id: string;
}

export interface UpdateFollowRecordResponse {
  data?: FollowRecordData;
  code: number;
  msg: string;
}

export interface DeleteFollowRecordRequest {
  space_id: string;
  follow_record_id: string;
}

export interface DeleteFollowRecordResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface ListProductsRequest {
  space_id: string;
  keyword?: string;
  status?: CRMStatus;
  created_at_start?: string;
  created_at_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListProductsResponse {
  data?: ProductListData;
  code: number;
  msg: string;
}

export interface GetProductRequest {
  space_id: string;
  product_id: string;
}

export interface GetProductResponse {
  data?: ProductData;
  code: number;
  msg: string;
}

export interface CreateProductRequest {
  space_id: string;
  product_name: string;
  product_code?: string;
  category?: string;
  unit_price?: string;
  status?: CRMStatus;
  remark?: string;
}

export interface CreateProductResponse {
  data?: ProductData;
  code: number;
  msg: string;
}

export interface UpdateProductRequest extends CreateProductRequest {
  product_id: string;
}

export interface UpdateProductResponse {
  data?: ProductData;
  code: number;
  msg: string;
}

export interface DeleteProductRequest {
  space_id: string;
  product_id: string;
}

export interface DeleteProductResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface ListSalesOrdersRequest {
  space_id: string;
  customer_id?: string;
  opportunity_id?: string;
  product_id?: string;
  sales_user_id?: string;
  keyword?: string;
  status?: CRMStatus;
  created_at_start?: string;
  created_at_end?: string;
  order_date_start?: string;
  order_date_end?: string;
  page?: number;
  page_size?: number;
}

export interface ListSalesOrdersResponse {
  data?: SalesOrderListData;
  code: number;
  msg: string;
}

export interface GetSalesOrderRequest {
  space_id: string;
  sales_order_id: string;
}

export interface GetSalesOrderResponse {
  data?: SalesOrderData;
  code: number;
  msg: string;
}

export interface CreateSalesOrderRequest {
  space_id: string;
  customer_id: string;
  opportunity_id?: string;
  product_id: string;
  product_name?: string;
  sales_user_id?: string;
  sales_user_name?: string;
  quantity?: string;
  amount?: string;
  order_date?: string;
  status?: CRMStatus;
  remark?: string;
}

export interface CreateSalesOrderResponse {
  data?: SalesOrderData;
  code: number;
  msg: string;
}

export interface UpdateSalesOrderRequest extends CreateSalesOrderRequest {
  sales_order_id: string;
}

export interface UpdateSalesOrderResponse {
  data?: SalesOrderData;
  code: number;
  msg: string;
}

export interface DeleteSalesOrderRequest {
  space_id: string;
  sales_order_id: string;
}

export interface DeleteSalesOrderResponse {
  data?: Record<string, never>;
  code: number;
  msg: string;
}

export interface GetDashboardOverviewRequest {
  space_id: string;
}

export interface GetDashboardOverviewResponse {
  data?: DashboardOverviewData;
  code: number;
  msg: string;
}

export const GetDashboardOverview = /*#__PURE__*/ createAPI<
  GetDashboardOverviewRequest,
  GetDashboardOverviewResponse
>({
  url: '/api/crm/dashboard/overview',
  method: 'GET',
  name: 'GetDashboardOverview',
  reqType: 'GetDashboardOverviewRequest',
  reqMapping: { query: ['space_id'] },
  resType: 'GetDashboardOverviewResponse',
  schemaRoot,
  service,
});

export const ListCustomers = /*#__PURE__*/ createAPI<
  ListCustomersRequest,
  ListCustomersResponse
>({
  url: '/api/crm/customer/list',
  method: 'GET',
  name: 'ListCustomers',
  reqType: 'ListCustomersRequest',
  reqMapping: {
    query: [
      'space_id',
      'keyword',
      'status',
      'owner_user_id',
      'created_at_start',
      'created_at_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListCustomersResponse',
  schemaRoot,
  service,
});

export const GetCustomer = /*#__PURE__*/ createAPI<
  GetCustomerRequest,
  GetCustomerResponse
>({
  url: '/api/crm/customer/get',
  method: 'GET',
  name: 'GetCustomer',
  reqType: 'GetCustomerRequest',
  reqMapping: { query: ['space_id', 'customer_id'] },
  resType: 'GetCustomerResponse',
  schemaRoot,
  service,
});

export const CreateCustomer = /*#__PURE__*/ createAPI<
  CreateCustomerRequest,
  CreateCustomerResponse
>({
  url: '/api/crm/customer/create',
  method: 'POST',
  name: 'CreateCustomer',
  reqType: 'CreateCustomerRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_name',
      'customer_code',
      'industry',
      'level',
      'owner_user_id',
      'owner_user_name',
      'status',
      'mobile',
      'email',
      'address',
      'remark',
    ],
  },
  resType: 'CreateCustomerResponse',
  schemaRoot,
  service,
});

export const UpdateCustomer = /*#__PURE__*/ createAPI<
  UpdateCustomerRequest,
  UpdateCustomerResponse
>({
  url: '/api/crm/customer/update',
  method: 'POST',
  name: 'UpdateCustomer',
  reqType: 'UpdateCustomerRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_id',
      'customer_name',
      'customer_code',
      'industry',
      'level',
      'owner_user_id',
      'owner_user_name',
      'status',
      'mobile',
      'email',
      'address',
      'remark',
    ],
  },
  resType: 'UpdateCustomerResponse',
  schemaRoot,
  service,
});

export const DeleteCustomer = /*#__PURE__*/ createAPI<
  DeleteCustomerRequest,
  DeleteCustomerResponse
>({
  url: '/api/crm/customer/delete',
  method: 'POST',
  name: 'DeleteCustomer',
  reqType: 'DeleteCustomerRequest',
  reqMapping: { body: ['space_id', 'customer_id'] },
  resType: 'DeleteCustomerResponse',
  schemaRoot,
  service,
});

export const ListContacts = /*#__PURE__*/ createAPI<
  ListContactsRequest,
  ListContactsResponse
>({
  url: '/api/crm/contact/list',
  method: 'GET',
  name: 'ListContacts',
  reqType: 'ListContactsRequest',
  reqMapping: {
    query: [
      'space_id',
      'customer_id',
      'keyword',
      'status',
      'created_at_start',
      'created_at_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListContactsResponse',
  schemaRoot,
  service,
});

export const GetContact = /*#__PURE__*/ createAPI<
  GetContactRequest,
  GetContactResponse
>({
  url: '/api/crm/contact/get',
  method: 'GET',
  name: 'GetContact',
  reqType: 'GetContactRequest',
  reqMapping: { query: ['space_id', 'contact_id'] },
  resType: 'GetContactResponse',
  schemaRoot,
  service,
});

export const CreateContact = /*#__PURE__*/ createAPI<
  CreateContactRequest,
  CreateContactResponse
>({
  url: '/api/crm/contact/create',
  method: 'POST',
  name: 'CreateContact',
  reqType: 'CreateContactRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_id',
      'contact_name',
      'mobile',
      'email',
      'title',
      'is_primary',
      'status',
      'remark',
    ],
  },
  resType: 'CreateContactResponse',
  schemaRoot,
  service,
});

export const UpdateContact = /*#__PURE__*/ createAPI<
  UpdateContactRequest,
  UpdateContactResponse
>({
  url: '/api/crm/contact/update',
  method: 'POST',
  name: 'UpdateContact',
  reqType: 'UpdateContactRequest',
  reqMapping: {
    body: [
      'space_id',
      'contact_id',
      'customer_id',
      'contact_name',
      'mobile',
      'email',
      'title',
      'is_primary',
      'status',
      'remark',
    ],
  },
  resType: 'UpdateContactResponse',
  schemaRoot,
  service,
});

export const DeleteContact = /*#__PURE__*/ createAPI<
  DeleteContactRequest,
  DeleteContactResponse
>({
  url: '/api/crm/contact/delete',
  method: 'POST',
  name: 'DeleteContact',
  reqType: 'DeleteContactRequest',
  reqMapping: { body: ['space_id', 'contact_id'] },
  resType: 'DeleteContactResponse',
  schemaRoot,
  service,
});

export const ListOpportunities = /*#__PURE__*/ createAPI<
  ListOpportunitiesRequest,
  ListOpportunitiesResponse
>({
  url: '/api/crm/opportunity/list',
  method: 'GET',
  name: 'ListOpportunities',
  reqType: 'ListOpportunitiesRequest',
  reqMapping: {
    query: [
      'space_id',
      'customer_id',
      'owner_user_id',
      'keyword',
      'status',
      'created_at_start',
      'created_at_end',
      'expected_close_date_start',
      'expected_close_date_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListOpportunitiesResponse',
  schemaRoot,
  service,
});

export const GetOpportunity = /*#__PURE__*/ createAPI<
  GetOpportunityRequest,
  GetOpportunityResponse
>({
  url: '/api/crm/opportunity/get',
  method: 'GET',
  name: 'GetOpportunity',
  reqType: 'GetOpportunityRequest',
  reqMapping: { query: ['space_id', 'opportunity_id'] },
  resType: 'GetOpportunityResponse',
  schemaRoot,
  service,
});

export const CreateOpportunity = /*#__PURE__*/ createAPI<
  CreateOpportunityRequest,
  CreateOpportunityResponse
>({
  url: '/api/crm/opportunity/create',
  method: 'POST',
  name: 'CreateOpportunity',
  reqType: 'CreateOpportunityRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_id',
      'opportunity_name',
      'stage',
      'amount',
      'expected_close_date',
      'owner_user_id',
      'owner_user_name',
      'status',
      'remark',
    ],
  },
  resType: 'CreateOpportunityResponse',
  schemaRoot,
  service,
});

export const UpdateOpportunity = /*#__PURE__*/ createAPI<
  UpdateOpportunityRequest,
  UpdateOpportunityResponse
>({
  url: '/api/crm/opportunity/update',
  method: 'POST',
  name: 'UpdateOpportunity',
  reqType: 'UpdateOpportunityRequest',
  reqMapping: {
    body: [
      'space_id',
      'opportunity_id',
      'customer_id',
      'opportunity_name',
      'stage',
      'amount',
      'expected_close_date',
      'owner_user_id',
      'owner_user_name',
      'status',
      'remark',
    ],
  },
  resType: 'UpdateOpportunityResponse',
  schemaRoot,
  service,
});

export const DeleteOpportunity = /*#__PURE__*/ createAPI<
  DeleteOpportunityRequest,
  DeleteOpportunityResponse
>({
  url: '/api/crm/opportunity/delete',
  method: 'POST',
  name: 'DeleteOpportunity',
  reqType: 'DeleteOpportunityRequest',
  reqMapping: { body: ['space_id', 'opportunity_id'] },
  resType: 'DeleteOpportunityResponse',
  schemaRoot,
  service,
});

export const ListFollowRecords = /*#__PURE__*/ createAPI<
  ListFollowRecordsRequest,
  ListFollowRecordsResponse
>({
  url: '/api/crm/follow_record/list',
  method: 'GET',
  name: 'ListFollowRecords',
  reqType: 'ListFollowRecordsRequest',
  reqMapping: {
    query: [
      'space_id',
      'customer_id',
      'contact_id',
      'owner_user_id',
      'keyword',
      'status',
      'created_at_start',
      'created_at_end',
      'next_follow_time_start',
      'next_follow_time_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListFollowRecordsResponse',
  schemaRoot,
  service,
});

export const GetFollowRecord = /*#__PURE__*/ createAPI<
  GetFollowRecordRequest,
  GetFollowRecordResponse
>({
  url: '/api/crm/follow_record/get',
  method: 'GET',
  name: 'GetFollowRecord',
  reqType: 'GetFollowRecordRequest',
  reqMapping: { query: ['space_id', 'follow_record_id'] },
  resType: 'GetFollowRecordResponse',
  schemaRoot,
  service,
});

export const CreateFollowRecord = /*#__PURE__*/ createAPI<
  CreateFollowRecordRequest,
  CreateFollowRecordResponse
>({
  url: '/api/crm/follow_record/create',
  method: 'POST',
  name: 'CreateFollowRecord',
  reqType: 'CreateFollowRecordRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_id',
      'contact_id',
      'follow_type',
      'content',
      'next_follow_time',
      'owner_user_id',
      'owner_user_name',
      'status',
    ],
  },
  resType: 'CreateFollowRecordResponse',
  schemaRoot,
  service,
});

export const UpdateFollowRecord = /*#__PURE__*/ createAPI<
  UpdateFollowRecordRequest,
  UpdateFollowRecordResponse
>({
  url: '/api/crm/follow_record/update',
  method: 'POST',
  name: 'UpdateFollowRecord',
  reqType: 'UpdateFollowRecordRequest',
  reqMapping: {
    body: [
      'space_id',
      'follow_record_id',
      'customer_id',
      'contact_id',
      'follow_type',
      'content',
      'next_follow_time',
      'owner_user_id',
      'owner_user_name',
      'status',
    ],
  },
  resType: 'UpdateFollowRecordResponse',
  schemaRoot,
  service,
});

export const DeleteFollowRecord = /*#__PURE__*/ createAPI<
  DeleteFollowRecordRequest,
  DeleteFollowRecordResponse
>({
  url: '/api/crm/follow_record/delete',
  method: 'POST',
  name: 'DeleteFollowRecord',
  reqType: 'DeleteFollowRecordRequest',
  reqMapping: { body: ['space_id', 'follow_record_id'] },
  resType: 'DeleteFollowRecordResponse',
  schemaRoot,
  service,
});

export const ListProducts = /*#__PURE__*/ createAPI<
  ListProductsRequest,
  ListProductsResponse
>({
  url: '/api/crm/product/list',
  method: 'GET',
  name: 'ListProducts',
  reqType: 'ListProductsRequest',
  reqMapping: {
    query: [
      'space_id',
      'keyword',
      'status',
      'created_at_start',
      'created_at_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListProductsResponse',
  schemaRoot,
  service,
});

export const GetProduct = /*#__PURE__*/ createAPI<
  GetProductRequest,
  GetProductResponse
>({
  url: '/api/crm/product/get',
  method: 'GET',
  name: 'GetProduct',
  reqType: 'GetProductRequest',
  reqMapping: { query: ['space_id', 'product_id'] },
  resType: 'GetProductResponse',
  schemaRoot,
  service,
});

export const CreateProduct = /*#__PURE__*/ createAPI<
  CreateProductRequest,
  CreateProductResponse
>({
  url: '/api/crm/product/create',
  method: 'POST',
  name: 'CreateProduct',
  reqType: 'CreateProductRequest',
  reqMapping: {
    body: [
      'space_id',
      'product_name',
      'product_code',
      'category',
      'unit_price',
      'status',
      'remark',
    ],
  },
  resType: 'CreateProductResponse',
  schemaRoot,
  service,
});

export const UpdateProduct = /*#__PURE__*/ createAPI<
  UpdateProductRequest,
  UpdateProductResponse
>({
  url: '/api/crm/product/update',
  method: 'POST',
  name: 'UpdateProduct',
  reqType: 'UpdateProductRequest',
  reqMapping: {
    body: [
      'space_id',
      'product_id',
      'product_name',
      'product_code',
      'category',
      'unit_price',
      'status',
      'remark',
    ],
  },
  resType: 'UpdateProductResponse',
  schemaRoot,
  service,
});

export const DeleteProduct = /*#__PURE__*/ createAPI<
  DeleteProductRequest,
  DeleteProductResponse
>({
  url: '/api/crm/product/delete',
  method: 'POST',
  name: 'DeleteProduct',
  reqType: 'DeleteProductRequest',
  reqMapping: { body: ['space_id', 'product_id'] },
  resType: 'DeleteProductResponse',
  schemaRoot,
  service,
});

export const ListSalesOrders = /*#__PURE__*/ createAPI<
  ListSalesOrdersRequest,
  ListSalesOrdersResponse
>({
  url: '/api/crm/sales_order/list',
  method: 'GET',
  name: 'ListSalesOrders',
  reqType: 'ListSalesOrdersRequest',
  reqMapping: {
    query: [
      'space_id',
      'customer_id',
      'opportunity_id',
      'product_id',
      'sales_user_id',
      'keyword',
      'status',
      'created_at_start',
      'created_at_end',
      'order_date_start',
      'order_date_end',
      'page',
      'page_size',
    ],
  },
  resType: 'ListSalesOrdersResponse',
  schemaRoot,
  service,
});

export const GetSalesOrder = /*#__PURE__*/ createAPI<
  GetSalesOrderRequest,
  GetSalesOrderResponse
>({
  url: '/api/crm/sales_order/get',
  method: 'GET',
  name: 'GetSalesOrder',
  reqType: 'GetSalesOrderRequest',
  reqMapping: { query: ['space_id', 'sales_order_id'] },
  resType: 'GetSalesOrderResponse',
  schemaRoot,
  service,
});

export const CreateSalesOrder = /*#__PURE__*/ createAPI<
  CreateSalesOrderRequest,
  CreateSalesOrderResponse
>({
  url: '/api/crm/sales_order/create',
  method: 'POST',
  name: 'CreateSalesOrder',
  reqType: 'CreateSalesOrderRequest',
  reqMapping: {
    body: [
      'space_id',
      'customer_id',
      'opportunity_id',
      'product_id',
      'product_name',
      'sales_user_id',
      'sales_user_name',
      'quantity',
      'amount',
      'order_date',
      'status',
      'remark',
    ],
  },
  resType: 'CreateSalesOrderResponse',
  schemaRoot,
  service,
});

export const UpdateSalesOrder = /*#__PURE__*/ createAPI<
  UpdateSalesOrderRequest,
  UpdateSalesOrderResponse
>({
  url: '/api/crm/sales_order/update',
  method: 'POST',
  name: 'UpdateSalesOrder',
  reqType: 'UpdateSalesOrderRequest',
  reqMapping: {
    body: [
      'space_id',
      'sales_order_id',
      'customer_id',
      'opportunity_id',
      'product_id',
      'product_name',
      'sales_user_id',
      'sales_user_name',
      'quantity',
      'amount',
      'order_date',
      'status',
      'remark',
    ],
  },
  resType: 'UpdateSalesOrderResponse',
  schemaRoot,
  service,
});

export const DeleteSalesOrder = /*#__PURE__*/ createAPI<
  DeleteSalesOrderRequest,
  DeleteSalesOrderResponse
>({
  url: '/api/crm/sales_order/delete',
  method: 'POST',
  name: 'DeleteSalesOrder',
  reqType: 'DeleteSalesOrderRequest',
  reqMapping: { body: ['space_id', 'sales_order_id'] },
  resType: 'DeleteSalesOrderResponse',
  schemaRoot,
  service,
});
