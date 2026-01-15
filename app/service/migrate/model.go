package migrate

import (
	"time"
)

type Seed struct {
	Name      string     `gorm:"primaryKey" json:"name" form:"name" binding:"-"`
	CreatedAt *time.Time `gorm:"index" json:"created_at,omitempty" form:"created_at" binding:"-" ignore:"true"`
}
