package validate

import (
	"admin_backend/pkg/common/xerr"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Validator 定义验证器结构
type Validator struct {
	validate *validator.Validate
}

// New 创建一个新的验证器实例
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Struct 验证结构体
func (v *Validator) Struct(s any) error {
	return v.validate.Struct(s)
}

func (v *Validator) Validate(r *http.Request, data any) error {
	if err := v.Struct(data); err != nil {
		return xerr.NewErrCodeMsg(xerr.ParamError, "参数错误:"+err.Error())
	}
	return nil
}

// 全局验证器实例
var DefaultValidator = New()
