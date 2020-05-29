package util

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type CommonSubrouter struct {
	Router *mux.Router
	Db     *gorm.DB
}
