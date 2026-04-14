package stdhttp

import (
	"encoding/json"
	"net/http"
)

// Server implements ServerInterface by echoing received parameters back as JSON.
type Server struct{}

var _ ServerInterface = (*Server)(nil)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) GetSimplePrimitive(w http.ResponseWriter, r *http.Request, param int32)        { writeJSON(w, param) }
func (s *Server) GetSimpleExplodePrimitive(w http.ResponseWriter, r *http.Request, param int32)  { writeJSON(w, param) }
func (s *Server) GetSimpleNoExplodeArray(w http.ResponseWriter, r *http.Request, param []int32)  { writeJSON(w, param) }
func (s *Server) GetSimpleExplodeArray(w http.ResponseWriter, r *http.Request)                   { w.WriteHeader(http.StatusOK) }
func (s *Server) GetSimpleNoExplodeObject(w http.ResponseWriter, r *http.Request, param Object)  { writeJSON(w, param) }
func (s *Server) GetSimpleExplodeObject(w http.ResponseWriter, r *http.Request)                  { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelPrimitive(w http.ResponseWriter, r *http.Request)                       { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelExplodePrimitive(w http.ResponseWriter, r *http.Request)                { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelNoExplodeArray(w http.ResponseWriter, r *http.Request)                  { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelExplodeArray(w http.ResponseWriter, r *http.Request)                    { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelNoExplodeObject(w http.ResponseWriter, r *http.Request)                 { w.WriteHeader(http.StatusOK) }
func (s *Server) GetLabelExplodeObject(w http.ResponseWriter, r *http.Request)                   { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixPrimitive(w http.ResponseWriter, r *http.Request)                      { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixExplodePrimitive(w http.ResponseWriter, r *http.Request)               { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixNoExplodeArray(w http.ResponseWriter, r *http.Request)                 { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixExplodeArray(w http.ResponseWriter, r *http.Request)                   { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixNoExplodeObject(w http.ResponseWriter, r *http.Request)                { w.WriteHeader(http.StatusOK) }
func (s *Server) GetMatrixExplodeObject(w http.ResponseWriter, r *http.Request)                  { w.WriteHeader(http.StatusOK) }
func (s *Server) GetContentObject(w http.ResponseWriter, r *http.Request, param string)          { writeJSON(w, param) }
func (s *Server) GetPassThrough(w http.ResponseWriter, r *http.Request, param string)            { writeJSON(w, param) }
func (s *Server) GetQueryForm(w http.ResponseWriter, r *http.Request, params GetQueryFormParams)  { writeJSON(w, params) }
func (s *Server) GetDeepObject(w http.ResponseWriter, r *http.Request, params GetDeepObjectParams) { writeJSON(w, params) }
func (s *Server) GetHeader(w http.ResponseWriter, r *http.Request, params GetHeaderParams)       { writeJSON(w, params) }
func (s *Server) GetCookie(w http.ResponseWriter, r *http.Request, params GetCookieParams)       { writeJSON(w, params) }
