package restful

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

type VaTest struct {
	Name  string `validate:"required"`
	Name2 string `validate:"required"`
}

func (d *VaTest) ValidateCtx(ctx context.Context, sl validator.StructLevel) {
	fmt.Println("run validator ValidateCtx")
}

func TestValidator(t *testing.T) {
	va := validator.New()
	valid := Validator{
		Validator: va,
	}
	vaTest := &VaTest{
		Name: "123",
	}
	valid.Register(vaTest)
	ctx := context.Background()
	err := valid.Validate(ctx, vaTest)
	if err == nil {
		t.Errorf("Validator.Validate fail, error=%v", err)
	}
	fmt.Println(err)
	if err := valid.ValidatePartial(ctx, vaTest, []string{"Name"}); err != nil {
		t.Errorf("Validator.ValidatePartial fail, error=%v", err)
	}

}
