package ws

import (
	c "Engee-Server/TestApp/common"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"Engee-Server/utils"

	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const url = "ws://localhost:8090/"

const timeout = 2 * time.Second

const endMsg = "EndMSG"

var conn *websocket.Conn

func connect(game utils.Game, user utils.User, incoming chan []byte, outgoing chan []byte) {

	ctx := context.Background()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dialer := websocket.Dialer{
		Subprotocols: []string{"json"},
	}

	var err error

	conn, _, err = dialer.DialContext(ctx, url+"games/"+user.UID, nil)
	if err != nil {
		incoming <- []byte(err.Error())
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer conn.Close()
		defer close(done)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				incoming <- []byte(err.Error())
				return
			}

			incoming <- msg
		}
	}()

	for {
		select {
		case o := <-outgoing:
			if string(o) == endMsg {
				conn.Close()
				return
			}
			time.Sleep(time.Millisecond * 5)
			if err := conn.WriteMessage(websocket.TextMessage, o); err != nil {
				log.Printf("writing: %v\n", err)
				return
			}
		case <-interrupt:
			log.Println("interrupting)")
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Printf("error closing: %v", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			conn.Close()
			return
		}
	}
}

func readyPlayers(t *testing.T, gid string, plrs []utils.User, count int, outgoing chan []byte, incoming chan []byte) {
	rMsg := utils.GameMsg{
		GID:     gid,
		Type:    "Status",
		Content: "Ready",
	}

	for i := 0; i < count; i++ {
		rMsg.UID = plrs[i].UID
		out, _ := json.Marshal(rMsg)

		outgoing <- out
	}

	//Allow replies and updates to come through
	time.Sleep(timeout)
}

func TestConnect(t *testing.T) {
	user, _ := c.PostUser(t, c.User)
	game := c.Game
	game.Leader = user.UID
	game, _ = c.PostGame(t, game)

	incoming := make(chan []byte)
	outgoing := make(chan []byte)

	go connect(game, user, incoming, outgoing)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case rec := <-incoming:
			want := "Info"

			var msg utils.GameMsg
			err := json.Unmarshal(rec, &msg)
			if msg.Type != want || err != nil {
				t.Fatalf(`TestConnection = %q, %v, want %q, "nil"`, msg, err, want)
			}

			outgoing <- []byte("test message")
			outgoing <- []byte(endMsg)

			return
		case <-timer.C:
			t.Fatalf("TestConect = Timed out")
		}
	}
}

func Setup(t *testing.T, pCount int) (utils.Game, []utils.User, chan []byte, chan []byte) {
	leader, err := c.PostUser(t, c.User)
	if err != nil {
		t.Fatalf("Could not create leader: %v", err)
	}

	game := c.Game
	game.Leader = leader.UID
	game, err = c.PostGame(t, game)
	if err != nil {
		t.Fatalf("Could not create game: %v", err)
	}

	leader.GID = game.GID

	var plrs []utils.User
	plrs = append(plrs, leader)

	for i := 1; i < pCount; i++ {
		user, err := c.PostUser(t, c.User)
		if err != nil {
			t.Fatalf("Could not create user(%d): %v", i, err)
		}

		user.GID = game.GID
		_, err = c.PutUser(t, user)
		if err != nil {
			t.Fatalf("Could not update user(%d): %v", i, err)
		}

		plrs = append(plrs, user)
		go connect(game, user, nil, nil)
	}

	incoming := make(chan []byte)
	outgoing := make(chan []byte)

	go connect(game, leader, incoming, outgoing)

	return game, plrs, incoming, outgoing
}

func Teardown(t *testing.T, game utils.Game, plrs []utils.User, incoming chan []byte, outgoing chan []byte) {

	outgoing <- []byte(endMsg)

	for _, p := range plrs {
		c.DeleteUser(t, p.UID)
	}

	close(incoming)
	close(outgoing)
}

