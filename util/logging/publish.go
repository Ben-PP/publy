package logging

import "log/slog"

func LogPublish(
	src_ip string,
	method string, path string,
	userAgent string,
	success bool,
	pub string,
	output string,
) {
	log(
		slog.LevelInfo,
		"Publish event triggered",
		"audit",
		slog.Group(
			"event",
			slog.String("eventtype", "audit:publish"),
			slog.String("src_ip", src_ip),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("user_agent", userAgent),
			slog.String("pub", pub),
			slog.Bool("success", success),
			slog.String("output", output),
		),
	)
}
