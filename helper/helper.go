package helper

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success string `json:"success"`
}

func APIResponse(message string, code int, success string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Success: success,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}
