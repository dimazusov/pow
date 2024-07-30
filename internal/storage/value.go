package storage

type Value struct {
	value []byte
}

func NewValue() *Value {
	return &Value{}
}

func (m *Value) Get() []byte {
	return m.value
}

func (m *Value) Set(val []byte) {
	m.value = val
}
