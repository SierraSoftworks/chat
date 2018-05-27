# Chat Protocol

Chat uses a text-based, UTF-8 encoded, line protocol in which messages are encoded as a command code, arguments tuple.
The protocol is intended to be easily visually debug-able and should be robust against most erroneous input. It is also
stateless at this time, meaning that there is no need for a client to perform an authentication handshake or jump through
hoops before sending messages.

Messages are received and sent over TCP with the same line protocol being used for both senders and receivers. In other
words: if a sender writes `.msg sender_name "This is an example"` then all other clients will receive the exact same
line from the server.

*Chat* also operates with a single, global, room - without support for direct messages or specific channels at this
stage. This simplifies implementation greatly and is one of the features which enables the protocol to be stateless.

## Line Format
Here is an EBNF grammar for the protocol, you probably won't need it but if you're curious as to what the exact format
should look like then this is a good place to check. It is written in PEG format and you can try it out online
[here](https://pegjs.org/online).

### Grammar
```peg
Commands =
	head:Command
	tail:("\n" Command)* { return [head, ...tail.map(x => x[1])]; }

Command =
  op:OperationCode args:(_+ ArgumentsList)? { return { operand: op, arguments: args && args[1] || [] } }

OperationCode =
    "." operand:[a-z]* { return operand.join(""); }

ArgumentsList =
    head:Argument
    tail:(_ Argument)* { return [head, ...tail.map(x => x[1])]; }

Argument =
    QuotedArgument / UnquotedArgument

UnquotedArgument =
    [a-zA-Z0-9]+ { return text(); }

QuotedArgument =
    '"' chars:QuotedChar* '"' { return chars.join(""); }

QuotedChar =
    !('"' / "\\") . { return text() }
    / "\\" escape:EscapeSequence { return escape; }

EscapeSequence =
    '"' /
    '\\' /
    "n" { return "\n"; }


_ "whitespace" =
    " " / "\n" / "\t"
```

### Example
```
.msg sender "message with a...\n...newline in it"
.ping
.msg sender hi
```

## Commands

### `.msg`
The `.msg` command is used by clients wishing to send a message on the *Chat* server. It accepts two arguments: the name
of the sender and the message to be sent.

```
.msg SENDER MESSAGE
```

In most cases, the `SENDER` argument will be an unquoted string, however this is not required. Similarly, the `MESSAGE`
argument will usually be a quoted string to allow the use of spaces and other special characters.

```
CLIENT: .msg ben "This is an example of a common message"
```

Of course, one may also use more complex `SENDER` names, or simpler `MESSAGE`s.

```
CLIENT: .msg "Benjamin Pannell" Hi
```

### `.ping`
The `.ping` command is a special command which the server uses to ensure that clients are still online and connected.
When a client receives a `.ping` message, it is expected to respond to the server with a `.ping` response or it will
be disconnected.

```
SERVER: .ping
CLIENT: .ping
```