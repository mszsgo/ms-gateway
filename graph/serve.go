package graph

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/mszsgo/hgraph"
	"github.com/rs/cors"
)

func ListenServe() {
	fmt.Println("******************************************************************************************")
	defer PrintStack()
	args()
	app.ListenAndServe()
}

func args() {
	var (
		// 服务名与端口号
		name string
		port int
	)
	flag.StringVar(&name, "name", app.Name, fmt.Sprintf("Set Application name. Default '%s'", app.Name))
	flag.IntVar(&port, "port", app.Port, fmt.Sprintf("Set Port. Default is %d", app.Port))

	flag.Parse()

	app.Name = name
	app.Port = port
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))

	if err := recover(); err != nil {
		log.Fatalf("** Main Fatalf-> %s", err)
	}
}

type Application struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

var app = &Application{Name: "gateway", Version: "0.0.1", Host: "", Port: 80}

func (app *Application) ListenAndServe() {
	Handles()
	log.Printf("MicroService: %s  ListenAndServe %s:%d   Start server http://127.0.0.1:%d", app.Name, app.Host, app.Port, app.Port)
	panic(http.ListenAndServe(fmt.Sprintf("%s:%d", app.Host, app.Port), nil))
}

func Handles() {

	http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "")
	})

	// 反向代理，访问服务
	http.Handle("/", cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path[1:]
		index := strings.Index(uri, "/")
		if index <= 0 {
			w.Write([]byte("path=''"))
			return
		}
		host := uri[0:index]
		path := uri[index:]
		log.Printf("host=%s  path=%s", host, path)
		proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = host
			req.URL.Path = path
		}}
		proxy.ServeHTTP(w, r)
	})))

	// Graphql服务
	http.Handle("/api/graphql", cors.Default().Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		if bytes == nil || len(bytes) == 0 {
			io.WriteString(w, "The request `Body` cannot be null")
			return
		}
		log.WithField("requestBody", bytes).Info("请求报文")
		rs := hgraph.Gateway(bytes)
		log.WithField("responseBody", bytes).Info("响应报文")
		w.Write(rs)
	})))

}
