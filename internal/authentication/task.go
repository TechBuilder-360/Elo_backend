package authentication

import (
	"time"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/queue"
)

const queueName = "wallet"

func NewDefaultWalletTask(u *ent.User) error {
	return queue.Enqueue(constant.TaskDafaultWallet, queue.TaskPayload{
		TaskID:    u.ID,
		QueueName: queueName,
		Retention: time.Hour,
		Retry:     3,
		Timeout:   time.Second * 30,
		WaitTime:  0,
		Data:      u.ID,
	})
}
