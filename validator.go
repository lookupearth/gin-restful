package restful

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/lookupearth/restful/response"
)

type ModelValidatorCtx interface {
	ValidateCtx(context.Context, validator.StructLevel)
}

type ModelValidator interface {
	Validate(validator.StructLevel)
}

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Register(model interface{}) {
	modelValidate, ok := model.(ModelValidator)
	if ok {
		v.Validator.RegisterStructValidation(modelValidate.Validate, model)
	}
	modelValidateCtx, ok := model.(ModelValidatorCtx)
	if ok {
		v.Validator.RegisterStructValidationCtx(modelValidateCtx.ValidateCtx, model)
	}
}

func (v *Validator) Validate(ctx context.Context, data interface{}) *response.Error {
	return v.convert(v.Validator.StructCtx(ctx, data))
}

func (v *Validator) ValidatePartial(ctx context.Context, data interface{}, fields []string) *response.Error {
	return v.convert(v.Validator.StructPartialCtx(ctx, data, fields...))
}

func (v *Validator) convert(err error) *response.Error {
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return response.NewError(500, err)
		}
		for _, err := range err.(validator.ValidationErrors) {
			return response.NewError(400, err)
		}
	}
	return nil
}
