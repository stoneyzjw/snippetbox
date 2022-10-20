package main 

import (
    "fmt"
    "net/http"
    "strconv"
    "errors" 
    "github.com/stoneyzjw/snippetbox/internal/models"
    "github.com/stoneyzjw/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

/*
 * Update our snippetCreateForm struct to include struct tags which tell the 
 * decoder how to map HTML form values into the different struct fields. So, for 
 * example, here we're telling the decoder to store the value from the HTML form 
 * input with the name "title" in the Title field. The struct tag `form:"-"` 
 * tells the decoder to completely ignore a field during decoding. 
 */
type snippetCreateForm struct {
    Title       string      `form:"title"`
    Content     string      `form:"content"`
    Expires     int         `form:"expires"`
    validator.Validator     `form:"-"`
}

// Create a new userSignupFor struct. 
type userSignupForm struct {
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
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

    // Initialize a new createSnippetForm instance and pass it to the template 
    // Notice how this is also a greate opportunity to set any default or 
    // 'initial' values for the form -- here we set the initial value for the snippet expire to 365
    // days

    data.Form = snippetCreateForm {
        Expires: 365, 
    }

    app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    // Declare a new empty instance of the snippetCreateForm struct. 
    var form snippetCreateForm

    /* 
     * Call the Decode() method of the form decoder, passing in the current 
     * request and *a pointer* to our snippetCreateForm struct. This will 
     * essentially fill our struct with the relevant values from the HTML form. 
     * If there is a problem, we return a 400 Bad Request response to the client.
     */
    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // The validate and use the data as normal ... 
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7 or 365")

    if !form.Valid() {
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

	// Use the Put() method to add a string value ("Snippet successfully 
	// created!") and the corresponding key ("flash") to the session data. 
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = userSignupForm{}
    app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
