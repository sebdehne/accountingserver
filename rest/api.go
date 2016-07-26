package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"github.com/sebdehne/accountingserver/domain"
	"strconv"
)

func New(s *storage.Storage) RestApi {
	return RestApi{s, &CategoryApi{s}, &AccountApi{s}, &PartApi{s}, &TransactionApi{s}}
}

type RestApi struct {
	store          *storage.Storage
	CategoryApi    *CategoryApi
	AccountApi     *AccountApi
	PartApi        *PartApi
	TransactionApi *TransactionApi
}

const DefaultLimit = 100
const MaxLimit = 1000

func ExtractPageFilter(c *iris.Context) domain.PageFilter {
	offset, limit := 0, DefaultLimit

	if i, found, _ := paramToDate(c, "offset"); found && i >= 0 {
		offset = int(i)
	}
	if i, found, _ := paramToDate(c, "limit"); found && i >= 0 && i < MaxLimit {
		limit = int(i)
	}

	return domain.PageFilter{Offset:offset, Limit:limit}
}

func ExtractDateFilter(c *iris.Context) (result domain.DateFilter, err error) {
	if i, found, e := paramToDate(c, "from_date"); e != nil {
		err = e
	} else if found && i > 0 {
		result.FromDate = &i

		if i, found, e := paramToDate(c, "to_date"); e != nil {
			err = e
		} else if found && i > 0 {
			result.ToDate = &i
		}
	}
	return
}

func paramToDate(c *iris.Context, paramKey string) (result int64, found bool, err error) {
	paramValue := c.Param(paramKey)
	if len(paramValue) > 0 {
		result, err = strconv.ParseInt(paramValue, 10, 0)
		if err == nil {
			found = true
		}
	}
	return
}