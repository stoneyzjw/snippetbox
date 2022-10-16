package main 

import (
    "fmt"
    "net/http"
    "strconv"
    // "log"
    "html/template"
    "errors" 
    "github.com/stoneyzjw/snippetbox/internal/models"
)

// Define a home handler function which writes a byte slice containing 
// "Hello from Snippetbox" as the response body. 

/* 
 * Change the signature of the home handler so it is defined as a method against 
 * *application
 */
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

    for _, snippet := range snippets {
        fmt.Fprintf(w, "%+v\n", snippet)
    }
    // w.Write([]byte("Hello from Snippetbox\n"))
    /*
     * Initialize a slice containing the paths to the two files. It's important 
     * to note that the file containing our base template must be the *first* 
     * file in the slice
     */
    // files := []string {
        // "./ui/html/base.tmpl",
        // "./ui/html/partials/nav.tmpl",
        // "./ui/html/pages/home.tmpl",
    // }
    /* 
     * Use the template.ParseFiles() functions to read the template file into a 
     * template set. If there's an error, we log the detailed error message and use 
     * the http.Error() funtion to send a generic 500 Internal Server Error 
     * response to the user.
     */
    // ts, err := template.ParseFiles(files...)
    // if err != nil {
        // log.Print(err.Error())
        // app.errorLog.Print(err.Error())
        // app.serverError(w, err) // Use the serverError helper
        // http.Error(w, "Internal Server Error", 500)
        // return
    // }
//
    /* 
     * We then use the Execute() method on the template set to write the 
     * template content as the response body. The last parameter to Execute() 
     * represents any dynamic data that we want to pass in, which for now we'll 
     * leave as nil.
     */
     /*
      * Use the ExecuteTemplate() method to write the content of the "base" 
      * template as the response body.
      */
    // err = ts.ExecuteTemplate(w, "base", nil)
    // if err != nil {
        // log.Print(err.Error())
        // app.errorLog.Print(err.Error())
        // app.serverError(w, err) // Use the serverError() helper
        // http.Error(w, "Internal Server Error", 500)
    // }
}

// Add a snippetView handler function. 
/*
 * Change the signature of the snippetView handler so it is defined as a method 
 * against *application. 
 */
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    /* Extract the value of the id parameter, from the query string and try to 
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
    /*
     * Initialize a slice containing the paths to the view.tmpl file, 
     * plus the base layout and navigation partial that we made earlier. 
     */ 

    files := []string {
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/view.tmpl",
    }
    // Parse the template files ...
    ts, err := template.ParseFiles(files...) 
    if err != nil {
        app.serverError(w, err)
        return
    }

    /* Create an instance of a templateData struct holding the snippet data. */ 
    data := &templateData {
        Snippet: snippet,
    }
    /*
     * And then execute them. Notice how we are passing in the snippet 
     * data (a models.Snippet struct) as the final parameter?
     */ 
    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, err)
    }
}

// Add a snippetCreate handler function 
/* 
 * Change the signature of the snippetCreate handler so it is defined as a method 
 * against *application
 */
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    // Use r.Method to check whether the request is using POST or not. 
    // if r.Method != "POST" {
    if r.Method != http.MethodPost {
        /* If it's not, use the w.WriteHeader() method to send a 405 status 
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

