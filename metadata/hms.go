package metadata

const (
	// HMS_HEADER_TOKEN 在 Header 里的Token名称
	HMS_HEADER_TOKEN = "goms-token"

	// HMS_HEADER_TRACING_EXTEND 在 Header 里的Tracing扩展名称
	HMS_HEADER_TRACING_EXTEND = "goms-tracing-extend"

	// HMS_TOKEN_TRACING 在 Token 里的 Tracing 成员名称，一般由API网关设置
	HMS_TOKEN_TRACING = "goms-tracing"

	// HMS_HEADER_USER 在 Header 里的User扩展信息，一般由API网关设置
	HMS_HEADER_USER = "goms-current-user"

	// HMS_TOKEN_USER 在 Token 里的 User 成员名称
	HMS_TOKEN_USER = "goms-user"

	// HMS_TOKEN_BE_USER 在 Token 里的 后端User 成员名称
	HMS_TOKEN_BE_USER = "goms-be-user"

	// HMS_QUERY_ACTION 请求参数Action名称
	HMS_QUERY_ACTION = "goms_action"

	// HMS_QUERY_TAG 请求服务标签
	HMS_QUERY_TAG = "goms_tag"

	// HMS_QUERY_USER 请求参数的用户信息名称
	HMS_QUERY_USER = "goms_user"

	// HMS_QUERY_BE_USER 请求参数的后端用户信息名称
	HMS_QUERY_BE_USER = "goms_be_user"

	//HMS_HEADER_BE_USER 在 Header 里的后端User扩展信息，一般由API网关设置
	HMS_HEADER_BE_USER = "goms-current-be-user"
)
