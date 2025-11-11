package httpclient

type Request struct {
	Method      string
	URL         string
	Body        any
	Headers     map[string]string
	QueryParams map[string]string
	Protected   bool
}
