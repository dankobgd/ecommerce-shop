package app

import (
	"net/http"
	"strconv"
	"strings"
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
func (a *App) IssueTokens(user *model.User) (*model.TokenMetadata, *model.AppErr) {
	settings := &a.Cfg().AuthSettings
	atID := uuid.New().String()
	atExp := time.Now().Add(time.Minute * 15)
	atClaims := model.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: atExp},
			ID:        atID,
			IssuedAt:  &jwt.Time{Time: time.Now()},
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	at, err := token.SignedString([]byte(settings.AccessTokenSecret))
	if err != nil {
		return nil, model.NewAppErr("App.GenerateTokens", model.ErrInternal, locale.GetUserLocalizer("en"), msgGenerateTokens, http.StatusInternalServerError, nil)
	}

	rtID := uuid.New().String()
	rtExp := time.Now().Add(time.Hour * 24 * 7)
	rtClaims := model.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: rtExp},
			ID:        rtID,
			IssuedAt:  &jwt.Time{Time: time.Now()},
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refresh.SignedString([]byte(settings.RefreshTokenSecret))
	if err != nil {
		return nil, model.NewAppErr("App.GenerateTokens", model.ErrInternal, locale.GetUserLocalizer("en"), msgGenerateTokens, http.StatusInternalServerError, nil)
	}

	meta := &model.TokenMetadata{
		AccessToken:    at,
		RefreshToken:   rt,
		AccessUUID:     atID,
		RefreshUUID:    rtID,
		AccessExpires:  atExp,
		RefreshExpires: rtExp,
		TokenType:      model.AccessTokenType,
	}

	return meta, nil
}

// AttachSessionCookies sets the token inside cookies
func (a *App) AttachSessionCookies(w http.ResponseWriter, meta *model.TokenMetadata) {
	secure, httpOnly := false, false
	if a.IsProd() {
		secure, httpOnly = true, true
		httpOnly = true
	}

	accessCookie := &http.Cookie{
		Name:     model.AccessCookieName,
		Value:    meta.AccessToken,
		Expires:  meta.AccessExpires,
		HttpOnly: httpOnly,
		Secure:   secure,
		Path:     "/",
	}

	refreshCookie := &http.Cookie{
		Name:     model.RefreshCookieName,
		Value:    meta.RefreshToken,
		Expires:  meta.RefreshExpires,
		HttpOnly: httpOnly,
		Secure:   secure,
		Path:     "/",
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

// DeleteSessionCookies deletes the cookies
func (a *App) DeleteSessionCookies(w http.ResponseWriter) {
	expiredAccess := expireCookie(model.AccessCookieName)
	expiredRefresh := expireCookie(model.RefreshCookieName)
	http.SetCookie(w, expiredAccess)
	http.SetCookie(w, expiredRefresh)
}

func expireCookie(cookieName string) *http.Cookie {
	return &http.Cookie{
		Name:    cookieName,
		Value:   "",
		Expires: time.Now().Add(-100 * 24 * time.Hour),
		MaxAge:  -1,
		Path:    "/",
	}
}

// ExtractAuthTokenFromRequest gets the auth token in few different ways
func ExtractAuthTokenFromRequest(r *http.Request) (string, model.AccessTokenLocation) {
	authHeader := r.Header.Get(model.HeaderBearer)

	// extract from cookie
	if cookie, err := r.Cookie(model.AccessCookieName); err == nil {
		return cookie.Value, model.TokenLocationCookie
	}

	// extract from auth headers
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == model.HeaderBearer {
		return authHeader[7:], model.TokenLocationHeader // default bearer
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == model.HeaderAuthorization {
		return authHeader[6:], model.TokenLocationHeader // oauth
	}

	// extract from query string
	if token := r.URL.Query().Get("access_token"); token != "" {
		return token, model.TokenLocationQueryString
	}

	return "", model.TokenLocationNotFound
}

// VerifyToken checks if token is valid
func VerifyToken(r *http.Request, settings *config.AuthSettings) (*jwt.Token, *model.AppErr) {
	tokenString, _ := ExtractAuthTokenFromRequest(r)
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.NewAppErr("VerifyToken", model.ErrUnauthorized, locale.GetUserLocalizer("en"), msgVerifyTokenMethod, http.StatusUnauthorized, nil)
		}
		return []byte(settings.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, model.NewAppErr("VerifyToken", model.ErrUnauthorized, locale.GetUserLocalizer("en"), msgVerifyToken, http.StatusUnauthorized, nil)
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
	claims, ok := token.Claims.(*model.Claims)
	if ok && token.Valid {
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return nil, model.NewAppErr("ExtractTokenMetadata", model.ErrInvalid, locale.GetUserLocalizer("en"), msgExtractTokenMeta, http.StatusBadRequest, nil)
		}

		ad := &model.AccessData{
			AccessUUID: claims.ID,
			UserID:     userID,
		}

		return ad, nil
	}
	return nil, err
}

// RefreshToken refreshes the token and returns the token detials
func (a *App) RefreshToken(rt *model.RefreshToken) (*model.TokenMetadata, *model.AppErr) {
	l := locale.GetUserLocalizer("en")
	token, err := jwt.ParseWithClaims(rt.RefreshToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
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
		udata := &model.User{Role: claims.Role, ID: userID}

		meta, err := a.IssueTokens(udata)
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
