package itemsHandler

import (
	"database/sql"
	"encoding/json"
	"example.com/chi-restful/pkg/handlers/model"
	itemStorage "example.com/chi-restful/pkg/handlers/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type itemResource struct {
	itemStorage *itemStorage.ItemStorage
}

func New(conn *sql.DB) chi.Router {
	i := itemResource{itemStorage: itemStorage.New(conn)}
	return i.routes()
}

func (ir itemResource) routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", ir.getAllItems)
	r.Post("/", ir.createItem)

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", ir.updateItem)
		r.Get("/", ir.getById)
		r.Delete("/", ir.deleteById)
	})

	return r
}

func (ir itemResource) getAllItems(w http.ResponseWriter, r *http.Request) {
	rows, err := ir.itemStorage.GetAll()
	if err != nil {
		body, _ := json.Marshal(errorResponse{Message: err.Error()})
		w.Write(body)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(rows)
	w.Write(res)
}

func (ir itemResource) createItem(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var newItem model.Item
	err = json.Unmarshal(body, &newItem)
	newItem.Id = uuid.NewString()

	_ = ir.itemStorage.CreateItem(newItem)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (ir itemResource) updateItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var item model.Item
	json.Unmarshal(body, &item)

	itemId := chi.URLParam(r, "id")
	_, err = getById(ir.itemStorage, itemId)
	if err == sql.ErrNoRows || err != nil {
		http.Error(w, fmt.Sprintf("item %s not found", itemId), http.StatusNotFound)
		return
	}

	ir.itemStorage.UpdateById(itemId, item)

}

func getById(itemStorage *itemStorage.ItemStorage, id string) (model.Item, error) {
	item, err := itemStorage.GetById(id)
	if err == sql.ErrNoRows || err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func (ir itemResource) getById(w http.ResponseWriter, r *http.Request) {
	itemId := chi.URLParam(r, "id")
	item, err := getById(ir.itemStorage, itemId)
	if err == sql.ErrNoRows || err != nil {
		body, _ := json.Marshal(errorResponse{Message: fmt.Sprintf("item %s not found", itemId)})
		w.Write(body)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	body, _ := json.Marshal(item)
	w.Write(body)
}

func (ir itemResource) deleteById(w http.ResponseWriter, r *http.Request) {
	itemId := chi.URLParam(r, "id")
	_, err := getById(ir.itemStorage, itemId)
	if err == sql.ErrNoRows || err != nil {
		body, _ := json.Marshal(errorResponse{Message: fmt.Sprintf("item %s not found", itemId)})
		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
		return
	}
	if err = ir.itemStorage.DeleteById(itemId); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
