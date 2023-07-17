package route

import (
	"github.com/dipankarupd/immigrant-management-system/pkg/controller"
	"github.com/gorilla/mux"
)

var RegisterNewroute = func(router *mux.Router) {

	router.HandleFunc("/immigrants", controller.CreateImmigrant).Methods("POST")
	router.HandleFunc("/immigrants", controller.GetImmigrants).Methods("GET")
	router.HandleFunc("/immigrants/{passportno}", controller.GetImmigrant).Methods("GET")
	router.HandleFunc("/immigrants/accept/{passportno}", controller.AcceptImmigrant).Methods("PUT")
}
