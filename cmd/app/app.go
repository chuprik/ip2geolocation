package main

import (
	"github.com/chuprik/ip2geolocation/internal/maxmind"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oschwald/geoip2-golang"
	"github.com/unrolled/render"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	filePath, err := maxmind.FindDBLocation()

	if err != nil {
		licenseKey := os.Getenv("MAXMIND_LICENSE_KEY")
		if licenseKey == "" {
			log.Fatalln("MAXMIND_LICENSE_KEY environment variable is not set")
		}
		err := maxmind.Download("GeoLite2-City", licenseKey)
		if err != nil {
			log.Fatalln("cannot download maxmind database", err)
		}
		filePath, err = maxmind.FindDBLocation()
		if err != nil {
			log.Fatalln("cannot find maxmind database after download", err)
		}
	}

	db, err := geoip2.Open(filePath)
	if err != nil {
		log.Fatal("cannot open geoip database", err)
	}
	defer db.Close()

	rndr := render.New(render.Options{})

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/geoip/city/{ip}", func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(chi.URLParam(r, "ip"))
		record, err := db.City(ip)
		if err != nil {
			log.Println("cannot get geoip record", err)
			rndr.JSON(w, http.StatusBadRequest, map[string]string{"err": err.Error()})
			return
		}
		rndr.JSON(w, http.StatusOK, record)
	})
	r.Get("/geoip/country/{ip}", func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(chi.URLParam(r, "ip"))
		record, err := db.Country(ip)
		if err != nil {
			log.Println("cannot get geoip record", err)
			rndr.JSON(w, http.StatusBadRequest, map[string]string{"err": err.Error()})
			return
		}
		rndr.JSON(w, http.StatusOK, record)
	})
	http.ListenAndServe(":3000", r)
}
