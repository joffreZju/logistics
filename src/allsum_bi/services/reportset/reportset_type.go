package reportset

import "time"

type Condition struct {
	TimeCondition
	//	IntCondition
	//	FloatCondition
	//	StringCondition
}

type TimeCondition struct {
	FieldName string
	Start     time.Time
	End       time.Time
}

//type IntCondition struct {
//	FieldName string
//	Smaller   int
//	Larger    int
//}
//
//type FloatCondition struct {
//	FieldName string
//	Smaller   float32
//	Larger    float32
//}
//
//type StringCondition struct {
//	FieldName string
//	Length    int
//}
