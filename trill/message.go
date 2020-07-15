package trill

type IMessage interface {
	GetDataLen() uint32
	SetDataLen(uint32)

	GetData() []byte
	SetData([]byte)

	GetMsgID() uint32
	SetMsgID(uint32)
}

type message struct {
	data    []byte
	dataLen uint32
	id      uint32
}

func NewMsgPacket(id uint32, data []byte) *message {
	m := &message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
	return m
}

func (m *message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}

func (m *message) GetData() []byte {
	return m.data
}

func (m *message) SetData(data []byte) {
	m.data = data
}

func (m *message) GetMsgID() uint32 {
	return m.id
}

func (m *message) SetMsgID(msgID uint32) {
	m.id = msgID
}
