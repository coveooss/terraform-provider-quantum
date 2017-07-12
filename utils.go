package main

import (
	"fmt"
)

func interfaceToString(source []interface{}) (result []string) {
	result = make([]string, len(source))
	for i := range source {
		result[i] = fmt.Sprintf("%v", source[i])
	}
	return
}

// func dump(v interface{}) {
// 	js, _ := json.Marshal(v)
// }
