package protocol

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	Convey("Parser", t, func() {
		p := &Parser{}

		Convey("Valid Examples", func() {
			examples := []struct {
				Line   string
				Parsed *RawCommand
			}{
				{".ping", &RawCommand{Operand: "ping", Arguments: []string{}}},
				{"      .ping    ", &RawCommand{Operand: "ping", Arguments: []string{}}},
				{".poke ben", &RawCommand{Operand: "poke", Arguments: []string{"ben"}}},
				{`.msg ben "This is an example"`, &RawCommand{Operand: "msg", Arguments: []string{"ben", "This is an example"}}},
				{`.msg ben "This is a \"quote\""`, &RawCommand{Operand: "msg", Arguments: []string{"ben", "This is a \"quote\""}}},
				{`.msg ben test\ spaces`, &RawCommand{Operand: "msg", Arguments: []string{"ben", "test spaces"}}},
				{`.msg ben "test\nnewlines"`, &RawCommand{Operand: "msg", Arguments: []string{"ben", "test\nnewlines"}}},
			}

			for _, example := range examples {
				ex := example
				Convey(ex.Line, func() {
					cmd, err := p.ParseLine(ex.Line)
					So(err, ShouldBeNil)
					So(cmd, ShouldResemble, ex.Parsed)
				})
			}
		})

		Convey("Invalid Lines", func() {
			examples := []struct {
				Line  string
				Error string
			}{
				{"ping", "expected whitespace or '.', got 'p' instead"},
				{".poke 'ben'", "expected alphanumeric character"},
				{".10", "expected alpha character"},
				{".poke ben\\b", "expected a valid escape sequence"},
				{".poke \"ben\\b\"", "expected a valid escape sequence"},
				{`.msg ben "this is unclosed`, "unable to finalize in state 'quoted_argument'"},
				{`.msg ben "this is unclosed\`, "unable to finalize in state 'quoted_argument_escape'"},
			}

			for _, example := range examples {
				ex := example
				Convey(ex.Line, func() {
					cmd, err := p.ParseLine(ex.Line)
					So(err, ShouldNotBeNil)
					So(cmd, ShouldBeNil)
					So(err.Error(), ShouldContainSubstring, ex.Error)
				})
			}
		})
	})
}
