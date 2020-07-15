package trill

type IMessage interface {
	GetDataLen() uint32
	SetDataLen(uint32)

	GetData() []byte
	SetData([]byte)

	GetMsgID() uint32
	SetMsgID(uint32)
}

type Message struct {
	Data    []byte
	DataLen uint32
	ID      uint32

}

func NewMsgPacket(id uint32, data []byte) *Message {
	m := &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	return m
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) SetMsgID(msgID uint32) {
	m.ID = msgID
}
