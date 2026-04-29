package tracker

import (
	"encoding/json"
	"log"
	"math"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// the number of seconds after which an RPC server is removed from the list
const expiryDuration = time.Second * 30

// the number of seconds to wait between announces
const interval = time.Second * 10

type TrackerSubscriber interface {
	OnNodeAdded(node RpcServerInfo)
	OnNodeRemoved(node RpcServerInfo)
}

type Tracker struct {
	sync.RWMutex
	subscribers map[TrackerSubscriber]struct{}

	RpcServers map[string]clientInfo
}

type clientInfo struct {
	RpcServerInfo
	expiryTimer *time.Timer
}

type RpcServerInfo struct {
	Ip            string    `json:"ip"`
	Port          int       `json:"port"`
	LastSeen      time.Time `json:"last_seen"`
	HardwareModel string    `json:"hardware_model"` // the hardware's model name
	MaxSize       int64     `json:"max_size"`
	Battery       float64   `json:"battery"`
	Temperature   float64   `json:"temperature"`
}

func NewTracker() *Tracker {
	return &Tracker{
		RpcServers:  make(map[string]clientInfo),
		subscribers: make(map[TrackerSubscriber]struct{}),
	}
}

func (t *Tracker) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/announce", t.Announce)
	mux.HandleFunc("/servers", t.ListServers)
}

func (t *Tracker) Announce(w http.ResponseWriter, r *http.Request) {
	log.Printf("Announce request from %s: %v", r.Host, r.URL)
	type response struct {
		Interval int `json:"interval"`
	}

	port := r.URL.Query().Get("port")
	if port == "" {
		http.Error(w, "missing port", http.StatusBadRequest)
		return
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		http.Error(w, "invalid port", http.StatusBadRequest)
		return
	}

	ip := r.URL.Query().Get("ip")
	if ip == "" {
		// Fill with the remote IP while handling IPv4 and IPv6 forms.
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = strings.Trim(host, "[]")
		} else {
			ip = strings.Trim(r.RemoteAddr, "[]")
		}
	}

	hardwareModel := r.URL.Query().Get("model")

	var maxSize int64 = -1
	if maxSizeStr := r.URL.Query().Get("max_size"); maxSizeStr != "" {
		maxSize, _ = strconv.ParseInt(maxSizeStr, 10, 64)
	}

	var battery float64 = math.NaN()
	if batteryStr := r.URL.Query().Get("battery"); batteryStr != "" {
		battery, _ = strconv.ParseFloat(batteryStr, 64)
	}

	var temperature float64 = math.NaN()
	if temperatureStr := r.URL.Query().Get("temperature"); temperatureStr != "" {
		temperature, _ = strconv.ParseFloat(temperatureStr, 64)
	}

	// validate ip
	if net.ParseIP(ip) == nil {
		http.Error(w, "invalid ip", http.StatusBadRequest)
		return
	}

	clientId := ip // + ":" + port

	func() {
		t.Lock()
		defer t.Unlock()
		notifyNew := false

		// avoid duplicate timers
		if existingTimer := t.RpcServers[clientId].expiryTimer; existingTimer != nil {
			existingTimer.Stop()
			log.Printf("Reannounce from %s", clientId)
		} else {
			notifyNew = true
			log.Printf("New announce from %s", clientId)
		}

		announceTime := time.Now()

		serverInfo := RpcServerInfo{
			LastSeen:      announceTime,
			Ip:            ip,
			Port:          portNum,
			HardwareModel: hardwareModel,
			MaxSize:       maxSize,
			Battery:       battery,
			Temperature:   temperature,
		}

		t.RpcServers[clientId] = clientInfo{
			RpcServerInfo: serverInfo,
			expiryTimer: time.AfterFunc(expiryDuration, func() {
				t.Lock()
				defer t.Unlock()

				// there's a possible race condition if the client announces just as the timer expires,
				// preventing the timer from being stopped. To prevent that, we verify that the last seen time
				// has not been changed.
				if t.RpcServers[clientId].LastSeen.Equal(announceTime) {
					serverInfo := t.RpcServers[clientId].RpcServerInfo
					delete(t.RpcServers, clientId)
					t.notifyNodeRemoved(serverInfo)
					log.Printf("Removed %s from tracker", clientId)
				}
			}),
		}

		if notifyNew {
			t.notifyNodeAdded(serverInfo)
		}
	}()

	// respond
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response{
		Interval: int(interval.Seconds()),
	})

	if err != nil {
		log.Printf("Failed to respond to announce: %v", err)
	}
}

func (t *Tracker) ListServers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Servers []RpcServerInfo `json:"servers"`
	}

	servers := t.GetServers()

	w.Header().Add("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response{
		Servers: servers,
	})

	if err != nil {
		log.Printf("Failed to respond to list servers: %v", err)
	}
}

func (t *Tracker) GetServers() []RpcServerInfo {
	t.RLock()
	defer t.RUnlock()

	servers := make([]RpcServerInfo, 0, len(t.RpcServers))
	for _, server := range t.RpcServers {
		servers = append(servers, server.RpcServerInfo)
	}

	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Ip < servers[j].Ip {
			return true
		}
		if servers[i].Ip > servers[j].Ip {
			return false
		}
		return servers[i].Port < servers[j].Port
	})

	return servers
}

func (t *Tracker) Subscribe(subscriber TrackerSubscriber) {
	t.Lock()
	defer t.Unlock()
	t.subscribers[subscriber] = struct{}{}
}

func (t *Tracker) Unsubscribe(subscriber TrackerSubscriber) {
	t.Lock()
	defer t.Unlock()
	delete(t.subscribers, subscriber)
}

func (t *Tracker) notifyNodeAdded(node RpcServerInfo) {
	for subscriber := range t.subscribers {
		subscriber.OnNodeAdded(node)
	}
}

func (t *Tracker) notifyNodeRemoved(node RpcServerInfo) {
	for subscriber := range t.subscribers {
		subscriber.OnNodeRemoved(node)
	}
}
