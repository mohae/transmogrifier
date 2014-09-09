package mogger


// Formatters describe the format of something.
type  formatter interface {
	FormatInfo(interface{})
	GetFormatInfo(interface{})
	HeaderInfo(interface{})
	GetHeaderInfo() interface{}
}
