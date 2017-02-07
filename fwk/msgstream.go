package fwk

import (
	"fmt"
	"io"
	"os"
)

// WriteSyncer is an io.Writer which can be sync'ed/flushed.
type WriteSyncer interface {
	io.Writer
	Sync() error
}

type msgstream struct {
	lvl Level
	w   WriteSyncer
	n   string
}

// NewMsgStream creates a new MsgStream value with name name and minimum
// verbosity level lvl.
// This MsgStream will print messages into w.
func NewMsgStream(name string, lvl Level, w WriteSyncer) MsgStream {
	return newMsgStream(name, lvl, w)
}

func newMsgStream(name string, lvl Level, w WriteSyncer) msgstream {
	if w == nil {
		w = os.Stdout
	}

	return msgstream{
		lvl: lvl,
		w:   w,
		n:   fmt.Sprintf("%-20s ", name),
	}
}

// Debugf displays a (formated) DBG message
func (msg msgstream) Debugf(format string, a ...interface{}) (int, error) {
	return msg.Msg(LvlDebug, format, a...)
}

// Infof displays a (formated) INFO message
func (msg msgstream) Infof(format string, a ...interface{}) (int, error) {
	return msg.Msg(LvlInfo, format, a...)
}

// Warnf displays a (formated) WARN message
func (msg msgstream) Warnf(format string, a ...interface{}) (int, error) {
	defer msg.flush()
	return msg.Msg(LvlWarning, format, a...)
}

// Errorf displays a (formated) ERR message
func (msg msgstream) Errorf(format string, a ...interface{}) (int, error) {
	defer msg.flush()
	return msg.Msg(LvlError, format, a...)
}

// Msg displays a (formated) message with level lvl.
func (msg msgstream) Msg(lvl Level, format string, a ...interface{}) (int, error) {
	if lvl < msg.lvl {
		return 0, nil
	}
	format = msg.n + msg.lvl.msgstring() + " " + format
	return fmt.Fprintf(msg.w, format, a...)
}

func (msg msgstream) flush() error {
	return msg.w.Sync()
}

// EOF
