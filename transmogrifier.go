// transmogrifier: transmogrifier package, is the main package for mog.
package transmogrifier

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	EnvMogTOML = "EnvMogTOML"
)

// AppConfig contains the application configuration settings.
var AppConfig appConfig

type appConfig struct {
	header          bool
	format          bool
	log             bool
	logfile         string
	compression     string
	destinationType string
}

// SetEnv sets the environment variables, if they do not already exist
func SetEnv() error {
	var err error
	var tmp string
	tmp = os.Getenv(EnvMogTOML)

	if tmp == "" {
		tmp = "mog.toml"
	}

	_, err = toml.DecodeFile(tmp, &AppConfig)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

/*
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
*/

// Append argstring handles the appending of a string
func appendArgString(a, s string) string {
	if s == "" {
		return a
	}

	if a == "" {
		return s
	}

	a += "," + s
	return a
}
