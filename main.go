package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
	. "gokitdemo/Services"
	"gokitdemo/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// 获取参数
	name := flag.String("name", "", "服务器名称")
	port := flag.Int("p", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请输入服务名")
	}
	if *port == 0 {
		log.Fatal("请输入端口")
	}

	util.SetServiceNameAndPort(*name, *port)
	user := UserService{}
	// api限流
	limit := rate.NewLimiter(1, 3)
	endp := RateLimit(limit)(GetUserEndpoint(user))
	//endp := GetUserEndpoint(user)
	// 错误处理
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(MyErrorEncoder),
	}
	serverHandler := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse,options...)
	r := mymux.NewRouter()
	//r.Handle(`/user/{uid:\d+}`,serverHandler)
	//r.Methods("GET").Path(`/user/{uid:\d+}`).Handler(serverHandler)
	//r = mymux.NewRouter()

	{
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})
	}
	errChan := make(chan error)
	go (func() {
		util.RegService()
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r)
		if err != nil {
			errChan <- err
		}
	})()
	go func() {
		sig_c := make(chan os.Signal)
		signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig_c)
	}()

	getErr := <-errChan
	util.Unregservice()
	log.Println(getErr)

}
