package post_category

import (
	. "github.com/axlle-com/blog/pkg/post"
	"gorm.io/gorm"
	"time"
)

type PostCategory struct {
	gorm.Model
	ID                 uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	PostCategoryID     uint           `gorm:"type:int;default:null;index;" json:"post_category_id"`
	Render             string         `gorm:"type:varchar(255);default:null" json:"render"`
	MetaTitle          string         `gorm:"type:varchar(100);default:null" json:"meta_title"`
	MetaDescription    string         `gorm:"type:varchar(200);default:null" json:"meta_description"`
	Alias              string         `gorm:"type:varchar(255);default:null;unique" json:"alias"`
	URL                string         `gorm:"type:varchar(1000);default:null;unique" json:"url"`
	IsPublished        uint8          `gorm:"type:smallint;default:1" json:"is_published"`
	IsFavourites       uint8          `gorm:"type:smallint;default:0" json:"is_favourites"`
	IsWatermark        uint8          `gorm:"type:smallint;default:1" json:"is_watermark"`
	IsSitemap          uint8          `gorm:"type:smallint;default:1" json:"is_sitemap"`
	Image              string         `gorm:"type:varchar(255);default:null" json:"image"`
	ShowImage          uint8          `gorm:"type:smallint;default:1" json:"show_image"`
	Title              string         `gorm:"type:varchar(255);not null;index" json:"title"`
	TitleShort         string         `gorm:"type:varchar(150);default:null" json:"title_short"`
	Description        string         `gorm:"type:text;default:null" json:"description"`
	PreviewDescription string         `gorm:"type:text;default:null" json:"preview_description"`
	Sort               uint           `gorm:"type:int;default:0" json:"sort"`
	Script             string         `gorm:"type:longtext;default:null" json:"script"`
	CSS                string         `gorm:"type:longtext;default:null" json:"css"`
	CreatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index;type:timestamp;default:null" json:"deleted_at"`
	Posts              []*Post
}
