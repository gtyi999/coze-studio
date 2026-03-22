include "../base.thrift"

namespace go im

struct ListIMPlatformsRequest {
    255: optional base.Base Base (api.none="true")
}

struct IMPlatformInfo {
    1: required string platform
    2: required string name
    3: required i64 connector_id (agw.js_conv="str", api.js_conv="true")
    4: required string callback_path
    5: required string callback_url
    6: required bool enabled
}

struct ListIMPlatformsResponse {
    1: required list<IMPlatformInfo> data

    253: i64 code
    254: string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct FeishuEventRequest {
    255: optional base.Base Base (api.none="true")
}

struct FeishuEventResponse {
    253: i64 code
    254: string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct DingTalkEventRequest {
    255: optional base.Base Base (api.none="true")
}

struct DingTalkEventResponse {
    253: i64 code
    254: string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

struct WeComEventRequest {
    255: optional base.Base Base (api.none="true")
}

struct WeComEventResponse {
    253: i64 code
    254: string msg
    255: optional base.BaseResp BaseResp (api.none="true")
}

service IMService {
    ListIMPlatformsResponse ListIMPlatforms(1: ListIMPlatformsRequest req)(api.get="/api/im/platforms", api.category="im", api.gen_path="im", agw.preserve_base="true")
    FeishuEventResponse FeishuEvent(1: FeishuEventRequest req)(api.post="/api/im/feishu/event", api.category="im", api.gen_path="im", agw.preserve_base="true")
    DingTalkEventResponse DingTalkEvent(1: DingTalkEventRequest req)(api.get="/api/im/dingtalk/event", api.post="/api/im/dingtalk/event", api.category="im", api.gen_path="im", agw.preserve_base="true")
    WeComEventResponse WeComEvent(1: WeComEventRequest req)(api.get="/api/im/wecom/event", api.post="/api/im/wecom/event", api.category="im", api.gen_path="im", agw.preserve_base="true")
}