func wantResponse(uid string, gid string, cause string, message string) utils.GameMsg {
	resp := utils.Response{
		Cause:   cause,
		Message: message,
	}

	content, _ := json.Marshal(resp)

	return utils.GameMsg{
		UID:     uid,
		GID:     gid,
		Type:    "Response",
		Content: string(content),
	}
}

type inOut struct {
	out utils.GameMsg
	in  utils.GameMsg
}

type pListOut struct {
	out   utils.GameMsg
	pList []utils.User
}

func TestInvalids(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	randID := uuid.NewString()

	var testCases = []inOut{
		// Invalid Message Type
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Invalid",
			},
			in: wantResponse(
				uid, gid, "Error", "Unsupported message type: Invalid",
			),
		},

		// Invalid GID
		{
			out: utils.GameMsg{
				UID: uid, GID: randID, Type: "Status", Content: "Ready",
			},
			in: wantResponse(
				uid, randID, "Error", "Invalid ID(s) provided",
			),
		},

		// Invalid UID
		{
			out: utils.GameMsg{
				UID: randID, GID: gid, Type: "Status", Content: "Ready",
			},
			in: wantResponse(
				randID, gid, "Error", "Invalid ID(s) provided",
			),
		},

		// Not Leader
		{
			out: utils.GameMsg{
				UID: plrs[1].UID, GID: gid, Type: "Reset",
			},
			in: wantResponse(
				plrs[1].UID, gid, "Error", "Must be a leader to Reset",
			),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),

			func(t *testing.T) {
				want := tc.in
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf("TestInvalid = %v", err)
						}
						if msg == want {
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestInvalid = time out")
					}
				}
			},
		)
	}
}

func TestHandlePause(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 1)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	var testCases = []inOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Pause",
			},
			in: utils.GameMsg{
				GID: gid, Type: "Status", Content: "Pause",
			},
		},
	}

	time.Sleep(1 * time.Second)

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),

			func(t *testing.T) {
				want := tc.in
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf("TestPause = %v", err)
						}
						if msg == want {
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestPause = time out")
					}
				}
			},
		)
	}
}

func TestHandleUnpause(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	var testCases = []inOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Pause",
			},
			in: utils.GameMsg{
				GID: gid, Type: "Status", Content: "Lobby",
			},
		},
	}

	time.Sleep(time.Second)

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),

			func(t *testing.T) {
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Pause:
				for {
					select {
					case in := <-incoming:
						_ = json.Unmarshal(in, &msg)
						if msg.Type == "Status" {
							break Pause
						}
					case <-timer.C:
						t.Fatalf("TestPause = time out")
					}
				}

				timer.Stop()

				want := tc.in

				outgoing <- out

				timer = time.NewTimer(timeout)
				defer timer.Stop()

			Unpause:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf("TestUnpause = %v", err)
						}
						if msg == want {
							break Unpause
						}
					case <-timer.C:
						t.Fatalf("TestUnpause = time out")
					}
				}
			},
		)
	}
}

func createPList(plrs []utils.User, low int, high int) []utils.User {
	pList := []utils.User{}
	for i, p := range plrs {
		if i >= low && i <= high {
			continue
		}
		p.Status = "Not Ready"

		pList = append(pList, p)
	}

	return pList
}

func comparePlrsList(rec []utils.User, want []utils.User) bool {
	if len(rec) != len(want) {
		return false
	}

	for _, r := range rec {
		found := false
		for _, w := range want {
			if r == w {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true

}

func TestHandleStatus(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	plrs[0].Status = "Ready"

	var testCases = []utils.GameMsg{
		{
			UID: uid, GID: gid, Type: "Status", Content: "Ready",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s, Content: %s",
			tc.UID, tc.GID, tc.Type, tc.Content),

			func(t *testing.T) {
				out, _ := json.Marshal(tc)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf(`TestStatus = %v`, err)
						}

						if msg.Type == "Players" {
							var recList []utils.User
							err = json.Unmarshal([]byte(msg.Content), &recList)
							if err != nil {
								t.Fatalf(`TestStatus = %v`, err)
							}

							if !comparePlrsList(recList, plrs) {
								t.Fatalf(`TestStatus = %v, want %v`, recList, plrs)
							}
							break Reply
						}

					case <-timer.C:
						t.Fatalf("TestStatus = time out")
					}
				}
			},
		)
	}
}

