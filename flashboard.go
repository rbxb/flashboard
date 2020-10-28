package flashboard

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	sync.Mutex
	postSize int
	lifetime time.Duration
	posts    [][]byte
	times    []int64
}

func NewServer(postSize, cap int, lifetime time.Duration) *Server {
	return &Server{
		Mutex:    sync.Mutex{},
		postSize: postSize,
		lifetime: lifetime,
		posts:    make([][]byte, 0, cap),
		times:    make([]int64, 0, cap),
	}
}

func (sv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		sv.Lock()
		defer sv.Unlock()
		sv.clean()
		enc := json.NewEncoder(w)
		if err := enc.Encode(sv.posts); err != nil {
			http.Error(w, "Internal error.", 500)
			return
		}
	case "POST":
		sv.Lock()
		defer sv.Unlock()
		if text, err := ioutil.ReadAll(req.Body); err != nil {
			http.Error(w, "Bad request.", 400)
			return
		} else if sv.postSize < len(text) || len(text) == 0 {
			http.Error(w, "Post size invalid.", 400)
			return
		} else {
			sv.clean()
			if len(sv.posts) == cap(sv.posts) {
				http.Error(w, "Server full.", 503)
				return
			}
			sv.posts = append(sv.posts, text)
			sv.times = append(sv.times, time.Now().Add(sv.lifetime).Unix())
			w.WriteHeader(200)
		}
	default:
		w.WriteHeader(405)
	}
}

func (sv *Server) clean() {
	now := time.Now().Unix()
	for i := 0; i < len(sv.times); {
		if sv.times[i] < now {
			sv.posts = append(sv.posts[:i], sv.posts[i+1:]...)
			sv.times = append(sv.times[:i], sv.times[i+1:]...)
		} else {
			i++
		}
	}
}
