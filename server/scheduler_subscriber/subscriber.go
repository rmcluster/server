package schedulersubscriber

import (
	"log"
	"strconv"

	"github.com/wk-y/rama-swap/microservices/scheduling"
	"github.com/wk-y/rama-swap/tracker"
)

type SchedulerSubscriber struct {
	scheduler scheduling.Scheduler
}

func NewSchedulerSubscriber(scheduler scheduling.Scheduler) *SchedulerSubscriber {
	return &SchedulerSubscriber{
		scheduler: scheduler,
	}
}

func (s *SchedulerSubscriber) OnNodeAdded(trackerNode tracker.RpcServerInfo) {
	log.Printf("Node added: %v", trackerNode)
	s.scheduler.OnNodeConnect(convertTrackerNode(trackerNode))
}

func (s *SchedulerSubscriber) OnNodeRemoved(trackerNode tracker.RpcServerInfo) {
	log.Printf("Node removed: %v", trackerNode)
	s.scheduler.OnNodeDisconnect(convertTrackerNode(trackerNode))
}

func convertTrackerNode(trackerNode tracker.RpcServerInfo) *node {
	return &node{
		id:      trackerNode.Ip + ":" + strconv.Itoa(trackerNode.Port),
		ip:      trackerNode.Ip,
		port:    trackerNode.Port,
		maxSize: trackerNode.MaxSize,
	}
}

var _ tracker.TrackerSubscriber = (*SchedulerSubscriber)(nil)
