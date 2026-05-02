package main

import "net/http"

type server struct{
	addr string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello from the server"))
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path{
		case "/":
			w.Write([]byte("Index page"))
			return 
		case "/users":
			w.Write([]byte("Users page"))
			return
		}
	default:
		w.Write([]byte("404 page not found"))
		return
	}
}


func main(){
	s := &server{addr: ":8080"}

	srv := &http.Server{
		Addr: s.addr,
		Handler: s,
	}
	//err := http.ListenAndServe(s.addr, s)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}


