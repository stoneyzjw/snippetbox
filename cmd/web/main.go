package main 

import (
    "flag"
    "log"
    "net/http"
)

func main() {
    /* 
     * Define a new command-line flag with the name 'addr', a default value of ":4000"
     * and some short help text explaining what the flag controls. The value of the 
     * falg will be stored in the addr variable at runtime.
     */
    addr := flag.String("addr", ":4000", "HTTP network address")
    /* 
     * Importantly, we use the flag.Parse() function to parse the command-line flag. 
     * This reads in the command-line flag value and assigns it to the addr 
     * variable. You need to call this *before* you ue the addr variable 
     * otherwise it will always contain the default value of ":4000". If any errors are 
     * encoutered during parsing the application will be terminated. 
     */
    flag.Parse()

    // Use the http.NewServeMux() function to initialize a new servemux, then
    // register the home function as the handler for the "/" URL pattern 
    mux := http.NewServeMux()
    /*
     * Create a file server which serves files out of the "./ui/static" directory 
     * Note that the path given to the http.Dir function is relative to the project 
     * directory root. 
     */
    fileServer := http.FileServer(http.Dir("./ui/static/"))

    /*
     * Use the mux.Handle() function toregister the file server as the handler for 
     * all URL paths that start with "/static/". For matching paths, we strip the 
     * "/static" prefix before the request reaches the file server.
     */
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    /* Register the other application routes as normal. */
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet/view", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)
    /* 
     * Use the http.listenAndServe() function to star a new web server. We pass in 
     * two parameters: the TCP network address to listen on (in this case ":4000") 
     * and the servemux we just created. If http.ListenAndServe() returns an error
     * we ues the log.Fatal() function to log the error message and exit. Note 
     * that any error returned by http.ListenAndServe() is always not-nil. 
     */
    /* 
     * The value returned from the flag.String() function is a pointer to the flag
     * value, not the value itself. So we need to dereference the pointer (i.e.
     * prefix it with the * symbol) before using it. Note that we're using the 
     * log.Print() function to interpolate the address with the log message. 
     */
    // log.Print("Starting server on :4000")
    log.Printf("Starting server on %s", *addr)
    err := http.ListenAndServe(*addr, mux) 
    log.Fatal(err)
}
