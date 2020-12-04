package main
import (
  "flag"
  "io"
  "log"
  "path/filepath"
  "net/http"
  "os"
  "strconv"
)
var (
  port     = flag.String( "port"     ,  "80"   ,  "port web server"                                       )
  dir      = flag.String( "dir"      ,  "."    ,  "root directory"                                        )
  status   = flag.Int   ( "status"   ,  0      ,  "force return code"                                     )
  username = flag.String( "user"     ,  ""     ,  "username for basic authentication (modification only)" )
  password = flag.String( "pass"     ,  ""     ,  "password for basic authentication (modification only)" )
  nocache  = flag.Bool  ( "nocache"  ,  false  ,  "force not to cache" )
)
func basicAuth(w http.ResponseWriter, r *http.Request) bool {
  if( *username!="" && *password!="" ) {
    if user, pass, ok := r.BasicAuth(); !ok || user!=*username || pass != *password { 
      log.Println("Wrong credential")
      w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
      returnCode(w,401)
      return false
    }
  }
  return true
}
func enableCors(w *http.ResponseWriter, r *http.Request) {
  if origin := r.Header.Get("Origin"); origin != "" {
    (*w).Header().Set("Access-Control-Allow-Origin", origin)
    if r.Method == "OPTIONS" {
      (*w).Header().Set("Access-Control-Allow-Credentials", "true")
      (*w).Header().Set("Access-Control-Allow-Methods", "HEAD POST, GET, OPTIONS, PUT, DELETE")
      (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }
  }
}
func returnCode(w http.ResponseWriter,code int) {
  w.WriteHeader(code)
  w.Write([]byte(http.StatusText(code)))
}
func fileHandler(w http.ResponseWriter, r *http.Request) {
  var fullpath string
  if stat, err := os.Stat(*dir+"/"+r.Host); err == nil && stat.IsDir() { fullpath = *dir+"/"+r.Host
  } else { fullpath = *dir
  }
  log.Println( r.Method, r.URL.Path )
  if( *nocache ) { w.Header().Set("Cache-Control","no-cache, no-store, must-revalidate"); w.Header().Set("Expires","0"); }
  if( (r.Method!="GET")&&(r.Method!="HEAD")&&(r.Method!="OPTIONS") ) { if !basicAuth(w,r) { return } }
  enableCors(&w,r)
  if( *status!=0 ) { // We force return code
    if str := http.StatusText(*status); str != "" {
      w.WriteHeader(*status)
      if( (*status==301)||(*status==302)||(*status==303) ) { w.Header().Set("Location","/") }
      if( *status!= 204 ) { w.Write([]byte(strconv.Itoa(*status)+" - "+str)) }
    } else { returnCode(w,http.StatusInternalServerError) }
  } else { // We serve files
    if( r.Method == "OPTIONS" ) {
      return
    } else if( (r.Method == "PUT") || (r.Method == "POST") ) { // Upload fle
      if _, err := os.Stat(filepath.Dir(fullpath+r.URL.Path)); err != nil { 
        if err := os.MkdirAll(filepath.Dir(fullpath+r.URL.Path),0755); err != nil { returnCode(w,http.StatusInternalServerError); return }
      }
      dst, err := os.Create(fullpath+r.URL.Path)
      if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
      defer dst.Close()
      defer r.Body.Close() 
      if _, err := io.Copy(dst, r.Body); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
      returnCode(w,http.StatusCreated)
    } else if( r.Method == "DELETE" ) { // Delete file
      if err:= os.Remove(fullpath+r.URL.Path); err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
      returnCode(w,http.StatusNoContent)
    } else if( (r.Method == "GET") || (r.Method == "HEAD") ) { // Download file or file info
      http.FileServer(http.Dir(fullpath)).ServeHTTP(w, r)
    } else { // Unknown method
      returnCode(w,http.StatusMethodNotAllowed)
    }
  }
}
func main() {
  flag.Parse()
  http.Handle("/", http.HandlerFunc(fileHandler))
  http.HandleFunc("/ping", func (w http.ResponseWriter, r *http.Request) { log.Println( r.Method, r.URL.Path ); w.Write([]byte("pong")) } )
  log.Println("Starting web server with port "+*port+" on directory "+*dir+" with status response "+strconv.Itoa(*status))
  log.Fatal(http.ListenAndServe(":"+*port, nil))
}
