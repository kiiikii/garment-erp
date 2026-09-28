package orders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	newID, err := h.service.CreateOrder(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Success creates Order #%d\n", newID)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderList, err := h.service.GetAllOrders()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderList)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	//! enforcing PATCH
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	//! get the order ID
	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID in URL", http.StatusBadRequest)
		return
	}

	//! decode JSON body
	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed Decode JSON", http.StatusBadRequest)
		return
	}

	//! send to service layer
	err = h.service.UpdateStatus(orderID, req.CurrentStatus, req.NewStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d successfully update to %s\n", orderID, req.NewStatus)
}

// ! handling upload image
func (h *OrderHandler) UploadImages(w http.ResponseWriter, r *http.Request) {
	//! enforce POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//! parsing multipart form data (with limit upload 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parsing file. File might be too large", http.StatusBadRequest)
		return
	}

	//! get the order id from form text field
	orderIDStr := r.FormValue("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}

	//! retreive the array of files
	files := r.MultipartForm.File["tech_packs"]
	if len(files) == 0 {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	//! loop through every upload
	for _, fileHeader := range files {
		//! open incoming file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		//! Generate a unique filename
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

		//! the path where it will be saved on our drive
		savePath := filepath.Join("..", "uploads", fileName)

		//! create empty file on our drive
		destinationFile, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to save file on server", http.StatusInternalServerError)
			return
		}
		defer destinationFile.Close()

		//! copy the binary data from the incoming req
		if _, err := io.Copy(destinationFile, file); err != nil {
			http.Error(w, "Failed to write file data", http.StatusInternalServerError)
			return
		}

		//! save the database record
		dbFileURL := fmt.Sprintf("/uploads/%s", fileName)
		err = h.service.SaveImages(orderID, dbFileURL)
		if err != nil {
			destinationFile.Close()
			os.Remove(savePath)

			http.Error(w, "Failed to link image to database", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Succes uploaded %d images for order #%d", len(files), orderID)
}
