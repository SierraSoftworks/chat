package protocol

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRawCommand(t *testing.T) {
	Convey("RawCommand", t, func() {
		Convey("String()", func() {
			examples := []struct {
				Cmd    *RawCommand
				String string
			}{
				{&RawCommand{}, ""},
				{&RawCommand{"poke", []string{}}, ".poke"},
				{&RawCommand{"msg", []string{"ben", "this is a test"}}, `.msg ben "this is a test"`},
			}

			for _, example := range examples {
				ex := example
				Convey(ex.String, func() {
					So(ex.Cmd.String(), ShouldEqual, ex.String)
				})
			}
		})
	})
}
