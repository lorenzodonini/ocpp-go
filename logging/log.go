package logging

// Logger is the adapter interface that needs to be implemented, if the library should internally print logs.
//
// This allows to hook up your logger of choice.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// VoidLogger is an empty implementation of the Logger interface, which doesn't actually process any logs.
// It may be used as a dummy implementation, if no logs should be visible.
type VoidLogger struct{}

func (l *VoidLogger) Debug(args ...interface{})                 {}
func (l *VoidLogger) Debugf(format string, args ...interface{}) {}
func (l *VoidLogger) Info(args ...interface{})                  {}
func (l *VoidLogger) Infof(format string, args ...interface{})  {}
func (l *VoidLogger) Error(args ...interface{})                 {}
func (l *VoidLogger) Errorf(format string, args ...interface{}) {}
