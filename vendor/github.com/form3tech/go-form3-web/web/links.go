package web

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Links struct {
	First string
	Last  string
	Next  string
	Prev  string
	Self  string
}

func BuildItemLinks(c *gin.Context, id string) Links {
	urlPath := c.Request.URL.Path

	forwardUrlPath := c.Request.Header.Get("X-Forwarded-Path")
	if forwardUrlPath != "" {
		urlPath = getUrlPathFromForwardedPath(urlPath, forwardUrlPath)
	}

	urlPath = strings.TrimSuffix(urlPath, "/")
	if !strings.HasSuffix(strings.ToLower(urlPath), strings.ToLower(id)) {
		urlPath = path.Join(urlPath, id)
	}

	return Links{
		Self: urlPath,
	}
}

func getUrlPathFromForwardedPath(urlPath, forwardUrlPath string) string {
	forwardParts := strings.Split(forwardUrlPath, "/")
	matchingPart := forwardParts[len(forwardParts)-1]
	idx := strings.Index(urlPath, matchingPart)
	if idx == -1 {
		return forwardUrlPath + urlPath[len("/v1"):]
	}
	return forwardUrlPath + urlPath[idx+len(matchingPart):]
}

func BuildListLinks(c *gin.Context, results PageResults) Links {
	urlPath := c.Request.URL.Path

	forwardUrlPath := c.Request.Header.Get("X-Forwarded-Path")
	if forwardUrlPath != "" {
		urlPath = getUrlPathFromForwardedPath(urlPath, forwardUrlPath)
	}

	urlPath = strings.TrimSuffix(urlPath, "/")

	next := ""
	if results.CurrentPage < lastPageNumber(results.TotalRecords, results.PageSize) {
		next = pageNumber(c.Request.URL, urlPath, strconv.FormatInt(int64(results.CurrentPage+1), 10))
	}

	prev := ""
	if results.CurrentPage > 0 {
		prev = pageNumber(c.Request.URL, urlPath, strconv.FormatInt(int64(results.CurrentPage-1), 10))
	}

	selfPath := urlPath
	if len(c.Request.URL.Query()) > 0 {
		selfPath = urlPath + "?" + c.Request.URL.Query().Encode()
	}

	return Links{
		Self:  selfPath,
		First: pageNumber(c.Request.URL, urlPath, "first"),
		Next:  next,
		Prev:  prev,
		Last:  pageNumber(c.Request.URL, urlPath, "last"),
	}
}

func pageNumber(url *url.URL, path string, number string) string {
	query := url.Query()
	query.Set("page[number]", number)
	return path + "?" + query.Encode()
}
