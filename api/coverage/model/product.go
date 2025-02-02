/*
Copyright (c) 2022-2025, webmaster@testandwin.net, Michael Schlottmann
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package model

type Product struct {
	Id   int64  `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}
