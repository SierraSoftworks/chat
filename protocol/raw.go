package protocol

import (
	"fmt"
	"strings"
)

type RawCommand struct {
	Operand   string
	Arguments []string
}

func (c *RawCommand) String() string {
	if c.Operand == "" {
		return ""
	}

	out := fmt.Sprintf(".%s", c.Operand)

	for _, arg := range c.Arguments {
		out = fmt.Sprintf("%s %s", out, c.formatArgument(arg))
	}

	return out
}

func (c *RawCommand) formatArgument(arg string) string {
	if c.isBasicString(arg) {
		return arg
	}

	b := &strings.Builder{}
	b.WriteRune('"')
	for _, c := range arg {
		switch c {
		case '"':
			b.WriteRune('\\')
			b.WriteRune(c)
		case '\n':
			b.WriteRune('\\')
			b.WriteRune('n')
		case '\\':
			b.WriteRune('\\')
			b.WriteRune('\\')
		default:
			b.WriteRune(c)
		}

	}
	b.WriteRune('"')

	return b.String()
}

func (c *RawCommand) isBasicString(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			continue
		}

		if c >= 'A' && c <= 'Z' {
			continue
		}

		if c >= '0' && c <= '9' {
			continue
		}

		return false
	}

	return true
}
