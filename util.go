package antgo

import (
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