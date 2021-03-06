package model

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"

	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/utils/is"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	bcryptCost = 14

	numbers          = "0123456789"
	symbols          = " !\"\\#$%&'()*+,-./:;<=>?@[]^_`|~"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	userEmailMaxLength    = 128
	userPasswordMaxLength = 72
	userUsernameMaxRunes  = 64
	userFirstnameMaxRunes = 64
	userLastnameMaxRunes  = 64
	userUsernameMaxLength = 64
	userUsernameMinLength = 1
	userLocaleMaxLength   = 5
	userDefaultLocale     = "en"
)

var reservedNames = []string{"app", "api", "admin", "signup", "login", "oauth", "error", "help"}
var restrictedUsernames = []string{"app", "api", "admin", "system"}
var validUsernameChars = regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)

// messages for localization
var (
	msgInvalidUser            = &i18n.Message{ID: "model.user.validate.app_error", Other: "invalid user data"}
	msgValidateUserID         = &i18n.Message{ID: "model.user.validate.id.app_error", Other: "uppercase letter required"}
	msgValidateUserCrAt       = &i18n.Message{ID: "model.user.validate.created_at.app_error", Other: "invalid created_at timestamp"}
	msgValidateUserUpAt       = &i18n.Message{ID: "model.user.validate.updated_at.app_error", Other: "invalid updated_at timestamp"}
	msgValidateUsername       = &i18n.Message{ID: "model.user.validate.username.app_error", Other: "invalid username"}
	msgValidateUserEmail      = &i18n.Message{ID: "model.user.validate.email.app_error", Other: "invalid email"}
	msgValidateUserFName      = &i18n.Message{ID: "model.user.validate.first_name.app_error", Other: "invalid first name"}
	msgValidateUserLName      = &i18n.Message{ID: "model.user.validate.last_name.app_error", Other: "invalid last name"}
	msgValidateUserPwd        = &i18n.Message{ID: "model.user.validate.password.app_error", Other: "invalid password"}
	msgValidateUserConfirmPwd = &i18n.Message{ID: "model.user.validate.confirm_password.app_error", Other: "invalid confirm password"}
	msgValidateUserGender     = &i18n.Message{ID: "model.user.validate.gender.app_error", Other: "invalid gender"}
	msgValidateUserLocale     = &i18n.Message{ID: "model.user.validate.locale.app_error", Other: "invalid locale"}
	msgValidatePwdLength      = &i18n.Message{ID: "model.user.validate.password_length.app_error", Other: "invalid password length"}
	msgValidatePwdUpper       = &i18n.Message{ID: "model.user.validate.password_uppercase.app_error", Other: "uppercase letter required"}
	msgValidatePwdLower       = &i18n.Message{ID: "model.user.validate.password_lowercase.app_error", Other: "lowercase letter required"}
	msgValidatePwdNumber      = &i18n.Message{ID: "model.user.validate.password_numbers.app_error", Other: "number required"}
	msgValidatePwdSymbol      = &i18n.Message{ID: "model.user.validate.password_symbols.app_error", Other: "symbol required"}
)

// User represents the shop user model
type User struct {
	TotalRecordsCount
	ID              int64      `json:"id" db:"id" schema:"-"`
	FirstName       string     `json:"first_name" db:"first_name" schema:"first_name"`
	LastName        string     `json:"last_name" db:"last_name" schema:"last_name"`
	Username        string     `json:"username" db:"username" schema:"username"`
	Email           string     `json:"email" db:"email" schema:"email"`
	Password        string     `json:"password,omitempty" db:"password" schema:"password"`
	ConfirmPassword string     `json:"confirm_password,omitempty" schema:"confirm_password"`
	Gender          *string    `json:"gender" db:"gender" schema:"gender"`
	Role            string     `json:"role" db:"role" schema:"role"`
	Locale          string     `json:"locale" db:"locale" schema:"locale"`
	AvatarURL       *string    `json:"avatar_url" db:"avatar_url" schema:"-"`
	AvatarPublicID  *string    `json:"avatar_public_id" db:"avatar_public_id" schema:"-"`
	Active          bool       `json:"active" db:"active" schema:"-"`
	EmailVerified   bool       `json:"email_verified" db:"email_verified" schema:"-"`
	FailedAttempts  int        `json:"failed_attempts,omitempty" db:"failed_attempts" schema:"-"`
	LastLoginAt     time.Time  `json:"last_login_at" db:"last_login_at" schema:"-"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at" schema:"-"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at" schema:"-"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at" schema:"-"`
	rawpw           string
}

// UserLogin is the data for creation
type UserLogin struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UserAddress is the contact address info
type UserAddress struct {
	UserID    int64 `json:"user_id" db:"user_id"`
	AddressID int64 `json:"address_id" db:"address_id"`
}