func TestHandleLeave(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	gid := game.GID

	uid := plrs[len(plrs)-1].UID

	var testCases = []pListOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Leave",
			},
			pList: plrs[:len(plrs)-1],
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s, Content: %s",
			tc.out.UID, tc.out.GID, tc.out.Type, tc.out.Content),

			func(t *testing.T) {
				out, _ := json.Marshal(tc.out)

				outgoing <- out

				var msg utils.GameMsg

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf(`TestLeave = %v`, err)
						}

						if msg.Type == "Players" {
							var recList []utils.User
							err = json.Unmarshal([]byte(msg.Content), &recList)
							if err != nil {
								t.Fatalf(`TestLeave = %v`, err)
							}

							if !comparePlrsList(recList, tc.pList) {
								t.Fatalf(`TestStatus = %v, want %v`, recList, plrs)
							}
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestLeave = time out")
					}
				}
			},
		)
	}
}

func TestStart(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	readyPlayers(t, gid, plrs, len(plrs), outgoing, incoming)

	var testCases = []inOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Start",
			},
			in: utils.GameMsg{
				GID: gid, Type: "Status", Content: "Play",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),
			func(t *testing.T) {
				want := tc.in
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout * 5)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf(`TestStart = %v`, err)
						}

						if msg == want {
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestStart = time out")
					}
				}
			},
		)
	}
}

func TestHandleEnd(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	var testCases = []inOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "End",
			},
			in: utils.GameMsg{
				GID: gid, Type: "End",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),

			func(t *testing.T) {
				want := tc.in
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf("Test End = %v", err)
						}

						if msg == want {
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestEnd = time out")
					}
				}
			},
		)
	}
}

func TestHandleReset(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	var testCases = []inOut{
		{
			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Reset",
			},
			in: utils.GameMsg{
				GID: gid, Type: "Status", Content: "Lobby",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s",
			tc.out.UID, tc.out.GID, tc.out.Type),

			func(t *testing.T) {
				want := tc.in
				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf("TestReset = %v", err)
						}

						if msg == want {
							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestReset = time out")
					}
				}
			},
		)
	}
}

func TestHandleRemove(t *testing.T) {
	game, plrs, incoming, outgoing := Setup(t, 4)
	defer Teardown(t, game, plrs, incoming, outgoing)

	uid := plrs[0].UID
	gid := game.GID

	var testCases = []pListOut{
		{

			out: utils.GameMsg{
				UID: uid, GID: gid, Type: "Remove", Content: plrs[1].UID,
			},
			pList: createPList(plrs, 1, 1),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("UID: %s, GID: %s, Type: %s, Content: %s",
			tc.out.UID, tc.out.GID, tc.out.Type, tc.out.Content),

			func(t *testing.T) {

				out, _ := json.Marshal(tc.out)

				var msg utils.GameMsg

				outgoing <- out

				timer := time.NewTimer(timeout)
				defer timer.Stop()

			Reply:
				for {
					select {
					case in := <-incoming:
						err := json.Unmarshal(in, &msg)
						if err != nil {
							t.Fatalf(`TestRemove = %v`, err)
						}
						if msg.Type == "Players" {
							var recList []utils.User
							err = json.Unmarshal([]byte(msg.Content), &recList)
							if err != nil {
								t.Fatalf(`TestLeave = %v`, err)
							}

							if !comparePlrsList(recList, tc.pList) {
								t.Fatalf(`TestStatus = %v, want %v`, recList, plrs)
							}

							break Reply
						}
					case <-timer.C:
						t.Fatalf("TestRemove = time out")
					}
				}
			},
		)
	}
}
