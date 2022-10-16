package main 

import (
    "net/http"
    "github.com/justinas/alice"
)

/* The routes() method returns a servemux containing our application routes. */ 
/* 
 * Update the signature for the routes() method so that it returns a 
 * http.Handler instead of *http.ServerMux.
 */
func (app *application) routes() http.Handler {
    mux := http.NewServeMux() 
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("/", app.home) 
    mux.HandleFunc("/snippet/view", app.snippetView) 
    mux.HandleFunc("/snippet/create", app.snippetCreate)

    /*
     * Pass the servermux as the 'next' parameter to the secureHeaders middleware
     * Because secureHeaders is just a function, and the function returns a 
     * http.Handle we don't need to do anything else. 
     */
    // return secureHeaders(mux)
    // Wrap the existing chain with the logRequest middleware 
    // return app.recoverPanic(app.logRequest(secureHeaders(mux)))
    /* 
     * Create a middleware chain containing our 'standard' middleware 
     * which will be used for every request our application receives.
     */ 
    standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
    // Return the 'standard' middleware chain followed by the servemux 
    return standard.Then(mux)
}
