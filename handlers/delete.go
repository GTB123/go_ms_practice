package handlers

import (
	"data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} products deleteProduct
// responses:
//	201: noContent

// DeleteProduct removes a product from the data store by id

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	// this will always convert because of router
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
