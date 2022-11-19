package struct_utils

import (
	"encoding/json"
	"strings"
)

func MergeStruct(src interface{}, data interface{}) (err error) {
	var dataBytes []byte
	if dataBytes, err = json.Marshal(data); err != nil {
		return
	}
	err = json.NewDecoder(strings.NewReader(string(dataBytes))).Decode(&src)
	return
}
