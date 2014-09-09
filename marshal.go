// marshal implements the marshal command
package transmogrifier


// a convenience struct to handle the arguments, easier this way
type arg struct {
	kvSeparator string
	str	string
	separator string
}

// appends the settings to the settings string.
func (a *arg) append(k, v string) {
	if v == "" {
		return
	}

	if a.str == "" {
		a.str =  a.getKVString(k,  v)
		return
	}

	a.str += a.getKVString(k, v)
	return
	
}

// getKVString creates a string version of the key value, separated by the kv
// separator. 
func (a *arg) getKVString(k,v string) string {
	if k == "" || v == "" {
		return ""
	}

	if k == "" {
		return a.kvSeparator + v
	}

	if v == "" {
		return k + a.kvSeparator
	}

	return k + a.kvSeparator + v
			
}

var Arg = &arg{kvSeparator: "=", separator: ", "}

// Marshal takes the arguments, generates a finalized transmogrification 
// template and hands the transmogrification over to the appropriate
// transmogrifier, which depends on the source type and the type it is to be
// transformed into.
func Marshal(a MarshalFilter) error {

	Arg.append("DestinationType", a.DestinationType)
	Arg.append("Format", a.Format)
	Arg.append("Header", a.Header)
	Arg.append("Log", a.Log)
	Arg.append("LogFile", a.LogFile)

}
