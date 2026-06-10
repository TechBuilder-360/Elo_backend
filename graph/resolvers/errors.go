package resolver

import (
	stderrors "errors"
	apperrors "github.com/Toflex/directory_v2/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func toGQLError(err error) error {
	if err == nil {
		return nil
	}
	var gqlErr *gqlerror.Error
	if stderrors.As(err, &gqlErr) {
		return gqlErr
	}
	var safe *apperrors.SafeError
	if stderrors.As(err, &safe) {
		return &gqlerror.Error{
			Message: safe.Message,
			Extensions: map[string]any{
				"code": string(safe.Code),
			},
		}
	}
	return &gqlerror.Error{Message: err.Error()}
}
