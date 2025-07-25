package util

import (
	"fmt"
	"time"

	"github.com/tj/go-spin"
)

func Spin(label string, close chan struct{}) {
	go func() {
		s := spin.New()
	LOOP:
		for {
			fmt.Printf("\r\033[m%s \033[m %s", label, s.Next())
			time.Sleep(100 * time.Millisecond)
			select {
			case _, ok := <-close:
				if !ok {
					break LOOP
				}
			default:
			}
		}
	}()
}
