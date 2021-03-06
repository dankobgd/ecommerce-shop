package app

import (
	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/payment"
	"github.com/dankobgd/ecommerce-shop/zlog"
)

// App represents the app struct
type App struct {
	srv             *Server
	cfg             *config.Config
	log             *zlog.Logger
	paymentProvider payment.Provider
}

// Option for the app
type Option func(*App) error

// OptionCreator for the app
type OptionCreator func() []Option

// New creates the new App
func New(options ...Option) *App {
	app := &App{}
	for _, option := range options {
		option(app)
	}
	return app
}

// SiteURL returs the site URL
func (a *App) SiteURL() string {
	url := "http://localhost:3000"

	if a.IsProd() {
		url = "http://www.xxx.xxx"
	}

	return url
}

// IsDev returs true if app is in development mode
func (a *App) IsDev() bool {
	return a.Cfg().ENV == "development"
}

// IsProd returs true if app is in production mode
func (a *App) IsProd() bool {
	return a.Cfg().ENV == "production"
}

// IsTest returs true if app is in test mode
func (a *App) IsTest() bool {
	return a.Cfg().ENV == "test"
}

// Srv retrieves the app server
func (a *App) Srv() *Server {
	return a.srv
}

// SetServer sets the app server
func (a *App) SetServer(srv *Server) {
	a.srv = srv
}

// Cfg retrieves the app config
func (a *App) Cfg() *config.Config {
	return a.cfg
}

// SetConfig sets the app config
func (a *App) SetConfig(cfg *config.Config) {
	a.cfg = cfg
}

// Log retrieves the app logger
func (a *App) Log() *zlog.Logger {
	return a.log
}

// SetLogger sets the app logger
func (a *App) SetLogger(logger *zlog.Logger) {
	a.log = logger
}

// PaymentProvider retrieves the app payment provider service
func (a *App) PaymentProvider() payment.Provider {
	return a.paymentProvider
}

// SetPaymentProvider option for the app
func SetPaymentProvider(provider payment.Provider) Option {
	return func(a *App) error {
		a.paymentProvider = provider
		return nil
	}
}

// SetConfig option for the app
func SetConfig(cfg *config.Config) Option {
	return func(a *App) error {
		a.cfg = cfg
		return nil
	}
}

// SetLogger option for the app
func SetLogger(logger *zlog.Logger) Option {
	return func(a *App) error {
		a.log = logger
		return nil
	}
}

// SetServer option for the app
func SetServer(server *Server) Option {
	return func(a *App) error {
		a.srv = server
		return nil
	}
}
