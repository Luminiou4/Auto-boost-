package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Spin struct {
	s    *spinner.Spinner
	msg  string
}

func NewSpinner(message string) *Spin {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = "  " + message
	s.Color("cyan", "bold")
	return &Spin{s: s, msg: message}
}

func (sp *Spin) Start() {
	sp.s.Start()
}

func (sp *Spin) UpdateMessage(msg string) {
	sp.s.Suffix = "  " + msg
	sp.msg = msg
}

func (sp *Spin) Stop(finalMessage string) {
	sp.s.Stop()
	if finalMessage != "" {
		fmt.Println(finalMessage)
	}
}

func (sp *Spin) StopSuccess(msg string) {
	sp.s.FinalMSG = "\r\033[92m✔\033[0m  " + msg + "\n"
	sp.s.Stop()
}

func (sp *Spin) StopFail(msg string) {
	sp.s.FinalMSG = "\r\033[91m✘\033[0m  " + msg + "\n"
	sp.s.Stop()
}
