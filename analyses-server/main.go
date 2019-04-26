/*
 * Copyright (c) 2019 XLAB d.o.o
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/private-predictions/analyses-server/framingham"
	"github.com/fentec-project/private-predictions/serialization"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/framingham", framingham.Routes())
	})

	return router
}

func DeriveKey() {
	r2 := big.NewInt(34362)
	r3 := big.NewInt(263588)
	r4 := big.NewInt(188030)
	r5 := big.NewInt(112673)
	r6 := big.NewInt(-90941)
	r7 := big.NewInt(59397)
	r8 := big.NewInt(52320)
	r9 := big.NewInt(68602)
	y1 := data.NewVector([]*big.Int{r2, r3, r4, r5, r6, r7, r8, r9})

	jsonValue, _ := json.Marshal(y1)
	response, err := http.Post("http://localhost:8080/v1/api/paillier",
		"application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Accessing response body failed with error %s\n", err)
	}
	key1 := string(data1)

	t2 := big.NewInt(48123)
	t3 := big.NewInt(339222)
	t4 := big.NewInt(139862)
	t5 := big.NewInt(-439)
	t6 := big.NewInt(16081)
	t7 := big.NewInt(99858)
	t8 := big.NewInt(19035)
	t9 := big.NewInt(49756)
	y2 := data.NewVector([]*big.Int{t2, t3, t4, t5, t6, t7, t8, t9})

	jsonValue, _ = json.Marshal(y2)
	response, err = http.Post("http://localhost:8080/v1/api/paillier",
		"application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data2, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Accessing response body failed with error %s\n", err)
	}
	key2 := string(data2)

	key1 = strings.TrimSpace(key1)
	key2 = strings.TrimSpace(key2)

	key1Int, _ := new(big.Int).SetString(key1, 10)
	key2Int, _ := new(big.Int).SetString(key2, 10)

	serialization.WriteGob("framingham30-FE-y1-key.gob", key1Int)
	serialization.WriteGob("framingham30-FE-y2-key.gob", key2Int)
}

func main() {
	_, err := os.Stat("framingham30-FE-y1-key.gob")
	if os.IsNotExist(err) {
		DeriveKey()
	}

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":8081", router))
}
