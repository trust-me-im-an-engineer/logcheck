package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func testSlog() {
	slog.Info("hello world")          // want "msg"
	slog.DebugContext(nil, "ctx msg") // want "msg"
	slog.Log(nil, 0, "log msg")       // want "msg"
}

func testZap(l *zap.Logger, s *zap.SugaredLogger) {
	l.Info("zap info")  // want "msg"
	l.Log(0, "zap log") // want "msg"

	s.Infof("formatted %s", "msg") // want "msg"
	s.Logf(0, "logf msg")          // want "msg"
}
