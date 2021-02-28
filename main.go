package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/zserge/webview"
)

func handler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("%s\n", "http://"+ln.Addr().String())

	go func() {
		http.HandleFunc("/", handler)
		http.ListenAndServe(ln.Addr().String(), nil)
		// Set up your http server here
		log.Fatal(http.Serve(ln, nil))
	}()
	fmt.Printf("vor webview\n")
	w := webview.New(webview.Settings{"Forumlauncher", "http://" + ln.Addr().String(), 400, 400, false, true, nil})
	//webview.Open("Hello", "http://"+ln.Addr().String(), 400, 300, false)
	fmt.Printf("vor run\n")
	go func() {
		fmt.Printf("vor dispatch\n")
		time.Sleep(time.Second)
		w.Dispatch(func() {
			for i := 1; i <= 100; i++ {
				w.Eval(fmt.Sprintf("setProgress(%d)", i))
			}
		})
	}()

	w.Run()
}
