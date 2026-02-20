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
