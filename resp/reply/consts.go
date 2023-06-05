package reply

//回复固定的消息

// PONG Reply

type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func (p PongReply) ToBytes() []byte {
	return pongBytes
}

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

// OK Reply

type OKReply struct {
}

var okBytes = []byte("+ok\r\n")

func (O OKReply) ToBytes() []byte {
	return okBytes
}

var theOKReply = new(OKReply)

func MakeOKReply() *OKReply {
	return theOKReply
}

// null bulk Reply 空字符串回复

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

var theNullBulkReply = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

// emptyMultiBulk Reply 空数组回复

type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (n EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

// No Reply

type NoReply struct {
}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}

var theNoReply = new(NoReply)

func MakeNoReply() *NoReply {
	return theNoReply
}
