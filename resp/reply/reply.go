package reply

import (
	"bytes"
	"ledis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type BulkReply struct {
	Arg []byte
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

// 数组

type MultiBulkReply struct {
	Args [][]byte
}

func (m MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(m.Args)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

var theMultiBulkReply = new(MultiBulkReply)

func MakeMultiBulkReply(arg [][]byte) *MultiBulkReply {
	return theMultiBulkReply
}

// status

type StatusReply struct {
	Status string
}

func (s StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

var theStatusReply = new(StatusReply)

func MakeStatusReply() *StatusReply {
	return theStatusReply
}

// Int

type IntReply struct {
	Code int64
}

func (i IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandardErrReply struct {
	Status string
}

func (s StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsError

func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
