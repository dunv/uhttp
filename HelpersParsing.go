package uhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseBody(r *http.Request, model interface{}) error {
	reader, err := ReaderHelper(r.Header, r.Body)
	if err != nil {
		return fmt.Errorf("err parsing request (err getting reader %s)", err)
	}

	err = json.NewDecoder(reader).Decode(model)
	if err != nil {
		return fmt.Errorf("err parsing request (err marshaling %s)", err)
	}
	defer r.Body.Close()

	return nil
}
