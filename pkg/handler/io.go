package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/EnMasseProject/maas-service-broker/pkg/broker"
	"github.com/EnMasseProject/maas-service-broker/pkg/errors"
	"github.com/op/go-logging"
	"strconv"
)

func readRequest(log *logging.Logger, r *http.Request, obj interface{}) error {
	// TODO: uncomment this when the service catalog controller starts setting the content-type properly
	//contentType := r.Header.Get("Content-Type")
	//if contentType != "application/json" {
	//	return errors.NewBadRequest("error: invalid content-type: " + contentType)
	//}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	reader := bytes.NewReader(buf.Bytes())
	log.Infof("Request body: ", buf.String())

	err := json.NewDecoder(reader).Decode(&obj)
	if err != nil {
		log.Info("Could not parse request body: " + err.Error())
		return errors.NewBadRequest("could not parse request body : " + err.Error())
	}

	return nil
}

func writeResponse(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// return json.NewEncoder(w).Encode(obj)

	// pretty-print for easier debugging
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	i := bytes.Buffer{}
	json.Indent(&i, b, "", "  ")
	i.WriteString("\n")
	_, err = w.Write(i.Bytes())
	return err
}

func writeDefaultResponse(w http.ResponseWriter, code int, resp interface{}, err error, log *logging.Logger) error {
	if err == nil {
		return writeResponse(w, code, resp)
	} else {
		return writeErrorResponse(w, err, log)
	}
}

func writeErrorResponse(w http.ResponseWriter, err error, log *logging.Logger) error {
	if brokerError, ok := err.(errors.BrokerError); ok {
		log.Warning("Sending broker error response: " + strconv.Itoa(brokerError.Status) + ", " + brokerError.Description)
		return writeResponse(w, brokerError.Status, broker.NewErrorResponse(brokerError.Description))
	} else {
		log.Warning("Sending internal error response: " + err.Error())
		return writeResponse(w, http.StatusInternalServerError, broker.NewErrorResponse("Unknown error: "+err.Error()))
	}
}
