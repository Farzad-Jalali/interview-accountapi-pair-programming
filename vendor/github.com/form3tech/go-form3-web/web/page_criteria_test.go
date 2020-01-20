package web

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetPageNumber_GetLastPageWhenRecordCountIsOddReturnsTheLastPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("last", "2", 3)

	assert.Equal(1, pageNumber)
}

func Test_GetPageNumber_GetLastPageWhenPageSizeIsGreaterThanRecordCountReturnsTheFirstPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("last", "10", 2)

	assert.Equal(0, pageNumber)
}

func Test_GetPageNumber_GetLastPageWhenRecordCountIsEvenReturnsTheLastPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("last", "2", 4)

	assert.Equal(1, pageNumber)
}

func Test_GetPageNumber_GetLastPageWhenPageSizeIsZeroReturnsTheLastPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("last", "0", 4)

	assert.Equal(0, pageNumber)
}

func Test_GetPageNumber_GetLastPageWhenRecordCountIsZeroReturnsTheFirstPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("last", "1", 0)

	assert.Equal(0, pageNumber)
}

func Test_GetPageNumber_GetFirstPageReturnsTheFirstPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("first", "2", 4)

	assert.Equal(0, pageNumber)
}

func Test_GetPageNumber_GetFirstPageWhenPageSizeIsZeroReturnsTheFirstPage(t *testing.T) {
	assert := assert.New(t)

	pageNumber := GetPageNumber("first", "0", 4)

	assert.Equal(0, pageNumber)
}

func GetPageNumber(pageNumber string, pageSize string, recordCount int) int {
	url, _ := url.Parse(fmt.Sprintf("/v1/mandates?page[number]=%s&page[size]=%s", pageNumber, pageSize))
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	pageCriteria := BuildPageCriteria(&ctx)

	return pageCriteria.GetPageNumber(recordCount)
}
