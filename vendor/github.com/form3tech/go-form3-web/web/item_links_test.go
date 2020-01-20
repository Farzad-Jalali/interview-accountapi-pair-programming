package web

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_BuildItemSelfLink(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildItemLinks(&ctx, "4c1c71e0-d0d8-4e89-bd92-d5e636f34045")
	assert.Equal("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045", link.Self)
}

func Test_BuildItemSelfLinkFromForwardHeader(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045")
	ctx := gin.Context{
		Request: &http.Request{
			URL:    url,
			Header: http.Header{"X-Forwarded-Path": []string{"/v1/transactions/mandates"}},
		},
	}
	link := BuildItemLinks(&ctx, "4c1c71e0-d0d8-4e89-bd92-d5e636f34045")
	assert.Equal("/v1/transactions/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045", link.Self)
}

func Test_BuildItemSelfComplexLinkFromForwardHeader(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045/reversals/")
	ctx := gin.Context{
		Request: &http.Request{
			URL:    url,
			Header: http.Header{"X-Forwarded-Path": []string{"/v1/transactions/mandates"}},
		},
	}
	link := BuildItemLinks(&ctx, "8e4ce432-becc-4c8b-8d92-403eb258f379")
	assert.Equal("/v1/transactions/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045/reversals/8e4ce432-becc-4c8b-8d92-403eb258f379", link.Self)
}

func Test_BuildItemSelfFromForwardHeader_NotPartsInCommon(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/sortcodes/123456")
	ctx := gin.Context{
		Request: &http.Request{
			URL:    url,
			Header: http.Header{"X-Forwarded-Path": []string{"/v1/validations/gbdsc"}},
		},
	}
	link := BuildItemLinks(&ctx, "123456")
	assert.Equal("/v1/validations/gbdsc/sortcodes/123456", link.Self)
}

func Test_BuildItemSelfLinkAppendsId(t *testing.T) {
	assert := assert.New(t)
	url, _ := url.Parse("/v1/mandates")
	ctx := gin.Context{
		Request: &http.Request{
			URL: url,
		},
	}
	link := BuildItemLinks(&ctx, "4c1c71e0-d0d8-4e89-bd92-d5e636f34045")
	assert.Equal("/v1/mandates/4c1c71e0-d0d8-4e89-bd92-d5e636f34045", link.Self)
}
