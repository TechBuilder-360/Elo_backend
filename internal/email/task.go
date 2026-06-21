package email

import (
	"time"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/queue"
	"github.com/Toflex/directory_v2/pkg/util"
)

const queueName = "email"

func NewEmailWelcomeTask(request WelcomeEmailPayload) error {
	return queue.Enqueue(constant.TaskTypeWelcomeEmail, queue.TaskPayload{
		TaskID:    util.GenerateUUID(),
		QueueName: queueName,
		Retention: time.Hour,
		Retry:     3,
		Timeout:   time.Second * 30,
		WaitTime:  0,
		Data:      request,
	})
}

func NewOTPTask(request OTPMailRequest) error {
	return queue.Enqueue(constant.TaskTypeOTPEmail, queue.TaskPayload{
		TaskID:    util.GenerateUUID(),
		QueueName: queueName,
		Retention: time.Hour,
		Retry:     3,
		Timeout:   time.Second * 30,
		WaitTime:  0,
		Data:      request,
	})
}

func NewUserVerificationTask(request VerificationMailPayload) error {
	return queue.Enqueue(constant.TaskUserVerification, queue.TaskPayload{
		TaskID:    util.GenerateUUID(),
		QueueName: queueName,
		Retention: time.Hour,
		Retry:     3,
		Timeout:   time.Second * 30,
		WaitTime:  0,
		Data:      request,
	})
}
