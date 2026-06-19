package verification

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/queue"
)

const queueName = "verification"

func QueueVerificationTask(request VerificationResult) error {
	byt, err := json.Marshal(request)
	if err != nil {
		return err
	}

	b64 := base64.RawStdEncoding.EncodeToString(byt)

	fmt.Println(string(b64))
	return queue.Enqueue(constant.TaskTypeIdentityVerification, queue.TaskPayload{
		TaskID:    request.ReferenceID,
		QueueName: queueName,
		Retention: time.Hour,
		Retry:     3,
		Timeout:   time.Minute,
		WaitTime:  0,
		Data:      b64,
	})
}
