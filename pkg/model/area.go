/*
Copyright (c) 2022, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

type Area struct {
	Id         int64   `json:"id"           db:"id"`
	ProductId  int64   `json:"product-id"   db:"product_id"`
	Name       string  `json:"name"         db:"name"`
	Total      int64   `json:"total"`
	FirstTotal int64   `json:"first-total"`
	Passes     int64   `json:"passes"`
	Pending    int64   `json:"pending"`
	Failures   int64   `json:"failures"`
	Skipped    int64   `json:"skipped"`
	ExplTests  int64   `json:"expl-tests"`
	ExplRating float64 `json:"expl-rating"`
}
