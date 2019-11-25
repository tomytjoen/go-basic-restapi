package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// function untuk mencetak data berformat JSON
func SetJSONResp(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func main() {

	//buat struct Product
	type Product struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	}
	//buat varibel Product dengan isi banyak(array)
	database := make(map[string]Product)
	// isi data variabelnya
	database["001"] = Product{ID: "001", Name: "Samsung Galaxy S1", Quantity: 10}
	database["002"] = Product{ID: "002", Name: "Samsung Galaxy S2", Quantity: 15}

	// handle url kalau di home
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		message := []byte(`{"message":"server is up"}`)
		SetJSONResp(res, message, http.StatusOK)
	})

	//handle url untuk mendapatkan semua data
	http.HandleFunc("/get-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`"message":"invalid http method"`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var products []Product

		for _, product := range database {
			products = append(products, product)
		}

		productJSON, err := json.Marshal(&products)

		if err != nil {
			message := []byte(`{"message":"Error While Parsing Data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)
	})

	// untuk mendapatkan 1 url product
	http.HandleFunc("/get-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`"message":"invalid http method"`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"Required Product ID"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]

		if !ok {
			message := []byte(`{"message":"Product Not Found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		productJSON, err := json.Marshal(&product)

		if err != nil {
			message := []byte(`{"message":"Error While Parsing Data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)
	})

	// url untuk hapus products
	http.HandleFunc("/delete-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			message := []byte(`"message":"invalid http method"`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"Required Product ID"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]

		if !ok {
			message := []byte(`{"message":"Product Not Found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		delete(database, id)

		productJSON, err := json.Marshal(&product)

		if err != nil {
			message := []byte(`{"message":"Error While Parsing Data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)
	})

	// menambah data product
	http.HandleFunc("/add-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`"message":"invalid http method"`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var product Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&product)

		if err != nil {
			message := []byte(`{"message":"Error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		database[product.ID] = product

		message := []byte(`{"message":"Success Create Product"}`)
		SetJSONResp(res, message, http.StatusCreated)
	})

	// mengupdate/merubah data product
	http.HandleFunc("/update-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			message := []byte(`"message":"invalid http method"`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"Required Product ID"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]

		if !ok {
			message := []byte(`{"message":"Product Not Found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		var newProduct Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&newProduct)

		if err != nil {
			message := []byte(`{"message":"Error when parsing products"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		product.Name = newProduct.Name
		product.Quantity = newProduct.Quantity

		database[product.ID] = product

		productJSON, err := json.Marshal(&product)

		if err != nil {
			message := []byte(`{"message":"Error While Parsing Products"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)
	})

	// membuat script dapat di handle/request dgn port 9000
	err := http.ListenAndServe(":9000", nil)

	// kalau ada problem, cetak problemnya
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
