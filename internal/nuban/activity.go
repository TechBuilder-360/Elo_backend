package nuban

import (
	"github.com/samber/do/v2"
)

type activityWrapper struct {
	svc IService
}

func NewActivityWrapper(i do.Injector) *activityWrapper {
	svc := do.MustInvoke[IService](i)
	return &activityWrapper{svc: svc}
}
