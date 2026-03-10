package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elliotchance/phpserialize"
)

func main() {
	payload := `a:10:{s:2:"id";i:1429;s:10:"salutation";s:3:"Mr.";s:11:"designation";a:3:{s:2:"id";i:4;s:4:"name";s:11:"SUPER ADMIN";s:9:"is_active";b:1;}s:4:"name";s:6:"Kartik";s:13:"mobile_number";s:10:"9316909534";s:5:"email";s:19:"kartik.b@medkart.in";s:13:"employee_code";s:4:"2199";s:9:"is_active";b:1;s:8:"pass_key";s:5:"42712";s:10:"last_login";O:25:"Illuminate\Support\Carbon":3:{s:4:"date";s:26:"2026-02-25 10:46:20.000000";s:13:"timezone_type";i:3;s:8:"timezone";s:12:"Asia/Kolkata";}}`

	result, err := phpserialize.UnmarshalAssociativeArray([]byte(payload))
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("JSON Marshal error: %v\n", err)
	}

	fmt.Println(string(data))
}
