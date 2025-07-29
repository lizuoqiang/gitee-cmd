package utils

import "github.com/mark3labs/mcp-go/mcp"

func CheckRequired(params map[string]interface{}, required ...string) (*mcp.CallToolResult, error) {
	for _, key := range required {
		if _, ok := params[key]; !ok {
			return mcp.NewToolResultError("Missing required parameter: " + key),
				NewParamError(key, "parameter is required")
		}
	}
	return nil, nil
}

func CombineOptions(options ...[]mcp.ToolOption) []mcp.ToolOption {
	var result []mcp.ToolOption
	for _, option := range options {
		result = append(result, option...)
	}
	return result
}
