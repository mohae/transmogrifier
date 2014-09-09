// Contains log related stuff.
package tomd 

import (
	"errors"
	"io"

	seelog "github.com/cihub/seelog"\
        "github.com/mohae/tranmogrifier/format"
        "github.com/mohae/tranmogrifier/mog"
	"github.com/mohae/tranmogrifier/tmog"
)

var logger seelog.LoggerInterface

func init() {
	//Disable logger by default
	DisableLog()
}

// DisableLog disables all package output
func DisableLog() {
	format.DisableLog()
	mog.DisableLog()
	tmog.DisableLog()
        logger = seelog.Disabled
}

// UseLoggers uses a specified seelog.LoggerInterface to output package to log.
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
	format.UseLogger(logger)
	mog.UseLogger(logger)
	tmog.UseLogger(tmog)
}

// SetLogWriter uses a specified io.Writer to output library log.
// Use this func if you are not using Seelog logging system in your app.
func SetLogWriter(writer io.Writer) error {
	if writer == nil {
		return errors.New("Nil writer")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevel(writer, seelog.TraceLvl)
	if err != nil {
		return err
	}

	UseLogger(newLogger)
	return nil
}

// FlushLog, call before app shutdown. This is called by realMain(). If a
// logger other than Seelog is going to be used, use the 
func FlushLog() {
	format.FlushLog()
	mog.FlushLog()
	tmog.FlushLog()
	logger.Flush()
}
