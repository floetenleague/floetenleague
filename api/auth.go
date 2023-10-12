package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/floetenleague/floetenleague/api/apigen"
	"github.com/floetenleague/floetenleague/database"
	"github.com/floetenleague/floetenleague/database/dbgen"
	"github.com/floetenleague/floetenleague/hclient"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

func (a *api) GetAuth(ctx echo.Context) error {
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	session, _ := a.store.Get(ctx.Request(), sessionID)
	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return ctx.JSON(http.StatusOK, &apigen.LoginState{
			Username:   "guest",
			Id:         -1,
			LoggedIn:   false,
			Permission: apigen.UserPermissionUnverified,
		})
	}

	user, err := db.GetUserById(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusOK, &apigen.LoginState{
			Username:   "guest",
			Id:         -1,
			LoggedIn:   false,
			Permission: apigen.UserPermissionUnverified,
		})
	}

	if err != nil {
		return ctx.JSON(http.StatusOK, &apigen.LoginState{
			Username:   "guest",
			Id:         -1,
			LoggedIn:   false,
			Permission: apigen.UserPermissionUnverified,
		})
	}
	return ctx.JSON(http.StatusOK, &apigen.LoginState{
		Username:   user.Username,
		Id:         user.ID,
		Permission: apigen.UserPermission(user.Permission),
		LoggedIn:   true,
	})
}
func (a *api) Logout(ctx echo.Context) error {
	session, _ := a.store.Get(ctx.Request(), sessionID)
	session.Options.MaxAge = -1
	if err := a.store.Save(ctx.Request(), ctx.Response().Writer, session); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
func (a *api) LoginPoe(ctx echo.Context) error {
	if a.cfg.POEClientID == "" {
		return ctx.JSON(http.StatusOK, "internal")
	}
	stateKey := randKey()

	session, _ := a.store.Get(ctx.Request(), sessionID)
	session.Values["oauth2State"] = stateKey
	_ = a.store.Save(ctx.Request(), ctx.Response().Writer, session)

	return ctx.JSON(http.StatusOK, a.auth.AuthCodeURL(stateKey))
}
func (a *api) LoginInternal(ctx echo.Context) error {
	var req apigen.LoginInternalJSONBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed body")
	}
	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()
	user, err := db.GetUserByName(ctx.Request().Context(), req.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not find a user with that name"+err.Error())
	}
	if user.Password == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "password login not allowed")
	}
	if bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "wrong password")
	}
	session, _ := a.store.Get(ctx.Request(), sessionID)
	session.Values["user_id"] = user.ID
	err = a.store.Save(ctx.Request(), ctx.Response().Writer, session)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type POEProfile struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (a *api) PoeCallback(ctx echo.Context) error {
	if ctx.QueryParam("error") != "" {
		log.Info().Interface("query", ctx.QueryParams()).Msg("OAuth Callback Error")
		return ctx.Redirect(http.StatusTemporaryRedirect, "/#/auth/denied")
	}

	session, _ := a.store.Get(ctx.Request(), sessionID)
	stateKey, ok := session.Values["oauth2State"].(string)
	delete(session.Values, "oauth2State")
	loginErr := func() error {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/#/auth/error")
	}

	if !ok {
		return loginErr()
	}

	if stateKey != ctx.QueryParam("state") {
		log.Warn().Str("queryKey", ctx.QueryParam("state")).Str("sessionKey", stateKey).Msg("Callback failed")
		return loginErr()
	}
	exchangeCtx := context.WithValue(context.Background(), oauth2.HTTPClient, hclient.Client)
	token, err := a.auth.Exchange(exchangeCtx, ctx.QueryParam("code"))
	if err != nil {
		log.Warn().Str("code", ctx.QueryParam("code")).Err(err).Msg("Oauth Exchange failed")
		return loginErr()
	}

	req, err := http.NewRequest("GET", "https://api.pathofexile.com/profile", nil)
	if err != nil {
		log.Error().Err(err).Msg("Build Login Profile Request")
		return loginErr()
	}
	req.Header.Add(echo.HeaderAuthorization, "Bearer "+token.AccessToken)

	res, err := hclient.Client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Send Profile Request")
		return loginErr()
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Error().Int("code", res.StatusCode).Str("body", string(body)).Msg("Send Profile Request NonOK")
		return loginErr()
	}
	var profile POEProfile
	err = json.NewDecoder(res.Body).Decode(&profile)
	res.Body.Close()
	if err != nil {
		log.Error().Err(err).Msg("Parsing Profile Response")
		return loginErr()
	}

	db, err := a.db.Aquire(ctx.Request().Context())
	if err != nil {
		return err
	}
	defer db.Close()

	user, err := db.InsertPOEUser(ctx.Request().Context(), dbgen.InsertPOEUserParams{
		Username:   profile.Name,
		PoeID:      profile.UUID,
		Permission: dbgen.UserPermissionUnverified,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error inserting user")
		return loginErr()
	}
	session.Values["user_id"] = user.ID
	err = a.store.Save(ctx.Request(), ctx.Response().Writer, session)
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, "/#/auth/success")
}

func (a *api) getUser(ctx echo.Context, db *database.Queries) (dbgen.User, error) {
	session, _ := a.store.Get(ctx.Request(), sessionID)
	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		return dbgen.User{}, echo.NewHTTPError(http.StatusUnauthorized, "no entry here m8")
	}

	user, err := db.GetUserById(ctx.Request().Context(), userID)
	if err != nil {
		return dbgen.User{}, echo.NewHTTPError(http.StatusInternalServerError, "invalid user")
	}

	switch user.Permission {
	case dbgen.UserPermissionBanned, dbgen.UserPermissionUnverified:
		return dbgen.User{}, echo.NewHTTPError(http.StatusForbidden, "unverified or banned")
	}

	return user, nil
}

func randKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
