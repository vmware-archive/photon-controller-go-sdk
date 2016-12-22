package logs

// Logger is an interface that extracts all the relevant methods of the stdlib logger.
// But it allows for different logging libraries to be used as well.
// There's no standard interface, this is the closest we get, unfortunately.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
