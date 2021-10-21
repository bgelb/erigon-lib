package dbg

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// PanicReplacer - does format panic to make it readable in our logs
var PanicReplacer = strings.NewReplacer("\n", " ", "\t", "", "\r", "")

// Recover - does save panic to datadir/crashreports, bud doesn't log to logger and doesn't stop the process
// it returns recovered panic as error in format friendly for our logger
// common pattern of use - assign to named output param:
//  func A() (err error) {
//	    defer func() { err = debug.Recover(err) }() // avoid crash because Erigon's core does many things
//  }
func Recover(err error) error {
	panicResult := recover()
	if panicResult == nil {
		return err
	}

	stack := string(debug.Stack())
	switch typed := panicResult.(type) {
	case error:
		err = fmt.Errorf("%w, trace: %s", typed, PanicReplacer.Replace(stack))
	default:
		err = fmt.Errorf("%+v, trace: %s", typed, PanicReplacer.Replace(stack))
	}
	return err
}