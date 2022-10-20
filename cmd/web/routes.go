package main 

import (
    "net/http"
    "github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"
)

/* The routes() method returns a servemux containing our application routes. */ 
/* 
 * Update the signature for the routes() method so that it returns a 
 * http.Handler instead of *http.ServerMux.
 */
func (app *application) routes() http.Handler {
	// Initialize the router 
	router := httprouter.New()
	// Create a handler function which wraps our notFound() helper, and then 
	// assign it as the custom handler for 404 Not Found responses. You can also 
	// set a custom handler for 405 Method Not Allowed responses by setting 
	// router.MethodNotAllowed in the same way too. 
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// Update the pattern for the route for the static files.
    fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	// And then create the routes using the appropriate methods, patterns and handlers.
	// Update these routes to use the new dynamic middleware chain followed by 
	// the appropriate handler function. Note that because the alice ThenFunc() 
	// method returns a http.Handler (rather than a http.HandleFunc) we also 
	// need to switch to registering the route using the router.Handler() method. 
	// router.HandlerFunc(http.MethodGet, "/", app.home)
	// router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	// router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	// router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	// Add the five new routes, all of which use our 'dynamic' middleware chain
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// Create the middleware chain as normal 
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders) 

	// Wrap the router with the middleware and return it as normal 
	return standard.Then(router)
}
