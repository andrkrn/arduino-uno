package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
)

var (
	humidity = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "arduino",
		Name:      "humidity",
	})
	temperature = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "arduino",
		Name:      "temperature",
	})
	// ops = promauto.NewCounter(prometheus.CounterOpts{
	// 	Namespace: "arduino",
	// 	Name:      "ops",
	// })
)

func recordMetrics() {
	go func() {
		config := &serial.Config{
			Name:        "/dev/cu.usbserial-1410",
			Baud:        9600,
			ReadTimeout: time.Second * 5,
			Size:        0,
		}

		stream, err := serial.OpenPort(config)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println(scanner.Text())
			s := strings.Fields(text)
			h := strings.Split(s[0], "=")[1]
			t := strings.Split(s[1], "=")[1]

			tfloat, _ := strconv.ParseFloat(t, 32)
			hfloat, _ := strconv.ParseFloat(h, 32)

			temperature.Set(tfloat)
			humidity.Set(hfloat)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
}

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":2112", nil); err != nil {
		panic(err)
	}
}
