package helper

func FormatResponse(message string, data any) map[string]any {
	var response = map[string]any{}
	response["message"] = message
	if data != nil {
		response["data"] = data
	}
	return response
}

func FormatPaginationResponse(message string, data any, pagination any) map[string]any {
	var responsePagination = map[string]any{}

	responsePagination["message"] = message
	responsePagination["data"] = data
	responsePagination["pagination"] = pagination

	return responsePagination
}
