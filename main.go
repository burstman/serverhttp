// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 195.

// Http4 is an e-commerce server that registers the /list and /price
// endpoint by calling http.HandleFunc.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gopl.io/ch7/http4/htmllist"
)

const tpl = `
<html>
<body>

<table>
	<tr>
		<th>item</th>
		<th>price</a></th>
	</tr>
{{range .}}
	<tr>
		<td>{{.Item}}</td>
		<td>{{.Price}}</td>
	</td>
{{end}}
</body>
</html>`

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//!-main

type dollars float64

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

//Data that will be transferred to template list
type Data struct {
	Item  string
	Price dollars
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	var d []Data
	for item, price := range db {
		//fmt.Fprintf(w, "%s: %s\n", item, price)
		d = append(d, Data{item, price})
	}
	htmllist.Templist(w, d, tpl, "listtemplate")

}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
func (db database) create(w http.ResponseWriter, r *http.Request) {

	item := r.FormValue("item")
	pricestr := r.FormValue("price")
	if item == "" {
		http.Error(w, "please input a name for the item", http.StatusBadRequest)
		return
	}
	if _, ok := db[item]; ok {
		http.Error(w, "item already existe in database", http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(pricestr, 2)
	if err != nil {
		fmt.Fprintf(w, "please input a number for the item: %q\n", item)
		return
	}
	if db == nil {
		db = make(map[string]dollars, 0)
	}
	db[item] = dollars(price)

}
func (db database) update(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")
	pricestr := r.FormValue("price")
	if item == "" {
		http.Error(w, "please input a name for the item", http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(pricestr, 2)
	if err != nil {
		fmt.Fprintf(w, "please input a number for the item: %q\n", item)
		return
	}
	if _, ok := db[item]; ok {
		db[item] = dollars(price)
	} else {
		http.Error(w, "item not found", http.StatusBadRequest)
	}

}
func (db database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")
	if _, ok := db[item]; ok {
		delete(db, item)
	} else {
		http.Error(w, "item not found", http.StatusBadRequest)
	}
}
