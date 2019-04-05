package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var help = flag.Bool("help", false, "show keys")
var static_dir = flag.String("static_dir", "static", "static files folder")
var http_port = flag.String("http_port", "3000", "http port")
var https = flag.Bool("https", false, " -https for use https protocol")
var crt = flag.String("crt", "domain.crt", "crt - file")
var key = flag.String("key", "domain.key", "key - file")

func init() {

}

func main() {
	flag.Parse()
	if *help == true {
		flag.VisitAll(func(flag *flag.Flag) {
			format := "\t-%s: %s (Default: %s) \n"
			fmt.Printf(format, flag.Name, flag.Usage, flag.DefValue)
		})
		fmt.Println("For example StaticServer.exe -run -http_port=3000 static_dir=static")
	} else {
		fs := http.FileServer(http.Dir(*static_dir))
		http.Handle("/", fs)
		http.HandleFunc("/Help", helpHF)
		http.HandleFunc("/About", helpHF)
		http.HandleFunc("/formHF", formHF)
		InitSessionManager()
		log.Println("Listening...")
		fmt.Println("Port :", *http_port, " Static files dirrectory:", *static_dir)
		fmt.Println("Press Ctrl+C for stop server.")
		var err error = nil
		httpServer := &http.Server{Addr: ":" + *http_port,
			Handler:        nil,
			ReadTimeout:    1000 * time.Second,
			WriteTimeout:   1000 * time.Second,
			MaxHeaderBytes: 1 << 20}

		if *https {
			err = httpServer.ListenAndServeTLS(*crt, *key)
		} else {
			err = httpServer.ListenAndServe()
		}

		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		log.Println("Stopped")
	}
}

func helpHF(res http.ResponseWriter, req *http.Request) {
	ses := GetSession(res, req)
	var i int64 = ses.GetInt64("IntVar")
	var f float64 = ses.GetFloat64("FloatVar")
	ses.Set("IntVar", i+1)
	ses.Set("FloatVar", f+3.14)
	ses.Set("DT", time.Now())
	//fmt.Println(i)
	fmt.Fprintf(res, "Simple http server by go.\n Yuor SESSION ID="+ses.GetId()+
		" \n IntVar = "+ses.GetString("IntVar")+
		" \n FloatVar = "+ses.GetString("FloatVar")+
		" \n DT = "+ses.GetString("DT"))
}

func formHF(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, "Err")
	} else {

		cb2_1 := req.FormValue("checkbox2")
		fmt.Println("cb2_1=", cb2_1)
		rb1 := req.FormValue("optionRadios")
		fmt.Println("rb1=", rb1)

		fmt.Println(req.PostForm)
		for k, v := range req.PostForm {
			fmt.Println(k, "=", v)
		}
		cb2_2 := req.PostForm["checkbox2"]
		fmt.Println("cb2_2=", cb2_2)

		fmt.Fprintf(res, "Ok")
	}

}
