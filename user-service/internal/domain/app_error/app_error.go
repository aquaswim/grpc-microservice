package appError

import "github.com/joomcode/errorx"

var appErrorNamespace = errorx.NewNamespace("app_error")

var (
	ErrInternal     = errorx.NewType(appErrorNamespace, "internal")
	ErrValidation   = appErrorNamespace.NewType("validation")
	ErrNotFound     = appErrorNamespace.NewType("not_found")
	ErrUnauthorized = appErrorNamespace.NewType("unauthorized")
)

func Switch(err error) *errorx.Type {
	return errorx.TypeSwitch(err,
		// register all error here
		ErrInternal,
		ErrValidation,
		ErrNotFound,
		ErrUnauthorized,
	)
}
