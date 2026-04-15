include "../base.thrift"
include "./common.thrift"

namespace go crm

struct CustomerInfo {
    1: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional string customer_name
    5: optional string customer_code
    6: optional string industry
    7: optional string level
    8: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string owner_user_name
    10: optional string status
    11: optional string mobile
    12: optional string email
    13: optional string address
    14: optional string remark
    15: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    16: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    17: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    18: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct ContactInfo {
    1: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional string contact_name
    6: optional string mobile
    7: optional string email
    8: optional string title
    9: optional bool is_primary
    10: optional string status
    11: optional string remark
    12: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    13: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    14: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    15: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct OpportunityInfo {
    1: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional string opportunity_name
    6: optional string stage
    7: optional string amount
    8: optional string expected_close_date
    9: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    10: optional string owner_user_name
    11: optional string status
    12: optional string remark
    13: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    14: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    15: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    16: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct FollowRecordInfo {
    1: optional i64 follow_record_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    6: optional string follow_type
    7: optional string content
    8: optional string next_follow_time
    9: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    10: optional string owner_user_name
    11: optional string status
    12: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    13: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    14: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    15: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct ProductInfo {
    1: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional string product_name
    5: optional string product_code
    6: optional string category
    7: optional string unit_price
    8: optional string status
    9: optional string remark
    10: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    11: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    12: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    13: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct SalesOrderInfo {
    1: optional i64 sales_order_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    6: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    7: optional string product_name
    8: optional i64 sales_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string sales_user_name
    10: optional string quantity
    11: optional string amount
    12: optional string order_date
    13: optional string status
    14: optional string remark
    15: optional i64 created_by (agw.js_conv="str", api.js_conv="true")
    16: optional i64 updated_by (agw.js_conv="str", api.js_conv="true")
    17: optional i64 created_at (agw.js_conv="str", api.js_conv="true")
    18: optional i64 updated_at (agw.js_conv="str", api.js_conv="true")
}

struct CreateCustomerRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional string customer_name
    4: optional string customer_code
    5: optional string industry
    6: optional string level
    7: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    8: optional string owner_user_name
    9: optional string status
    10: optional string mobile
    11: optional string email
    12: optional string address
    13: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct CreateCustomerResponse {
    1: optional CustomerInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateCustomerRequest {
    1: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional string customer_name
    5: optional string customer_code
    6: optional string industry
    7: optional string level
    8: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string owner_user_name
    10: optional string status
    11: optional string mobile
    12: optional string email
    13: optional string address
    14: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct UpdateCustomerResponse {
    1: optional CustomerInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteCustomerRequest {
    1: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteCustomerResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetCustomerDetailRequest {
    1: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetCustomerDetailResponse {
    1: optional CustomerInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListCustomersRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListCustomersData {
    1: optional list<CustomerInfo> list
    2: optional common.PageInfo page_info
}

struct ListCustomersResponse {
    1: optional ListCustomersData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateContactRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    4: optional string contact_name
    5: optional string mobile
    6: optional string email
    7: optional string title
    8: optional bool is_primary
    9: optional string status
    10: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct CreateContactResponse {
    1: optional ContactInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateContactRequest {
    1: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional string contact_name
    6: optional string mobile
    7: optional string email
    8: optional string title
    9: optional bool is_primary
    10: optional string status
    11: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct UpdateContactResponse {
    1: optional ContactInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteContactRequest {
    1: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteContactResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetContactDetailRequest {
    1: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetContactDetailResponse {
    1: optional ContactInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListContactsRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListContactsData {
    1: optional list<ContactInfo> list
    2: optional common.PageInfo page_info
}

struct ListContactsResponse {
    1: optional ListContactsData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateOpportunityRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    4: optional string opportunity_name
    5: optional string stage
    6: optional string amount
    7: optional string expected_close_date
    8: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string owner_user_name
    10: optional string status
    11: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct CreateOpportunityResponse {
    1: optional OpportunityInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateOpportunityRequest {
    1: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional string opportunity_name
    6: optional string stage
    7: optional string amount
    8: optional string expected_close_date
    9: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    10: optional string owner_user_name
    11: optional string status
    12: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct UpdateOpportunityResponse {
    1: optional OpportunityInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteOpportunityRequest {
    1: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteOpportunityResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetOpportunityDetailRequest {
    1: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetOpportunityDetailResponse {
    1: optional OpportunityInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListOpportunitiesRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListOpportunitiesData {
    1: optional list<OpportunityInfo> list
    2: optional common.PageInfo page_info
}

struct ListOpportunitiesResponse {
    1: optional ListOpportunitiesData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateFollowRecordRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    5: optional string follow_type
    6: optional string content
    7: optional string next_follow_time
    8: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string owner_user_name
    10: optional string status
    255: optional base.Base Base (api.none="true")
}

struct CreateFollowRecordResponse {
    1: optional FollowRecordInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateFollowRecordRequest {
    1: optional i64 follow_record_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional i64 contact_id (agw.js_conv="str", api.js_conv="true")
    6: optional string follow_type
    7: optional string content
    8: optional string next_follow_time
    9: optional i64 owner_user_id (agw.js_conv="str", api.js_conv="true")
    10: optional string owner_user_name
    11: optional string status
    255: optional base.Base Base (api.none="true")
}

struct UpdateFollowRecordResponse {
    1: optional FollowRecordInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteFollowRecordRequest {
    1: optional i64 follow_record_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteFollowRecordResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetFollowRecordDetailRequest {
    1: optional i64 follow_record_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetFollowRecordDetailResponse {
    1: optional FollowRecordInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListFollowRecordsRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListFollowRecordsData {
    1: optional list<FollowRecordInfo> list
    2: optional common.PageInfo page_info
}

struct ListFollowRecordsResponse {
    1: optional ListFollowRecordsData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateProductRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional string product_name
    4: optional string product_code
    5: optional string category
    6: optional string unit_price
    7: optional string status
    8: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct CreateProductResponse {
    1: optional ProductInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateProductRequest {
    1: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional string product_name
    5: optional string product_code
    6: optional string category
    7: optional string unit_price
    8: optional string status
    9: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct UpdateProductResponse {
    1: optional ProductInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteProductRequest {
    1: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteProductResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetProductDetailRequest {
    1: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetProductDetailResponse {
    1: optional ProductInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListProductsRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListProductsData {
    1: optional list<ProductInfo> list
    2: optional common.PageInfo page_info
}

struct ListProductsResponse {
    1: optional ListProductsData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct CreateSalesOrderRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    5: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    6: optional string product_name
    7: optional i64 sales_user_id (agw.js_conv="str", api.js_conv="true")
    8: optional string sales_user_name
    9: optional string quantity
    10: optional string amount
    11: optional string order_date
    12: optional string status
    13: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct CreateSalesOrderResponse {
    1: optional SalesOrderInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct UpdateSalesOrderRequest {
    1: optional i64 sales_order_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    4: optional i64 customer_id (agw.js_conv="str", api.js_conv="true")
    5: optional i64 opportunity_id (agw.js_conv="str", api.js_conv="true")
    6: optional i64 product_id (agw.js_conv="str", api.js_conv="true")
    7: optional string product_name
    8: optional i64 sales_user_id (agw.js_conv="str", api.js_conv="true")
    9: optional string sales_user_name
    10: optional string quantity
    11: optional string amount
    12: optional string order_date
    13: optional string status
    14: optional string remark
    255: optional base.Base Base (api.none="true")
}

struct UpdateSalesOrderResponse {
    1: optional SalesOrderInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DeleteSalesOrderRequest {
    1: optional i64 sales_order_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DeleteSalesOrderResponse {
    1: optional common.DeleteResult data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct GetSalesOrderDetailRequest {
    1: optional i64 sales_order_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    3: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct GetSalesOrderDetailResponse {
    1: optional SalesOrderInfo data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct ListSalesOrdersRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    3: optional i32 page_no
    4: optional i32 page_size
    5: optional string keyword
    6: optional list<common.CrmFilter> filters
    7: optional list<common.CrmSort> sorts
    255: optional base.Base Base (api.none="true")
}

struct ListSalesOrdersData {
    1: optional list<SalesOrderInfo> list
    2: optional common.PageInfo page_info
}

struct ListSalesOrdersResponse {
    1: optional ListSalesOrdersData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DashboardOrderTrendInfo {
    1: optional string date
    2: optional i64 order_count
    3: optional string order_amount
}

struct GetDashboardOverviewRequest {
    1: optional i64 tenant_id (agw.js_conv="str", api.js_conv="true")
    2: optional i64 space_id (agw.js_conv="str", api.js_conv="true")
    255: optional base.Base Base (api.none="true")
}

struct DashboardOverviewData {
    1: optional i64 customer_total
    2: optional i64 new_customers_this_month
    3: optional string opportunity_total_amount
    4: optional i64 new_opportunities_this_month
    5: optional string sales_order_total_amount
    6: optional list<DashboardOrderTrendInfo> recent_order_trend
}

struct GetDashboardOverviewResponse {
    1: optional DashboardOverviewData data
    253: optional i64 code
    254: optional string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

service CrmCustomerService {
    CreateCustomerResponse CreateCustomer(1: CreateCustomerRequest req)(api.post="/api/crm/customer/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateCustomerResponse UpdateCustomer(1: UpdateCustomerRequest req)(api.post="/api/crm/customer/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteCustomerResponse DeleteCustomer(1: DeleteCustomerRequest req)(api.post="/api/crm/customer/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetCustomerDetailResponse GetCustomerDetail(1: GetCustomerDetailRequest req)(api.get="/api/crm/customer/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListCustomersResponse ListCustomers(1: ListCustomersRequest req)(api.get="/api/crm/customer/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmContactService {
    CreateContactResponse CreateContact(1: CreateContactRequest req)(api.post="/api/crm/contact/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateContactResponse UpdateContact(1: UpdateContactRequest req)(api.post="/api/crm/contact/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteContactResponse DeleteContact(1: DeleteContactRequest req)(api.post="/api/crm/contact/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetContactDetailResponse GetContactDetail(1: GetContactDetailRequest req)(api.get="/api/crm/contact/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListContactsResponse ListContacts(1: ListContactsRequest req)(api.get="/api/crm/contact/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmOpportunityService {
    CreateOpportunityResponse CreateOpportunity(1: CreateOpportunityRequest req)(api.post="/api/crm/opportunity/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateOpportunityResponse UpdateOpportunity(1: UpdateOpportunityRequest req)(api.post="/api/crm/opportunity/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteOpportunityResponse DeleteOpportunity(1: DeleteOpportunityRequest req)(api.post="/api/crm/opportunity/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetOpportunityDetailResponse GetOpportunityDetail(1: GetOpportunityDetailRequest req)(api.get="/api/crm/opportunity/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListOpportunitiesResponse ListOpportunities(1: ListOpportunitiesRequest req)(api.get="/api/crm/opportunity/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmFollowRecordService {
    CreateFollowRecordResponse CreateFollowRecord(1: CreateFollowRecordRequest req)(api.post="/api/crm/follow_record/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateFollowRecordResponse UpdateFollowRecord(1: UpdateFollowRecordRequest req)(api.post="/api/crm/follow_record/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteFollowRecordResponse DeleteFollowRecord(1: DeleteFollowRecordRequest req)(api.post="/api/crm/follow_record/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetFollowRecordDetailResponse GetFollowRecordDetail(1: GetFollowRecordDetailRequest req)(api.get="/api/crm/follow_record/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListFollowRecordsResponse ListFollowRecords(1: ListFollowRecordsRequest req)(api.get="/api/crm/follow_record/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmProductService {
    CreateProductResponse CreateProduct(1: CreateProductRequest req)(api.post="/api/crm/product/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateProductResponse UpdateProduct(1: UpdateProductRequest req)(api.post="/api/crm/product/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteProductResponse DeleteProduct(1: DeleteProductRequest req)(api.post="/api/crm/product/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetProductDetailResponse GetProductDetail(1: GetProductDetailRequest req)(api.get="/api/crm/product/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListProductsResponse ListProducts(1: ListProductsRequest req)(api.get="/api/crm/product/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmSalesOrderService {
    CreateSalesOrderResponse CreateSalesOrder(1: CreateSalesOrderRequest req)(api.post="/api/crm/sales_order/create", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    UpdateSalesOrderResponse UpdateSalesOrder(1: UpdateSalesOrderRequest req)(api.post="/api/crm/sales_order/update", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    DeleteSalesOrderResponse DeleteSalesOrder(1: DeleteSalesOrderRequest req)(api.post="/api/crm/sales_order/delete", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    GetSalesOrderDetailResponse GetSalesOrderDetail(1: GetSalesOrderDetailRequest req)(api.get="/api/crm/sales_order/detail", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
    ListSalesOrdersResponse ListSalesOrders(1: ListSalesOrdersRequest req)(api.get="/api/crm/sales_order/list", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}

service CrmDashboardService {
    GetDashboardOverviewResponse GetDashboardOverview(1: GetDashboardOverviewRequest req)(api.get="/api/crm/dashboard/overview", api.category="crm", api.gen_path="crm", agw.preserve_base="true")
}
