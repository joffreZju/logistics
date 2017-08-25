package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	DBErrStrDuplicateKey  = "duplicate key value violates unique constraint"
	DBErrStrAlreadyExists = "already exists"
)

//封装好的类型:Str数组,Int数组,Json
//用户更方便的操作postgres中的对应类型

type JsonMap map[string]interface{}

func (m JsonMap) Value() (driver.Value, error) {
	b, e := json.Marshal(m)
	if e != nil {
		return "", e
	}
	return string(b), nil
}
func (m *JsonMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("assert pg json failed")
	}
	e := json.Unmarshal(b, m)
	if e != nil {
		return e
	}
	return nil
}

type StrSlice []string

func NewStrSlice(args ...string) StrSlice {
	s := make(StrSlice, 0)
	for _, arg := range args {
		s = append(s, arg)
	}
	return s
}

func (s StrSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return "{" + strings.Join(s, ",") + "}", nil
}

func (s *StrSlice) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New("assert pg array failed")
	}
	str := string(b)
	str = strings.TrimLeft(str, "{")
	str = strings.TrimRight(str, "}")
	for _, v := range strings.Split(str, ",") {
		*s = append(*s, v)
	}
	return nil
}

type IntSlice []int

func NewIntSlice(args ...int) IntSlice {
	s := make(IntSlice, 0)
	for _, arg := range args {
		s = append(s, arg)
	}
	return s
}

func (s IntSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	str := "{"
	for _, v := range s {
		str += fmt.Sprintf("%d,", v)
	}
	return strings.TrimRight(str, ",") + "}", nil
}

func (s *IntSlice) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New("assert pg array failed")
	}
	str := string(b)
	str = strings.TrimLeft(str, "{")
	str = strings.TrimRight(str, "}")
	for _, v := range strings.Split(str, ",") {
		if i, e := strconv.Atoi(v); e == nil {
			*s = append(*s, i)
		} else {
			return errors.New("pg array strToint failed")
		}
	}
	return nil
}
