package connections

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	db "Engee-Server/database"
	"Engee-Server/utils"

	"github.com/gorilla/websocket"
)

type connPool struct {
	mu sync.Mutex
	v  map[string]*Conn
}

type Conn struct {
	Mu sync.Mutex
	V  *websocket.Conn
}

var poolMap = map[string]*connPool{}

var LOCALTEST = false

func SETLOCALTEST(on bool) {
	LOCALTEST = on
}

func CheckConnection(gid string, uid string) bool {
	if LOCALTEST {
		return true
	}
	pool := poolMap[gid]
	pool.mu.Lock()
	defer pool.mu.Unlock()

	_, found := pool.v[uid]
	return found
}

func AddConnectionPool(gid string) error {
	if LOCALTEST {
		return nil
	}
	_, k := poolMap[gid]
	if k {
		return fmt.Errorf("connection pool already exists: %v", gid)
	}

	pool := new(connPool)
	pool.mu = sync.Mutex{}
	pool.v = map[string]*Conn{}

	poolMap[gid] = pool

	return nil
}

func AddConnection(gid string, uid string, c *websocket.Conn) error {
	if LOCALTEST {
		return nil
	}
	pool, k := poolMap[gid]
	if !k {
		return fmt.Errorf("no connection pool found for given gid: %v", gid)
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	_, k = pool.v[uid]
	if k {
		return fmt.Errorf("connection already exists in pool: %v", uid)
	}

	nuConn := new(Conn)
	nuConn.Mu = sync.Mutex{}
	nuConn.V = c

	pool.v[uid] = nuConn

	return nil
}

func GetConnections(gid string) (map[string]*Conn, error) {
	pool, k := poolMap[gid]
	if !k {
		return nil, fmt.Errorf("no connection pool found for given gid: %v", gid)
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.v) == 0 {
		return nil, fmt.Errorf("no connections found in the pool")
	}

	return pool.v, nil
}

func Broadcast(msg utils.GameMsg) error {
	if LOCALTEST {
		return nil
	}

	pool, k := poolMap[msg.GID]
	if !k {
		log.Printf("Deleting game, no connection pool")
		db.RemoveGame(msg.GID)
		return fmt.Errorf("no connection pool found for given gid: %v", msg.GID)
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.v) == 0 {
		log.Printf("Deleting game, no connection pool")
		db.RemoveGame(msg.GID)
		return fmt.Errorf("no connections found in the pool")
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	for _, conn := range pool.v {
		conn.Mu.Lock()
		conn.V.WriteMessage(websocket.TextMessage, message)
		conn.Mu.Unlock()
	}

	return nil
}

func SingleMessage(msg utils.GameMsg) error {
	if LOCALTEST {
		return nil
	}
	pool, k := poolMap[msg.GID]
	if !k {
		return fmt.Errorf("no connection pool found for given gid: %v", msg.GID)
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.v) == 0 {
		return fmt.Errorf("no connections found in the pool")
	}

	conn, k := pool.v[msg.UID]
	if !k {
		return fmt.Errorf("no connection found for given uid in pool: %v", msg.UID)
	}

	conn.Mu.Lock()
	defer conn.Mu.Unlock()

	message, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	conn.V.WriteMessage(websocket.TextMessage, message)

	return nil
}

func RemoveConnection(gid string, uid string) error {
	if LOCALTEST {
		return nil
	}
	pool, err := GetConnections(gid)
	if err != nil {
		return err
	}

	conn, k := pool[uid]
	if !k {
		return fmt.Errorf("no connection found for given uid: %v", uid)
	}

	conn.V.Close()

	delete(pool, uid)

	return nil
}

func RemoveConnectionPool(gid string) error {
	if LOCALTEST {
		return nil
	}
	pool, k := poolMap[gid]
	if !k {
		return fmt.Errorf("connection pool not found: %v", gid)
	}

	pool.mu.Lock()

	for i := range pool.v {
		err := RemoveConnection(gid, i)
		if err != nil {
			pool.mu.Unlock()
			return fmt.Errorf("could not remove all connections from pool :%v", err)
		}
	}

	pool.mu.Unlock()
	delete(poolMap, gid)

	return nil
}
