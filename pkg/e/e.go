package e

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	type envelope map[string]any
	env := envelope{"error": message}

	js, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
}
