package trill

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}

type request struct {
	conn IConnection
	data []byte
}

func (r *request) GetConnection() IConnection {
	return r.conn
}

func (r *request) GetData() []byte {
	return r.data
}
