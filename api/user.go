package api

import (
	"net/http"

	"github.com/floetenleague/floetenleague/api/apigen"
	"github.com/floetenleague/floetenleague/database/dbgen"
	"github.com/labstack/echo/v4"
)

func (a *api) GetUsers(ctx echo.Context) error {
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

	users, err := db.GetUsers(ctx.Request().Context())
	if err != nil {
		return err
	}
	resp := a.conv.ConvertUsers(users)
	return ctx.JSON(http.StatusOK, resp)
}

func (a *api) SetUserPermission(ctx echo.Context, userId int64, permission apigen.UserPermission) error {
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

	err = db.SetUserPermission(ctx.Request().Context(), dbgen.SetUserPermissionParams{
		ID:         userId,
		Permission: dbgen.UserPermission(permission),
	})
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
