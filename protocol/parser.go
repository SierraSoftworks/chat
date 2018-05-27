package protocol

import (
	"fmt"

	"github.com/pkg/errors"
)

type Parser struct {
}

func (p *Parser) ParseLine(line string) (*RawCommand, error) {
	ctx := &parseContext{
		cmd: &RawCommand{
			Operand:   "",
			Arguments: []string{},
		},
		state: "start",
	}

	for i, c := range line {
		if err := ctx.Next(c, newParseLoc(0, i)); err != nil {
			return nil, errors.Wrap(err, "parser: failed to parse line")
		}
	}

	if err := ctx.Finalize(); err != nil {
		return nil, errors.Wrap(err, "parser: failed to parse line")
	}

	return ctx.cmd, nil
}

type parseLoc struct {
	line int
	col  int
}

func newParseLoc(line, col int) *parseLoc {
	return &parseLoc{line, col}
}

func (l *parseLoc) String() string {
	return fmt.Sprintf("%d, %d", l.line, l.col)
}

type parseContext struct {
	cmd *RawCommand

	state   string
	current string
}

func (c *parseContext) Next(r rune, loc *parseLoc) error {
	switch c.state {
	case "start":
		if c.isWhitespace(r) {
			// Gobble whitespace
			return nil
		}

		if r != '.' {
			return errors.Errorf("expected whitespace or '.', got '%c' instead (%s)", r, loc)
		}

		c.state = "operand"
		c.current = ""

	case "operand":
		if c.isWhitespace(r) {
			c.cmd.Operand = c.current
			c.current = ""
			c.state = "argument"
		} else if c.isAlphaLowercase(r) {
			c.current = fmt.Sprintf("%s%c", c.current, r)
		} else {
			return errors.Errorf("expected alpha character (a-z), got '%c' instead (%s)", r, loc)
		}

	case "argument":
		if c.current == "" && c.isWhitespace(r) {
			// Gobble extra whitespace at the start of arguments
			return nil
		}

		if c.current == "" && r == '"' {
			c.state = "quoted_argument"
			return nil
		}

		if r == '\\' {
			c.state = "unquoted_argument_escape"
			return nil
		}

		if c.isAlphanumeric(r) {
			c.current = fmt.Sprintf("%s%c", c.current, r)
			return nil
		}

		if c.isWhitespace(r) {
			c.cmd.Arguments = append(c.cmd.Arguments, c.current)
			c.current = ""
			c.state = "argument"
			return nil
		}

		return errors.Errorf("expected alphanumeric character (a-zA-Z0-9), got '%c' instead (%s)", r, loc)

	case "quoted_argument":
		if r == '"' {
			c.cmd.Arguments = append(c.cmd.Arguments, c.current)
			c.current = ""
			c.state = "argument"
			return nil
		}

		if r == '\\' {
			c.state = "quoted_argument_escape"
			return nil
		}

		c.current = fmt.Sprintf("%s%c", c.current, r)
		return nil

	case "unquoted_argument_escape":
		switch r {
		case 'n':
			c.current = fmt.Sprintf("%s\n", c.current)
			c.state = "argument"
		case ' ':
			c.current = fmt.Sprintf("%s ", c.current)
			c.state = "argument"
		case '\\':
			c.current = fmt.Sprintf("%s\\", c.current)
			c.state = "argument"
		default:
			return errors.Errorf("expected a valid escape sequence (\\n, \\ , \\\\), got '\\%c' instead (%s)", r, loc)
		}

	case "quoted_argument_escape":
		switch r {
		case 'n':
			c.current = fmt.Sprintf("%s\n", c.current)
			c.state = "quoted_argument"
		case '"':
			c.current = fmt.Sprintf("%s\"", c.current)
			c.state = "quoted_argument"
		case '\\':
			c.current = fmt.Sprintf("%s\\", c.current)
			c.state = "quoted_argument"
		default:
			return errors.Errorf("expected a valid escape sequence (\\n, \\\", \\\\), got '\\%c' instead (%s)", r, loc)
		}

	default:
		return errors.Errorf("encountered unexpected state '%s' (%s)", c.state, loc)
	}

	return nil
}

func (c *parseContext) Finalize() error {
	switch c.state {
	case "operand":
		c.cmd.Operand = c.current
	case "argument":
		if c.current != "" {
			c.cmd.Arguments = append(c.cmd.Arguments, c.current)
		}
	default:
		return errors.Errorf("unable to finalize in state '%s', line is likely not complete", c.state)
	}

	return nil
}

func (c *parseContext) isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func (c *parseContext) isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func (c *parseContext) isAlphaLowercase(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func (c *parseContext) isAlphaUppercase(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func (c *parseContext) isAlpha(r rune) bool {
	return c.isAlphaUppercase(r) || c.isAlphaLowercase(r)
}

func (c *parseContext) isAlphanumeric(r rune) bool {
	return c.isAlpha(r) || c.isNumeric(r)
}
