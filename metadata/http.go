package metadata

const (
	HttpTraceId          = "trace_id" // HttpTraceId 链路ID
	HttpShreqid          = "Shreqid"
	HttpFrom             = "from" // 来自http中的_userAgent
	HTTP_CIRCUIT_BREAKER = "cb"   // HTTP_CIRCUIT_BREAKER 降级状态

	// 中间件拦截获取的一些参数
	HTTP_MIDDLEWARE_AB_TEST      = "ABTest"      // ABTest abtest的map
	HTTP_MIDDLEWARE_APP_PLATFORM = "appPlatform" //
	HTTP_MIDDLEWARE_CLIENT_CODE  = "clientCode"
	HTTP_MIDDLEWARE_VERSION      = "v"            // 版本
	HTTP_MIDDLEWARE_GOMS_USER    = "goms_user"    // 前台用户信息
	HTTP_MIDDLEWARE_GOMS_BE_USER = "goms_be_user" // 前台用户信息
	HTTP_MIDDLEWARE_NETWORK      = "network"      // 网络
	HTTP_MIDDLEWARE_PLATFORM     = "platform"     // 手机平台
	HTTP_MIDDLEWARE_MOBILE_BRAND = "mobile_brand" // 手机品牌
	HTTP_MIDDLEWARE_MIN_VERSION  = "min_version"  // minVersion
	HTTP_MIDDLEWARE_CHANNEL      = "channel"      // 渠道

	HTTP_COLOR = "color" // HTTP_COLOR 流量染色的颜色
)
