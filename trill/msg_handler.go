package trill

import (
	"fmt"
	"strconv"
)

type IMsgHandle interface {
	DoMsgHandle(request IRequest)
	AddRouter (msgID uint32, router IRouter)
}

type msgHandle struct {
	APIs map[uint32]IRouter
}

func NewMsgHandle() *msgHandle {
	return &msgHandle{
		APIs: make(map[uint32]IRouter),
	}
}

func (m *msgHandle) DoMsgHandle(request IRequest) {
	handler, ok := m.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API msgID = ", request.GetMsgID(), " is not FOUND!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *msgHandle) AddRouter(msgID uint32, router IRouter) {
	if _, ok := m.APIs[msgID]; ok {
		panic("repeated api, msgID = " + strconv.Itoa(int(msgID)))
	}
	m.APIs[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}
