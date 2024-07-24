package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, data JSONResponse) {
	w.WriteHeader(data.Status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	data := JSONResponse{
		Status:  statusCode,
		Message: "",
		Data:    nil,
	}
	if err != nil {
		data.Message = err.Error()
	}
	JSON(w, data)
}
