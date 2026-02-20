package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func testSlog() {
	slog.Info("hello world")
	slog.DebugContext(nil, "ctx msg")
	slog.Log(nil, 0, "log msg")
}

func testZap(l *zap.Logger) {
	l.Info("zap info")
	l.Log(0, "Zap log")
}
