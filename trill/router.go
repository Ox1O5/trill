package trill

type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}

type baseRouter struct{}

func (b *baseRouter) PreHandle(request IRequest) {}

func (b *baseRouter) Handle(request IRequest) {}

func (b *baseRouter) PostHandle(request IRequest) {}
