package dict

type Consumer func(key string, val interface{}) bool

// 字典类型 Map

type Dict interface {
	Get(key string) (val interface{}, exist bool)
	Len() int64
	Put(key string, val interface{}) (result int)
	PutIfAbsent(key string, val interface{}) (result int) // 如果不存在
	PutIfExists(key string, val interface{}) (result int) // 如果存在
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	clear()
}
