package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {

	err := a.initialize(dbuser, dbpwd, "test")
	if err != nil {
		log.Fatal("Error occured while initialising the database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    quantity int,
	price float(10,7),
    PRIMARY KEY (id)
	);`

	_, err := a.db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.db.Exec("DELETE from products")
	a.db.Exec("ALTER table products AUTO_INCREMENT=1")
	log.Println("clearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chairs", "quantity":1, "price":200}`)
	req, _ := http.NewRequest("POST", "/addnewproduct", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Name"] != "chairs" {
		t.Errorf("Expected name: %v, Got: %v", "cybertruck", m["Name"])
	}
	if m["Quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, Got: %v", 1.0, m["Quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/products/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/products/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(req)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name":"connector", "quantity":90, "price":10}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["Id"] != newValue["Id"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["Id"], oldValue["Id"])
	}

	if oldValue["Name"] != newValue["Name"] {
		t.Errorf("Expected name: %v, Got: %v", newValue["Name"], oldValue["Name"])
	}

	if oldValue["Price"] != newValue["Price"] {
		t.Errorf("Expected price: %v, Got: %v", newValue["Price"], oldValue["Price"])
	}

	if oldValue["Quantity"] == newValue["Quantity"] {
		t.Errorf("Expected quantity: %v, Got: %v", newValue["Quantity"], oldValue["Quantity"])
	}
}
