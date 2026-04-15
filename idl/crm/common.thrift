namespace go crm

enum CrmFilterOperator {
    CrmFilterOperatorEq = 1
    CrmFilterOperatorIn = 2
    CrmFilterOperatorLike = 3
    CrmFilterOperatorGte = 4
    CrmFilterOperatorLte = 5
}

enum CrmSortOrder {
    CrmSortOrderAsc = 1
    CrmSortOrderDesc = 2
}

struct CrmFilter {
    1: optional string field_name
    2: optional CrmFilterOperator operator
    3: optional list<string> values
}

struct CrmSort {
    1: optional string field_name
    2: optional CrmSortOrder order
}

struct PageInfo {
    1: optional i32 page_no
    2: optional i32 page_size
    3: optional i64 total (agw.js_conv="str", api.js_conv="true")
}

struct DeleteResult {
    1: optional bool success
}
