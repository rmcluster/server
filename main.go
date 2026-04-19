package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-systemd/v22/activation"

	"github.com/wk-y/rama-swap/llama"
	"github.com/wk-y/rama-swap/microservices/dashboard"
	"github.com/wk-y/rama-swap/microservices/homepage"
	"github.com/wk-y/rama-swap/microservices/scheduling"
	"github.com/wk-y/rama-swap/server"
	schedulersubscriber "github.com/wk-y/rama-swap/server/scheduler_subscriber"
	"github.com/wk-y/rama-swap/tracker"
)

const EX_USAGE = 64

func main() {
	mux := http.NewServeMux()

	args, rest, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(EX_USAGE)
	}

	if len(rest) > 0 {
		fmt.Fprintf(os.Stderr, "%s: unexpected positional argument %v\n", os.Args[0], rest[0])
		os.Exit(EX_USAGE)
	}

	setDefaults(&args)

	ramalama := llama.Llama{
		Command: args.Ramalama,
	}
	tracker := tracker.NewTracker()
	tracker.AddRoutes(mux)
	scheduler := scheduling.NewPartitioningScheduler(scheduling.NewInstanceFactory(&ramalama, 49170), 3)
	tracker.Subscribe(schedulersubscriber.NewSchedulerSubscriber(scheduler))
	server := server.NewServer(ramalama, scheduler)
	dashboard := dashboard.NewDashboard(tracker)
	dashboard.RegisterHandlers(mux)
	homepage := homepage.NewHomepage()
	homepage.RegisterHandlers(mux)

	server.ModelNameMangler = func(s string) string {
		return strings.ReplaceAll(s, "/", "_")
	}

	// serve on the configured host/port
	log.Printf("Listening on http://%s:%d\n", *args.Host, *args.Port)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *args.Host, *args.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	server.HandleHttp(mux)
	err = http.Serve(l, mux)

	log.Fatalf("Failed to serve: %v", err)
}
