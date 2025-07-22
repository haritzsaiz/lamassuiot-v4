package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"
)

type FormatterHandler struct {
	out       io.Writer
	level     Leveler
	formatter *Formatter
	goas      []groupOrAttrs
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

type Formatter struct {
	FieldsOrder           []string
	TimestampFormat       string
	HideKeys              bool
	NoColors              bool
	NoFieldsColors        bool
	NoFieldsSpace         bool
	ShowFullLevel         bool
	NoUppercaseLevel      bool
	TrimMessages          bool
	CallerFirst           bool
	CustomCallerFormatter func(pc *runtime.Frame) string
	ReportCaller          bool
}

func NewFormatterHandler(w io.Writer, l Leveler, f *Formatter) slog.Handler {
	if f.TimestampFormat == "" {
		f.TimestampFormat = time.StampMilli
	}
	return &FormatterHandler{out: w, level: l, formatter: f}
}

func (h *FormatterHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.Level(h.level.Level())
}

func (h *FormatterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *FormatterHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *FormatterHandler) withGroupOrAttrs(goa groupOrAttrs) *FormatterHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *FormatterHandler) Handle(ctx context.Context, r slog.Record) error {
	f := h.formatter
	var b bytes.Buffer

	// Time
	b.WriteString(r.Time.Format(f.TimestampFormat))

	// Caller
	if f.CallerFirst {
		writeCaller(&b, r, f)
	}

	// Level
	level := r.Level.String()
	if !f.NoUppercaseLevel {
		level = strings.ToUpper(level)
	}
	if !f.NoColors {
		fmt.Fprintf(&b, "\x1b[%dm", getColorByLevel(r.Level))
	}
	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")
	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// Attributes
	fields := make(map[string]any)
	for _, goa := range h.goas {
		for _, attr := range goa.attrs {
			fields[attr.Key] = attr.Value.Any()
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	if f.FieldsOrder != nil {
		writeOrderedFields(&b, fields, f)
	} else {
		writeSortedFields(&b, fields, f)
	}
	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// Message
	msg := r.Message
	if f.TrimMessages {
		msg = strings.TrimSpace(msg)
	}
	b.WriteString(msg)

	if !f.CallerFirst {
		writeCaller(&b, r, f)
	}
	b.WriteByte('\n')

	_, err := h.out.Write(b.Bytes())
	return err
}

func writeSortedFields(b *bytes.Buffer, fields map[string]any, f *Formatter) {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		writeField(b, key, fields[key], f)
	}
}

func writeOrderedFields(b *bytes.Buffer, fields map[string]any, f *Formatter) {
	seen := map[string]bool{}
	for _, key := range f.FieldsOrder {
		if val, ok := fields[key]; ok {
			seen[key] = true
			writeField(b, key, val, f)
		}
	}
	var rest []string
	for k := range fields {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	for _, key := range rest {
		writeField(b, key, fields[key], f)
	}
}

func writeField(b *bytes.Buffer, key string, val any, f *Formatter) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", val)
	} else {
		fmt.Fprintf(b, "[%s:%v]", key, val)
	}
	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

func writeCaller(b *bytes.Buffer, r slog.Record, f *Formatter) {
	if f.ReportCaller {
		if pc := r.PC; pc != 0 {
			fn := runtime.FuncForPC(pc)
			file, line := fn.FileLine(pc)
			frame := &runtime.Frame{Function: fn.Name(), File: file, Line: line}

			if f.CustomCallerFormatter != nil {
				b.WriteString(f.CustomCallerFormatter(frame))
			} else {
				fmt.Fprintf(b, " (%s:%d %s)", path.Base(file), line, frame.Function)
			}
		}
	}
}

func getColorByLevel(level slog.Level) int {
	switch {
	case level <= slog.LevelDebug-4:
		return 37 // gray
	case level == slog.LevelDebug:
		return 37 // gray
	case level == slog.LevelInfo:
		return 36 // blue
	case level == slog.LevelWarn:
		return 33 // yellow
	case level >= slog.LevelError:
		return 31 // red
	default:
		return 36 // fallback blue
	}
}
