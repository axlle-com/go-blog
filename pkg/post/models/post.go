package models

import (
	"encoding/json"
	"fmt"
	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/common/logger"
	common "github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/file"
	"net/http"
	"time"
)

type Post struct {
	ID                 uint                   `gorm:"primaryKey" json:"id" form:"id" binding:"-"`
	UserID             *uint                  `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         *uint                  `gorm:"index" json:"template_id" form:"template_id" binding:"omitempty"`
	PostCategoryID     *uint                  `gorm:"index" json:"post_category_id" form:"post_category_id" binding:"omitempty"`
	MetaTitle          *string                `gorm:"size:100" json:"meta_title" form:"meta_title" binding:"omitempty,max=100"`
	MetaDescription    *string                `gorm:"size:200" json:"meta_description" form:"meta_description" binding:"omitempty,max=200"`
	Alias              string                 `gorm:"size:255;unique" json:"alias" form:"alias" binding:"omitempty,max=255"`
	URL                string                 `gorm:"size:1000;unique" json:"url" form:"url" binding:"omitempty,max=1000"`
	IsPublished        bool                   `gorm:"not null;default:false" json:"is_published" form:"is_published" binding:"omitempty"`
	IsFavourites       bool                   `gorm:"not null;default:false" json:"is_favourites" form:"is_favourites" binding:"omitempty"`
	HasComments        bool                   `gorm:"not null;default:false" json:"has_comments" form:"has_comments" binding:"omitempty"`
	ShowImagePost      bool                   `gorm:"not null;default:false" json:"show_image_post" form:"show_image_post"`
	ShowImageCategory  bool                   `gorm:"not null;default:false" json:"show_image_category" form:"show_image_category" binding:"omitempty"`
	InSitemap          bool                   `gorm:"not null;default:false" json:"in_sitemap" form:"in_sitemap" binding:"omitempty"`
	Media              *string                `gorm:"size:255" json:"media" form:"media" binding:"omitempty,max=255"`
	Title              string                 `gorm:"size:255;not null" json:"title" form:"title" binding:"required,max=255"`
	TitleShort         *string                `gorm:"size:155;default:null" json:"title_short" form:"title_short" binding:"omitempty,max=155"`
	DescriptionPreview *string                `gorm:"type:text" json:"description_preview" form:"description_preview" binding:"omitempty"`
	Description        *string                `gorm:"type:text" json:"description" form:"description" binding:"omitempty"`
	ShowDate           bool                   `gorm:"not null;default:false" json:"show_date" form:"show_date" binding:"omitempty"`
	DatePub            *time.Time             `json:"date_pub,date,omitempty" time_format:"02.01.2006" form:"date_pub" binding:"omitempty"`
	DateEnd            *time.Time             `json:"date_end,date,omitempty" time_format:"02.01.2006" form:"date_end" binding:"omitempty"`
	Image              *string                `gorm:"size:255" json:"image" form:"image" binding:"omitempty,max=255"`
	Hits               uint                   `gorm:"not null;default:0" json:"hits" form:"hits" binding:"-"`
	Sort               int                    `gorm:"not null;default:0" json:"sort" form:"sort" binding:"omitempty"`
	Stars              float32                `gorm:"not null;default:0.0" json:"stars" form:"stars" binding:"-"`
	CreatedAt          *time.Time             `json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty" form:"updated_at" binding:"-" ignore:"true"`
	DeletedAt          *time.Time             `gorm:"index" json:"deleted_at" form:"deleted_at" binding:"-" ignore:"true"`
	Galleries          []contracts.Gallery    `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`
	dirty              map[string]interface{} `ignore:"true"`
	original           *Post                  `ignore:"true"`
	*common.Field      `ignore:"true"`
}

func (p *Post) GetCategoryID() uint {
	var categoryID uint
	if p.PostCategoryID != nil {
		categoryID = *p.PostCategoryID
	}
	return categoryID
}

func (p *Post) GetTemplateID() uint {
	var templateID uint
	if p.TemplateID != nil {
		templateID = *p.TemplateID
	}
	return templateID
}

func (p *Post) GetID() uint {
	return p.ID
}

func (p *Post) GetResource() string {
	return "posts"
}

func (p *Post) Creating() {
	p.Saving()
}

func (p *Post) Updating() {
	p.Saving()
}

func (p *Post) Deleting() bool {
	err := p.DeleteImageFile()
	if err != nil {
		return false
	}
	return true
}

func (p *Post) Saving() {
	p.SetDirty()
	//logger.Print(p.GetDirty())
	p.setTitleShort()
	p.setAlias()
	p.setURL()
	p.setDate()
	p.SetDirty()
	//logger.Print(p.GetDirty())
}

func (p *Post) DeleteImageFile() error {
	if p.Image == nil {
		return nil
	}
	err := file.DeleteFile(*p.Image)
	if err != nil {
		return err
	}
	p.Image = nil
	return nil
}

func (p *Post) UploadImageFile(r *http.Request) error {
	_, img, _ := r.FormFile("file")
	if img != nil {
		newFileName := fmt.Sprintf("%s/%d", p.GetResource(), p.ID)
		path, err := file.SaveUploadedFile(img, newFileName)
		if err != nil {
			logger.Error(err)
			return err
		}
		if p.Image != nil {
			err := p.DeleteImageFile()
			if err != nil {
				return err
			}
		}
		p.Image = &path
	}
	return nil
}

func (p *Post) setAlias() {
	if p.Title == "" {
		return
	}
	if !p.isDirty("Alias") && p.Alias != "" {
		return
	}

	if p.Alias == "" {
		p.Alias = alias.Generate(p, p.Title)
	} else {
		p.Alias = alias.Generate(p, p.Alias)
	}
}

func (p *Post) setURL() {
	if p.Alias != "" {
		p.URL = "/" + p.Alias
	}
}

func (p *Post) setTitleShort() {
	if p.TitleShort == nil {
		return
	}
	if *p.TitleShort == "" {
		p.TitleShort = nil
	}
}

func (p *Post) setDate() {
	if p.DatePub == nil || p.DatePub.IsZero() {
		p.DatePub = nil
	}
	if p.DateEnd == nil || p.DateEnd.IsZero() {
		p.DateEnd = nil
	}
}

func (p *Post) SetOriginal(o *Post) {
	p.original = o
}

func (p *Post) SetDirty() {
	p.dirty = p.GetChangedFields(p.original, p)
}

func (p *Post) GetDirty() string {
	if len(p.dirty) > 0 {
		jsonData, err := json.Marshal(p.dirty)
		if err != nil {
			logger.Fatal(err)
		}
		return string(jsonData)
	}
	return ""
}

func (p *Post) isDirty(s string) bool {
	_, ok := p.dirty[s]
	return ok
}
