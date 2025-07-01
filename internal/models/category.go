package models

type Category struct {
	ID       int64      `json:"id" db:"id"`
	Name     string     `json:"name" db:"name"`
	ParentID *int64     `json:"parent_id" db:"parent_id"`
	Children []Category `json:"children" db:"-"`
}
