package model

type Feature struct {
	Id          int64  `db:"id"          json:"id"`
	AreaId      int64  `db:"area_id"     json:"area-id"`
	Name        string `db:"name"        json:"name"`
	Description string `db:"description" json:"description"`
	Importance  string `db:"importance"  json:"importance"`
	Total       int64  `json:"total"`
	Passes      int64  `json:"passes"`
	Pending     int64  `json:"pending"`
	Failures    int64  `json:"failures"`
	Skipped     int64  `json:"skipped"`
	Tests       []Test `json:"tests"`
}

const CREATE_FEATURE = `CREATE TABLE IF NOT EXISTS features(
	id INT AUTO_INCREMENT PRIMARY KEY,
	area_id INT, 
	name VARCHAR(255),
	description TEXT,
	importance VARCHAR(255),
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (area_id) REFERENCES areas(id)
	)`

const INSERT_FEATURE = "INSERT INTO features(area_id, name, description, importance) VALUES (?,?,?,?)"

const SELECT_FEATURES = "SELECT id, area_id, name, description, importance FROM features WHERE area_id = ?;"

const UPDATE_FEATURE = "UPDATE features SET name = ?, description = ?, importance = ? WHERE id = ?"

const DELETE_FEATURE = "DELETE FROM features WHERE id = ?"
