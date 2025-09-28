package models

import (
	"fmt"
	"time"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PostCategory struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UUID               uuid.UUID  `gorm:"type:uuid;index,using:hash" json:"uuid" form:"uuid" binding:"-"`
	UserID             *uint      `gorm:"index" json:"user_id" form:"user_id" binding:"omitempty"`
	TemplateID         *uint      `gorm:"index" json:"template_id,omitempty"`
	PostCategoryID     *uint      `gorm:"index" json:"post_category_id,omitempty"`
	Path               string     `gorm:"index" json:"-"`
	MetaTitle          *string    `gorm:"size:100" json:"meta_title,omitempty"`
	MetaDescription    *string    `gorm:"size:200" json:"meta_description,omitempty"`
	Alias              string     `gorm:"size:255;unique" json:"alias"`
	URL                string     `gorm:"size:1000;unique" json:"url"`
	IsPublished        *bool      `gorm:"index;default:true" json:"is_published,omitempty"`
	IsFavourites       *bool      `gorm:"default:false" json:"is_favourites,omitempty"`
	InSitemap          *bool      `gorm:"index;default:true" json:"in_sitemap,omitempty"`
	Image              *string    `gorm:"size:255" json:"image,omitempty"`
	ShowImage          *bool      `gorm:"default:true" json:"show_image,omitempty"`
	Title              string     `gorm:"size:255;not null" json:"title"`
	TitleShort         *string    `gorm:"size:150" json:"title_short,omitempty"`
	Description        *string    `gorm:"type:text" json:"description,omitempty"`
	DescriptionPreview *string    `gorm:"type:text" json:"description_preview,omitempty"`
	Sort               *uint      `gorm:"index;default:0" json:"sort,omitempty"`
	CreatedAt          *time.Time `gorm:"index" json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	GalleriesSnapshot  datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'::jsonb" json:"galleries_snapshot"`
	InfoBlocksSnapshot datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'::jsonb" json:"info_blocks_snapshot"`

	Galleries  []contracts.Gallery   `gorm:"-" json:"galleries" form:"galleries" binding:"-" ignore:"true"`
	InfoBlocks []contracts.InfoBlock `gorm:"-" json:"info_blocks" form:"info_blocks" binding:"-" ignore:"true"`
	Category   *PostCategory         `gorm:"-" json:"category" form:"category" binding:"-" ignore:"true"`
	Template   contracts.Template    `gorm:"-" json:"template" form:"template" binding:"-" ignore:"true"`
	User       contracts.User        `gorm:"-" json:"user" form:"user" binding:"-" ignore:"true"`
}

func (c *PostCategory) GetID() uint {
	return c.ID
}

func (c *PostCategory) GetTable() string {
	return "post_categories"
}

func (c *PostCategory) GetUUID() uuid.UUID {
	return c.UUID
}

func (c *PostCategory) GetName() string {
	return c.GetTable()
}

func (c *PostCategory) GetURL() string {
	return c.URL
}

func (c *PostCategory) GetTitle() string {
	if c.TitleShort != nil && *c.TitleShort != "" {
		return *c.TitleShort
	}
	return c.Title
}

func (c *PostCategory) GetTemplateName() string {
	if c.Template != nil {
		return fmt.Sprintf("%s.%s", c.GetTable(), c.Template.GetName())
	}
	return fmt.Sprintf("%s.default", c.GetTable())
}

func (c *PostCategory) AdminURL() string {
	if c.ID == 0 {
		return "/admin/post/categories"
	}
	return fmt.Sprintf("/admin/post/categories/%d", c.ID)
}

func (c *PostCategory) SetUUID() {
	if c.UUID == uuid.Nil {
		c.UUID = uuid.New()
	}
}

func (c *PostCategory) SetAlias() {
	if c.Alias != "" {
		return
	}
	if c.UUID == uuid.Nil {
		c.SetUUID()
	}
	c.Alias = c.UUID.String()
}

func (c *PostCategory) GetCategoryID() uint {
	var categoryID uint
	if c.PostCategoryID != nil {
		categoryID = *c.PostCategoryID
	}
	return categoryID
}

func (c *PostCategory) GetTemplateID() uint {
	var templateID uint
	if c.TemplateID != nil {
		templateID = *c.TemplateID
	}
	return templateID
}

func (c *PostCategory) GetCategoryTitleShort() string {
	var titleShort string
	if c.Category != nil {
		titleShort = *c.Category.TitleShort
	}
	return titleShort
}

func (c *PostCategory) GetTemplateTitle() string {
	var title string
	if c.Template != nil {
		title = c.Template.GetTitle()
	}
	return title
}

func (c *PostCategory) UserLastName() string {
	var lastName string
	if c.User != nil {
		lastName = c.User.GetLastName()
	}
	return lastName
}

func (c *PostCategory) Creating() {
	c.Saving()
}

func (c *PostCategory) Updating() {
	c.Saving()
}

func (c *PostCategory) Deleting() bool {
	return true
}

func (c *PostCategory) Saving() {
	c.SetUUID()
	c.SetAlias()
	c.setTitleShort()
	c.setURL()
}

func (c *PostCategory) setURL() {
	if c.Alias != "" {
		c.URL = "/" + c.Alias
	}
}

func (c *PostCategory) setTitleShort() {
	if c.TitleShort == nil {
		return
	}
	if *c.TitleShort == "" {
		c.TitleShort = nil
	}
}

func (c *PostCategory) Date() string {
	if c.CreatedAt == nil {
		return ""
	}
	return c.CreatedAt.Format("02.01.2006 15:04:05")
}
