package request

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func Decode(ctx context.Context, r *http.Request, decoder DecoderInterface) error {

	bite, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.Unmarshal(bite, decoder)
	if err != nil {
		return nil
	}

	return nil
}
