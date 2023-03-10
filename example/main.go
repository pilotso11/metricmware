package main

/*
 * Copyright (c) 2023. Seth Osher.  All Rights Reserved.
 * MIT License
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import (
	"expvar"
	"flag"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	expvarmw "github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/pilotso11/metricmware"
	"github.com/zserge/metric"
)

var metrics = make(map[string]metric.Metric)

func main() {
	useExpvar := flag.Bool("e", false, "Use expvar")
	flag.Parse()

	app := fiber.New()

	// Register the middleware
	app.Use(metricmware.New(metricmware.Config{Exposed: &metrics}))
	if *useExpvar {
		app.Use(expvarmw.New())
	}

	// Add my static page
	app.Static("/", "./public")

	// Create the metrics
	if *useExpvar {
		expvar.Publish("my_counter", metric.NewCounter("5m1s", "15m30s", "1h1m"))
		expvar.Publish("my_stat", metric.NewGauge("5m1s", "15m30s", "1h1m"))
		expvar.Publish("my_latency", metric.NewHistogram("5m1s", "15m30s", "1h1m"))
	} else {
		metrics["my_counter"] = metric.NewCounter("5m1s", "15m30s", "1h1m")
		metrics["my_stat"] = metric.NewGauge("5m1s", "15m30s", "1h1m")
		metrics["my_latency"] = metric.NewHistogram("5m1s", "15m30s", "1h1m")

	}

	// Start the random generator
	go randomStatsGenerator(*useExpvar)

	_ = app.Listen("127.0.0.1:8000")
}

func randomStatsGenerator(useExpVar bool) {
	for {
		start := time.Now()
		delay := rand.Intn(500 * int(time.Millisecond))
		time.Sleep(time.Duration(delay))
		increment := rand.Float64() * 50

		if useExpVar {
			expvar.Get("my_counter").(metric.Metric).Add(1)                           // Increase counter each loop
			expvar.Get("my_stat").(metric.Metric).Add(increment)                      // Increase counter each loop
			expvar.Get("my_latency").(metric.Metric).Add(time.Since(start).Seconds()) // Increase counter each loop
		} else {
			metrics["my_counter"].Add(1)                           // Increase counter each loop
			metrics["my_stat"].Add(increment)                      // Increase counter each loop
			metrics["my_latency"].Add(time.Since(start).Seconds()) // Increase counter each loop
		}
	}
}
