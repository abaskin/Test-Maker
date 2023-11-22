package quizconfig

import "log"

// logger implements the paho.Logger interface
type Logger struct {
	Prefix string
}

// Println is the library provided NOOPLogger's
// implementation of the required interface function()
func (l Logger) Println(v ...interface{}) {
	log.Println(append([]interface{}{l.Prefix + ":"}, v...)...)
}

// Printf is the library provided NOOPLogger's
// implementation of the required interface function(){}
func (l Logger) Printf(format string, v ...interface{}) {
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format = format + "\n" // some log calls in paho do not add \n
	}
	log.Printf(l.Prefix+":"+format, v...)
}
