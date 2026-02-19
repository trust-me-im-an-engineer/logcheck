package zap

// Field is a dummy for the Zap structured logger
type Field struct{}

// Functions
func NewStdLog(l *Logger) {}

// Logger Methods
type Logger struct{}

func (log *Logger) Debug(msg string, fields ...Field)        {}
func (log *Logger) Info(msg string, fields ...Field)         {}
func (log *Logger) Warn(msg string, fields ...Field)         {}
func (log *Logger) Error(msg string, fields ...Field)        {}
func (log *Logger) DPanic(msg string, fields ...Field)       {}
func (log *Logger) Panic(msg string, fields ...Field)        {}
func (log *Logger) Fatal(msg string, fields ...Field)        {}
func (log *Logger) Log(lvl any, msg string, fields ...Field) {}

// SugaredLogger Methods
type SugaredLogger struct{}

func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{})        {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})         {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})         {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{})        {}
func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{})       {}
func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{})        {}
func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{})        {}
func (s *SugaredLogger) Logw(lvl any, msg string, keysAndValues ...interface{}) {}

func (s *SugaredLogger) Debugf(template string, args ...interface{})        {}
func (s *SugaredLogger) Infof(template string, args ...interface{})         {}
func (s *SugaredLogger) Warnf(template string, args ...interface{})         {}
func (s *SugaredLogger) Errorf(template string, args ...interface{})        {}
func (s *SugaredLogger) DPanicf(template string, args ...interface{})       {}
func (s *SugaredLogger) Panicf(template string, args ...interface{})        {}
func (s *SugaredLogger) Fatalf(template string, args ...interface{})        {}
func (s *SugaredLogger) Logf(lvl any, template string, args ...interface{}) {}
