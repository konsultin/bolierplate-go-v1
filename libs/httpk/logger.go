package httpk

// Logger interface for httpk logging
// This allows integration with any logging framework
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultLogger is a no-op logger
type DefaultLogger struct{}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {}
func (l *DefaultLogger) Infof(format string, args ...interface{})  {}
func (l *DefaultLogger) Warnf(format string, args ...interface{})  {}
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {}
