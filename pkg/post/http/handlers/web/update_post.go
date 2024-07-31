package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller *controller) UpdatePost(c *gin.Context) {
	id := controller.getID(c)
	if id == 0 {
		return
	}

	post, err := repository.NewRepository().GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		c.Abort()
		return
	}

	err = c.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		c.Abort()
		return
	}
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
	c.Set("title", post.Title)
	galleries := service.SaveFromForm(c)
	for _, gallery := range galleries {
		err := service.Attach(post, gallery)
		if err != nil {
			logger.New().Error(err)
		}
	}
}
