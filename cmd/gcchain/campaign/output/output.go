package output

import (
	"fmt"
	"os"
	"text/template"

	cm "github.com/gcchains/chain/cmd/gcchain/campaign/common"
	"github.com/gcchains/chain/commons/log"
)

// LogOutput output data in log
type LogOutput struct {
	logger *log.Logger
}

// NewLogOutput new a object
func NewLogOutput() LogOutput {
	logger := log.New()
	return LogOutput{logger}
}

// Status shows the status of node
func (l *LogOutput) Status(status *cm.Status) {
	outTmpl := `--------------------------

Mining:           {{.Mining}}

RNode:            {{.RNode}}

Proposer:         {{.Proposer}}
--------------------------
`
	tmpl, err := template.New("status").Parse(outTmpl)
	if err != nil {
		l.Error(err.Error())
	}
	err = tmpl.Execute(os.Stdout, status)
	if err != nil {
		l.Error(err.Error())
	}
	fmt.Println()
}

// Info log
func (l *LogOutput) Info(msg string, params ...interface{}) {
	l.logger.Info(msg, params...)
}

// Error log
func (l *LogOutput) Error(msg string, params ...interface{}) {
	l.logger.Info("error:"+msg, params...)
}

// Fatal log
func (l *LogOutput) Fatal(msg string, params ...interface{}) {
	l.logger.Info("fatal:"+msg, params...)
}

// Warn log
func (l *LogOutput) Warn(msg string, params ...interface{}) {
	l.logger.Warn("warn:"+msg, params...)
}
