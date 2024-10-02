package alias

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

func Generate(r contracts.Resource, s string) string {
	alias := Create(s)
	aliasNew := alias
	counter := 1
	repo := Repo()

	for {
		err := repo.GetByAlias(r.GetResource(), aliasNew, r.GetID())
		if err == gorm.ErrRecordNotFound {
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

func Create(title string) string {
	title = strings.ToLower(title)
	alias := transliterate(title)
	alias = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(alias, "-")
	alias = strings.Trim(alias, "-")
	return alias
}

func transliterate(input string) string {
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
