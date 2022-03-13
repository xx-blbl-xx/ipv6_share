package main

import (
	"fmt"
	xlog "ipv6_share/log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type peer struct {
	Ip   net.IP `json:"ip"`
	Port int64  `json:"port"`
}

var (
	Peers = make(map[string]map[string]*peer, 50)
	mutex sync.Mutex
)

func checkPeerAlive() {
	for {
		for fileHash, m := range Peers {
			if app.shareFiles[fileHash] == nil {
				mutex.Lock()
				delete(Peers, fileHash)
				mutex.Unlock()
				continue
			}

			for s, p := range m {
				if !p.Alive() {
					mutex.Lock()
					delete(Peers[fileHash], s)
					mutex.Unlock()
					xlog.Warn("not alive", p)
				}
			}
		}
		time.Sleep(time.Minute)
	}
}

func newPeer(peerInfo string) *peer {
	if peerInfo == "" {
		return nil
	}
	r, err := regexp.Compile(`\[(.*)]:(\d+)`)
	if err != nil {
		xlog.Error("Compile err", err)
		return nil
	}
	rr := r.FindStringSubmatch(peerInfo)
	if len(rr) != 3 {
		return nil
	}

	ip := net.ParseIP(rr[1])
	if ip == nil || ip.To4() != nil || !ip.IsGlobalUnicast() {
		return nil
	}
	port, err := strconv.ParseInt(rr[2], 10, 64)
	if err != nil {
		xlog.Error("ParseInt err", err)
		return nil
	}
	if port >= 65535 || port <= 0 {
		return nil
	}
	return &peer{
		Ip:   ip,
		Port: port,
	}
}

func addPeer(peerInfo, fileHash string) {
	p := newPeer(peerInfo)
	if p == nil {
		return
	}
	mutex.Lock()
	if _, ok := Peers[fileHash]; !ok {
		Peers[fileHash] = make(map[string]*peer, 100)
	}
	Peers[fileHash][peerInfo] = p
	mutex.Unlock()
}

func (p *peer) Alive() bool {
	rsp, err := http.Get(fmt.Sprintf("[%s]:%d/ping", p.Ip, p.Port))
	if err != nil {
		xlog.Error("checkPing err", p, err)
		return false
	}
	if rsp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
