package reportset

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_condition(b *testing.T) {
	jsonstr := "[{ \"field_name\": \"a\",  \"start\": \"2017-07-31T12:07:41.709073657+08:00\",\"end\": \"2017-07-31T12:08:41.709073657+08:00\"}, {  \"field_name\": \"b\",  \"smaller\": 0, \"larger\": 100  }]"
	jsonstrsimple := "{  \"field_name\": \"b\",  \"smaller\": 0, \"larger\": 100  }"
	var condition []TimeCondition
	err := json.Unmarshal([]byte(jsonstr), &condition)
	fmt.Println("condition", condition, err)
	var conditionsimple Condition
	err = json.Unmarshal([]byte(jsonstrsimple), &conditionsimple)
	fmt.Println("condition simple", conditionsimple, err)
	conditionstruct := TimeCondition{
		FieldName: "test",
		//	Start:     time.Time("2006-01-02T15:04:05Z07:00"),
		Start: time.Now(),
		//	End:       time.Time("2009-01-02T15:04:05Z07:00"),
		End: time.Now(),
	}
	btstr, err := json.Marshal(conditionstruct)
	fmt.Println("btstr ", string(btstr), err)
	b.Log(err)
}
