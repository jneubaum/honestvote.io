package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jneubaum/honestvote/core/core-crypto/crypto"
	"github.com/jneubaum/honestvote/core/core-database/database"
	corehttp "github.com/jneubaum/honestvote/core/core-http/http"
)

var CollectionPrefix = "test_"
var MongoConnection = database.MongoConnect("localhost")

func TestPostRequestAdminPrivileges(t *testing.T) {
	url := "http://127.0.0.1:7003"

	private_key, public_key := crypto.GenerateKeyPair()
	message := []byte("requesting administrator privileges")
	signature, err := crypto.Sign(message, private_key)
	if err != nil {
		return
	}

	var request database.RequestAdminPrivileges = database.RequestAdminPrivileges{
		PublicKey:   public_key,
		Domain:      "bizylife.com",
		Institution: "BizyLife",
		Signature:   signature,
		Message:     message,
	}

	json_request, err := json.Marshal(&request)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_request))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	database.CollectionPrefix = CollectionPrefix
	database.MongoDB = MongoConnection

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(corehttp.PostVoteHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//error because it must write a message to the administrator
func TestPostRegisterHandler(t *testing.T) {
	// jsonRegistration, _ := json.Marshal(database.AwaitingRegistration{})
	// req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonRegistration))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// database.CollectionPrefix = CollectionPrefix
	// database.MongoDB = MongoConnection

	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(corehttp.PostRegisterHandler)

	// handler.ServeHTTP(rr, req)

}

func TestPostVoteHandler(t *testing.T) {

}

func TestPostElectionHandler(t *testing.T) {

}

func TestGetElectionsHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	database.CollectionPrefix = CollectionPrefix
	database.MongoDB = MongoConnection

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(corehttp.GetElectionsHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	response := corehttp.Payload{}
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	expected := "OK"
	if response.Status != expected {
		t.Errorf("handler returned unexpected status: got %v want %v",
			response.Status, expected)
	}

	responseData, _ := json.Marshal(&response.Data)
	elections := []database.ElectionInfo{}
	json.Unmarshal(responseData, &elections)

	// expectedElections := []database.ElectionInfo{
	// 	database.ElectionInfo{
	// 		ElectionName: "Vote for Charity",
	// 		Institution:  "Honestvote",
	// 		Description:  "Whichever charities get the most votes, will be donated $50 each by Honestvote",
	// 		Start:        "Fri, 27 Feb 2020 08:00:00 EST",
	// 		End:          "Fri, 27 Feb 2020 22:00:00 EST",
	// 		Signature:    "3045022034466fa37fac0368c342705c342bda5e381a8ad92b0209161ca7dc310dfcef8b022100c6249adc9a0d690c4d4b76f32c8fbe226d77e5965d6ed6b234a5733ba76d3504",
	// 	},
	// }
	// if !reflect.DeepEqual(elections, expectedElections) {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		responseData, expectedElections)
	// }

}

func TestGetElectionHandler(t *testing.T) {
	// req, err := http.NewRequest("GET", "/election/{electionid:3045022034466fa37fac0368c342705c342bda5e381a8ad92b0209161ca7dc310dfcef8b022100c6249adc9a0d690c4d4b76f32c8fbe226d77e5965d6ed6b234a5733ba76d3504}", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// database.CollectionPrefix = CollectionPrefix
	// database.MongoDB = MongoConnection

	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(corehttp.GetElectionHandler)

	// handler.ServeHTTP(rr, req)

	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// response := corehttp.Payload{}
	// _ = json.Unmarshal(rr.Body.Bytes(), &response)

	// expected := "OK"
	// if response.Status != expected {
	// 	t.Errorf("handler returned unexpected status: got %v want %v",
	// 		response.Status, expected)
	// }

	// responseData, _ := json.Marshal(&response.Data)
	// elections := []database.ElectionInfo{}
	// json.Unmarshal(responseData, &elections)

	// expectedElections := []database.ElectionInfo{
	// 	database.ElectionInfo{
	// 		ElectionName: "Vote for Charity",
	// 		Institution:  "Honestvote",
	// 		Description:  "Whichever charities get the most votes, will be donated $50 each by Honestvote",
	// 		Start:        "Fri, 27 Feb 2020 08:00:00 EST",
	// 		End:          "Fri, 27 Feb 2020 22:00:00 EST",
	// 		Signature:    "3045022034466fa37fac0368c342705c342bda5e381a8ad92b0209161ca7dc310dfcef8b022100c6249adc9a0d690c4d4b76f32c8fbe226d77e5965d6ed6b234a5733ba76d3504",
	// 	},
	// }
	// if !reflect.DeepEqual(elections, expectedElections) {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		responseData, expectedElections)
	// }
}

func TestGetVotesHandler(t *testing.T) {

}
func TestGetPositionsHandler(t *testing.T) {

}

func TestGetPermissions(t *testing.T) {

}
