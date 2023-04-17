package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"no-q-solution/domain/interfaces"
	"no-q-solution/domain/usecases"
	"no-q-solution/http/error"
	"no-q-solution/http/transport/request"
	"no-q-solution/http/transport/request/decoders"
	"no-q-solution/http/transport/response"
	"no-q-solution/http/validators"
	"no-q-solution/utils/container"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type QueueController struct {
	usecase   usecases.QueuetUsecase
	validator validators.Validator
	repo      interfaces.MerchantRepository
}

func NewQueueController(ctr container.Containers) QueueController {
	ctl := QueueController{
		usecase:   usecases.NewQueuetUsecase(ctr.Repositories.Queue),
		validator: validators.NewValidator(),
		repo:      ctr.Repositories.Merchant,
	}

	return ctl
}

func (ctl QueueController) GetByMerchant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	merchant_id, err := strconv.Atoi(vars["merchant_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	queues, err := ctl.usecase.GetByMerchant(ctx, int64(merchant_id))
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(queues, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl QueueController) GetSlotsByDate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	date, ok := vars["date"]
	if !ok {
		err := errors.New("date not provided")
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	givenDate, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		err := errors.New("given date is invalid")
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	queue, err := ctl.usecase.GetSlotsByDate(ctx, int64(queue_id), givenDate)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(queue, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl QueueController) MakeItAvailable(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchantID, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	merchant, err := ctl.usecase.MakeItAvailable(ctx, merchantID, int64(queue_id))
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchant, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl QueueController) MakeItUnAvailable(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchantID, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	merchant, err := ctl.usecase.MakeItUnAvailable(ctx, merchantID, int64(queue_id))
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchant, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl QueueController) MakeDatesAvailable(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchantID, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	decoder := decoders.Dates{}

	err = request.Decode(ctx, r, &decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = ctl.validator.Validate(ctx, decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	dates, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.MakeDatesAvailable(ctx, merchantID, int64(queue_id), dates)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl QueueController) MakeDatesUnAvailable(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchantID, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	decoder := decoders.Dates{}

	err = request.Decode(ctx, r, &decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = ctl.validator.Validate(ctx, decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	dates, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.MakeDatesUnAvailable(ctx, merchantID, int64(queue_id), dates)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusAccepted)
}

func (ctl QueueController) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")

	fmt.Println(authHeader)
	fmt.Println(r.Header)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchant_id, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	decoder := decoders.Queue{}

	err = request.Decode(ctx, r, &decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = ctl.validator.Validate(ctx, decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	queue, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	queue.MerchantID = merchant_id

	id, err := ctl.usecase.Create(ctx, queue)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl QueueController) ReserveSlot(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	decoder := decoders.ReserveSlot{}

	err := request.Decode(ctx, r, &decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = ctl.validator.Validate(ctx, decoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	reserveSlot, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.ReserveSlot(ctx, reserveSlot)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl QueueController) UnReserveSlot(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	_, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	token_no, err := strconv.Atoi(vars["token_no"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.UnReserveSlot(ctx, int64(token_no))
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl QueueController) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := errors.New("authorization token not found")
		log.Println(err.Error())

		error.HandleError(w, errors.New("authorization token not found"), http.StatusBadRequest)
		return
	}

	token := authHeader[len("Bearer "):]

	merchant_id, err := ctl.repo.ValidateToken(ctx, token)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	queue_id, err := strconv.Atoi(vars["queue_id"])
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.Delete(ctx, merchant_id, int64(queue_id))
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}
