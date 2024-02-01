package utils

import (
	"log"
	"time"
)

const heartbeatPeriod = 3 * time.Second
const heartbeatThreshold = 12 * time.Second

func MonitorHeartbeats(heartbeats *map[string]time.Time, Delete func(uid string) error) {
	for {
		if heartbeats == nil {
			return
		}

		now := time.Now()

		for uid, lastBeat := range *heartbeats {
			if (now.Sub(lastBeat) * time.Second) < heartbeatThreshold {
				err := Delete(uid)
				if err != nil {
					log.Printf("[Error] Failed to delete after heartbeat failure: %v", err)
				}
			}
		}

		time.Sleep(heartbeatPeriod)
	}
}
