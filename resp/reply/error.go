package reply

// unknownErrBytes
type UnknowErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u UnknowErrReply) Error() string {
	return string(unknownErrBytes)
}

func (u UnknowErrReply) ToBytes() []byte {
	return unknownErrBytes
}

func MakeUnknowErrReply() *UnknowErrReply {
	return &UnknowErrReply{}
}

// ArgNumErrReply
type ArgNumErrReply struct {
	Cmd string
}

func (a ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" +
		a.Cmd +
		"' command"
}

func (a ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" +
		a.Cmd +
		"' command\r\n")
}

var theArgNumErrReply = new(ArgNumErrReply)

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// SyntaxErrReply
type SyntaxErrReply struct {
}

func (s SyntaxErrReply) Error() string {
	return "ERR Syntax error"
}

func (s SyntaxErrReply) ToBytes() []byte {
	return []byte("-ERR Syntax error\r\n")
}

var theSyntaxErrReply = new(SyntaxErrReply)

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

// WrongTypeErrReply
type WrongTypeErrReply struct {
}

func (w WrongTypeErrReply) Error() string {
	return "WrongType Operation against a key holding the wrong kind of value"
}

func (w WrongTypeErrReply) ToBytes() []byte {
	return []byte("-WrongType Operation against a key holding the wrong kind of value\r\n")
}

var theWrongTypeErrReply = new(WrongTypeErrReply)

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return theWrongTypeErrReply
}

// ProtocolErrReply
type ProtocolErrReply struct {
	Msg string
}

func (p ProtocolErrReply) Error() string {
	return "Err Protocol error" + p.Msg
}

func (p ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err Protocol error '" +
		p.Msg +
		"' \r\n")
}

var theProtocolErrReply = new(ProtocolErrReply)

func MakeProtocolErrReply() *ProtocolErrReply {
	return theProtocolErrReply
}
