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
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/fentec-project/gofe/innerprod/fullysec"
	"github.com/fentec-project/private-predictions/key-server/keys"
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
		r.Mount("/api", keys.Routes())
	})

	return router
}

func GenerateMasterKeys() {
	l := 8
	boundX := new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)
	boundY := new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)

	bitLength := 512
	lambda := 128

	paillier, err := fullysec.NewPaillier(l, lambda, bitLength, boundX, boundY)
	if err != nil {
		fmt.Errorf("Error during simple inner product creation: %v", err)
	}

	masterSecKey, masterPubKey, err := paillier.GenerateMasterKeys()
	if err != nil {
		fmt.Errorf("Error during master key generation: %v", err)
	}

	serialization.WriteGob("secKey.gob", masterSecKey)
	serialization.WriteGob("pubKey.gob", masterPubKey)

	serialization.WriteGob("paillier.gob", paillier.Params)
}

func main() {
	_, err := os.Stat("paillier.gob")
	if os.IsNotExist(err) {
		GenerateMasterKeys()
	}

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
