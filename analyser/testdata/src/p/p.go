package p

import (
	"context"
	"go.uber.org/zap"
	"log/slog"
)

type User struct {
	Token string
}

type CustomStruct struct {
	APIKey string
	Secret string
	Data   string
}

func testLowercaseRule() {
	var l *slog.Logger
	var z *zap.Logger
	ctx := context.Background()

	// slog
	l.Info("Starting server")           // want "log message should start with a lowercase letter"
	l.InfoContext(ctx, "Bad message")   // want "log message should start with a lowercase letter"
	l.Log(ctx, slog.LevelInfo, "Upper") // want "log message should start with a lowercase letter"
	slog.Info("Global logger")          // want "log message should start with a lowercase letter"

	// zap
	z.Error("Failed to connect") // want "log message should start with a lowercase letter"
	z.DPanic("Panic")            // want "log message should start with a lowercase letter"

	// ✅ OK
	l.Info("server started")
	z.Debug("connection established")
}

func testCharacterRules() {
	var l *slog.Logger
	var z *zap.Logger

	// Rule 2 & 3: English letters, numbers, and spaces only
	l.Info("запуск сервера") // want "log message should only contain english letters, numbers and spaces"
	l.Warn("warning!")       // want "log message should only contain english letters, numbers and spaces"
	l.Info("started 🚀")      // want "log message should only contain english letters, numbers and spaces"
	l.Info("status: 100%")   // want "log message should only contain english letters, numbers and spaces"
	l.Info("wait...")        // want "log message should only contain english letters, numbers and spaces"

	z.Fatal("Crash!!!") // want "log message should only contain english letters, numbers and spaces"

	// ✅ OK
	l.Info("server started on port 8080")
}

func testSensitiveDataRule() {
	var l *slog.Logger
	password := "12345"
	token := "abc"
	u := User{Token: "abc"}
	s := CustomStruct{APIKey: "key", Secret: "shh", Data: "public"}

	// Variable name leaks
	l.Info("user password " + password) // want "potential sensitive data leak: argument contains 'password'"
	l.Debug("auth " + u.Token)          // want "potential sensitive data leak: argument contains 'token'"
	l.Debug("auth", u.Token)            // want "potential sensitive data leak: argument contains 'token'"

	// Struct field leaks
	l.Info("data", s.APIKey) // want "potential sensitive data leak: argument contains 'key'"
	l.Info("data", s.Secret) // want "potential sensitive data leak: argument contains 'secret'"

	// Complex concatenation
	l.Error("failed with " + "token " + token) // want "potential sensitive data leak: argument contains 'token'"

	// ✅ OK
	l.Info("token validated") // "token" inside a string literal is fine
	l.Info("data", s.Data)    // non-sensitive field
}

func testEdgeCases() {
	var l *slog.Logger

	// Empty strings should not trigger the lowercase rule (usually len < 1)
	l.Info("") // ✅ OK
}
