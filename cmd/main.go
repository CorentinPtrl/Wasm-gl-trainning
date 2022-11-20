package main

import "net/http"

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	http.Handle("/wasm/", http.StripPrefix("/wasm/", http.FileServer(http.Dir("web/wasm/"))))
	http.Handle("/", &noCache{Handler: http.FileServer(http.Dir("web"))})
	http.ListenAndServe(":8080", nil)
}

type noCache struct {
	http.Handler
}

func (h *noCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	h.Handler.ServeHTTP(w, r)
}
