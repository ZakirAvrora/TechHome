package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const PageIdKey = "page"

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageId := r.URL.Query().Get(string(PageIdKey))
		page := 0
		var err error

		if pageId != "" {
			page, err = strconv.Atoi(pageId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), PageIdKey, page)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
