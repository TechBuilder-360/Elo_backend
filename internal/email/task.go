package email

import (
	"encoding/json"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/queue"
)

func NewEmailWelcomeTask(request WelcomeEmailPayload) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return queue.Enqueue(constant.TaskTypeWelcomeEmail, payload, 0)
}

func NewOTPTask(request OTPMailRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return queue.Enqueue(constant.TaskTypeOTPEmail, payload, 0)
}
