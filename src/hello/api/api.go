package api

import "github.com/julienschmidt/httprouter"

type API interface {
	Bind(r *httprouter.Router) error
}
