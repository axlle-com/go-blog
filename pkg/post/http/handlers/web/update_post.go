package web

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/gallery/service"
	"github.com/axlle-com/blog/pkg/post/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *controller) UpdatePost(ctx *gin.Context) {
	id := c.getID(ctx)
	if id == 0 {
		return
	}

	post, err := repository.NewPostRepository().GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		ctx.Abort()
		return
	}

	err = ctx.Request.ParseMultipartForm(32 << 20) // Устанавливаем максимальный размер для multipart/form-data
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"title":   "404 Not Found",
			"content": "errors.404.gohtml",
		})
		ctx.Abort()
		return
	}
	// Извлечение данных первого уровня
	post.PostCategoryID = db.IDStrPtr(ctx.PostForm("post_category_id"))
	post.TemplateID = db.IDStrPtr(ctx.PostForm("template_id"))
	post.Alias = ctx.PostForm("alias")
	post.URL = ctx.PostForm("url")
	post.Title = ctx.PostForm("title")
	post.TitleShort = db.StrPtr(ctx.PostForm("title_short"))
	post.DescriptionPreview = db.StrPtr(ctx.PostForm("description_preview"))
	post.MetaTitle = db.StrPtr(ctx.PostForm("meta_title"))
	post.MetaDescription = db.StrPtr(ctx.PostForm("meta_description"))
	post.Description = db.StrPtr(ctx.PostForm("description"))
	post.ShowImagePost = db.CheckStr(ctx.PostForm("show_image_post"))
	post.ShowImageCategory = db.CheckStr(ctx.PostForm("show_image_category"))
	post.HasComments = db.CheckStr(ctx.PostForm("has_comments"))
	post.DatePub = db.ParseDate(ctx.PostForm("date_pub"))
	post.IsPublished = db.CheckStr(ctx.PostForm("is_published"))
	post.ShowDate = db.CheckStr(ctx.PostForm("show_date"))
	post.DateEnd = db.ParseDate(ctx.PostForm("date_end"))
	post.IsFavourites = db.CheckStr(ctx.PostForm("is_favourites"))
	post.Sort = db.IntStr(ctx.PostForm("sort"))
	ctx.Set("title", post.Title)
	galleries := service.SaveFromForm(ctx)
	for _, gallery := range galleries {
		err := service.Attach(post, gallery)
		if err != nil {
			logger.New().Error(err)
		}
	}
}
