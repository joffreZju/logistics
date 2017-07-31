package reportset

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_condition(b *testing.T) {
	jsonstr := "[{ \"field_name\": \"a\",  \"start\": \"2006-01-02T15:04:05Z07:00\",\"end\": \"2019-01-02T15:04:05Z07:00\"}, {  \"field_name\": \"b\",  \"smaller\": 0, \"larger\": 100  }]"
	var condition []Condition
	err := json.Unmarshal([]byte(jsonstr), &condition)
	fmt.Println("condition", condition, err)
	b.Log(err)
}
