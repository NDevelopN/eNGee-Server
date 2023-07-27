package REST

import (
	c "Engee-Server/TestApp/common"
	"Engee-Server/utils"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestCreateUserValid(t *testing.T) {
	testCases := []utils.User{
		c.User,
	}

	fmt.Print("CreateUser(valid)\n")

	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				rUser, err := c.PostUser(t, user)
				if rUser.UID == "" || err != nil {
					t.Fatalf(`Received: %q, %v, want UID, "nil"`, rUser.UID, err)
				}

				user.UID = rUser.UID

				if rUser != user {
					t.Fatalf(`Received user: %v, expected user: %v`, rUser, user)
				}
			})
	}
}

func TestCreateUserErrors(t *testing.T) {
	testCases := []utils.User{
		{GID: "", UID: uuid.NewString(), Name: "TestLeader", Status: ""},
		{GID: uuid.NewString(), UID: "", Name: "TestLeader", Status: ""},
		{GID: "", UID: "", Name: "", Status: ""},
		{GID: "", UID: "", Name: "TestLeader", Status: "Test"},
	}

	fmt.Print("CreateUser(error)\n")

	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				rUser, err := c.PostUser(t, user)
				if rUser.UID != "" || err == nil {
					t.Fatalf(`Received: %q, %v, want "", ERROR`, rUser.UID, err)
				}
			})
	}

}

func TestGetUserValid(t *testing.T) {
	testCases := []utils.User{
		c.User,
	}

	fmt.Print("GetUser(valid)\n")

	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				rUser, _ := c.PostUser(t, user)

				nUser, err := c.GetUser(t, rUser.UID)
				if nUser != rUser || err != nil {
					t.Fatalf(`Received: %v, %q, want: %v, "nil"`, nUser, err, rUser)
				}
			})
	}
}

func TestUpdateUserValid(t *testing.T) {

	rUser, _ := c.PostUser(t, c.User)
	rGame, _ := c.PostGame(t, c.Game)

	testCases := []utils.User{
		{UID: rUser.UID, GID: rUser.GID, Name: rUser.Name, Status: "NewStatus"},
		{UID: rUser.UID, GID: rUser.GID, Name: "New Name", Status: "NewStatus"},
		{UID: rUser.UID, GID: rGame.GID, Name: "New Name", Status: "NewStatus"},
	}

	fmt.Print("UpdateUser(valid)\n")
	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				want := "Accept"

				reply, err := c.PutUser(t, user)
				if reply.Cause != want || err != nil {
					t.Fatalf(`Received: %q, %v, want: %q, "nil"`, reply.Cause, err, want)
				}
			})
	}
}

func TestDeleteUserValid(t *testing.T) {
	testCases := []utils.User{
		c.User,
	}

	fmt.Print("DeleteUser(valid)\n")
	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				rUser, _ := c.PostUser(t, user)

				want := "Accept"

				reply, err := c.DeleteUser(t, rUser.UID)
				if reply.Cause != want || err != nil {
					t.Fatalf(`Received: %q, %v, want: %q, "nil"`, reply.Cause, err, want)
				}
			})
	}
}

func TestCreateGameValid(t *testing.T) {
	testCases := []utils.Game{
		c.Game,
	}

	fmt.Print("CreateGame(valid)\n")

	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				rGame, err := c.PostGame(t, game)
				if rGame.GID == "" || err != err {
					t.Fatalf(`Received: %q, %v, want GID, "nil"`, rGame.GID, err)
				}

				game.GID = rGame.GID

				if rGame != game {
					t.Fatalf(`Received game: %v, expected game: %v`, rGame, game)
				}
			})
	}
}

func TestGetGamesValid(t *testing.T) {
	testCases := []utils.Game{
		c.Game,
	}

	fmt.Print("GetGame(valid)\n")

	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				rGame, _ := c.PostGame(t, game)

				nGame, err := c.GetGame(t, rGame.GID)
				if nGame != rGame || err != nil {
					t.Fatalf(`Received: %v, %q, want: %v, ""`, nGame, err, rGame)
				}
			})
	}
}

func TestUpdateGameValid(t *testing.T) {
	rGame, _ := c.PostGame(t, c.Game)

	testCases := []utils.Game{
		{
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			GID: rGame.GID, Name: "New Name", Type: rGame.Type, Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: 6, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: 6, CurPlrs: rGame.CurPlrs, AdditionalRules: "{sample: 'NewRule'}",
		},
	}

	fmt.Print("UpdateGame(valid)\n")
	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				want := "Accept"

				reply, err := c.PutGame(t, game)
				if reply.Cause != want || err != nil {
					t.Fatalf(`Received: %q, %v, want: %q, "nil"`, reply.Cause, err, want)
				}
			})
	}
}

func TestDeleteGameValid(t *testing.T) {
	testCases := []utils.Game{
		c.Game,
	}

	fmt.Print("DeleteGame(valid)\n")
	for _, game := range testCases {

		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				rGame, _ := c.PostGame(t, game)

				want := "Accept"

				reply, err := c.DeleteGame(t, rGame.GID)
				if reply.Cause != want || err != nil {
					t.Fatalf(`Received: %q, want: %q`, err, want)
				}
			})
	}
}
