package trill

import (
	"fmt"
	"github.com/Ox1O5/trill/utils"
	"strconv"
)

type IMsgHandle interface {
	DoMsgHandle(request IRequest)
	AddRouter (msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}

type msgHandle struct {
	APIs map[uint32]IRouter
	TaskQueue []chan IRequest
	WorkerPoolSize uint32
}

func NewMsgHandle() *msgHandle {
	return &msgHandle{
		APIs: make(map[uint32]IRouter),
		TaskQueue: make([]chan IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
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

func (m *msgHandle) startOneWorker(workerID int, taskQueue chan IRequest) {
	fmt.Println("WorkerID = ", workerID, " is started")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandle(request)
		}
	}
}

func (m *msgHandle) StartWorkerPool() {
	for i:= 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

func (m *msgHandle) SendMsgToTaskQueue(request IRequest) {
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID = ",
	request.GetMsgID(), "to workerID = ",workerID)
	m.TaskQueue[workerID] <- request
}
