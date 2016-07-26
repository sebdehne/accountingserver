package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"strconv"
	"github.com/sebdehne/accountingserver/domain"
	"encoding/json"
)

type CategoryDto struct {
	Name string `json:"name"`
}

type CategoryApi struct {
	store storage.Storage
}

func (cApi *CategoryApi) ListCategories(c *iris.Context) {
	root, err := cApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
	} else {
		c.SetHeader("ETag", strconv.Itoa(root.Version))
		c.JSON(200, root.Categories)
	}
}

func (cApi *CategoryApi) DeleteCategory(c *iris.Context) {
	// get the existing data
	root, err := cApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// validate the ETag header
	expectedVersion, err := strconv.Atoi(c.RequestHeader("ETag"))
	if err != nil {
		c.Error("Invalid ETag header", iris.StatusBadRequest)
		return
	}
	if expectedVersion != root.Version {
		c.Error("Invalid ETag header", iris.StatusConflict)
		return
	}

	categoriesInUse := root.GetTransactionsByCategory(domain.DateFilter{})
	if _, found := categoriesInUse[c.Param("id")]; found {
		c.Error("Category is currently in use", iris.StatusConflict)
		return
	} else {
		if root.RemoveCategory(c.Param("id")) {
			c.SetStatusCode(iris.StatusNoContent)
		} else {
			c.SetStatusCode(iris.StatusNotFound)
		}
	}

	root.Version++
	err = cApi.store.Save(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	c.SetHeader("ETag", strconv.Itoa(root.Version))
}

func (cApi *CategoryApi) PutCategory(c *iris.Context) {
	// try to unmarshall request body
	in := CategoryDto{}
	err := json.Unmarshal(c.Request.Body(), &in)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	inCat := domain.Category{Id:c.Param("id"), Name:in.Name}

	// get the existing data
	root, err := cApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// validate the ETag header
	expectedVersion, err := strconv.Atoi(c.RequestHeader("ETag"))
	if err != nil {
		c.Error("Invalid ETag header", iris.StatusBadRequest)
		return
	}
	if expectedVersion != root.Version {
		c.Error("Invalid ETag header", iris.StatusConflict)
		return
	}

	// all good, update the category now
	_, i, found := root.GetCategory(inCat.Id)
	if !found {
		root.Categories = append(root.Categories, inCat)
	} else {
		root.Categories[i] = inCat
	}
	root.Version++
	err = cApi.store.Save(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// prepare a response
	c.SetHeader("ETag", strconv.Itoa(root.Version))
	if found {
		c.JSON(iris.StatusOK, inCat)
	} else {
		c.JSON(iris.StatusCreated, inCat)
	}
}
