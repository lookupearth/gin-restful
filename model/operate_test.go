package model

import (
	"reflect"
	"testing"
)

type OperateCase struct {
	ID      int64  `gorm:"column:id;primaryKey;->" json:"id"`                    // 自增id
	Status  *int64 `gorm:"column:status;default:2" json:"status" operate:">="`   // 状态
	Status2 int64  `gorm:"column:status2;default:2" json:"status2" operate:"in"` // 状态
	Name    string `gorm:"column:name;default:app" json:"name" operate:"eq"`     // 名称
}

func TestOperate(t *testing.T) {
	opCase := &OperateCase{}
	at := reflect.TypeOf(opCase).Elem()
	f0, _ := at.FieldByName("ID")
	op0 := NewOperate(f0)
	if op0.Operate != "=" {
		t.Errorf("Name fail")
	}
	f1, _ := at.FieldByName("Status")
	op1 := NewOperate(f1)
	if op1.Operate != ">=" {
		t.Errorf("Status fail")
	}
	f2, _ := at.FieldByName("Status2")
	op2 := NewOperate(f2)
	if op2.Operate != "IN" {
		t.Errorf("Status2 fail")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Name fail")
		}
	}()
	f3, _ := at.FieldByName("Name")
	_ = NewOperate(f3)
}
