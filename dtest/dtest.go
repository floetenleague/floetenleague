package dtest

import (
	"context"
	"fmt"

	"github.com/floetenleague/floetenleague/database"
	"github.com/floetenleague/floetenleague/database/dbgen"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var id int

const pass = "hunter2"

func Test(pool *database.DB) {
	q, err := pool.Aquire(context.Background())
	check(err)

	_, err = q.GetUserByName(context.Background(), "mod1")
	if err == nil {
		return
	}

	addUser(q, "mod1", pass, dbgen.UserPermissionModerator)
	addUser(q, "mod2", pass, dbgen.UserPermissionModerator)
	addUser(q, "user1", pass, dbgen.UserPermissionUser)
	addUser(q, "user2", pass, dbgen.UserPermissionUser)
	addUser(q, "user3", pass, dbgen.UserPermissionUnverified)
	addUser(q, "user4", pass, dbgen.UserPermissionUnverified)
	addUser(q, "user5", pass, dbgen.UserPermissionBanned)

	b, err := q.AddBingo(context.Background(), dbgen.AddBingoParams{
		Name: "Cool Board",
		Size: 6,
	})
	check(err)
	for i := 0; i < 36; i++ {
		err := q.AddBingoField(context.Background(), dbgen.AddBingoFieldParams{
			BingoID: b.ID,
			Text:    fmt.Sprintf("Field %d", i),
		})
		check(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
}

func addUser(db *database.Queries, username, pass string, perm dbgen.UserPermission) {
	id++
	u, err := db.InsertPOEUser(context.Background(), dbgen.InsertPOEUserParams{
		Username:   username,
		Permission: perm,
		PoeID:      fmt.Sprint(id),
	})
	check(err)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), 8)
	hashedString := string(hashed)

	check(db.SetUserPass(context.Background(), dbgen.SetUserPassParams{
		ID:       u.ID,
		Password: &hashedString,
	}))
}
