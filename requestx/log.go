package requestx

import "fmt"

type Logger interface {
	Debug(msg string, v ...any)
	Info(msg string, v ...any)
}

type ZapLog struct{}

func (*ZapLog) Debug(msg string, v ...any) {
	fmt.Printf(msg, v...)
}
func (*ZapLog) Info(msg string, v ...any) {
	fmt.Printf(msg, v...)
}
