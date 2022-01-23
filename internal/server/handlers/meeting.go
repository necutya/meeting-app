package handlers

import (
	"net/http"

	"server/internal/models"
	"server/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type MeetingHandler struct {
	service *service.Service
	upgrader websocket.Upgrader
}

func NewMeetingHandler(service *service.Service) *MeetingHandler {

	return &MeetingHandler{
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// swagger:operation GET /assignments/{id}/history assignments get_assignments_history
//
// ---
// parameters:
// - name: id
//   in: path
//   required: true
//   type: string
// responses:
//   '200':
//     description: Fetched
//     schema:
//       "$ref": "#/definitions/AssigmentHistory"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (ms *MeetingHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	room, err := ms.service.CreateRoom(r.Context())
	if err != nil {
		SendHTTPError(w, err)
		return
	}

	SendResponse(w, http.StatusOK, models.CreateRoomHTTPResponse{RoomID: room.UUID})
}

// swagger:operation GET /assignments/{id}/history assignments get_assignments_history
//
// ---
// parameters:
// - name: id
//   in: path
//   required: true
//   type: string
// responses:
//   '200':
//     description: Fetched
//     schema:
//       "$ref": "#/definitions/AssigmentHistory"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
func (ms *MeetingHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	log.Info("Joining room...")
	vars := mux.Vars(r)
	roomIDstr, ok := vars["roomID"]
	if !ok {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDstr)
	if err != nil {
		SendHTTPError(w, err)
		return
	}

	names, ok := r.URL.Query()["username"]
	if !ok {
		SendEmptyResponse(w, http.StatusBadRequest)
		return
	}

	// var req models.JoinRoomHTPPRequest
	// err = UnmarshalRequest(r, &req)
	// if err != nil {
	// 	SendHTTPError(w, err)
	// }

	ws, err := ms.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Web Socket Upgrade Error", err)
	}

	log.Info(names[0])
	ms.service.JoinRoom(r.Context(), roomID, models.RoomParticipant{
		Username: names[0],
		Conn: ws,
	})
}
