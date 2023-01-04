package pankat

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

func timeElapsed(reason string) func() {
	fmt.Println(color.GreenString(">>>>>>>>>>>>>>>>>>>>>>>>>>"), reason, "start!", color.GreenString(">>>>>>>>>>>>>>>>>>>>>>>>>>"))

	start := time.Now()
	return func() {
		timeElapsed := time.Since(start)
		fmt.Println(color.GreenString("<<<<<<<<<<<<<<<<<<<<<<<<<<"), reason, "finished after", timeElapsed, color.GreenString("<<<<<<<<<<<<<<<<<<<<<<<<<<"))
	}
}
