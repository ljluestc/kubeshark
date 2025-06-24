package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kubeshark/kubeshark/internal/services"
)

type ConnectionHandler struct {
	connectionService services.ConnectionService
}

func NewConnectionHandler(connectionService services.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: connectionService,
	}
}

func (h *ConnectionHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	includeHalf := false
	halfParam := r.URL.Query().Get("half")
	if halfParam != "" {
		var err error
		includeHalf, err = strconv.ParseBool(halfParam)
		if err != nil {
			http.Error(w, "Invalid 'half' parameter, must be true or false", http.StatusBadRequest)
			return
		}
	}

	connections, err := h.connectionService.GetConnections(includeHalf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connections)
}
