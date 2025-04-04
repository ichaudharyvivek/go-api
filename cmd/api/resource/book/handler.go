package book

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type API struct {
	repository *Repository
}

func New(db *gorm.DB) *API {
	return &API{
		repository: NewRepository(db),
	}
}

// List godoc
//
//	@summary		List books
//	@description	List books
//	@tags			books
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	err.Error
//	@router			/books [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	books, err := a.repository.List()
	if err != nil {
		// Handle later
		return
	}

	if len(books) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(books.ToDto()); err != nil {
		// Handle later
		return
	}
}

// Create godoc
//
//	@summary		Create book
//	@description	Create book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Book form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/books [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		// handle later
		return
	}

	newBook := form.ToModel()
	newBook.ID = uuid.New()

	_, err := a.repository.Create(newBook)
	if err != nil {
		// handle later
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Get book with id
//	@description	Read book with id
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Book ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [get]
func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	book, err := a.repository.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// handle later
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		// handle later
		return
	}
}

// Update godoc
//
//	@summary		Update book
//	@description	Update book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"Book ID"
//	@param			body	body	Form	true	"Book form"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		// handle later
		return
	}

	book := form.ToModel()
	book.ID = id

	rows, err := a.repository.Update(book)
	if err != nil {
		// handle later
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Delete godoc
//
//	@summary		Delete book
//	@description	Delete book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Book ID"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		// handle later
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
