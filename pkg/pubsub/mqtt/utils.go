package mqtt

import "encoding/json"

func ConvertForMQTT(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
