package test_handler

import "log"

type EchoReq struct {
	Contant string
}
type EchoReply struct {
	Contant string
}

type Echo struct {
}

func (e *Echo) Ping(req EchoReq, reply *EchoReply) (err error) {
	log.Println("recv:", req.Contant)
	reply.Contant = req.Contant
	return
}
