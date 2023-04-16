package controllers

import (
	"encoding/json"
	"errors"
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
	"strings"

	"github.com/gorilla/mux"
)

type MerchantController struct {
	usecase   usecases.MerchantUsecase
	validator validators.Validator
	repo      interfaces.MerchantRepository
}

func NewMerchantController(ctr container.Containers) MerchantController {
	ctl := MerchantController{
		usecase:   usecases.NewMerchantUsecase(ctr.Repositories.Merchant),
		validator: validators.NewValidator(),
		repo:      ctr.Repositories.Merchant,
	}

	return ctl
}

func (ctl MerchantController) GetAll(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	param := r.FormValue("paginator")

	pageDecoder := decoders.Paginator{}

	err := json.Unmarshal([]byte(param), &pageDecoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	err = ctl.validator.Validate(ctx, pageDecoder)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	paginator, err := pageDecoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	merchants, err := ctl.usecase.GetAll(ctx, paginator)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchants, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) GetCategories(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	categories, err := ctl.usecase.GetCategories(ctx)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(categories, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) GetByCategory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	category, ok := vars["category"]
	if !ok {
		err := errors.New("category not found")

		log.Println(err.Error())
		error.HandleError(w, err, http.StatusBadRequest)
	}

	merchants, err := ctl.usecase.GetByCategory(ctx, category)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchants, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) GetSingle(w http.ResponseWriter, r *http.Request) {

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

	merchant, err := ctl.usecase.GetSingle(ctx, merchantID)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchant, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) Search(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	vars := mux.Vars(r)

	input := vars["input"]

	merchants, err := ctl.usecase.Search(ctx, input)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(merchants, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	decoder := decoders.Merchant{}

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

	merchant, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := ctl.usecase.Create(ctx, merchant)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(id, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl MerchantController) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	decoder := decoders.Login{}

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

	login, err := decoder.Validate()
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusBadRequest)
		return
	}

	token, err := ctl.usecase.Login(ctx, login)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(token, nil, "true")

	response.Send(w, payload, http.StatusCreated)
}

func (ctl MerchantController) Logout(w http.ResponseWriter, r *http.Request) {

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

	done, err := ctl.usecase.Logout(ctx, merchantID)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(done, nil, "true")

	response.Send(w, payload, http.StatusOK)
}

func (ctl MerchantController) Delete(w http.ResponseWriter, r *http.Request) {

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

	done, err := ctl.usecase.Delete(ctx, merchantID)
	if err != nil {
		log.Println(err.Error())

		error.HandleError(w, err, http.StatusUnprocessableEntity)
		return
	}

	payload := response.Encode(done, nil, "true")

	response.Send(w, payload, http.StatusOK)
}
