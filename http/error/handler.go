package error

import (
	"net/http"
	"no-q-solution/http/transport/response"
)

func HandleError(w http.ResponseWriter, err error, code int) {

	payload := response.Encode(nil, err.Error(), "false")

	response.Send(w, payload, code)
}
