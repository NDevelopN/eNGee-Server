package consequences

import (
	"Engee-Server/utils"
	"log"
	"testing"

	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/user"
)

func prepareConGame(t *testing.T, testName string) (string, string) {
	db.InitDB()

	gid, err := g.CreateGame(defGame)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (creating game): %v`, testName, err)
	}

	plr := defPlr
	plr.GID = gid

	lid, _ := u.CreateUser(defPlr)
	plr.UID = lid
	err = u.UpdateUser(plr)
	if err != nil {
		t.Fatalf(`%v = failed to prepare conGame (updating leader): %v`, testName, err)
	}

	for i := 0; i < 3; i++ {
		uid, _ := u.CreateUser(defPlr)
		plr.UID = uid
		err = u.UpdateUser(plr)
		if err != nil {
			t.Fatalf(`%v = failed to prepare conGame (updating user): %v`, testName, err)
		}
	}

	return gid, lid
}

func prepareInitTest(t *testing.T, testName string) utils.GameMsg {
	gid, lid := prepareConGame(t, testName)

	return utils.GameMsg{
		Type:    "Init",
		UID:     lid,
		GID:     gid,
		Content: string(ts),
	}
}

func TestInit(t *testing.T) {

	gMsg := prepareInitTest(t, "TestInit(Valid)")
	log.Printf("[DEBUG] gMsg: %v", gMsg)

}

func TestPause(t *testing.T) {

}

func TestUnpause(t *testing.T) {

}

func TestEnd(t *testing.T) {

}

func TestReset(t *testing.T) {

}

func TestRules(t *testing.T) {

}

func TestRemove(t *testing.T) {

}

func TestStatus(t *testing.T) {

}

func TestLeave(t *testing.T) {

}
