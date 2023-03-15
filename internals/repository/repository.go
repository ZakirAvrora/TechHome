package repository

import "github.com/ZakirAvrora/TechHome/internals/models"

type Repository interface {
	ListLinks(page int) ([]models.Link, error)
	GetLink(id int) (models.Link, error)
	CreateLink(link models.Link) (int64, error)
	UpdateLink(link models.Link, id int) error
	DeleteLink(id int) error
	FindLink(link string) (models.Link, error)
}
