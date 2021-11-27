package crud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"hack-change-api/db"
	"hack-change-api/models/entities"
	u "hack-change-api/muxutil"
	"net/http"
	"strconv"
)

var ThreadCommentCreate = func(w http.ResponseWriter, r *http.Request) {
	ThreadComment := &entities.ThreadComment{}
	err := json.NewDecoder(r.Body).Decode(ThreadComment)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	db := db.GetDB()
	err = db.Create(ThreadComment).Error

	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		res, _ := json.Marshal(ThreadComment)
		u.RespondJSON(w, res)
	}
}

var ThreadCommentRetrieve = func(w http.ResponseWriter, r *http.Request) {
	ThreadComment := &entities.ThreadComment{}

	params := mux.Vars(r)
	id := params["id"]

	db := db.GetDB()
	err := db.Preload("Author").First(&ThreadComment, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleInternalError(w, err)
		}
		return
	}

	res, err := json.Marshal(ThreadComment)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.RespondJSON(w, res)
	}
}

var ThreadCommentUpdate = func(w http.ResponseWriter, r *http.Request) {
	ThreadComment := &entities.ThreadComment{}

	params := mux.Vars(r)
	id := params["id"]

	db := db.GetDB()
	err := db.First(&ThreadComment, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.HandleNotFound(w)
		} else {
			u.HandleInternalError(w, err)
		}
		return
	}

	newThreadComment := &entities.ThreadComment{}
	err = json.NewDecoder(r.Body).Decode(newThreadComment)

	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	err = db.Model(&ThreadComment).Updates(newThreadComment).Error

	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.Respond(w, u.Message(true, "OK"))
	}
}

var ThreadCommentDelete = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db := db.GetDB()
	err := db.Delete(&entities.ThreadComment{}, id).Error

	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		u.Respond(w, u.Message(true, "OK"))
	}
}

var ThreadCommentQuery = func(w http.ResponseWriter, r *http.Request) {
	var agroModels []entities.ThreadComment
	var count string

	order := r.FormValue("_order")
	sort := r.FormValue("_sort")
	end, err1 := strconv.Atoi(r.FormValue("_end"))
	start, err2 := strconv.Atoi(r.FormValue("_start"))

	if err1 != nil || err2 != nil {
		u.HandleBadRequest(w, errors.New("bad _start or _end parameter value"))
		return
	}
	u.CheckOrderAndSortParams(&order, &sort)

	db := db.GetDB()
	err := db.Preload("Author").Order(fmt.Sprintf("%s %s", sort, order)).
		Offset(start).Limit(end - start).Find(&agroModels).Error

	if err != nil {
		u.HandleInternalError(w, err)
		return
	}

	res, err := json.Marshal(agroModels)

	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		db.Model(&entities.ThreadComment{}).Count(&count)
		u.SetTotalCountHeader(w, count)
		u.RespondJSON(w, res)
	}
}
