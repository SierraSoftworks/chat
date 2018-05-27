package protocol

func NewPing() *RawCommand {
	return &RawCommand{
		Operand:   "ping",
		Arguments: []string{},
	}
}

func NewMessage(sender, message string) *RawCommand {
	return &RawCommand{
		Operand: "msg",
		Arguments: []string{
			sender,
			message,
		},
	}
}

type MessageCommand RawCommand

func (m *MessageCommand) Sender() string {
	return m.Arguments[0]
}

func (m *MessageCommand) Message() string {
	return m.Arguments[1]
}
