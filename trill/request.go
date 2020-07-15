package trill

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
}

type request struct {
	conn IConnection
	msg  IMessage
}

func (r *request) GetConnection() IConnection {
	return r.conn
}

func (r *request) GetData() []byte {
	return r.msg.GetData()
}

func (r *request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
