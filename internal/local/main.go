package local

import (
	"embed"
	"io/fs"

	"github.com/gin-gonic/gin"
	goI18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed i18n/*
var translationFS embed.FS

type localizer struct {
	L *goI18n.Localizer
}

var bundle *goI18n.Bundle

// init initializes language configuration.
func init() {
	bundle = goI18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load embedded message files
	enFile, _ := fs.ReadFile(translationFS, "i18n/active.en.toml")
	bundle.ParseMessageFileBytes(enFile, "active.en.toml")

	zhFile, _ := fs.ReadFile(translationFS, "i18n/active.zh-CN.toml")
	bundle.ParseMessageFileBytes(zhFile, "active.zh-CN.toml")
}

// New returns localizer instance pointer.
func New(ctx *gin.Context) *localizer {
	accept := ctx.GetHeader("Accept-Language")
	var instance = localizer{
		L: goI18n.NewLocalizer(bundle, accept),
	}
	return &instance
}

// Message returns a localized message.
func (l *localizer) Message(id string) string {
	msg, _ := l.L.Localize(&goI18n.LocalizeConfig{
		MessageID: id,
	})

	return msg
}
