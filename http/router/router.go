package router

import (
	"net/http"
	"no-q-solution/http/controllers"
	"no-q-solution/http/transport/response"
	"no-q-solution/utils/container"

	"github.com/gorilla/mux"
)

func Init(ctr container.Containers) *mux.Router {

	r := mux.NewRouter()

	merchant := controllers.NewMerchantController(ctr)
	queue := controllers.NewQueueController(ctr)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.Send(w, []byte("No-Q Solution"), http.StatusOK)
	}).Methods(http.MethodGet)

	r.HandleFunc("/merchant/get_all", merchant.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/merchant/get_categories", merchant.GetCategories).Methods(http.MethodGet)
	r.HandleFunc("/merchant/get_by_category/{category}", merchant.GetByCategory).Methods(http.MethodGet)
	r.HandleFunc("/merchant/get_single", merchant.GetSingle).Methods(http.MethodGet)
	r.HandleFunc("/merchant/search/{input}", merchant.Search).Methods(http.MethodGet)
	r.HandleFunc("/merchant/create", merchant.Create).Methods(http.MethodPost)
	r.HandleFunc("/merchant/login", merchant.Login).Methods(http.MethodPost)
	r.HandleFunc("/merchant/logout", merchant.Logout).Methods(http.MethodGet)
	r.HandleFunc("/merchant/delete", merchant.Delete).Methods(http.MethodDelete)

	r.HandleFunc("/queue/get_by_merchant/{merchant_id}", queue.GetByMerchant).Methods(http.MethodGet)
	r.HandleFunc("/queue/get_slots_by_date/{queue_id}/{date}", queue.GetSlotsByDate).Methods(http.MethodGet)
	r.HandleFunc("/queue/make_it_available/{queue_id}", queue.MakeItAvailable).Methods(http.MethodPatch)
	r.HandleFunc("/queue/make_it_un_available/{queue_id}", queue.MakeItUnAvailable).Methods(http.MethodPatch)
	r.HandleFunc("/queue/make_dates_available/{queue_id}", queue.MakeDatesAvailable).Methods(http.MethodPost)
	r.HandleFunc("/queue/make_dates_un_available/{queue_id}", queue.MakeDatesUnAvailable).Methods(http.MethodDelete)
	r.HandleFunc("/queue/create", queue.Create).Methods(http.MethodPost)
	r.HandleFunc("/queue/reserve_slot", queue.ReserveSlot).Methods(http.MethodPost)
	r.HandleFunc("/queue/un_reserve_slot/{token_no}", queue.UnReserveSlot).Methods(http.MethodDelete)
	r.HandleFunc("/queue/delete/{queue_id}", queue.Delete).Methods(http.MethodDelete)

	return r
}
