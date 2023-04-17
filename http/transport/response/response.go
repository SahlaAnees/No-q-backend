package response

import "net/http"

func Send(w http.ResponseWriter, payload []byte, code int) {

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	w.Write(payload)
}
