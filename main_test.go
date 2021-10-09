package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPost(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/6161b199c19104130f7e82bd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"_id":"6161b199c19104130f7e82bd","Caption":"dis","ImageUrl":"sheesh","Timestamp":"2021-10-09 20:43:29.7963941 +0530 IST m=+327.985190701","User_id":"6161b45d4d008f1382ee45ce"}]`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreatePost(t *testing.T) {
	var jsonStr = []byte(`{"Caption":"testpost","Image":"test.png","User_id":"6161b45d4d008f1382ee45ce"}`)
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"InsertedID":"61617b6a91e6558a630edc60"}`
	if len(strings.TrimSpace(rr.Body.String())) != len(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"Name":"testuser","Email":"testuser.com","Password":"testpass"}`)
	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"InsertedID":"61617eaf8d9c98800c05de5f"}`
	if len(strings.TrimSpace(rr.Body.String())) != len(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/6161b45d4d008f1382ee45ce", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"_id":"6161b45d4d008f1382ee45ce","Name":"siddhartha","Email":"crazysid278@gmail.com","Password":"49472fa2771466002c0afa5f954169a6a6a45118132a1309fdd207743ace9165518b075f166f52"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostsByUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/user/6161b45d4d008f1382ee45ce?limit=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserPost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"_id":"6161b199c19104130f7e82bd","Caption":"dis","ImageUrl":"sheesh","Timestamp":"2021-10-09 20:43:29.7963941 +0530 IST m=+327.985190701","User_id":"6161b45d4d008f1382ee45ce"},{"_id":"6161b65ce66ad60ff55893aa","Caption":"Caption goes here","ImageUrl":"photobomb","Timestamp":"2021-10-09 21:03:48.351589 +0530 IST m=+217.869099901","User_id":"6161b45d4d008f1382ee45ce"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
