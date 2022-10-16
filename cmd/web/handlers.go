package main 

import (
    "fmt"
    "net/http"
    "strconv"
    "errors" 
    "github.com/stoneyzjw/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    snippets, err := app.snippets.Latest() 
    if err != nil {
        app.serverError(w, err) 
        return
    }
    // Use the new render helper. 
    app.render(w, http.StatusOK, "home.tmpl", &templateData {Snippets: snippets})
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    /*
     * Extract the value of the id parameter, from the query string and try to 
     * convert it to an integer using the strconv.Atoi() function. If it can't 
     * be convert to an integer, or the value is less than 1, we return a 404 page 
     * not found reponse
     */
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
    // Use the new render helper 
    app.render(w, http.StatusOK, "view.tmpl", &templateData{Snippet: snippet})
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    // Use r.Method to check whether the request is using POST or not. 
    // if r.Method != "POST" {
    if r.Method != http.MethodPost {
        /*
         * If it's not, use the w.WriteHeader() method to send a 405 status 
         * code and the w.Write() method to write a "Method Not Allowed" 
         * response body. We then return from the function so that the 
         * subsequent code is not executed 
         */ 
        w.Header().Set("Allow", http.MethodPost)
        // http.Error(w, "Method Not Allowed\n", http.StatusMethodNotAllowed)
        app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper
        // w.WriteHeader(405)
        // w.Write([]byte("Method Not Allowed\n"))
        return
    }
    /*
     * Create some variables holding dummy data. We'll remove these later on 
     * during the build. 
     */ 
    title := "0 snail"
    content := "0 snail\nClimb Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa" 
    expires := 7
    /*
     * Pass the data to the SnippetModel.Insert() method, receiving the 
     * ID of the new record back. 
     */ 
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err) 
        return
    }
    /* Redirect the user to the relevant page for the snippet. */ 
    http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
    // w.Write([]byte("Write a new snippet\n"))
}

