package response

const (
	RateLimitError = 401
	Invalidate     = 10001
)

var mapMsg = map[int]string{
	Invalidate:     "不合法的请求参数",
	RateLimitError: "接口请求太频繁，请稍后重试",
}
