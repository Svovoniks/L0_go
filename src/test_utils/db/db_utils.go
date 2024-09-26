package test_utils_db

import (
	_db "l0/types/db"
)



func DBCleanup(db *_db.DB) {
	db.Db.Exec(`delete from "orders"`)
}
