package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"github.com/sebdehne/accountingserver/domain"
	"strconv"
	"encoding/json"
)

type PartApi struct {
	store *storage.Storage
}

type PartDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (pApi *PartApi) DeleteParty(c *iris.Context) {
	// get the existing data
	root, err := pApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	id := c.Param("id")

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

	if root.IsPartyInUse(id) {
		c.Error("Party is currently in use", iris.StatusConflict)
		return
	}

	if !root.RemoveParty(id) {
		c.SetStatusCode(iris.StatusNotFound)
		return
	}

	root.Version++
	err = pApi.store.Save(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetStatusCode(iris.StatusNoContent)
	c.SetHeader("ETag", strconv.Itoa(root.Version))

}

func (pApi *PartApi) PutParty(c *iris.Context) {
	// try to unmarshall request body
	in := PartDto{}
	err := json.Unmarshal(c.Request.Body(), &in)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	in.Id = c.Param("id")
	inPart := domain.Party{Id:in.Id, Name:in.Name}

	// get the existing data
	root, err := pApi.store.Get()
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

	// all good, update the party now
	_, i, found := root.GetPart(inPart.Id)
	if !found {
		root.Parties = append(root.Parties, inPart)
	} else {
		root.Parties[i] = inPart
	}
	root.Version++
	err = pApi.store.Save(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// prepare a response
	c.SetHeader("ETag", strconv.Itoa(root.Version))
	if found {
		c.JSON(iris.StatusOK, MapParty(inPart))
	} else {
		c.JSON(iris.StatusCreated, MapParty(inPart))
	}
}

func (pApi *PartApi) ListParties(c *iris.Context) {
	root, err := pApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
	} else {
		c.SetHeader("ETag", strconv.Itoa(root.Version))
		c.JSON(iris.StatusOK, MapParties(root.Parties))
	}
}

func MapParties(in []domain.Party) []PartDto {
	result := make([]PartDto, 0)
	for _, p := range in {
		result = append(result, MapParty(p))
	}
	return result
}
func MapParty(in domain.Party) PartDto {
	return PartDto{Id:in.Id, Name:in.Name}
}
