package locale

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// localizers holds the list of all i18n localizers that can be fetched at each request
var localizers map[string]*i18n.Localizer = make(map[string]*i18n.Localizer)

// InitTranslations loads the translation files from the locales directory
func InitTranslations() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	localesDir := fmt.Sprintf("./locales")
	files, _ := ioutil.ReadDir(localesDir)

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			if !strings.HasPrefix(f.Name(), "active.") {
				continue
			}
			if f.Name() == "active.en.json" {
				continue
			}
			bundle.MustLoadMessageFile(fmt.Sprintf("./") + filepath.Join(localesDir, f.Name()))
		}
	}

	localizers[language.English.String()] = i18n.NewLocalizer(bundle, language.English.String())
	localizers[language.Serbian.String()] = i18n.NewLocalizer(bundle, language.Serbian.String())
}

// GetLocalizer gets the localizer by the specified language key
func GetLocalizer(lang string) (*i18n.Localizer, error) {
	val, ok := localizers[lang]
	if !ok {
		return nil, errors.New("could not get localizer from the map")
	}
	return val, nil
}

// GetSupportedLocales shows the list of supported locales
func GetSupportedLocales() map[string]*i18n.Localizer {
	return localizers
}

// GetUserLocalizer gets localizer based on user's preference lang
func GetUserLocalizer(locale string) *i18n.Localizer {
	if _, ok := localizers[locale]; !ok {
		locale = language.English.String()
	}
	return localizers[locale]
}

// LocalizeDefaultMessage localizes the provided message
func LocalizeDefaultMessage(l *i18n.Localizer, m *i18n.Message) string {
	s, err := l.LocalizeMessage(m)
	if err != nil {
		zlog.Warn("could not localize message", zlog.String("messageID", m.ID), zlog.Err(err))
		return ""
	}
	return s
}

// LocalizeWithConfig localizes with the provided config
func LocalizeWithConfig(l *i18n.Localizer, lc *i18n.LocalizeConfig) string {
	s, err := l.Localize(lc)
	if err != nil {
		zlog.Warn("could not localize message with config", zlog.Err(err))
		return ""
	}
	return s
}
