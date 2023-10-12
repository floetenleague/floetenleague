package conv

import (
	"time"

	"github.com/floetenleague/floetenleague/api/apigen"
	"github.com/floetenleague/floetenleague/database/dbgen"
)

// goverter:converter
// goverter:extend TimeToTime IntX
type Converter interface {
	ConvertUsers([]dbgen.GetUsersRow) []apigen.User
	// goverter:matchIgnoreCase
	ConvertUser(dbgen.GetUsersRow) apigen.User

	ConvertReviews([]dbgen.GetUnconfirmedUserFieldsRow) []apigen.FLBingoFieldReview
	ConvertFields([]dbgen.BingoField) []apigen.FLBingoField

	// goverter:map Text Label
	// goverter:map ID Id
	ConvertField(dbgen.BingoField) apigen.FLBingoField
}

func TimeToTime(t time.Time) time.Time {
	return t
}
func IntX(t int32) int64 {
	return int64(t)
}
