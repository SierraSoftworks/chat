package protocol

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommands(t *testing.T) {
	Convey("Commands", t, func() {
		Convey("NewPing()", func() {
			So(NewPing(), ShouldNotBeNil)
			So(NewPing(), ShouldResemble, &RawCommand{
				Operand:   "ping",
				Arguments: []string{},
			})
		})

		Convey("NewMessage()", func() {
			So(NewMessage("test", "this is a test"), ShouldNotBeNil)
			So(NewMessage("test", "this is a test"), ShouldResemble, &RawCommand{
				Operand:   "msg",
				Arguments: []string{"test", "this is a test"},
			})
		})
	})
}
