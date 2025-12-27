package saferoutine

import "github.com/Toflex/directory_v2/pkg/log"

func Run(process func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered from panic: %v", r)
		}
	}()

	go process()
}
