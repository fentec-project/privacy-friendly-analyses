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

package keys

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/fullysec"
	"github.com/fentec-project/private-predictions/serialization"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/paillier", DerivePaillierKey)
	return router
}

func DerivePaillierKey(w http.ResponseWriter, r *http.Request) {
	params := new(fullysec.PaillierParams)
	err := serialization.ReadGob("paillier.gob", params)
	if err != nil {
		fmt.Errorf("Error during Paillier params reading: %v", err)
	}
	paillier := fullysec.NewPaillierFromParams(params)

	masterSecKey := new(data.Vector)
	err = serialization.ReadGob("secKey.gob", masterSecKey)
	if err != nil {
		fmt.Errorf("Error during key reading: %v", err)
	}

	y1 := new(data.Vector)
	err = json.NewDecoder(r.Body).Decode(&y1)
	if err != nil {
		panic(err)
	}

	key1, err := paillier.DeriveKey(*masterSecKey, *y1)
	if err != nil {
		fmt.Errorf("Error during key derivation: %v", err)
	}

	render.JSON(w, r, key1)
}
