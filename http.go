package main

import (
    "math/rand"
    "time"
	"log"
    "net/http"
    "net/http/httputil"
    "net/url"
	"sync"
	"strconv"
	"os"
	"encoding/json"
    "os/exec"
	"fmt"
)

type Request struct {
    Code struct {
        PayloadCode string `json:"payload_code"`
    } `json:"code"`

    Env map[string]string `json:"env"`
}

var (
workdir string
p int
sandboxPid *exec.Cmd
sandboxOn bool
sandboxLock sync.Mutex
appfile string
data[] byte
User string
Pass string
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(data)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "not json", http.StatusUnsupportedMediaType)
        return
    }
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
	var vset string
	for k, v := range req.Env { 
	if vset!="" {vset+=","}
	vset+=fmt.Sprintf("(name=\"%s\",text=\"%s\")",k,v)
	}
	if vset!="" {vset+=","}
	vset+=fmt.Sprintf("(name=\"v_runtime\",text=\"v0.1\")")
	err := makeruntime(req.Code.PayloadCode,vset) 
	if err!=nil{
		log.Println(err)
		fmt.Fprintf(w,`{"status":"fail"}`)
	}else{
	time.Sleep(time.Millisecond*612)
	if sandboxOn {
	fmt.Fprintf(w,`{"status":"ok"}`)
	}else{
	fmt.Fprintf(w,`{"status":"fail"}`)
	}
	}
}

func makeruntime(payload,e string) error {
	err := os.WriteFile(workdir+"/main.js", []byte(payload), 0600)
	if err != nil {
		return err
	}
	p=rand.Intn(60000-10000+1) + 10000
	t:=`using Workerd = import "/workerd/workerd.capnp";const helloWorldExample :Workerd.Config = (services = [ (name = "main", worker = .helloWorld) ],sockets = [ ( name = "http", address = "127.0.0.1:`+strconv.Itoa(p)+`", http = (), service = "main" ) ]);const helloWorld :Workerd.Worker = (modules = [(name = "worker", esModule = embed "main.js")],compatibilityDate = "2025-11-17",compatibilityFlags = ["enable_web_file_system", "experimental"],bindings = [`+e+"]);"
	err = os.WriteFile(workdir+"/sandbox_"+strconv.Itoa(p), []byte(t), 0600)
	if err != nil {
		return err
	}
	return start_sandbox(p)
}

func start_sandbox(sandbox_Num int) error {
	sandboxLock.Lock()
	defer sandboxLock.Unlock()
	if sandboxOn {
		sandboxPid.Process.Kill()
		sandboxPid.Wait()
	}
	try := exec.Command(appfile,"serve",workdir+"/sandbox_"+strconv.Itoa(sandbox_Num),"--experimental")
	try.Dir = workdir
	try.Stdout = os.Stdout
	try.Stderr = os.Stderr
	if err := try.Start(); err != nil {
		log.Println("ERR: start up sandbox:",err)
		os.Remove(workdir+"/sandbox_"+strconv.Itoa(sandbox_Num))
		return err
	}
	sandboxPid = try
	sandboxOn = true
	go func(){
		sandboxPid.Wait()
		sandboxOn = false
	}()
	go func(){
	time.Sleep(time.Second*3)
	os.Remove(workdir+"/sandbox_"+strconv.Itoa(sandbox_Num))
	}()
	return nil
}

func devHandler(w http.ResponseWriter, r *http.Request) {
    user, pass, ok := r.BasicAuth()
    if !ok || user != User || pass != Pass {
        w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return;}
    handler(w, r)
}

func workerThread(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(404)
		return
	}
    var target *url.URL
    var err error
	if p < 8000 {
		log.Println("worker not start")
        http.Error(w, "worker not start", 500)
        return
	}
    target, err = url.Parse("http://127.0.0.1:"+strconv.Itoa(p))
    if err != nil {
		log.Println(err)
        http.Error(w, "Internal Server Error 500", 500)
        return
    }
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.Director = func(req *http.Request) {
        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host
        req.Host = target.Host
    }
    proxy.ServeHTTP(w, r)
}

func init_Nano() {
	rand.Seed(time.Now().UnixNano())
}

func sandbox_service(){
	argc:=len(os.Args)
	if argc <=1{
		log.Fatalln("Usage:",os.Args[0],"[config.json]")
	}
	conf,err:=checkconf(os.Args[1])
	if err!=nil{
		log.Fatalln("ERR: read config err:",err)
	}
	workdir=conf.WorkDir
	appfile=conf.AppFile
	User=conf.User
	Pass=conf.Pass
	data=[]byte(datea)
    http.HandleFunc("/", workerThread)
	http.HandleFunc("/dev_page", devHandler)
    log.Println("listen:",conf.Listen,"ssl:",conf.EnableSSL)
	if !conf.EnableSSL{
    log.Fatal(http.ListenAndServe(conf.Listen, nil))
	}else{
	log.Fatal(http.ListenAndServeTLS(conf.Listen, conf.SSLConfig.Crt, conf.SSLConfig.Key, nil))
	}
}