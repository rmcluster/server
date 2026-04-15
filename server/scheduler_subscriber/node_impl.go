package schedulersubscriber

import "github.com/wk-y/rama-swap/microservices/scheduling"

type node struct {
	id   string
	ip   string
	port int
}

// Id implements [scheduling.Node].
func (n *node) Id() string {
	return n.id
}

// Ip implements [scheduling.Node].
func (n *node) Ip() string {
	return n.ip
}

// Port implements [scheduling.Node].
func (n *node) Port() int {
	return n.port
}

var _ scheduling.Node = (*node)(nil)