// UserPatch is the user patch model
type UserPatch struct {
	FirstName string  `json:"first_name,omitempty" schema:"first_name"`
	LastName  string  `json:"last_name,omitempty" schema:"last_name"`
	Username  string  `json:"username,omitempty" schema:"username"`
	Email     string  `json:"email,omitempty" schema:"email"`
	Gender    *string `json:"gender,omitempty" schema:"gender"`
	Locale    string  `json:"locale,omitempty" schema:"locale"`
}

// Patch patches the user fields that are provided
func (u *User) Patch(patch *UserPatch) {
	if patch.FirstName != "" {
		u.FirstName = patch.FirstName
	}
	if patch.LastName != "" {
		u.LastName = patch.LastName
	}
	if patch.Username != "" {
		u.Username = patch.Username
	}
	if patch.Email != "" {
		u.Email = patch.Email
	}
	if patch.Gender != nil {
		u.Gender = patch.Gender
	}
	if patch.Locale != "" {
		u.Locale = patch.Locale
	}
}

// UserPatchFromJSON decodes the input and returns the UserPatch
func UserPatchFromJSON(data io.Reader) (*UserPatch, error) {
	var up *UserPatch
	err := json.NewDecoder(data).Decode(&up)
	return up, err
}

// ToJSON converts user to json string
func (u *User) ToJSON() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// UserFromJSON decodes the input and return the User
func UserFromJSON(data io.Reader) (*User, error) {
	var user *User
	err := json.NewDecoder(data).Decode(&user)
	return user, err
}

// ToJSON converts UserLogin to json string
func (u *UserLogin) ToJSON() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// UserLoginFromJSON decodes the input and return the User
func UserLoginFromJSON(data io.Reader) (*UserLogin, error) {
	var u *UserLogin
	err := json.NewDecoder(data).Decode(&u)
	return u, err
}

// SetAvatarDetails sets the avatar img and public_id
func (u *User) SetAvatarDetails(details *gocloudinary.ResourceDetails) {
	u.AvatarURL = NewString(details.SecureURL)
	u.AvatarPublicID = NewString(details.PublicID)
}

// IsValidEmail checks if email is valid
func IsValidEmail(email string) bool {
	return is.ValidEmail(email)
}

// IsValidUsername validates whether username matches the criteria
func IsValidUsername(username string) bool {
	if len(username) < userUsernameMinLength || len(username) > userUsernameMaxLength {
		return false
	}
	if !validUsernameChars.MatchString(username) {
		return false
	}
	for _, ru := range restrictedUsernames {
		if username == ru {
			return false
		}
	}
	return true
}

// IsValidGender checks if gender is valid
func IsValidGender(gender string) bool {
	if gender == "m" || gender == "f" {
		return true
	}
	return false
}

// IsValidLocale checks if locale is valid
func IsValidLocale(locale string) bool {
	if locale != "" {
		if len(locale) > userLocaleMaxLength {
			return false
		} else if _, err := language.Parse(locale); err != nil {
			return false
		}
	}
	return true
}

// IsValidPasswordCriteria checks if password fulfills the criteria
func IsValidPasswordCriteria(password string, settings *config.PasswordSettings) *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if len(password) < settings.MinLength || len(password) > settings.MaxLength {
		errs.Add(Invalid("password", l, msgValidatePwdLength))
	}
	if settings.Lowercase {
		if !strings.ContainsAny(password, lowercaseLetters) {
			errs.Add(Invalid("password", l, msgValidatePwdLower))
		}
	}
	if settings.Uppercase {
		if !strings.ContainsAny(password, uppercaseLetters) {
			errs.Add(Invalid("password", l, msgValidatePwdUpper))
		}
	}
	if settings.Number {
		if !strings.ContainsAny(password, numbers) {
			errs.Add(Invalid("password", l, msgValidatePwdNumber))
		}
	}
	if settings.Symbol {
		if !strings.ContainsAny(password, symbols) {
			errs.Add(Invalid("password", l, msgValidatePwdSymbol))
		}
	}

	if !errs.IsZero() {
		return NewValidationError("User", msgInvalidUser, "", errs)
	}
	return nil
}

