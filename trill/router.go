package trill

type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}

type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request IRequest) {}

func (b *BaseRouter) Handle(request IRequest) {}

func (b *BaseRouter) PostHandle(request IRequest) {}
