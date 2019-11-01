package types

type APIError struct {
	Code    int16       `json:"code"`
	Key     string      `json:"key"`
	Message string      `json:"message"`
	Extras  interface{} `json:"extras,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

// APIErrorExtras is a shortcut for map[string]interface{}
type APIErrorExtras map[string]interface{}
