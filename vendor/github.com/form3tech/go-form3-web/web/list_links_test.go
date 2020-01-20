package web

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_BuildListSelfLink(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 1, PageSize: 1})
	assert.Equal("/v1/mandates", link.Self)
}

func Test_BuildListSelfLinkFromForwardHeader(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL:    url,
			Header: http.Header{"X-Forwarded-Path": []string{"/v1/transactions/mandates"}},
		},
	}
	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 1, PageSize: 1})
	assert.Equal("/v1/transactions/mandates", link.Self)
}

func Test_BuildListSelfComplexLinkFromForwardHeader(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045/reversals/")
	ctx := gin.Context{
		Request: &http.Request{
			URL:    url,
			Header: http.Header{"X-Forwarded-Path": []string{"/v1/transactions/mandates"}},
		},
	}

	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 1, PageSize: 1})
	assert.Equal("/v1/transactions/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045/reversals", link.Self)
}

func Test_BuildListNextLink(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 4, PageSize: 2})
	assert.Equal("/v1/mandates?page%5Bnumber%5D=1", link.Next)
}

func Test_BuildListNextLinkWhenNextPageHasOneElement(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 3, PageSize: 2})
	assert.Equal("/v1/mandates?page%5Bnumber%5D=1", link.Next)
}

func Test_BuildListNextLinkHideWhenItsTheLastPage(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildListLinks(&ctx, PageResults{CurrentPage: 0, TotalRecords: 2, PageSize: 2})
	assert.Equal("", link.Next)
}
