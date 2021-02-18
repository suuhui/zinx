package znet

import "github.com/suuhui/v0.3/ziface"

type BaseRouter struct {}

func (baseRouter *BaseRouter) PreHandle(request ziface.IRequest) {}

func (baseRouter *BaseRouter) Handle(request ziface.IRequest) {}

func (baseRouter *BaseRouter) PostHandle(request ziface.IRequest) {}



