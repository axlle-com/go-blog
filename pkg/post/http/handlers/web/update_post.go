package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	postRepo := repository.NewRepository()
	post, err := postRepo.GetByID(id)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "admin.404", gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		c.Abort()
		return
	}

	err = c.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		return
	}
	//form := c.Request.MultipartForm

	// Извлечение данных первого уровня
	post.PostCategoryID = db.IDStrPtr(c.PostForm("post_category_id"))
	post.TemplateID = db.IDStrPtr(c.PostForm("template_id"))
	post.Alias = c.PostForm("alias")
	post.URL = c.PostForm("url")
	post.Title = c.PostForm("title")
	post.TitleShort = db.StrPtr(c.PostForm("title_short"))
	post.DescriptionPreview = db.StrPtr(c.PostForm("description_preview"))
	post.MetaTitle = db.StrPtr(c.PostForm("meta_title"))
	post.MetaDescription = db.StrPtr(c.PostForm("meta_description"))
	post.Description = db.StrPtr(c.PostForm("description"))
	post.ShowImagePost = db.CheckStr(c.PostForm("show_image_post"))
	post.ShowImageCategory = db.CheckStr(c.PostForm("show_image_category"))
	post.HasComments = db.CheckStr(c.PostForm("has_comments"))
	post.DatePub = db.ParseDate(c.PostForm("date_pub"))
	post.IsPublished = db.CheckStr(c.PostForm("is_published"))
	post.ShowDate = db.CheckStr(c.PostForm("show_date"))
	post.DateEnd = db.ParseDate(c.PostForm("date_end"))
	post.IsFavourites = db.CheckStr(c.PostForm("is_favourites"))
	post.Sort = db.IntStr(c.PostForm("sort"))

	service.SaveFromForm(c)
	//
	//galleries := make(map[string]Gallery)
	//re := regexp.MustCompile(`^galleries\[(.+?)\]\[images\]\[(.+?)\]\[(.+)\]$`)
	//for key, values := range form.Value {
	//	if matches := re.FindStringSubmatch(key); matches != nil {
	//		galleryID := matches[1]
	//		imageID := matches[2]
	//		field := matches[3]
	//		// Инициализация галереи, если необходимо
	//		if _, ok := galleries[galleryID]; !ok {
	//			galleries[galleryID] = Gallery{Images: make(map[string]Image)}
	//		}
	//		gallery := galleries[galleryID]
	//
	//		// Инициализация изображения, если необходимо
	//		if _, ok := gallery.Images[imageID]; !ok {
	//			gallery.Images[imageID] = Image{}
	//		}
	//		image := gallery.Images[imageID]
	//
	//		// Заполнение соответствующего поля изображения
	//		switch field {
	//		case "title":
	//			image.Title = values[0]
	//		case "description":
	//			image.Description = values[0]
	//		case "sort":
	//			image.Sort, _ = strconv.Atoi(values[0])
	//		}
	//		gallery.Images[imageID] = image
	//		galleries[galleryID] = gallery
	//	}
	//}
	//
	//for _, headers := range form.File {
	//	for _, header := range headers {
	//		contentDisposition := header.Header.Get("Content-Disposition")
	//		re := regexp.MustCompile(`name="(galleries\[(.+?)\]\[images\]\[(.+?)\]\[file\])"`)
	//		if matches := re.FindStringSubmatch(contentDisposition); matches != nil {
	//			galleryID := matches[2]
	//			imageID := matches[3]
	//			// Сохраняем файл
	//			filePath := fmt.Sprintf("uploads/%s", header.Filename)
	//			if err := c.SaveUploadedFile(header, filePath); err != nil {
	//				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//				return
	//			}
	//
	//			// Инициализация галереи, если необходимо
	//			if _, ok := galleries[galleryID]; !ok {
	//				galleries[galleryID] = Gallery{Images: make(map[string]Image)}
	//			}
	//			gallery := galleries[galleryID]
	//
	//			// Инициализация изображения, если необходимо
	//			if _, ok := gallery.Images[imageID]; !ok {
	//				gallery.Images[imageID] = Image{}
	//			}
	//			image := gallery.Images[imageID]
	//			image.File = filePath
	//
	//			gallery.Images[imageID] = image
	//			galleries[galleryID] = gallery
	//		}
	//	}
	//}
	//log.Println(post)
	//
	//err = postRepo.Update(post)
	//log.Println(err)
	//return
}