// Validate validates the user and returns an error if it doesn't pass criteria
func (u *User) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if u.ID != 0 {
		errs.Add(Invalid("id", l, msgValidateUserID))
	}
	if u.CreatedAt.IsZero() {
		errs.Add(Invalid("created_at", l, msgValidateUserCrAt))
	}
	if u.UpdatedAt.IsZero() {
		errs.Add(Invalid("updated_at", l, msgValidateUserUpAt))
	}
	if !IsValidUsername(u.Username) {
		errs.Add(Invalid("username", l, msgValidateUsername))
	}
	if len(u.Email) == 0 || len(u.Email) > userEmailMaxLength || !is.ValidEmail(u.Email) {
		errs.Add(Invalid("email", l, msgValidateUserEmail))
	}
	if utf8.RuneCountInString(u.Username) > userUsernameMaxRunes {
		errs.Add(Invalid("username", l, msgValidateUsername))
	}
	if utf8.RuneCountInString(u.FirstName) > userFirstnameMaxRunes {
		errs.Add(Invalid("first_name", l, msgValidateUserFName))
	}
	if utf8.RuneCountInString(u.LastName) > userLastnameMaxRunes {
		errs.Add(Invalid("last_name", l, msgValidateUserLName))
	}
	if len(u.rawpw) == 0 || len(u.rawpw) > userPasswordMaxLength {
		errs.Add(Invalid("password", l, msgValidateUserPwd))
	}
	if len(u.ConfirmPassword) == 0 || len(u.ConfirmPassword) > userPasswordMaxLength || u.ConfirmPassword != u.rawpw {
		errs.Add(Invalid("confirm_password", l, msgValidateUserConfirmPwd))
	}
	if !IsValidLocale(u.Locale) {
		errs.Add(Invalid("locale", l, msgValidateUserLocale))
	}

	if !errs.IsZero() {
		return NewValidationError("User", msgInvalidUser, "", errs)
	}
	return nil
}

// Validate validates the UserLogin and returns an error if it doesn't pass criteria
func (u *UserLogin) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if len(u.Email) == 0 || len(u.Email) > userEmailMaxLength || !is.ValidEmail(u.Email) {
		errs.Add(Invalid("email", l, msgValidateUserEmail))
	}
	if len(u.Password) == 0 || len(u.Password) > userPasswordMaxLength {
		errs.Add(Invalid("password", l, msgValidateUserPwd))
	}

	if !errs.IsZero() {
		return NewValidationError("User", msgInvalidUser, "", errs)
	}
	return nil
}

// Validate validates the user patch and returns an error if it doesn't pass criteria
func (up *UserPatch) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if up.Username != "" {
		if !IsValidUsername(up.Username) {
			errs.Add(Invalid("username", l, msgValidateUsername))
		}
	}
	if up.Email != "" {
		if len(up.Email) == 0 || len(up.Email) > userEmailMaxLength || !is.ValidEmail(up.Email) {
			errs.Add(Invalid("email", l, msgValidateUserEmail))
		}
	}
	if up.Username != "" {
		if utf8.RuneCountInString(up.Username) > userUsernameMaxRunes {
			errs.Add(Invalid("username", l, msgValidateUsername))
		}
	}
	if up.FirstName != "" {
		if utf8.RuneCountInString(up.FirstName) > userFirstnameMaxRunes {
			errs.Add(Invalid("first_name", l, msgValidateUserFName))
		}
	}
	if up.LastName != "" {
		if utf8.RuneCountInString(up.LastName) > userLastnameMaxRunes {
			errs.Add(Invalid("last_name", l, msgValidateUserLName))
		}
	}
	if up.Gender != nil {
		if !IsValidGender(*up.Gender) {
			errs.Add(Invalid("gender", l, msgValidateUserGender))
		}
	}
	if up.Locale != "" {
		if !IsValidLocale(up.Locale) {
			errs.Add(Invalid("locale", l, msgValidateUserLocale))
		}
	}

	if !errs.IsZero() {
		return NewValidationError("User", msgInvalidUser, "", errs)
	}
	return nil
}

// Sanitize removes any private data from the user object
func (u *User) Sanitize(options map[string]bool) {
	u.rawpw = ""
	u.Password = ""
	u.ConfirmPassword = ""
	u.FailedAttempts = 0
	if len(options) != 0 && !options["email"] {
		u.Email = ""
	}
}

// PreSave will set missing defaults and fill CreatedAt and UpdatedAt times
// It will also hash the password and it should be called before saving the user to the db
func (u *User) PreSave(isSeed ...bool) {
	u.Email = NormalizeEmail(u.Email)
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	u.LastLoginAt = u.CreatedAt
	u.Active = true

	if u.Role == "" {
		u.Role = UserRole
	}
	if u.Locale == "" {
		u.Locale = userDefaultLocale
	}

	if len(isSeed) == 0 || len(isSeed) > 0 && isSeed[0] == false {
		if len(u.Password) > 0 {
			u.rawpw = u.Password
			u.Password = HashPassword(u.Password)
		}
	}
}

// PreUpdate should be called before updating the user in the db
func (u *User) PreUpdate() {
	u.Email = NormalizeEmail(u.Email)
	u.UpdatedAt = time.Now()
}

// NormalizeEmail trims space and returns lowercase email
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// HashPassword generates a hash using bcrypt
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// ComparePassword compares the hash
func ComparePassword(hash string, password string) bool {
	if len(password) == 0 || len(hash) == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
