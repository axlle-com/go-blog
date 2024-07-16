package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model                        // adds ID, created_at etc.
	ID                 uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	PostCategoryID     uint           `gorm:"type:int;default:null;index;" json:"post_category_id"`
	Render             string         `gorm:"type:varchar(255);default:null" json:"render"`
	MetaTitle          string         `gorm:"type:varchar(100);default:null" json:"meta_title"`
	MetaDescription    string         `gorm:"type:varchar(200);default:null" json:"meta_description"`
	Alias              string         `gorm:"type:varchar(255);default:null;unique" json:"alias"`
	URL                string         `gorm:"type:varchar(1000);default:null;unique" json:"url"`
	IsPublished        uint8          `gorm:"type:smallint;default:1" json:"is_published"`
	IsFavourites       uint8          `gorm:"type:smallint;default:0" json:"is_favourites"`
	IsComments         uint8          `gorm:"type:smallint;default:0" json:"is_comments"`
	IsImagePost        uint8          `gorm:"type:smallint;default:1" json:"is_image_post"`
	IsImageCategory    uint8          `gorm:"type:smallint;default:1" json:"is_image_category"`
	IsWatermark        uint8          `gorm:"type:smallint;default:1" json:"is_watermark"`
	IsSitemap          uint8          `gorm:"type:smallint;default:1" json:"is_sitemap"`
	Media              string         `gorm:"type:varchar(255);default:null" json:"media"`
	Title              string         `gorm:"type:varchar(255);not null;index" json:"title"`
	TitleShort         string         `gorm:"type:varchar(155);default:null" json:"title_short"`
	DescriptionPreview string         `gorm:"type:text;default:null" json:"description_preview"`
	Description        string         `gorm:"type:text;default:null" json:"description"`
	ShowDate           uint8          `gorm:"type:smallint;default:1" json:"show_date"`
	DatePub            time.Time      `gorm:"type:timestamp;default:null" json:"date_pub"`
	DateEnd            time.Time      `gorm:"type:timestamp;default:null" json:"date_end"`
	ControlDatePub     uint8          `gorm:"type:smallint;default:0" json:"control_date_pub"`
	ControlDateEnd     uint8          `gorm:"type:smallint;default:0" json:"control_date_end"`
	Image              string         `gorm:"type:varchar(255);default:null" json:"image"`
	Hits               uint           `gorm:"type:int;default:0" json:"hits"`
	Sort               int            `gorm:"type:int;default:0" json:"sort"`
	Stars              float32        `gorm:"type:real;default:0.0" json:"stars"`
	Script             string         `gorm:"type:text;default:null" json:"script"`
	CSS                string         `gorm:"type:text;default:null" json:"css"`
	CreatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index;type:timestamp;default:null" json:"deleted_at"`
}
