package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Endereco struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

func getEnderecoCEP(cep string) (Endereco, error) {
	var endereco Endereco

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return endereco, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return endereco, err
	}

	err = json.Unmarshal(body, &endereco)
	if err != nil {
		return endereco, err
	}

	return endereco, nil
}

func buscarEnderecoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]

	endereco, err := getEnderecoCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endereco)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cep/{cep}", buscarEnderecoHandler).Methods("GET")

	fmt.Println("Servidor iniciado na porta 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
