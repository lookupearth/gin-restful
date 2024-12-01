package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Operate struct {
	Operate    string
	RawOperate string
}

func getRealOp(op string) string {
	opMap := map[string]string{
		"START": "LIKE",
		"END":   "LIKE",
	}
	value, ok := opMap[op]
	if ok {
		return value
	}
	return op
}

func isLegal(op string) bool {
	ops := []string{
		"=", ">", ">=", "<", "<=", "<>", "!=", "IN", "LIKE", "NOT IN", "NOT LIKE", "START", "END",
	}
	for _, o := range ops {
		if o == op {
			return true
		}
	}
	return false
}

func NewOperate(field reflect.StructField) *Operate {
	op := field.Tag.Get("operate")
	op = strings.TrimSpace(op)
	if op == "" {
		op = "="
	}
	op = strings.ToUpper(op)
	if !isLegal(op) {
		panic(fmt.Sprintf("operate <%s> is inlegal in field <%s>", op, field.Name))
	}
	f := &Operate{
		Operate:    getRealOp(op),
		RawOperate: op,
	}
	return f
}

func (operate *Operate) Value(value interface{}) interface{} {
	if operate.RawOperate == "LIKE" || operate.RawOperate == "NOT LIKE" {
		if s, ok := value.(string); ok {
			s = "%" + s + "%"
			return s
		}
	} else if operate.RawOperate == "START" {
		if s, ok := value.(string); ok {
			s = s + "%"
			return s
		}
	} else if operate.RawOperate == "END" {
		if s, ok := value.(string); ok {
			s = "%" + s
			return s
		}
	}
	return value
}
