package parser

import (
	"io"
	"ledis/interface/resp"
)

type Payload struct {
	Data resp.Reply // 发送和返回的数据结构一致，所以用reply
	Err  error
}

type readState struct {
	readingMultiLine  bool // 单行数据还是多行数据
	expectedArgsCount int  // 参数个数
	msgType           byte
	args              [][]byte // 已经解析的长度
	bulkLen           int64
}

// 解析器是否解析完成
func (r *readState) finished() bool {
	// 已经解析的数量和目标参数数量
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

// 解析器
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}
func parse0(reader io.Reader, ch chan *Payload) {

}
