// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignHandlerNoData(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/sign", nil)
	http.HandlerFunc(SignHandler).ServeHTTP(w, r)

	// Check the JSON response
	result := SignResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signed response: %v", err)
	}
	if result.Success {
		t.Error("Expected an error, got success response")
	}
}

func TestSignHandler(t *testing.T) {
	const assertions = `
  {
	  "brand-id": "System",
    "model":"聖誕快樂",
    "serial":"A1234/L",
		"revision": 2,
    "device-key":"ssh-rsa NNhqloxPyIYXiTP+3JTPWV/mNoBar2geWIf"
  }`

	Config = &ConfigSettings{PrivateKeyPath: "../TestKey.asc"}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/1.0/sign", bytes.NewBufferString(assertions))
	http.HandlerFunc(SignHandler).ServeHTTP(w, r)

	// Check the JSON response
	result := SignResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signed response: %v", err)
	}
	if !result.Success {
		t.Errorf("Error generated in signing the device: %s", result.ErrorMessage)
	}
	if result.Signature == "" {
		t.Errorf("Empty signed data returned.")
	}
}

func TestVersionHandler(t *testing.T) {

	Config = &ConfigSettings{Version: "1.2.5"}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/1.0/version", nil)
	http.HandlerFunc(VersionHandler).ServeHTTP(w, r)

	// Check the JSON response
	result := VersionResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the version response: %v", err)
	}
	if result.Version != Config.Version {
		t.Errorf("Incorrect version returned. Expected '%s' got: %v", Config.Version, result.Version)
	}

}
