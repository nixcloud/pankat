package pankat

import (
	"fmt"
	"time"
)

func timeElapsed(reason string) func() {
	start := time.Now()
	return func() {
		timeElapsed := time.Since(start)
		fmt.Println(reason, " took", timeElapsed, "time")
	}
}
