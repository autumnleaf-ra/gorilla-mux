package productcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/autumnleaf-ra/gorilla-mux/helper"
	"github.com/autumnleaf-ra/gorilla-mux/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ResponseJson = helper.ResponseJson
var ResponseError = helper.ResponseError

// Menampilkan semua data
func Index(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	if err := models.DB.Find(&products).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJson(w, http.StatusOK, products)
}

// Menampilkan data berdasarkan ID
func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// parse ke int
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	// id terdeteksi
	var product models.Product
	if err := models.DB.First(&product, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, "Product tidak ditemukan !")
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	ResponseJson(w, http.StatusOK, product)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// request dari client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	// simpan database
	if err := models.DB.Create(&product).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJson(w, http.StatusCreated, product)

}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// parse ke int
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	var product models.Product

	// request dari client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	// update data
	if models.DB.Where("id = ?", id).Updates(product).RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "Tidak dapat mengupdate Product")
		return
	}

	product.Id = id

	ResponseJson(w, http.StatusOK, product)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// input dari client
	input := map[string]string{"id": ""}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var product models.Product
	if models.DB.Delete(&product, input["id"]).RowsAffected == 0 {
		ResponseError(w, http.StatusBadRequest, "Tidak dapat menghapus Product")
		return
	}

	response := map[string]string{"message": "Product berhasil dihapus"}
	ResponseJson(w, http.StatusOK, response)
}

// func Delete(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)

// 	// parse ke int
// 	id, err := strconv.ParseInt(vars["id"], 10, 64)
// 	if err != nil {
// 		ResponseError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	var product models.Product

// 	// id terdeteksi
// 	if models.DB.Where("id = ?", id).Delete(&product).RowsAffected == 0 {
// 		ResponseError(w, http.StatusInternalServerError, "Tidak dapat menghapus produk !")
// 		return
// 	}

// 	ResponseJson(w, http.StatusOK, "Produk berhasil dihapus !")

// }
