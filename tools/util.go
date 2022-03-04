package tools

import "encoding/json"

func FastMarshal(input interface{}) string {
	b, _ := json.Marshal(input)
	return string(b)
}
