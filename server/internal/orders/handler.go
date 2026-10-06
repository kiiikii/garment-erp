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

type LayoutAssignmentRequest struct {
	LayoutID int `json:"layout_id"`
}

type SampleCompleteRequest struct {
	FinishedAt string `json:"finished_at"`
}

type WaitingReasonRequest struct {
	Reason string `json:"reason"`
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
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		h.GetByID(w, r)
		return
	}

	// Otherwise, return the full list for the dashboard
	orders, err := h.service.GetAllOrders()
	if err != nil {
		fmt.Println("GetAllOrders Handler Failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
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
	err = h.service.UpdateStatus(orderID, req.NewStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d successfully update to %s\n", orderID, req.NewStatus)
}

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

func (h *OrderHandler) AssignLayout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID or URL", http.StatusBadRequest)
		return
	}

	var req LayoutAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.service.AssignLayout(orderID, req.LayoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d successfully assign to layout #%d\n", orderID, req.LayoutID)
}

func (h *OrderHandler) CompleteSampling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID in URL", http.StatusBadRequest)
		return
	}

	finishedTime := time.Now().Format("2006-01-02 15:04:05")

	err = h.service.CompleteSampling(orderID, finishedTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d sampling marked as finished at %s\n", orderID, finishedTime)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Order Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) SetWaitingReason(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID in URL", http.StatusBadRequest)
		return
	}

	var req WaitingReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	if req.Reason == "" {
		http.Error(w, "Waiting reason cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.service.SetWaitingForMaterials(orderID, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d flagged as WAITING_FOR_MATERRIALS, Reason: %s\n", orderID, req.Reason)
}

func (h *OrderHandler) ResumeOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = h.service.ResumeOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order #%d resumed and is now READY_FOR_PRODUCTION\n", orderID)
}

func (h *OrderHandler) HoldProduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.service.HoldProduction(orderID, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produciton for order #%d halted successfully\n", orderID)
}

func (h *OrderHandler) ResumeProduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}

	err = h.service.ResumeProduction(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produciton for order #%d resumed successfully\n", orderID)
}
