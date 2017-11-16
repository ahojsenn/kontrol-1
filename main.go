package main

import (
	"bitbucket.org/rwirdemann/kontrol/rest"

	"bitbucket.org/rwirdemann/kontrol/processing"

	"bitbucket.org/rwirdemann/kontrol/parser"
)

func main() {
	bookings := parser.Import("2017-Buchungen-KG - Buchungen 2017.csv")
	for _, p := range bookings {
		processing.Process(p)
	}


	rest.StartService()
}
