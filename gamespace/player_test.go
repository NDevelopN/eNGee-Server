package gamespace

import (
	"testing"

	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"
	utils "Engee-Server/utils"
)

func TestStatusChangeReady(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	status := "Ready"

	err := ChangeStatus(pid, gid, status)
	if err != nil {
		t.Fatalf(`TestStatusChange(Ready) = %q, want "nil"`, err)
	}
}

func TestStatusChangeNotReady(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	status := "Not Ready"

	err := ChangeStatus(pid, gid, status)
	if err != nil {
		t.Fatalf(`TestStatusChange(Ready) = %q, want "nil"`, err)
	}
}

func TestStatusChangeInvalid(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	status := "Invalid"

	err := ChangeStatus(pid, gid, status)
	if err == nil {
		t.Fatalf(`TestStatusChange(Ready) = %q, want ERROR`, err)
	}
}

func TestStatusChangeEmpty(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	status := ""

	err := ChangeStatus(pid, gid, status)
	if err == nil {
		t.Fatalf(`TestStatusChange(Ready) = %q, want ERROR`, err)
	}
}

func TestLeaveValid(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	user := utils.DefUser
	user.UID = pid
	user.GID = gid

	u.UpdateUser(user)

	err := Leave(pid, gid)
	if err != nil {
		t.Fatalf(`TestLeave(Valid) = %q, want "nil"`, err)
	}
}

func TestLeaveNotInGame(t *testing.T) {
	db.InitDB()

	pid, _ := u.CreateUser(utils.DefUser)
	gid, _ := g.CreateGame(utils.DefGame)

	err := Leave(pid, gid)
	if err == nil {
		t.Fatalf(`TestLeave(Valid) = %q, want ERROR`, err)
	}
}
