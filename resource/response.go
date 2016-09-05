package resource

// Implementation of api2go.Responder
type Response struct {
	Res  interface{}
	Code int
}

func (r Response) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"version": "0",
	}
}

// Result returns the actual payload
func (r Response) Result() interface{} {
	return r.Res
}

// StatusCode returns the HTTP status code
func (r Response) StatusCode() int {
	return r.Code
}
