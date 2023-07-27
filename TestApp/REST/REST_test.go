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
		{
			//Create User with existing UID
			GID: "", UID: uuid.NewString(), Name: "TestLeader", Status: "",
		},
		{
			//Create user with existing GID
			GID: uuid.NewString(), UID: "", Name: "TestLeader", Status: "",
		},
		{
			//Create user with Name as empty string
			GID: "", UID: "", Name: "", Status: "",
		},
		{
			//Create User with existing Status
			GID: "", UID: "", Name: "TestLeader", Status: "Test",
		},
	}

	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				rUser, err := c.PostUser(t, user)
				if rUser.UID != "" || err == nil {
					t.Fatalf(`Received: %q, %v, want "", ERROR`, rUser.UID, err)
				}
				t.Log(err)
			})
	}

}

func TestGetUserValid(t *testing.T) {
	testCases := []utils.User{
		c.User,
	}

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

func TestGetUserError(t *testing.T) {
	_, _ = c.PostUser(t, c.User)

	testCases := []string{
		//Empty UID
		"",
		//Invalid String
		"User ID",
		//Invalid UID
		uuid.NewString(),
	}

	for _, uid := range testCases {
		t.Run(fmt.Sprintf("UID: %v", uid),
			func(t *testing.T) {
				rUser, err := c.GetUser(t, uid)
				if err == nil || rUser.UID != "" {
					t.Fatalf(`Received: %v, %v, want: "", Error`, rUser.UID, err)

				}
				t.Log(err)
			})
	}
}

func TestUpdateUserValid(t *testing.T) {

	rUser, _ := c.PostUser(t, c.User)
	rGame, _ := c.PostGame(t, c.Game)

	testCases := []utils.User{
		{
			//Change Status
			UID: rUser.UID, GID: rUser.GID, Name: rUser.Name, Status: "NewStatus",
		},
		{
			//Change Name
			UID: rUser.UID, GID: rUser.GID, Name: "New Name", Status: "NewStatus",
		},
		{
			//Change GID
			UID: rUser.UID, GID: rGame.GID, Name: "New Name", Status: "NewStatus",
		},
	}

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

func TestUpdateUserError(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	rGame, _ := c.PostGame(t, c.Game)

	rUser.GID = rGame.GID

	_, _ = c.PutUser(t, rUser)

	testCases := []utils.User{
		{
			//Change Name to empty string
			UID: rUser.UID, GID: rUser.GID, Name: "", Status: rUser.Status,
		},
		{
			//Change GID to invalid string
			UID: rUser.UID, GID: "GID", Name: rUser.Name, Status: rUser.Status,
		},
		{
			//Change GID to invalid GID
			UID: rUser.UID, GID: uuid.NewString(), Name: rUser.Name, Status: rUser.Status,
		},
		{
			//Change UID
			UID: uuid.NewString(), GID: rUser.GID, Name: rUser.Name, Status: rUser.Status,
		},
	}

	for _, user := range testCases {
		t.Run(fmt.Sprintf("GID: %s, UID:  %s, Name: %s, Status: %s",
			user.GID, user.UID, user.Name, user.Status),

			func(t *testing.T) {
				want := ""

				reply, err := c.PutUser(t, user)
				if err == nil || reply.Cause != want {
					t.Fatalf(`Received: %q, %v, want: %q, ERROR`, reply.Cause, err, want)
				}
				t.Log(err)
			})
	}
}

func TestDeleteUserValid(t *testing.T) {
	testCases := []utils.User{
		c.User,
	}

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

func TestDeleteUserError(t *testing.T) {
	c.PostUser(t, c.User)

	testCases := []string{
		//Empty String
		"",
		//Invalid String
		"User ID",
		//Invalid UID
		uuid.NewString(),
	}

	for _, uid := range testCases {
		t.Run(fmt.Sprintf("UID:  %s", uid),
			func(t *testing.T) {
				want := ""

				reply, err := c.DeleteUser(t, uid)
				if reply.Cause != want || err == nil {
					t.Fatalf(`Received: %q, %v, want: %q, ERROR`, reply.Cause, err, want)
				}
				t.Log(err)
			})
	}
}
func TestCreateGameValid(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)

	game := c.Game
	game.Leader = rUser.UID

	testCases := []utils.Game{
		game,
	}

	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				rGame, err := c.PostGame(t, game)
				if rGame.GID == "" || err != nil {
					t.Fatalf(`Received: %q, %v, want GID, "nil"`, rGame.GID, err)
				}

				game.GID = rGame.GID

				if rGame != game {
					t.Fatalf(`Received game: %v, expected game: %v`, rGame, game)
				}
			})
	}
}

