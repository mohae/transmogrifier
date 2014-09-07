package tomd

// MD constants
var (
	// Pipe is the MD column separator
	mdPipe []byte = []byte("|")
	
	// LeftJustify is the MD for left justification of columns.
	mdLeftJustify []byte = []byte(":--")

	// RightJustify is the Md for right justification of columns,
	mdRightJustify []byte = []byte("--:")
	mdCentered []byte = []byte(":--:")
	mdDontJustify []byte = []byte("--")
)

func GetMDPipe() []byte {
	return mdPipe
}

func GetMDLeftJustify() []byte {
	return mdLeftJustify
}

func GetMDRightJustify() []byte {
	return mdRightJustify
}

func GetMDCentered() []byte {
	return mdCentered
}

func GetMDDontJustify() []byte {
	return mdDontJustify
}

