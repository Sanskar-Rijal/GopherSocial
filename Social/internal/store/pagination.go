package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedQuery struct {
	Limit int `json:"limit" validate:"gte=1,lte=20"` //how much page to show
	Offset int `json:"offset" validate:"gte=0"` //from which page to start
	Sort string `json:"sort" validate:"oneof=asc desc"` //either ascending or descending order
	Search string `json:"search"` //search query
	Tags []string `json:"tags"` //filter by tags
	Since string `json:"since"` 
	Until string `json:"until"`
}

//Read query params from URL 
func (pq *PaginatedQuery) Parse(r *http.Request) error {
	queryString := r.URL.Query() //reads ?limit=10&offset=0&search=hello from URL

	limit :=queryString.Get("limit")
	if limit != ""{
		l , err := strconv.Atoi(limit)

		if err != nil {
			return err

		}
		pq.Limit = l
	}

	//off set, how many to skip
	offset := queryString.Get("offset")

	if offset != ""{
		o, err := strconv.Atoi(offset)
		if err != nil {
			return  err
		}
		pq.Offset = o
	}

	//Sort asc or desc 
	sort := queryString.Get("sort")
	if sort != ""{
		pq.Sort = sort
	}

	search := queryString.Get("search")
	if search != ""{
		pq.Search = search
	}
	

	//Filter by tags 
	// URL: ?tags=go,backend,docker
	tags := queryString.Get("tags")
	if tags != ""{
		pq.Tags = strings.Split(tags,",")
	}

	since := queryString.Get("since")
	if since != "" {
		pq.Since= since
	}

	until := queryString.Get("until")
	if until != "" {
		pq.Until= until
	}

	return nil 
}