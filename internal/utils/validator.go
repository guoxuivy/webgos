package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Validate 使用全局验证器实例验证DTO
// 先绑定数据，再进行验证
// 注意：dto必须是指向结构体的指针，以便能够修改其字段值
func Validate(c *gin.Context, dto any) error {
	if c.Request.Method == "GET" {
		return validateUri(c, dto)
	}
	// 数据绑定
	if err := c.ShouldBind(dto); err != nil {
		return err
	}

	// 验证数据
	validate := GetValidator()
	if err := validate.Struct(dto); err != nil {
		// 使用 ValidationError 函数处理验证错误，会自动使用 label 标签
		return ValidationError(err)
	}
	return nil
}

// validateUri 使用全局验证器实例验证URI参数DTO
// 先绑定URI参数，再进行验证
// 注意：dto必须是指向结构体的指针，以便能够修改其字段值
func validateUri(c *gin.Context, dto any) error {
	// URI参数绑定
	if err := c.ShouldBindUri(dto); err != nil {
		return err
	}

	// 验证数据
	validate := GetValidator()
	if err := validate.Struct(dto); err != nil {
		// 使用 ValidationError 函数处理验证错误，会自动使用 label 标签
		return ValidationError(err)
	}
	return nil
}

// GetValidator 获取全局验证器实例
// 可以手动使用，但不建议直接使用，推荐通过 Validate 函数来验证数据
func GetValidator() *validator.Validate {
	once.Do(func() {
		// 获取Gin默认的validator实例
		validate = validator.New()
		addDefaultValidations(validate)
	})

	return validate
}

// 这里可以添加一些默认的验证规则
func addDefaultValidations(v *validator.Validate) {

	// 注册标签名称函数，支持 label 标签
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 从结构体字段的tag中获取label的值
		label := fld.Tag.Get("label")
		if label != "" {
			return label // 有label标签则使用label
		}
		return fld.Name // 无label标签则使用默认字段名
	})

	// 注册自定义验证规则：手机号验证
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true // 允许空值（配合omitempty使用）
		}
		// 简单手机号正则：11位数字，以1开头
		match, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
		return match
	})
}

// ValidationError 处理验证错误，使用label标签给出错误提示
func ValidationError(err error) error {
	// 直接返回验证器的错误，因为已经通过 RegisterTagNameFunc 配置了使用 label 标签
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err // 如果不是验证错误，直接返回原始错误信息
	}

	// 构建更友好的错误信息
	var errs []string
	for _, e := range validationErrs {
		// 由于我们已经通过 RegisterTagNameFunc 注册了标签名称函数，
		// 验证错误会自动使用 label 标签作为字段名
		errs = append(errs, fmt.Sprintf("%s %s", e.Field(), getValidationMessage(e)))
	}
	return fmt.Errorf("参数验证失败: %s", strings.Join(errs, "; "))
}

// getValidationMessage 根据验证标签返回对应的错误消息
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "为必填项"
	case "min":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("长度不能少于%s个字符", e.Param())
		}
		return fmt.Sprintf("不能小于%s", e.Param())
	case "max":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("长度不能超过%s个字符", e.Param())
		}
		return fmt.Sprintf("不能大于%s", e.Param())
	case "email":
		return "格式不正确"
	case "gte":
		return fmt.Sprintf("必须大于等于%s", e.Param())
	case "lte":
		return fmt.Sprintf("必须小于等于%s", e.Param())
	case "oneof":
		return fmt.Sprintf("必须是%s中的一个", e.Param())
	default:
		return fmt.Sprintf("验证失败（%s）", e.Tag())
	}
}
