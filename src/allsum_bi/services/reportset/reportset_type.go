package reportset

import "time"

type Condition struct {
	TimeCondition
	IntCondition
	FloatCondition
	StringCondition
	FieldName string
}

type TimeCondition struct {
	Start time.Time
	End   time.Time
}

type IntCondition struct {
	Smaller int
	Larger  int
}

type FloatCondition struct {
	Smaller float32
	Larger  float32
}

type StringCondition struct {
	Length int
}
