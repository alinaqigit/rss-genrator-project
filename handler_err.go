package main

import "net/http"

func handlerErr(res http.ResponseWriter, req *http.Request) {
	responseWithError(res, 500, "Something Went wrong")
}