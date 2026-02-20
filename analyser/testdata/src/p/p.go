package p

import (
	"go.uber.org/zap"
	"log/slog"
)

func warnings() {
	var l *slog.Logger
	var z *zap.Logger
	password := "12345"

	// Rule 1: Lowercase
	l.Info("Starting server")    // want "log message should start with a lowercase letter"
	z.Error("Failed to connect") // want "log message should start with a lowercase letter"

	// Rule 2 & 3: English, Special chars & Emoji
	l.Info("запуск сервера") // want "log message should only contain english letters, numbers and spaces"
	l.Warn("warning!")       // want "log message should only contain english letters, numbers and spaces"
	l.Info("started 🚀")      // want "log message should only contain english letters, numbers and spaces"

	// Rule 4: Sensitive Data
	l.Info("user password " + password) // want "potential sensitive data leak: argument contains 'password'"

	type User struct {
		Token string
	}
	u := User{Token: "abc"}
	l.Debug("auth " + u.Token) // want "potential sensitive data leak: argument contains 'token'"

	// ✅ Correct cases
	l.Info("server started")
	z.Debug("connection established")
	l.Info("token validated") // "token" в строке разрешен, если это не имя переменной/поля (согласно вашему checkLogArg)
}
