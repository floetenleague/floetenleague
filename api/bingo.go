package api

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/floetenleague/floetenleague/api/apigen"
	"github.com/floetenleague/floetenleague/database/dbgen"
	"github.com/labstack/echo/v4"
)

func (a *api) JoinBingo(ctx echo.Context, bingoID int64) error {
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	user, err := a.getUser(ctx, db)
	if err != nil {
		return err
	}

	err = db.JoinBingo(ctx.Request().Context(), dbgen.JoinBingoParams{
		UserID:  user.ID,
		BingoID: bingoID,
	})
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (a *api) SetFieldStatus(ctx echo.Context, bingoId int64, fieldId int64, userId int64, status apigen.FLBingoFieldStatus) error {
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	user, err := a.getUser(ctx, db)
	if err != nil {
		return err
	}
	forCurrentUser := userId == user.ID
	isMod := user.Permission == dbgen.UserPermissionModerator

	if !isMod && !forCurrentUser {
		return echo.NewHTTPError(http.StatusForbidden, "nox")
	}

	err = db.BeginTxFunc(ctx.Request().Context(), func(tx *dbgen.Queries) error {
		field, _ := tx.GetBingoUserField(ctx.Request().Context(), dbgen.GetBingoUserFieldParams{
			UserID:       userId,
			BingoID:      bingoId,
			BingoFieldID: fieldId,
		})

		set := dbgen.SetBingoUserFieldStatusParams{
			UserID:       userId,
			BingoID:      bingoId,
			BingoFieldID: fieldId,
			DoneAt:       field.DoneAt,
			ConfirmedAt:  field.ConfirmedAt,
		}
		now := time.Now()
		switch status {
		case apigen.Blank:
			set.DoneAt = nil
			set.ConfirmedAt = nil
		case apigen.DoneInReview:
			set.DoneAt = &now
			set.ConfirmedAt = nil
		case apigen.Done:
			if user.Permission != dbgen.UserPermissionModerator {
				return echo.NewHTTPError(http.StatusBadRequest, "no")
			}
			if set.DoneAt == nil {
				set.DoneAt = &now
			}
			set.ConfirmedAt = &now
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "no")
		}

		err = tx.SetBingoUserFieldStatus(ctx.Request().Context(), set)
		if err != nil {
			return fmt.Errorf("could not update field status")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (a *api) GetReviews(ctx echo.Context) error {
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	user, err := a.getUser(ctx, db)
	if err != nil {
		return err
	}
	if user.Permission != dbgen.UserPermissionModerator {
		return echo.NewHTTPError(http.StatusForbidden, "no")
	}
	fields, err := db.GetUnconfirmedUserFields(ctx.Request().Context())
	if err != nil {
		return err
	}
	reviews := a.conv.ConvertReviews(fields)

	return ctx.JSON(http.StatusOK, reviews)
}
func (a *api) GetOverview(ctx echo.Context) error {
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	dbCtx := ctx.Request().Context()
	bingos, err := db.GetBingos(dbCtx)
	if err != nil {
		return err
	}

	var result apigen.FLOverview
	for _, bingo := range bingos {

		fields, err := db.GetBingoFields(dbCtx, bingo.ID)
		if err != nil {
			return err
		}

		users, err := db.GetBingoUsers(dbCtx, bingo.ID)
		if err != nil {
			return err
		}
		rBingo := apigen.FLBingo{}
		rBingo.Boards = []apigen.FLBingoUserBoard{}

		for _, user := range users {
			userFields, err := db.GetBingoUserFields(dbCtx, dbgen.GetBingoUserFieldsParams{
				UserID:  user.UserID,
				BingoID: bingo.ID,
			})
			if err != nil {
				return err
			}
			status := map[string]apigen.FLBingoUserBoardField{}

			for _, f := range userFields {
				s := apigen.FLBingoUserBoardField{
					Status: apigen.Blank,
				}
				switch {
				case f.ConfirmedAt != nil:
					s.Status = apigen.Done
					s.At = *f.DoneAt
				case f.DoneAt != nil:
					s.Status = apigen.DoneInReview
					s.At = *f.DoneAt
				}
				status[fmt.Sprint(f.BingoFieldID)] = s
			}

			rBingo.Boards = append(rBingo.Boards, apigen.FLBingoUserBoard{
				Username: user.Username,
				UserId:   user.UserID,
				Fields:   status,
				Id:       user.UserID,
			})
		}

		rBingo.Id = bingo.ID
		rBingo.Name = bingo.Name
		rBingo.Fields = a.conv.ConvertFields(fields)
		rBingo.Size = int(bingo.Size)
		result.Bingos = append(result.Bingos, rBingo)
	}

	fillAndSort(&result)

	return ctx.JSON(http.StatusOK, result)
}

func fillAndSort(overview *apigen.FLOverview) {
	for i, bingo := range overview.Bingos {
		for j, user := range bingo.Boards {
			user.Score = calcScore(bingo.Fields, user.Fields)
			max := findBingos(bingo.Size, bingo.Fields, user.Fields)
			for _, i := range max.Bingos {
				f := bingo.Fields[i]
				user.Fields[fmt.Sprint(f.Id)] = apigen.FLBingoUserBoardField{
					Status: apigen.Bingo,
					At:     user.Fields[fmt.Sprint(f.Id)].At,
				}
			}
			user.Score += int64(len(max.Bingos)/bingo.Size) * 5
			user.LastAt = max.Time
			bingo.Boards[j] = user
		}
		sort.Slice(bingo.Boards, func(i, j int) bool {
			left := bingo.Boards[i]
			right := bingo.Boards[j]
			if left.Score > right.Score {
				return true
			}
			if left.Score == right.Score {
				return left.LastAt.Before(right.LastAt)
			}
			return false
		})
		overview.Bingos[i] = bingo
	}
}

func calcScore(fields []apigen.FLBingoField, user map[string]apigen.FLBingoUserBoardField) int64 {
	score := int64(0)
	for _, field := range fields {
		if userField, ok := user[fmt.Sprint(field.Id)]; ok && (userField.Status == apigen.Done) {
			score += field.Score
		}
	}
	return score
}

func findBingos(size int, fields []apigen.FLBingoField, user map[string]apigen.FLBingoUserBoardField) *Maxer {
	maxer := &Maxer{Size: size}

	test := func(cell int) {
		if field, ok := user[fmt.Sprint(fields[cell].Id)]; ok && (field.Status == apigen.Done) {
			maxer.Record(cell, field.At)
		} else {
			maxer.Done()
		}
	}

	for i := 0; i < size; i++ {
		// check column
		for j := 0; j < size; j++ {
			cell := i + (j * size)
			test(cell)
		}
		maxer.Done()

		// check row
		for j := 0; j < size; j++ {
			cell := (i * size) + j
			test(cell)
		}

		maxer.Done()
	}

	// check diagonal, left -> right
	for i := 0; i < size; i++ {
		cell := i*size + i
		test(cell)
	}
	maxer.Done()

	// check diagonal, right -> left
	for i := 0; i < size; i++ {
		cell := size*i + size - i - 1
		test(cell)
	}
	maxer.Done()

	return maxer
}

type Maxer struct {
	Size    int
	current []int
	Bingos  []int
	Time    time.Time
}

type Match struct {
	Cells []int
	Time  time.Time
}

func (m *Maxer) Done() {
	if len(m.current) == m.Size {
		m.Bingos = append(m.Bingos, m.current...)
	}
	m.current = []int{}
}

func (m *Maxer) Record(cell int, t time.Time) {
	if t.After(m.Time) {
		m.Time = t
	}
	m.current = append(m.current, cell)
}
