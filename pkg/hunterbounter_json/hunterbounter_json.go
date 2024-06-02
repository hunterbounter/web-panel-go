package hunterbounter_json

import (
	"encoding/json"
	"log"
)

func ToString(v interface{}) string {
	return string(ToJson(v))
}

func ToInt(v interface{}) int {
	log.Println("gelen veri: ", v)
	switch v.(type) {
	// check is string
	case string:
		return int(ToInt64(v.(string)))
	}

	return int(v.(float64))

}

func ToInt64(v interface{}) int64 {
	return int64(v.(float64))
}

func ToStringBeautify(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func ToJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func StructToMap(v interface{}) map[string]interface{} {
	var data map[string]interface{}
	inrec, _ := json.Marshal(v)
	json.Unmarshal(inrec, &data)
	return data
}

func ToBase64(v string) string {
	return v
}

func UnmarshalJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
