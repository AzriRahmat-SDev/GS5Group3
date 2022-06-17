package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	var p Plot
	PopulateData(OpenVenueDB())
	p.Address = "Aljuinied Road, Happy Garden Estate, 389842"
	p.VenueName = "Aljunied Park"
	m := map[string]Plot{"ALJ027": p}

	jsValue, _ := json.Marshal(m)
	fmt.Println(m)
	request, _ := http.NewRequest(http.MethodPost, baseURL+"/plots/ALJ027", bytes.NewBuffer(jsValue))
	client := &http.Client{}

	resp, _ := client.Do(request)
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &p)
	fmt.Println("\nDATA =", p, m)
	resp.Body.Close()
}
