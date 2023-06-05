package resp

type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	Select(int) int
}
