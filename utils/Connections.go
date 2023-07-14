package utils

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var connPools = map[string]map[string]*websocket.Conn{}

func AddConnectionPool(gid string) error {
	_, k := connPools[gid]
	if k {
		return fmt.Errorf("connection pool already exists: %v", gid)
	}
	connPools[gid] = map[string]*websocket.Conn{}
	return nil
}

func GetConnections(gid string) (map[string]*websocket.Conn, error) {
	pool, k := connPools[gid]
	if !k {
		return nil, fmt.Errorf("no connection pool found for given gid: %v", gid)
	}

	if len(pool) == 0 {
		return nil, fmt.Errorf("no connections found in the pool")
	}

	return pool, nil
}

func AddConnection(gid string, uid string, conn *websocket.Conn) error {
	pool, k := connPools[gid]
	if !k {
		return fmt.Errorf("no connection pool found for given gid: %v", gid)
	}

	_, k = pool[uid]
	if k {
		return fmt.Errorf("connection already exists in pool: %v", uid)
	}

	pool[uid] = conn
	connPools[gid] = pool

	return nil
}

func RemoveConnection(gid string, uid string) error {
	pool, err := GetConnections(gid)
	if err != nil {
		return err
	}

	_, k := pool[uid]
	if !k {
		return fmt.Errorf("no connection found for given uid: %v", uid)
	}

	delete(pool, uid)
	connPools[gid] = pool

	return nil

}
