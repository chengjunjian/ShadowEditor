package category

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/tengge1/shadoweditor/helper"
	"github.com/tengge1/shadoweditor/server"
)

func init() {
	category := Category{}
	server.Mux.UsingContext().Handle(http.MethodGet, "/api/Category/List", category.List)
	server.Mux.UsingContext().Handle(http.MethodPost, "/api/Category/Save", category.Save)
}

// Category 类别控制器
type Category struct {
}

// List 获取列表
func (Category) List(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typ := strings.TrimSpace(r.FormValue("type"))

	db, err := server.Mongo()
	if err != nil {
		helper.WriteJSON(w, err.Error())
		return
	}

	docs := []Model{}

	if server.Config.Authority.Enabled {
		user, _ := server.GetCurrentUser(r)

		if user != nil {
			filter1 := bson.M{
				"UserID": user.ID,
			}

			if user.Name == "Administrator" {
				filter1 = bson.M{
					"$or": bson.A{
						filter1,
						bson.M{
							"UserID": bson.M{
								"$exists": 0,
							},
						},
					},
				}
			}

			if typ != "" {
				filter1 = bson.M{
					"$and": bson.A{
						filter1,
						bson.M{
							"Type": typ,
						},
					},
				}
			}
			db.FindMany(server.CategoryCollectionName, filter1, &docs)
		}
	} else {
		if typ != "" {
			filter1 := bson.M{
				"Type": typ,
			}
			db.FindMany(server.CategoryCollectionName, filter1, &docs)
		} else {
			db.FindAll(server.CategoryCollectionName, &docs)
		}
	}

	list := []map[string]string{}

	for _, i := range docs {
		obj := map[string]string{
			"ID":   i.ID,
			"Name": i.Name,
		}
		list = append(list, obj)
	}

	helper.WriteJSON(w, server.Result{
		Code: 200,
		Msg:  "Get Successfully!",
		Data: list,
	})
}

// Save 保存
func (Category) Save(w http.ResponseWriter, r *http.Request) {
}