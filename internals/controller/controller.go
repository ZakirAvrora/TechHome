package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZakirAvrora/TechHome/internals/models"
	rep "github.com/ZakirAvrora/TechHome/internals/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type repository interface {
	ListLinks(page int) ([]models.Link, error)
	GetLink(id int) (models.Link, error)
	CreateLink(link models.Link) (int64, error)
	UpdateLink(link models.Link, id int) error
	DeleteLink(id int) error
	FindLink(link string) (models.Link, error)
}

type cache interface {
	Add(key string, link string)
	Get(key string) (string, bool)
	Len() int
}

type controller struct {
	repo  repository
	cache cache
}

func New(repo repository, cache cache) *controller {
	return &controller{repo: repo, cache: cache}
}

func (c *controller) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	pageId, ok := r.Context().Value("page").(int)
	if !ok {
		http.Error(w, "error in parsing page id", http.StatusInternalServerError)
		return
	}

	links, err := c.repo.ListLinks(pageId)

	if len(links) == 0 || errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JsonSend(links, w)
}

func (c *controller) GetLink(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link, err := c.repo.GetLink(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	JsonSend(link, w)
}

func (c *controller) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.repo.CreateLink(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("New link object with id %d was added", id)))

}

func (c *controller) UpdateLink(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var link models.Link
	err = json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.repo.UpdateLink(link, id)
	if err != nil {
		if errors.Is(err, rep.ErrNoRowAffected) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte("Successfully updated"))
}

func (c *controller) DeleteLink(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.repo.DeleteLink(id); err != nil {
		if errors.Is(err, rep.ErrNoRowAffected) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *controller) Redirect(w http.ResponseWriter, r *http.Request) {
	linkKey := r.URL.Query().Get("link")

	var actualLink string
	var status int = http.StatusOK
	var flag bool = false

	if cacheVal, ok := c.cache.Get(linkKey); ok {
		actualLink = cacheVal
		flag = true
	} else {
		link, err := c.repo.FindLink(linkKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		actualLink = link.ActiveLink
	}

	if actualLink != linkKey {
		status = http.StatusMovedPermanently
	}

	c.cache.Add(linkKey, actualLink)

	if flag {
		actualLink = "from the cache: " + actualLink
	}

	w.WriteHeader(status)
	JsonSend(actualLink, w)
}

func JsonSend(data any, w http.ResponseWriter) {
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
