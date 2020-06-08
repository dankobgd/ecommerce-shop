package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgGenerateTokens     = &i18n.Message{ID: "app.generate_tokens.app_error", Other: "could not generate token"}
	msgVerifyToken        = &i18n.Message{ID: "app.verify_token.app_error", Other: "invalid token"}
	msgVerifyTokenMethod  = &i18n.Message{ID: "app.verify_token.app_error", Other: "invalid token signin method"}
	msgExtractTokenMeta   = &i18n.Message{ID: "app.extract_token_meta.app_error", Other: "could not extract token meta data"}
	msgRefreshToken       = &i18n.Message{ID: "app.refresh_token.app_error", Other: "invalid refresh token"}
	msgRefreshTokenMethod = &i18n.Message{ID: "app.refresh_token.app_error", Other: "invalid refresh token signing method"}
	msgDeleteToken        = &i18n.Message{ID: "app.refresh_token.app_error", Other: "could not delete old token"}
)

// IsValidPassword checks if user password is valid
func (a *App) IsValidPassword(password string) *model.AppErr {
	return model.IsValidPasswordCriteria(password, &a.Cfg().PasswordSettings)
}

// CheckUserPassword checks if password matches the hashed version
func (a *App) CheckUserPassword(user *model.User, password string) *model.AppErr {
	if !model.ComparePassword(user.Password, password) {
		return model.NewAppErr("App.ComparePassword", model.ErrConflict, locale.GetUserLocalizer("en"), model.MsgComparePwd, http.StatusBadRequest, nil)
	}
	return nil
}

// IssueTokens returns the token pair
func (a *App) IssueTokens(userID int64) (*model.TokenMetadata, *model.AppErr) {
	settings := &a.Cfg().AuthSettings
	accessID := uuid.New().String()
	accessExpires := time.Now().Add(time.Minute * 15)
	accessClaims := &model.Claims{
		Authorized: true,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: accessExpires},
			ID:        accessID,
			IssuedAt:  &jwt.Time{Time: time.Now()},
			Subject:   strconv.FormatInt(userID, 10),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString([]byte(settings.AccessTokenSecret))
	if err != nil {
		return nil, model.NewAppErr("App.GenerateTokens", model.ErrInternal, locale.GetUserLocalizer("en"), msgGenerateTokens, http.StatusInternalServerError, nil)
	}

	refreshID := uuid.New().String()
	refreshExpires := time.Now().Add(time.Hour * 24 * 7)
	refreshClaims := &model.Claims{
		Authorized: true,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: refreshExpires},
			ID:        refreshID,
			IssuedAt:  &jwt.Time{Time: time.Now()},
			Subject:   strconv.FormatInt(userID, 10),
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString([]byte(settings.RefreshTokenSecret))
	if err != nil {
		return nil, model.NewAppErr("App.GenerateTokens", model.ErrInternal, locale.GetUserLocalizer("en"), msgGenerateTokens, http.StatusInternalServerError, nil)
	}

	meta := &model.TokenMetadata{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		AccessUUID:     accessID,
		RefreshUUID:    refreshID,
		AccessExpires:  accessExpires,
		RefreshExpires: refreshExpires,
		TokenType:      model.AccessTokenType,
	}

	return meta, nil
}

// AttachSessionCookies sets the token inside cookies
func (a *App) AttachSessionCookies(w http.ResponseWriter, meta *model.TokenMetadata) {
	accessCookie := &http.Cookie{
		Name:     model.AccessCookieName,
		Value:    meta.AccessToken,
		Expires:  meta.AccessExpires,
		HttpOnly: false,
		Secure:   false,
	}

	refreshCookie := &http.Cookie{
		Name:     model.RefreshCookieName,
		Value:    meta.RefreshToken,
		Expires:  meta.RefreshExpires,
		HttpOnly: false,
		Secure:   false,
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

// ExtractToken from the request
func ExtractToken(r *http.Request) string {
	c, err := r.Cookie(model.AccessCookieName)

	if err != nil {
		return ""
	}

	return c.Value
}

// VerifyToken checks if token is valid
func VerifyToken(r *http.Request, settings *config.AuthSettings) (*jwt.Token, *model.AppErr) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.NewAppErr("VerifyToken", model.ErrInvalid, locale.GetUserLocalizer("en"), msgVerifyTokenMethod, http.StatusUnauthorized, nil)
		}
		return []byte(settings.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, model.NewAppErr("VerifyToken", model.ErrInvalid, locale.GetUserLocalizer("en"), msgVerifyToken, http.StatusUnauthorized, nil)
	}
	return token, nil
}

// TokenValid returns error if token is not valid
func (a *App) TokenValid(r *http.Request) *model.AppErr {
	token, err := VerifyToken(r, &a.Cfg().AuthSettings)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*model.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata extracts the token meta details
func (a *App) ExtractTokenMetadata(r *http.Request) (*model.AccessData, *model.AppErr) {
	token, err := VerifyToken(r, &a.Cfg().AuthSettings)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["jti"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			return nil, model.NewAppErr("ExtractTokenMetadata", model.ErrInvalid, locale.GetUserLocalizer("en"), msgExtractTokenMeta, http.StatusBadRequest, nil)
		}

		ad := &model.AccessData{
			AccessUUID: accessUUID,
			UserID:     userID,
		}

		return ad, nil
	}
	return nil, err
}

// RefreshToken refreshes the token and returns the token detials
func (a *App) RefreshToken(rt *model.RefreshToken) (*model.TokenMetadata, *model.AppErr) {
	l := locale.GetUserLocalizer("en")
	token, err := jwt.Parse(rt.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshTokenMethod, http.StatusUnauthorized, nil)
		}
		return []byte(a.Cfg().AuthSettings.RefreshTokenSecret), nil
	})

	if err != nil {
		return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshToken, http.StatusUnauthorized, nil)
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshToken, http.StatusUnauthorized, nil)
	}

	claims, ok := token.Claims.(*model.Claims)
	if ok && token.Valid {
		deleted, err := a.DeleteAuth(claims.ID)
		if err != nil || deleted == 0 {
			return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgDeleteToken, http.StatusUnauthorized, nil)
		}

		userID, _ := strconv.ParseInt(claims.Subject, 10, 64)
		meta, err := a.IssueTokens(userID)
		if err != nil {
			return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshToken, http.StatusUnauthorized, nil)
		}

		if err := a.SaveAuth(userID, meta); err != nil {
			return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshToken, http.StatusUnauthorized, nil)
		}

		return meta, nil
	}

	return nil, model.NewAppErr("RefreshToken", model.ErrInvalid, l, msgRefreshToken, http.StatusUnauthorized, nil)
}
