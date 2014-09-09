package mogger

// tmogger contains common information for transmogrifiers. Producer and
// consumer are the terms used here because they have different start chars
// whereas source and sink don't
type tmogger struct {
	// producer is where the data is coming from.
	producer string

	// pType identifies what type of producer it is, e.g. file, *io.Reader.
	pType
	
	// consumer is where the data is going to
	consumer  string


	// cType is the type of the consumer receiving the output.
	//	[]byte	no destination needed
	//	file	destination optional, if not set the output will be
	//		`sourceFilename.md` instead of `sourceFilename.csv`.
	destinationType string
}

