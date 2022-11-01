package master

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	masterServer struct {
		hserver *http.Server

		clients map[string]*masterClient
	}
	masterClient struct {
		conn *websocket.Conn

		id string
	}
	apiError struct {
		Code    int
		Message string
	}
	apiRequest struct {
		Id      string
		Method  apiMethod
		Payload interface{}
	}
	apiResponse struct {
		RequestId string

		Ok      bool
		Error   *apiError
		Payload interface{}
	}
)

type apiMethod uint8

const (
	apiMethodRegistration apiMethod = iota
)

const (
	apiErrWorkerId     = "Worker-Id not specified"
	apiErrUpgrade      = "got an internal error in proto upgrading process"
	apiErrUnauthorized = "secret key is invalid"
)

var (
	wsUpgrader = websocket.Upgrader{}
)

func newApiError(code int, desc string) *apiError {
	return &apiError{
		Code:    code,
		Message: desc,
	}
}

func (m *apiError) httpErrorRespond(w http.ResponseWriter) {
	defaultRespond := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if e := json.NewEncoder(w).Encode(m); e != nil {
		gLog.Error().Err(e).Msg("got an error in json marshaling of the api error")
		defaultRespond(w)
		return
	}

	w.WriteHeader(m.Code)
}

func newMasterServer() *masterServer {
	return &masterServer{}
}

func (*masterServer) bootstrap() (e error) {
	return
}

func (*masterServer) serve() {

}

func (*masterServer) getApiError(status int, error string) *apiError {
	return &apiError{}
}

func (m *masterServer) handleIncomingRequest(w http.ResponseWriter, r *http.Request) {
	var id string
	if id = r.Header.Get("X-Worker-Id"); id == "" {
		newApiError(http.StatusBadRequest, apiErrWorkerId).httpErrorRespond(w)
		return
	}

	gLog.Info().Str("remote_addr", r.RemoteAddr).Str("user-agent", r.UserAgent()).Str("worker_id", r.Header.Get("X-Worker-Id")).
		Msg("handling new incomig connection...")
	if ok, e := m.authorizeIncomingRequest(id, r.Header.Get("Authorization")); e != nil {
		return
	} else if !ok {
		newApiError(http.StatusUnauthorized, apiErrUnauthorized).httpErrorRespond(w)
	}

	gLog.Debug().Str("worker_id", id).Msg("upgrading proto to ws...")
	conn, e := wsUpgrader.Upgrade(w, r, nil)
	if e != nil {
		newApiError(http.StatusInternalServerError, e.Error()).httpErrorRespond(w)
		return
	}

	m.clients[id] = &masterClient{
		id:   id,
		conn: conn,
	}

	//
}

func (*masterServer) authorizeIncomingRequest(id, payload string) (_ bool, e error) {
	gLog.Debug().Str("worker_id", id).Msg("authorizing new incoming request...")

	//
	return
}
