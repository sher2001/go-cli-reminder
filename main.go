package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "GOLANG_CLI_REMINDER"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <text message>\n", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()

	when := when.New(nil)
	when.Add(en.All...)
	when.Add(common.All...)
	res, err := when.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if res == nil {
		fmt.Println("Unable to parse time")
		os.Exit(2)
	}

	if now.After(res.Time) {
		fmt.Println("Set a future time !")
		os.Exit(3)
	}

	tInterval := res.Time.Sub(now)

	if os.Getenv(markName) == markValue {
		time.Sleep(tInterval)
		if err := beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png"); err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		fmt.Println("Reminder will be displayed after", tInterval.Round(time.Second))
		os.Exit(0)
	}
}
