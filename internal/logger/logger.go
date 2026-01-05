// Copyright 2026 Keith Marshall
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"context"
	"log/slog"
	"os"
)

type MultiHandler struct {
	stdoutHandler slog.Handler
	stderrHandler slog.Handler
	threshold     slog.Level
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.stdoutHandler.Enabled(ctx, level) || h.stderrHandler.Enabled(ctx, level)
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= h.threshold {
		return h.stderrHandler.Handle(ctx, record)
	}
	return h.stdoutHandler.Handle(ctx, record)
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		stderrHandler: h.stderrHandler.WithAttrs(attrs),
		threshold:     h.threshold,
	}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	return &MultiHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		stderrHandler: h.stderrHandler.WithGroup(name),
	}
}

func SetupLogger(silent, verbose bool) *slog.Logger {
	var level slog.Level

	switch {
	case silent:
		level = slog.LevelError
	case verbose:
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	handler := &MultiHandler{
		stdoutHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
		stderrHandler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}),
		threshold: slog.LevelWarn,
	}

	return slog.New(handler)
}
