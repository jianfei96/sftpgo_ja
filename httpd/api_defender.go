package httpd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/render"

	"github.com/drakkan/sftpgo/v2/common"
	"github.com/drakkan/sftpgo/v2/dataprovider"
)

func getDefenderHosts(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	hosts, err := common.GetDefenderHosts()
	if err != nil {
		sendAPIResponse(w, r, err, "", getRespStatus(err))
		return
	}
	if hosts == nil {
		render.JSON(w, r, make([]dataprovider.DefenderEntry, 0))
		return
	}
	render.JSON(w, r, hosts)
}

func getDefenderHostByID(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	ip, err := getIPFromID(r)
	if err != nil {
		sendAPIResponse(w, r, err, "", http.StatusBadRequest)
		return
	}
	host, err := common.GetDefenderHost(ip)
	if err != nil {
		sendAPIResponse(w, r, err, "", getRespStatus(err))
		return
	}
	render.JSON(w, r, host)
}

func deleteDefenderHostByID(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	ip, err := getIPFromID(r)
	if err != nil {
		sendAPIResponse(w, r, err, "", http.StatusBadRequest)
		return
	}
	if !common.DeleteDefenderHost(ip) {
		sendAPIResponse(w, r, nil, "Not found", http.StatusNotFound)
		return
	}

	sendAPIResponse(w, r, nil, "OK", http.StatusOK)
}

func getIPFromID(r *http.Request) (string, error) {
	decoded, err := hex.DecodeString(getURLParam(r, "id"))
	if err != nil {
		return "", errors.New("invalid host id")
	}
	ip := string(decoded)
	err = validateIPAddress(ip)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func validateIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("ip address %#v is not valid", ip)
	}
	return nil
}
