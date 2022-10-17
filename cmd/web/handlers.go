package main 

import (
    "fmt"
    "net/http"
    "strconv"
    "errors" 
    "github.com/stoneyzjw/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the 
	// manual check of r.URL.Path != "/" from this handler 

    snippets, err := app.snippets.Latest() 
    if err != nil {
        app.serverError(w, err) 
        return
    }
    /* 
     * Call the newTemplateData() helper to get a templateData struct containing 
     * the 'default' data (which for now is just the current year), and add the 
     * snippets slice to it. 
     */ 
    data := app.newTemplateData(r)
    data.Snippets = snippets
    // Use the new render helper. 
    app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    /*
	 * When httprouter is parsing a request, the values of any named parameter 
	 * will be stored in the request context. We'll talk about request context 
	 * in detail later in the book, but for now it's enough to know that you can 
	 * use the ParamsFromContext() function to retrieve a slice containing these 
	 * parameter names and values like so:
     */
	params := httprouter.ParamsFromContext(r.Context())
	// We can then use the ByName() method to get the value of the "id" named
	// parameer from the slice and validate it as normal.
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        // http.NotFound(w, r)
        app.notFound(w)
        return
    }
    /*
     * Use the SnippetModel object's Get method to retrieve the data for a 
     * specific record based on its ID. If no matching record is found, 
     * return a 404 Not Found response. 
     */ 
    snippet, err := app.snippets.Get(id) 
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return 
    }
    /* 
     * Call the newTemplateData() helper to get a templateData struct containing 
     * the 'default' data (which for now is just the current year), and add the 
     * snippets slice to it. 
     */ 
    data := app.newTemplateData(r)
    data.Snippet = snippet
    // Use the new render helper 
    app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Checking if the request method is a POST is now superfluous and can be 
	// removed, because this is done automatically by httprouter. 
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa" 
	expires := 7
	id, err := app.snippets.Insert(title, content, expires) 
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Update the redirect path to use the new clean URL format. 
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
