package antgo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JsonDecode(str []byte) map[string]interface{} {
	var dict map[string]interface{}
	err := json.Unmarshal(str, &dict)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return dict
}

func JsonEncode(dict interface{}) []byte {
	str, err := json.Marshal(dict)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return str
}

func Fastjoin(strargs ...string) string {
	var buffer bytes.Buffer
	for _, str := range strargs {
		buffer.WriteString(str)
	}
	return buffer.String()
}

func MapKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key, _ := range data {
		keys = append(keys, key)
	}
	return keys
}

func MapVals(data map[string]interface{}) []interface{} {
	vals := make([]interface{}, 0, len(data))
	for _, val := range data {
		vals = append(vals, val)
	}
	return vals
}
