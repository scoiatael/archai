package http

type Context interface {
	HandleErr(error)
}

type HttpContext interface {
	SendJson(interface{})
	GetSegment(string) string
	ServerErr(error)
}

type GetContext interface {
	HttpContext
	StringParam(string) string
	IntParam(string, int) int
}

type PostContext interface {
	HttpContext
	JsonBodyParams() (map[string]interface{}, error)
}
