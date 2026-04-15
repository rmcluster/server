package scheduling

import (
	"context"
	"os/exec"
	"sync"

	"github.com/wk-y/rama-swap/llama"
)

func NewInstanceFactory(llmService *llama.Llama, lowestPort int) InstanceFactory {
	return &instanceFactoryImpl{
		llmService: llmService,
		lowestPort: lowestPort,
		usedPorts:  make(map[int]struct{}),
	}
}

type instanceFactoryImpl struct {
	sync.Mutex
	llmService *llama.Llama
	lowestPort int              // the lowest port to use
	usedPorts  map[int]struct{} // ports that are currently in use
}

// StartInstance implements [InstanceFactory].
func (i *instanceFactoryImpl) StartInstance(model string, nodes []Node) (Instance, error) {
	// build list of rpc nodes
	rpcNodes := make([]llama.RpcNode, len(nodes))
	for idx, node := range nodes {
		rpcNodes[idx] = llama.RpcNode{
			Ip:   node.Ip(),
			Port: node.Port(),
		}
	}

	// find the lowest port that is not used
	port := i.lowestPort
	for {
		if _, ok := i.usedPorts[port]; !ok {
			break
		}
		port++
	}

	// start the llama server
	cmd, err := func() (*exec.Cmd, error) {
		// guard critical section
		i.Lock()
		defer i.Unlock()

		cmd := i.llmService.ServeCommand(context.Background(), llama.ServeArgs{
			Model:    model,
			RpcNodes: rpcNodes,
			Port:     port,
		})

		err := cmd.Start()
		if err != nil {
			return nil, err
		}

		// mark port as used
		i.usedPorts[port] = struct{}{}

		return cmd, err
	}()

	if err != nil {
		return nil, err
	}

	dead := make(chan struct{})

	// wait for instance to die, then free port
	go func() {
		cmd.Wait()
		close(dead)

		i.Lock()
		delete(i.usedPorts, port)
		i.Unlock()
	}()

	// return the new instance
	return &instanceImpl{
		process: cmd.Process,
		port:    port,
		dead:    dead,
		model:   model,
	}, nil
}

var _ InstanceFactory = (*instanceFactoryImpl)(nil)
