package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	db     *sql.DB
}

func checkerror(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func (app *App) initialize(dbuser string, dbpwd string, dbname string) error {
	datasourcename := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", dbuser, dbpwd, dbname)
	var err error
	app.db, err = sql.Open("mysql", datasourcename)
	checkerror(err)

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handlerequst()
	return nil
}

func (app *App) run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendresponse(w http.ResponseWriter, status_code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status_code)
	w.Write(response)
}

func senderror(w http.ResponseWriter, statuscode int, err string) {
	error_message := map[string]string{"error": err}
	sendresponse(w, statuscode, error_message)
}

func (app *App) mainpage(w http.ResponseWriter, r *http.Request) {
	products, err := Getproducts(app.db)
	if err != nil {
		senderror(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendresponse(w, http.StatusOK, products)
}

func (app *App) findspecificID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slice_urlpath, err := strconv.Atoi(vars["ID"])
	if err != nil {
		senderror(w, http.StatusBadRequest, "invalid product id")
		return
	}
	prod := product{Id: slice_urlpath}
	err = prod.GetProduct(app.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			senderror(w, http.StatusNotFound, "product not found")
		default:
			senderror(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendresponse(w, http.StatusOK, prod)
}

func (app *App) createproduct(w http.ResponseWriter, r *http.Request) {
	var produ product
	err := json.NewDecoder(r.Body).Decode(&produ)
	if err != nil {
		senderror(w, http.StatusBadRequest, err.Error())
		return
	}
	err = produ.CreateProduct(app.db)
	if err != nil {
		senderror(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendresponse(w, http.StatusCreated, produ)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slice_urlpath, err := strconv.Atoi(vars["ID"])
	if err != nil {
		senderror(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	var p product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		senderror(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.Id = slice_urlpath
	err = p.UpdateProduct(app.db)
	if err != nil {
		senderror(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendresponse(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slice_urlpath, err := strconv.Atoi(vars["ID"])
	if err != nil {
		senderror(w, http.StatusBadRequest, "invalid product ID")
		return
	}
	p := product{Id: slice_urlpath}
	err = p.DeleteProduct(app.db)
	if err != nil {
		senderror(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendresponse(w, http.StatusOK, map[string]string{"result": "successful deletion"})
}

func (app *App) handlerequst() {
	app.Router.HandleFunc("/products", app.mainpage).Methods("GET")
	app.Router.HandleFunc("/products/{ID}", app.findspecificID).Methods("GET")
	app.Router.HandleFunc("/addnewproduct", app.createproduct).Methods("POST")
	app.Router.HandleFunc("/products/{ID}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/products/{ID}", app.deleteProduct).Methods("DELETE")
}
