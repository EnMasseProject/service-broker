package handler

import (
	"net/http"

	"github.com/EnMasseProject/maas-service-broker/pkg/broker"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/op/go-logging"
)

// TODO: implement asynchronous operations
// TODO: authentication / authorization

type handler struct {
	log *logging.Logger
	router mux.Router
	broker broker.Broker
}

func NewHandler(log *logging.Logger, b broker.Broker) http.Handler {
	h := handler{log: log, broker: b}

	// TODO: handle X-Broker-API-Version header, currently poorly defined
	//root := h.router.Headers("X-Broker-API-Version", "2.9").Subrouter()
	root := h.router.PathPrefix("/").Subrouter()

	root.HandleFunc("/v2/bootstrap", h.bootstrap).Methods(http.MethodPost)
	root.HandleFunc("/v2/catalog", h.catalog).Methods(http.MethodGet)
	root.HandleFunc("/v2/service_instances/{instance_uuid}", h.provision).Methods(http.MethodPut)
	root.HandleFunc("/v2/service_instances/{instance_uuid}", h.update).Methods(http.MethodPatch)
	root.HandleFunc("/v2/service_instances/{instance_uuid}", h.deprovision).Methods(http.MethodDelete)
	root.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}", h.bind).Methods(http.MethodPut)
	root.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}", h.unbind).Methods(http.MethodDelete)

	// TODO NotFoundHandler (must return json!)

	return h
}

func (h handler) bootstrap(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp, err := h.broker.Bootstrap()
	writeDefaultResponse(w, http.StatusOK, resp, err, h.log)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h handler) catalog(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp, err := h.broker.Catalog()
	writeDefaultResponse(w, http.StatusOK, resp, err, h.log)
}

func (h handler) provision(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	h.log.Info("Received provision request: " + r.RequestURI)

	instanceUUID := uuid.Parse(mux.Vars(r)["instance_uuid"])
	if instanceUUID == nil {
		h.log.Info("Invalid instance_uuid in request")
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	var req *broker.ProvisionRequest
	err := readRequest(h.log, r, &req)

	if err != nil {
		writeErrorResponse(w, err, h.log)
		return
	}

	resp, err := h.broker.Provision(instanceUUID, req)
	writeDefaultResponse(w, http.StatusCreated, resp, err, h.log)
}

func (h handler) update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	instanceUUID := uuid.Parse(mux.Vars(r)["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	var req *broker.UpdateRequest
	if err := readRequest(h.log, r, &req); err != nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: err.Error()})
		return
	}

	resp, err := h.broker.Update(instanceUUID, req)

	writeDefaultResponse(w, http.StatusOK, resp, err, h.log)
}

func (h handler) deprovision(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	h.log.Info("Received deprovision request: " + r.RequestURI)

	instanceUUIDstring := mux.Vars(r)["instance_uuid"]
	instanceUUID := uuid.Parse(instanceUUIDstring)
	if instanceUUID == nil {
		h.log.Info("Invalid instance_uuid in request: %s", instanceUUIDstring)
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	serviceId := r.FormValue("service_id")
	if serviceId == "" {
		h.log.Info("Missing service_id parameter")
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "missing service_id parameter"})
		return
	}

	planId := r.FormValue("plan_id")
	if planId == "" {
		h.log.Info("Missing plan_id parameter")
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "missing plan_id parameter"})
		return
	}

	resp, err := h.broker.Deprovision(instanceUUID, serviceId, planId)

	//if errors.IsNotFound(err) {
	//	writeResponse(w, http.StatusGone, broker.DeprovisionResponse{})
	//} else {
		writeDefaultResponse(w, http.StatusOK, resp, err, h.log)
	//}
}

func (h handler) bind(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	instanceUUID := uuid.Parse(mux.Vars(r)["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	bindingUUID := uuid.Parse(mux.Vars(r)["binding_uuid"])
	if bindingUUID == nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid binding_uuid"})
		return
	}

	var req *broker.BindRequest
	if err := readRequest(h.log, r, &req); err != nil {
		writeResponse(w, http.StatusInternalServerError, broker.ErrorResponse{Description: err.Error()})
		return
	}

	resp, err := h.broker.Bind(instanceUUID, bindingUUID, req)

	writeDefaultResponse(w, http.StatusCreated, resp, err, h.log)
}

func (h handler) unbind(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	instanceUUID := uuid.Parse(mux.Vars(r)["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	bindingUUID := uuid.Parse(mux.Vars(r)["binding_uuid"])
	if bindingUUID == nil {
		writeResponse(w, http.StatusBadRequest, broker.ErrorResponse{Description: "invalid binding_uuid"})
		return
	}

	err := h.broker.Unbind(instanceUUID, bindingUUID)

	//if errors.IsNotFound(err) {
	//	writeResponse(w, http.StatusGone, struct{}{})
	//} else {
		writeDefaultResponse(w, http.StatusOK, struct{}{}, err, h.log)
	//}
	return
}
