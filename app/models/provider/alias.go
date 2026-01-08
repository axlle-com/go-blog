package provider

import "github.com/axlle-com/blog/app/models/contract"

type AliasProvider interface {
	Generate(publisher contract.Publisher, aliasOld string) string
	Create(title string) string
	Transliterate(input string) string
}
