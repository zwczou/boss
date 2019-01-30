package opclient

import (
	em "zwczou/gobase/middleware"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	extra.RegisterFuzzyDecoders()
	em.SetNamingStrategy(em.LowerCaseWithUnderscores)
	em.RegisterTimeAsFormatCodec("2006-01-02 15:04:05")
}
