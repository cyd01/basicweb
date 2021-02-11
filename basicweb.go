package main
import (
  "bufio"; "flag"; "io"; "log"; "net/http"; "path/filepath"; "os"; "os/exec"; "strconv"; "strings"
)
var (
  command  = flag.String( "cmd"      ,  ""     ,  "external command" )
  dir      = flag.String( "dir"      ,  "."    ,  "root directory"                                        )
  nocache  = flag.Bool  ( "nocache"  ,  false  ,  "force not to cache" )
  password = flag.String( "pass"     ,  ""     ,  "password for basic authentication (modification only)" )
  port     = flag.String( "port"     ,  "80"   ,  "port web server"                                       )
  status   = flag.Int   ( "status"   ,  0      ,  "force return code"                                     )
  username = flag.String( "user"     ,  ""     ,  "username for basic authentication (modification only)" )
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
func setnocache(w http.ResponseWriter) {
  w.Header().Set("Cache-Control","no-cache, no-store, must-revalidate"); 
  w.Header().Set("Expires","0");
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
  if( *nocache ) { setnocache(w) }
  if( (r.Method!="GET")&&(r.Method!="HEAD")&&(r.Method!="OPTIONS") ) { if !basicAuth(w,r) { return } }
  if origin := r.Header.Get("Origin"); origin != "" {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    if r.Method == "OPTIONS" {
      w.Header().Set("Access-Control-Allow-Credentials", "true")
      w.Header().Set("Access-Control-Allow-Methods", "HEAD POST, GET, OPTIONS, PUT, DELETE")
      w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }
  }
  if( *status!=0 ) { // We force return code
    if str := http.StatusText(*status); str != "" {
      w.WriteHeader(*status)
      if( (*status==301)||(*status==302)||(*status==303) ) { w.Header().Set("Location","/") }
      if( *status!=204 ) { w.Write([]byte(strconv.Itoa(*status)+" - "+str)) }
    } else { returnCode(w,http.StatusInternalServerError) }
  } else { // We serve files
    if( r.Method == "OPTIONS" ) { return
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
    } else { returnCode(w,http.StatusMethodNotAllowed)
    }
  }
}
func cmdHandler(w http.ResponseWriter, r *http.Request) {
  log.Println( r.Method, r.URL.Path )
  commands := strings.Split(*command+" 2>&1"," ")
  setnocache(w)
  cmd := exec.Command(commands[0],commands[1:]...)
  stdinPipe, _ := cmd.StdinPipe() ; defer stdinPipe.Close()
  stdoutPipe, _ := cmd.StdoutPipe() ; defer stdoutPipe.Close()
  r.ParseForm()
  cmd.Env = append(os.Environ(),"REQUEST_METHOD="+r.Method,"REQUEST_URI="+r.URL.Path,"SCRIPT_NAME="+r.URL.Path,"HTTP_HOST="+r.Host,"SERVER_PROTOCOL="+r.Proto,"REMOTE_ADDR="+r.RemoteAddr,"CONTENT_TYPE="+r.Header.Get("Content-type"),"CONTENT_LENGTH="+r.Header.Get("Content-length"),"QUERY_STRING="+r.URL.RawQuery)
  for key, val := range r.Header { cmd.Env = append(cmd.Env, "HTTP_" + strings.ReplaceAll( strings.ToUpper( key ), "-", "_" )+"="+val[0]) }
  cmd.Start()
  if l,_ := strconv.Atoi(r.Header.Get("Content-Length")) ; l>0 { go func() { io.Copy(stdinPipe, r.Body) ; stdinPipe.Close() }() }
  scanner := bufio.NewScanner(stdoutPipe)
  for scanner.Scan() {
    out := scanner.Text() ; 
    if( len(out)==0 ) { break }
    if( !strings.Contains(out,":") ) { w.Write([]byte(out)); break }
    head := strings.SplitN( out, ":", 2)
    if( strings.EqualFold(head[0],"Status") ) {
      if s,err := strconv.Atoi(strings.TrimSpace(head[1])) ; err==nil { w.WriteHeader(s) }
    }
    w.Header().Set( head[0], strings.TrimSpace(head[1]) )
  }
  for scanner.Scan() {
    out := scanner.Text()
    w.Write( []byte(out+"\r\n") )
  }
  cmd.Wait()
}
func main() {
  flag.Parse()
  log.Println("Starting web server with port "+*port+" on directory "+*dir+" with status response "+strconv.Itoa(*status))
  if( len(*command)>0 ) { log.Println("Add dynamic command <"+*command+"> to /cmd path"); http.Handle("/cmd", http.HandlerFunc(cmdHandler)) }
  http.HandleFunc("/ping", func (w http.ResponseWriter, r *http.Request) { log.Println( r.Method, r.URL.Path ); w.Write([]byte("pong")) } )
  http.Handle("/", http.HandlerFunc(fileHandler))
  log.Fatal(http.ListenAndServe(":"+*port, nil))
}
