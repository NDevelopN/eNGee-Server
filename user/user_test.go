package user

import (
	"testing"

	db "Engee-Server/database"
	g "Engee-Server/game"
	u "Engee-Server/utils"

	"github.com/google/uuid"
)

func TestCreateUserValid(t *testing.T) {
	db.InitDB()
	msg, err := CreateUser(u.DefUser)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateUser(valid) = %q, %q, want "uuid", "nil"`, msg, err)
	}
}

func TestCreateUserMutli(t *testing.T) {
	db.InitDB()
	_, _ = CreateUser(u.DefUser)
	msg, err := CreateUser(u.DefUser)
	_, pe := uuid.Parse(msg)
	if pe != nil || err != nil {
		t.Fatalf(`CreateUser(multi) = %q, %q, want "uuid", "nil"`, msg, err)
	}
}

func TestCreateUserEmptyName(t *testing.T) {
	db.InitDB()
	user := u.DefUser
	user.Name = ""

	msg, err := CreateUser(user)
	if err == nil {
		t.Fatalf(`CreateUser(Empty name) = %q, %q, want "", ERROR`, msg, err)
	}
}

func TestCreateGameInjection(t *testing.T) {
	db.InitDB()
	//TODO
}

func TestGetUserValid(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	want := u.DefUser
	want.UID = uid

	user, err := GetUser(uid)
	if user != want || err != nil {
		t.Fatalf(`GetUser(valid) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestGetUserMulti(t *testing.T) {
	db.InitDB()

	_, _ = CreateUser(u.DefUser)
	uid, _ := CreateUser(u.DefUser)

	want := u.DefUser
	want.UID = uid

	user, err := GetUser(uid)
	if user != want || err != nil {
		t.Fatalf(`GetUser(multi) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestGetUserInvalidGID(t *testing.T) {
	db.InitDB()

	_, _ = CreateUser(u.DefUser)

	user, err := GetUser(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetUser(InvalidUID) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestGetUserEmptyUID(t *testing.T) {
	db.InitDB()

	_, _ = CreateUser(u.DefUser)

	user, err := GetUser("")
	if err == nil {
		t.Fatalf(`GetUser(EmptyUID) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestGetUserEmptyDB(t *testing.T) {
	db.InitDB()

	user, err := GetUser(uuid.NewString())
	if err == nil {
		t.Fatalf(`GetUser(EmptyDB) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestGetUserInjection(t *testing.T) {
	db.InitDB()
	//TODO
}

func TestUpdateUserChangeName(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = uid
	user.Name = "Name Test"

	err := UpdateUser(user)
	if err != nil {
		t.Fatalf(`UpdateUser(Name) = %q, want "nil"`, err)
	}

	want := user

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(Name) = %q, %q, want %q, "nil"`, user, err, want)
	}

}

func TestUpdateUserChangeStatus(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = uid
	user.Status = "Status Test"

	err := UpdateUser(user)
	if err != nil {
		t.Fatalf(`UpdateUser(Status) = %q, want "nil"`, err)
	}

	want := user

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(Status) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserChangeGID(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)
	gid, _ := g.CreateGame(u.DefGame)

	user := u.DefUser
	user.UID = uid
	user.GID = gid

	err := UpdateUser(user)
	if err != nil {
		t.Fatalf(`UpdateUser(GID) = %q, want "nil"`, err)
	}

	want := user

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(GID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserChangeInvalidGID(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = uid
	user.GID = uuid.NewString()

	err := UpdateUser(user)
	if err == nil {
		t.Fatalf(`UpdateUser(InvalidGID) = %q, want ERROR`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(InvalidGID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserAll(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)
	gid, _ := g.CreateGame(u.DefGame)

	user := u.DefUser
	user.UID = uid
	user.GID = gid
	user.Name = "Name Test"
	user.Status = "Status Test"

	err := UpdateUser(user)
	if err != nil {
		t.Fatalf(`UpdateUser(all) = %q, want "nil"`, err)
	}

	want := user

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(all) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserInvalidUID(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = uuid.NewString()
	user.Name = "Name Test"

	err := UpdateUser(user)
	if err == nil {
		t.Fatalf(`UpdateUser(InvalidUID) = %q, want ERROR`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(InvalidUID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserEmptyUID(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = ""
	user.Name = "Name Test"

	err := UpdateUser(user)
	if err == nil {
		t.Fatalf(`UpdatedUser(EmptyUID) = %q, want ERROR`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(EmptyUID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestUpdateUserEmptyDB(t *testing.T) {
	db.InitDB()

	user := u.DefUser
	user.UID = uuid.NewString()
	user.Name = "Name Test"

	err := UpdateUser(user)
	if err == nil {
		t.Fatalf(`UpdateUser(EmptyDB) = %q, want ERROR`, err)
	}

	user, err = GetUser(user.UID)
	if err == nil {
		t.Fatalf(`UpdateUser(EmptyDB) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestUpdateUserNoChange(t *testing.T) {
	db.InitDB()

	uid, _ := CreateUser(u.DefUser)

	user := u.DefUser
	user.UID = uid

	err := UpdateUser(user)
	if err != nil {
		t.Fatalf(`UpdateUser(NoChange) = %q, want "nil"`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err = GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`UpdateUser(NoChange) = %q, want %q, "nil"`, err, want)
	}
}

func TestDeleteUserValid(t *testing.T) {
	db.InitDB()
	uid, _ := CreateUser(u.DefUser)
	err := DeleteUser(uid)
	if err != nil {
		t.Fatalf(`DeleteUser(Valid) = %q, want "nil"`, err)
	}

	user, err := GetUser(uid)
	if err == nil {
		t.Fatalf(`DeletUser(Valid) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestDeleteUserMulti(t *testing.T) {
	db.InitDB()
	_, _ = CreateUser(u.DefUser)
	uid, _ := CreateUser(u.DefUser)
	err := DeleteUser(uid)
	if err != nil {
		t.Fatalf(`DeleteUser(Multi) = %q, want "nil"`, err)
	}

	user, err := GetUser(uid)
	if err == nil {
		t.Fatalf(`DeleteUser(Multi) = %q, %q, want "nil", ERROR`, user, err)
	}
}

func TestDeleteUserInvalidUID(t *testing.T) {
	db.InitDB()
	uid, _ := CreateUser(u.DefUser)

	err := DeleteUser(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteUser(InvalidGID) = %q, want ERROR`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err := GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`DeletUser(InvalidGID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestDeletUserEmptyGID(t *testing.T) {
	db.InitDB()
	uid, _ := CreateUser(u.DefUser)

	err := DeleteUser("")
	if err == nil {
		t.Fatalf(`DeleteUser(EmptyGID) = %q, want ERROR`, err)
	}

	want := u.DefUser
	want.UID = uid

	user, err := GetUser(uid)
	if want != user || err != nil {
		t.Fatalf(`DeletUser(InvalidGID) = %q, %q, want %q, "nil"`, user, err, want)
	}
}

func TestDeleteUserEmptyDB(t *testing.T) {
	db.InitDB()

	err := DeleteUser(uuid.NewString())
	if err == nil {
		t.Fatalf(`DeleteUser(EmptyDB) = %q, want ERROR`, err)
	}
}

func TestDeleteUserRepeat(t *testing.T) {
	db.InitDB()
	uid, _ := CreateUser(u.DefUser)
	_ = DeleteUser(uid)

	err := DeleteUser(uid)
	if err == nil {
		t.Fatalf(`DeleteUser(Repeat) = %q, want ERROR`, err)
	}
}

func TestDeleteUserInjection(t *testing.T) {
	db.InitDB()
	//TODO
}
