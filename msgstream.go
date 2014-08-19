package fwk

import (
	"fmt"
	"io"
	"os"
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type msgstream struct {
	lvl Level
	w   WriteSyncer
	n   string
}

func NewMsgStream(name string, lvl Level, w WriteSyncer) msgstream {
	if w == nil {
		w = os.Stdout
	}

	return msgstream{
		lvl: lvl,
		w:   w,
		n:   fmt.Sprintf("%-20s ", name),
	}
}

func (msg msgstream) Debugf(format string, a ...interface{}) (int, error) {
	return msg.Msg(LvlDebug, format, a...)
}

func (msg msgstream) Infof(format string, a ...interface{}) (int, error) {
	return msg.Msg(LvlInfo, format, a...)
}

func (msg msgstream) Warnf(format string, a ...interface{}) (int, error) {
	defer msg.flush()
	return msg.Msg(LvlWarning, format, a...)
}

func (msg msgstream) Errorf(format string, a ...interface{}) (int, error) {
	defer msg.flush()
	return msg.Msg(LvlError, format, a...)
}

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
