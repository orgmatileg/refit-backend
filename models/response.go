package models

// ResponseDefault models
type ResponseDefault struct {
	Code int           `json:"code"`
	List interface{}   `json:"_list,omitempty"`
	Data []interface{} `json:"data"`
}

// IsList ..
func (r *ResponseDefault) IsList(backURL, nextURL, count, limit, offset, order string) {
	r.List = struct {
		BackURL string `json:"back"`
		NextURL string `json:"next"`
		Count   string `json:"count"`
		Limit   string `json:"limit"`
		Offset  string `json:"offset"`
		Order   string `json:"order"`
	}{backURL, nextURL, count, limit, offset, order}
}

// RequestQueryParam models
type RequestQueryParam map[string]string
