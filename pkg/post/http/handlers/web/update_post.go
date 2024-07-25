package web

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/db"
	. "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Image struct {
	Title       string
	Description string
	Sort        int
	File        string // путь к файлу на сервере
}

type Gallery struct {
	Title  string
	Images map[string]Image
}

type FormData struct {
	PostCategoryID     int
	Render             int
	Alias              string
	URL                string
	Title              string
	TitleShort         string
	PreviewDescription string
	MetaTitle          string
	MetaDescription    string
	Description        string
	Files              []byte
	IsImagePost        bool
	IsImageCategory    bool
	IsComments         bool
	DatePub            string
	IsPublished        bool
	ShowDate           bool
	DateEnd            string
	IsFavourites       bool
	Sort               int
	Galleries          map[string]Gallery
}

func (controller *controller) UpdatePost(c *gin.Context) {
	id := controller.getID(c)
	if id == 0 {
		return
	}
	var formData Post

	err := c.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		return
	}
	form := c.Request.MultipartForm

	// Извлечение данных первого уровня
	//formData.ID = id
	formData.PostCategoryID = db.IDStrPtr(c.PostForm("post_category_id"))
	formData.TemplateID = db.IDStrPtr(c.PostForm("template_id"))
	formData.Alias = c.PostForm("alias")
	formData.URL = c.PostForm("url")
	formData.Title = c.PostForm("title")
	formData.TitleShort = db.StrPtr(c.PostForm("title_short"))
	formData.DescriptionPreview = db.StrPtr(c.PostForm("description_preview"))
	formData.MetaTitle = db.StrPtr(c.PostForm("meta_title"))
	formData.MetaDescription = db.StrPtr(c.PostForm("meta_description"))
	formData.Description = db.StrPtr(c.PostForm("description"))
	formData.ShowImagePost = db.CheckStr(c.PostForm("show_image_post"))
	formData.ShowImageCategory = db.CheckStr(c.PostForm("show_image_category"))
	formData.HasComments = db.CheckStr(c.PostForm("has_comments"))
	formData.DatePub = db.ParseDate(c.PostForm("date_pub"))
	formData.IsPublished = db.CheckStr(c.PostForm("is_published"))
	formData.ShowDate = db.CheckStr(c.PostForm("show_date"))
	formData.DateEnd = db.ParseDate(c.PostForm("date_end"))
	formData.IsFavourites = db.CheckStr(c.PostForm("is_favourites"))
	formData.Sort = db.IntStr(c.PostForm("sort"))

	galleries := make(map[string]Gallery)
	re := regexp.MustCompile(`^galleries\[(.+?)\]\[images\]\[(.+?)\]\[(.+)\]$`)
	for key, values := range form.Value {
		if matches := re.FindStringSubmatch(key); matches != nil {
			galleryID := matches[1]
			imageID := matches[2]
			field := matches[3]
			// Инициализация галереи, если необходимо
			if _, ok := galleries[galleryID]; !ok {
				galleries[galleryID] = Gallery{Images: make(map[string]Image)}
			}
			gallery := galleries[galleryID]

			// Инициализация изображения, если необходимо
			if _, ok := gallery.Images[imageID]; !ok {
				gallery.Images[imageID] = Image{}
			}
			image := gallery.Images[imageID]

			// Заполнение соответствующего поля изображения
			switch field {
			case "title":
				image.Title = values[0]
			case "description":
				image.Description = values[0]
			case "sort":
				image.Sort, _ = strconv.Atoi(values[0])
			}
			gallery.Images[imageID] = image
			galleries[galleryID] = gallery
		}
	}

	for _, headers := range form.File {
		for _, header := range headers {
			contentDisposition := header.Header.Get("Content-Disposition")
			re := regexp.MustCompile(`name="(galleries\[(.+?)\]\[images\]\[(.+?)\]\[file\])"`)
			if matches := re.FindStringSubmatch(contentDisposition); matches != nil {
				galleryID := matches[2]
				imageID := matches[3]
				// Сохраняем файл
				filePath := fmt.Sprintf("uploads/%s", header.Filename)
				if err := c.SaveUploadedFile(header, filePath); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Инициализация галереи, если необходимо
				if _, ok := galleries[galleryID]; !ok {
					galleries[galleryID] = Gallery{Images: make(map[string]Image)}
				}
				gallery := galleries[galleryID]

				// Инициализация изображения, если необходимо
				if _, ok := gallery.Images[imageID]; !ok {
					gallery.Images[imageID] = Image{}
				}
				image := gallery.Images[imageID]
				image.File = filePath

				gallery.Images[imageID] = image
				galleries[galleryID] = gallery
			}
		}
	}
	log.Println(formData)

	postRepo := repository.NewRepository().Update(&formData)
	log.Println(postRepo)
	return
}
