package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type StrSlice []string

func NewStrSlice(args ...string) StrSlice {
	s := make(StrSlice, 0)
	for _, arg := range args {
		s = append(s, arg)
	}
	return s
}

func (s StrSlice) Value() (driver.Value, error) {
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
