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
const MaxLimit = 10000

func ExtractPageFilter(c *iris.Context) domain.PageFilter {
	offset, limit := 0, DefaultLimit

	if i, found, _ := paramToInt64(c, "offset"); found && i >= 0 {
		offset = int(i)
	}
	if i, found, _ := paramToInt64(c, "limit"); found && i >= 0 && i < MaxLimit {
		limit = int(i)
	}

	return domain.PageFilter{Offset:offset, Limit:limit}
}

func ExtractDateFilter(c *iris.Context) (result domain.DateFilter, err error) {
	if i, found, e := paramToInt64(c, "from_date"); e != nil {
		err = e
	} else if found && i > 0 {
		result.FromDate = &i

		if i, found, e := paramToInt64(c, "to_date"); e != nil {
			err = e
		} else if found && i > 0 {
			result.ToDate = &i
		}
	}
	return
}

func paramToInt64(c *iris.Context, paramKey string) (result int64, found bool, err error) {
	paramValue := string(c.QueryArgs().Peek(paramKey))
	if len(paramValue) > 0 {
		result, err = strconv.ParseInt(paramValue, 10, 0)
		if err == nil {
			found = true
		}
	}
	return
}