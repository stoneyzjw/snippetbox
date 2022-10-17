package main 

import (
    "fmt"
    "net/http"
    "strconv"
    "errors" 
    "strings"
    "unicode/utf8"
    "github.com/stoneyzjw/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
    Title       string 
    Content     string 
    Expires     int 
    FieldErrors map[string]string
}

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
    data := app.newTemplateData(r)

    app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    /* 
     * First we call r.ParseForm() which adds any ata in POST request bodies 
     * to the r.PostForm map. This also works in the same way for PUT and PATCH 
     * requests. If there are any errors, we use our app.ClientError() helper to 
     * send a 400 Bad Request response to the user. 
     */ 
    err := r.ParseForm() 
    if err != nil {
        app.clientError(w, http.StatusBadRequest) 
        return
    }
    // Use the r.PostForm.Get() method to retrieve the title and content 
    // from the r.PostForm map 
    // title := r.PostForm.Get("title")
    // content := r.PostForm.Get("content")
    /* 
     * The r.PostForm.Get() method always returns the form data as a *string*. 
     * However, we're expecting our expires to be a number, and want to 
     * represent it in our Go code as an integer. So we need to manually convert 
     * the form data to an integer using strconv.Atoi(), and we send a 400 Bad 
     * Request response if the conversion fails
     */
    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // Initialize a map to hold any validation errors for the form fields. 

    form := snippetCreateForm {
        Title:          r.PostForm.Get("title"), 
        Content:        r.PostForm.Get("content"),
        Expires:        expires, 
        FieldErrors:    map[string]string{},
    }
    /*
     * Check that the title is not blank and is not more than 100 
     * characters long. If it fails either of those checks, add a message to the 
     * errors map using the field name as the key.
     */
    if strings.TrimSpace(form.Title) == "" {
        form.FieldErrors["title"] = "This field cannot be blank" 
    } else if utf8.RuneCountInString(form.Title) > 100 {
        form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
    }

    // Check that the Content value isn't blank 
    if strings.TrimSpace(form.Content) == "" {
        form.FieldErrors["content"] = "This field cannot be blank"
    }

    // Check the expires value matches one of the permitted values (1, 7 or 365)
    if expires != 1 && expires != 7 && expires != 365 {
        form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
    }

    // If there are any erros, dump them in a plain text HTTP response and 
    // return from the handler 
    if len(form.FieldErrors) > 0 {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        return 
    }

    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires) 
    if err != nil {
        app.serverError(w, err)
        return
    }

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
