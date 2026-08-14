package quote

import (
	"strings"
	"unicode/utf8"

	"github.com/Equationzhao/g/internal/request"
)

func Name(name string, mode request.QuoteMode, zero bool) string {
	if zero {
		return strings.ReplaceAll(name, "\x00", "?")
	}
	switch mode {
	case request.QuoteAlways:
		return always(name)
	case request.QuoteLiteral:
		return literal(name)
	default:
		if needsQuote(name) {
			return single(name)
		}
		return name
	}
}

func needsQuote(s string) bool {
	if s == "" {
		return true
	}
	if !utf8.ValidString(s) {
		return true
	}
	return strings.ContainsAny(s, " \t'\"\\\n\r")
}

func single(s string) string {
	var b strings.Builder
	b.WriteByte('\'')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'':
			b.WriteString(`'\''`)
		case c == '\n':
			b.WriteString(`\n`)
		case !utf8.ValidString(s[i : i+1]):
			b.WriteString(`\x`)
			const hex = "0123456789abcdef"
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&0xf])
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

func always(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\', '"', '$', '`':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('"')
	return b.String()
}

func literal(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == 0x7f {
			b.WriteByte('?')
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

func Classify(kind byte) string {
	switch kind {
	case 'd':
		return "/"
	case 'x':
		return "*"
	case 'l':
		return "@"
	case 'p':
		return "|"
	case 's':
		return "="
	default:
		return ""
	}
}
