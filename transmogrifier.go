// transmogrifier: transmogrifier package, is the main package for mog.
package transmogrifier

import (
	"os"

	"github.com/BurntSushi/toml"
	seelog "github.com/cihub/seelog"
	"github.com/mohae/transmogrifier/format"
	"github.com/mohae/transmogrifier/mog"
	"github.com/mohae/transmogrifier/tmog"
)

var logger seelog.LoggerInterface

const (
	EnvMogTOML	= "EnvMogTOML"
)

const (
	logging = "logging"
)

func init() {
	// Disable logger by default
	DisableLog()
}

// AppConfig contains the application configuration settings.
var AppConfig appConfig

type appConfig struct {
	header bool
	format bool
	log bool
	logfile string
	compression string
	destinationType string
}

// DisableLog disables all package log output
func DisableLog() {
	logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output lobrary log.
// Use this func if you are using Seelog logging system in your app--I am
// so I'm done
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
}

// Call this before app shutdown
func FlushLog() {
	logger.Flush()
}

// SetEnv sets the environment variables, if they do not already exist
func SetEnv() error {
	var err error
	var tmp string
	tmp = os.Getenv(EnvMogTOML)

	if tmp == "" {
		tmp = "mog.toml"
	}

	_, err  = toml.DecodeFile(tmp, &AppConfig)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

// Set sets the values for the keys within the passed settings. A key that does
// not exist is considered an error. Any archiver setting is a legitimate map
// key. Unmatched keys will result in an error.
// We handle everything as a string because we won't know when to override the
// boolean values if they were of type bool. 
func (t *Transmogrifier) Set(settings map[string]string) error {

	if len(settings) == 0 {
		return errors.New("Unable to initialize Tranmogrifier: no settings were received")
	}

	for k, v := settings {
		switch "k" {
		case "compression":
			if v != "" {
				c.compression = v
			}
		case "destinationtype":
			if v != "" {
				a.destinationType = v
			}
		case "format":
			if v:= "" {
				a.format = v.(bool)
			}
		case "header":
			if v != "" {
				a.header = v.(bool)
			}

		case "log":
			if v != "" {
				a.log = v.(bool)
			}
		case "logfile": 
			a.logfile = v
		default:
			return errors.New("Unsupported setting received " + k + ":" + v.(string))
		}		
	}

	return nil
}

// Append argstring handles the appending of a string 
func appendArgString(a, s  string) string {
	if s == "" {
		return a
	}

	if a == "" {
		return s
	}

	a += "," + s
	return a
}