func TestCreateGameErrors(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)

	testCases := []utils.Game{
		{
			//Create Game with Name as empty string
			GID: "", Name: "", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with Type as empty string
			GID: "", Name: "TestGame", Type: "", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with Invalid Type
			GID: "", Name: "TestGame", Type: "Invalid", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with Old Status
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "Lobby",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with invalid Leader ID
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: uuid.NewString(), MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with empty Leader
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: "", MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with Invalid MinPlrs (-1)
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: -1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with MinPlrs > MaxPlrs
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 6, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
		{
			//Create Game with CurPlrs != 0
			GID: "", Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 1, AdditionalRules: "",
		},
		{
			//Create Game with existing GID
			GID: uuid.NewString(), Name: "TestGame", Type: "test", Status: "Lobby", OldStatus: "",
			Leader: rUser.UID, MinPlrs: 1, MaxPlrs: 5, CurPlrs: 0, AdditionalRules: "",
		},
	}

	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				rGame, err := c.PostGame(t, game)
				if rGame.GID != "" || err == nil {
					t.Fatalf(`Received: %q, %v, want "", ERROR`, rGame.GID, err)
				}

				game.GID = rGame.GID
				if rGame != game {
					t.Fatalf(`Received game: %v, expected game: %v`, rGame, game)
				}
				t.Log(err)
			})
	}
}

func TestGetGamesValid(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	g := c.Game
	g.Leader = rUser.UID

	testCases := []utils.Game{
		g,
	}

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
					t.Fatalf(`Received: %v, %v, want: %v, "nil"`, nGame, err, rGame)
				}
			})
	}
}

func TestGetGamesErrors(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)

	g := c.Game
	g.Leader = rUser.UID

	c.PostGame(t, g)

	testCases := []string{
		"", "GID", uuid.NewString(),
	}

	for _, gid := range testCases {
		t.Run(fmt.Sprintf("GID: %s", gid),
			func(t *testing.T) {
				want := ""
				nGame, err := c.GetGame(t, gid)
				if nGame.GID != want || err == nil {
					t.Fatalf(`Received: %q, %v, want: %q, ERROR`, nGame, err, want)
				}
				t.Log(err)
			})
	}
}

func TestUpdateGameValid(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	g := c.Game
	g.Leader = rUser.UID

	rGame, _ := c.PostGame(t, g)

	testCases := []utils.Game{
		{
			//Change Status
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Name
			GID: rGame.GID, Name: "New Name", Type: rGame.Type, Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Type
			GID: rGame.GID, Name: "New Name", Type: "consequences", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change MinPlrs
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change MaxPlrs
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: 6, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change AdditionalRules
			GID: rGame.GID, Name: "New Name", Type: "New Type", Status: "New Status",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: 2,
			MaxPlrs: 6, CurPlrs: rGame.CurPlrs, AdditionalRules: "{sample: 'NewRule'}",
		},
	}

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
				t.Log(err)
			})
	}
}

func TestUpdateGameErrors(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	g := c.Game
	g.Leader = rUser.UID

	rGame, _ := c.PostGame(t, g)

	testCases := []utils.Game{
		{
			//Change Name to empty string
			GID: rGame.GID, Name: "", Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Type to invalid type
			GID: rGame.GID, Name: rGame.Name, Type: "Invalid", Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Type to empty string
			GID: rGame.GID, Name: rGame.Name, Type: "", Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Status to empty string
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: "",
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Leader to empty string
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: "", MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change Leader to invalid string
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: "Leader", MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change leader to invalid UUID
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: uuid.NewString(), MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change MinPlrs to invalid value (-1)
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: -1,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change MinPlrs to value larger than MaxPlrs
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MaxPlrs + 1,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: rGame.CurPlrs, AdditionalRules: rGame.AdditionalRules,
		},
		{
			//Change CurPlrs
			GID: rGame.GID, Name: rGame.Name, Type: rGame.Type, Status: rGame.Status,
			OldStatus: rGame.OldStatus, Leader: rGame.Leader, MinPlrs: rGame.MinPlrs,
			MaxPlrs: rGame.MaxPlrs, CurPlrs: 3, AdditionalRules: rGame.AdditionalRules,
		},
	}

	for _, game := range testCases {
		t.Run(fmt.Sprintf("GID: %s, Name:  %s, Type: %s, Status: %s,"+
			" OldStatus: %s, Leader: %s, MinPlrs: %d,"+
			" MaxPlrs: %d, CurPlrs: %d, AdditionalRules: %s",
			game.GID, game.Name, game.Type, game.Status, game.OldStatus, game.Leader,
			game.MinPlrs, game.MaxPlrs, game.CurPlrs, game.AdditionalRules),

			func(t *testing.T) {
				want := ""

				reply, err := c.PutGame(t, game)
				if reply.Cause != want || err == nil {
					t.Fatalf(`Received: %q, %v, want: %q, ERROR`, reply.Cause, err, want)
				}
				t.Log(err)
			})
	}
}

func TestDeleteGameValid(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	g := c.Game
	g.Leader = rUser.UID

	testCases := []utils.Game{
		g,
	}

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
					t.Fatalf(`Received: %q, %v, want: %q, "nil"`, reply.Cause, err, want)
				}
			})
	}
}

func TestDeleteGameErrors(t *testing.T) {
	rUser, _ := c.PostUser(t, c.User)
	g := c.Game
	g.Leader = rUser.UID

	c.PostGame(t, g)

	testCases := []string{
		//Empty GID
		"",
		//Invalid String
		"GID",
		//Invalid GID
		uuid.NewString(),
	}

	for _, gid := range testCases {
		t.Run(fmt.Sprintf("GID: %s", gid),
			func(t *testing.T) {
				want := ""

				reply, err := c.DeleteGame(t, gid)
				if reply.Cause != want || err == nil {
					t.Fatalf(`Received: %q, %v, want: %q, ERROR`, reply.Cause, err, want)
				}
				t.Log(err)
			})
	}
}
