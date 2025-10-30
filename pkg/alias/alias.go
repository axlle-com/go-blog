package alias

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"gorm.io/gorm"
)

type AliasProvider interface {
	Generate(r contract.Publisher, s string) string
	Create(title string) string
	transliterate(input string) string
}

func NewAliasProvider(aliasRepo AliasRepository) AliasProvider {
	return &provider{
		aliasRepo: aliasRepo,
	}
}

type provider struct {
	aliasRepo AliasRepository
}

func (p *provider) Generate(publisher contract.Publisher, aliasOld string) string {
	alias := p.Create(aliasOld)
	aliasNew := alias
	counter := 1

	for {
		err := p.aliasRepo.GetByAlias(publisher.GetID(), publisher.GetTable(), aliasNew)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		} else if err != nil {
			logger.Fatal(err)
			break
		}
		aliasNew = fmt.Sprintf("%s-%d", alias, counter)
		counter++
	}

	return aliasNew
}

func (p *provider) Create(title string) string {
	title = strings.ToLower(title)
	alias := p.transliterate(title)
	alias = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(alias, "-")
	alias = strings.Trim(alias, "-")
	return alias
}

func (p *provider) transliterate(input string) string {
	translitMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
		'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
		'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ы': "y", 'э': "e", 'ю': "yu", 'я': "ya",
		'ь': "", 'ъ': "",
	}

	var result strings.Builder
	for _, char := range input {
		if val, ok := translitMap[char]; ok {
			result.WriteString(val)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}
