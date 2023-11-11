package main

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	eridaAddr   = os.Getenv("ERIDA_ADDR")
	from        = os.Getenv("ERIDA_FROM")
	to          = strings.Split(os.Getenv("ERIDA_TO"), ",")
	bodyPattern = os.Getenv("ERIDA_BODY_PATTERN")
)

// following code is a stress test for erida
// it will send 10 messages to erida server
//
// they will be launched in the same cluster
func main() {
	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			issueRequest(i)
		}(i)
	}

	wg.Wait()
}

func issueRequest(reqNumber int) {
	erida, err := smtp.Dial(eridaAddr)
	if err != nil {
		slog.Error("failed to dial erida server", "error", err)
		return
	}

	defer func() {
		if err = erida.Quit(); err != nil {
			slog.Error("failed to quit erida server", "error", err)
		}
	}()

	var body = strings.NewReader(fmt.Sprintf(bodyPattern, reqNumber))
	if err = erida.SendMail(from, to, body); err != nil {
		slog.Error("failed to send email", "error", err)
		return
	}

	slog.Info("email sent", "reqNumber", reqNumber)
}
