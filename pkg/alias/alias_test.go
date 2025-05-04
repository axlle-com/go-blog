package alias

import (
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Простой пример с кириллицей",
			args: args{title: "Привет, мир!"},
			want: "privet-mir",
		},
		{
			name: "Слова с ё",
			args: args{title: "Ёлка"},
			want: "yolka",
		},
		{
			name: "Пример с цифрами",
			args: args{title: "Москва 2023"},
			want: "moskva-2023",
		},
		{
			name: "Сложное слово с щ",
			args: args{title: "Щука на озере"},
			want: "shchuka-na-ozere",
		},
		{
			name: "Текст на латинице",
			args: args{title: "Hello World"},
			want: "hello-world",
		},
		{
			name: "Текст с дефисами и пробелами",
			args: args{title: "Это - тест!"},
			want: "eto-test",
		},
		{
			name: "Пустая строка",
			args: args{title: ""},
			want: "",
		},
		{
			name: "Только цифры",
			args: args{title: "123"},
			want: "123",
		},
		{
			name: "Символы",
			args: args{title: "-123~+>)]@>_:^^\\$@–_{@&\\).[\\:[}@_aliaS--"},
			want: "123-alias",
		},
	}
	newDB, _ := db.SetupDB(config.Config())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAliasProvider(NewAliasRepo(newDB)).Create(tt.args.title); got != tt.want {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransliterate(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Простой пример с кириллицей",
			args: args{input: "привет"},
			want: "privet",
		},
		{
			name: "Слово с ё",
			args: args{input: "ёлка"},
			want: "yolka",
		},
		{
			name: "Слово с мягким знаком",
			args: args{input: "Мягкий"},
			want: "Мyagkiy",
		},
		{
			name: "Слово с твердым знаком",
			args: args{input: "объект"},
			want: "obekt",
		},
		{
			name: "Слово с щ",
			args: args{input: "щука"},
			want: "shchuka",
		},
		{
			name: "Текст на латинице",
			args: args{input: "test"},
			want: "test",
		},
		{
			name: "Пустая строка",
			args: args{input: ""},
			want: "",
		},
		{
			name: "Только цифры",
			args: args{input: "123"},
			want: "123",
		},
		{
			name: "Смешанный текст",
			args: args{input: "Москва2023"},
			want: "Мoskva2023",
		},
		{
			name: "Символы",
			args: args{input: "~+>)]@>_:^^\\$@–_{@&\\).[\\:[}@_"},
			want: "~+>)]@>_:^^\\$@–_{@&\\).[\\:[}@_",
		},
	}
	newDB, _ := db.SetupDB(config.Config())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAliasProvider(NewAliasRepo(newDB)).transliterate(tt.args.input); got != tt.want {
				t.Errorf("transliterate() = %v, want %v", got, tt.want)
			}
		})
	}
}
