package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

type Benchmark struct {
	Connections int
	Block       int
	ServerAddr  *url.URL
	Dialer      *websocket.Dialer
	CreatedAt   time.Time
}

func (client *Benchmark) Bench() (err error) {
	taskAdd(func() { addWs(client.ServerAddr.String(), 4) })
	data := genRandBytes(client.Block)
	for {
		if connPool.Count() < 10 {
			c := &muxConn{
				id: genRandBytes(3),
				ws: wsPool.getWs(),
			}
			_, err = c.bench(data)
			if err != nil {
				log.Warnf(err.Error())
			}
			_ = c.Close()
		}
	}
}

func benchStats() {
	var ou, od int64
	for {
		time.Sleep(time.Second)
		speedUp := uploaded - ou
		speedDown := downloaded - od
		log.Infof("stats: uploaded %s | %s/s, downloaded %s | %s/s",
			ByteCountSI(uploaded), ByteCountSI(speedUp), ByteCountSI(downloaded), ByteCountSI(speedDown))
		ou = uploaded
		od = downloaded
	}
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
