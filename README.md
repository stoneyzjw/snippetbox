# snippetbox

## Chapter02 

Use go mod init command to create the module: github.com/stoneyzjw/snippetbox. the file go.mod is the
result.

## Web application basics 

Now that everything is set up correctly let's make the first iteration of our web application. We'll
begin with the three absolute essentials: 
1. The first thing we need is a handler. If you're coming from an MVC-background, you can think of
   handlers as being a bit like controllers. They're responsible for excuting your application logic
   and for writing HTTP response headers and bodies.

2. The second component is a router 9or servermux in Go terminology). This stores a mapping between the
   URL patterns for your application and the corresponding handlers. Usually you have one servemux for
   your application containing all your routes. 

3. The last thing we need is a web server. One of the great things about Go is that you can establish a
   web server and listen for incoming request as part of your application itself. You don't need an
   external third-party server like Nginx or Apache. 

Let's put these components together in the main.go file to make a working application. 

## Routing requests 

Having a web application which just one route isn't very exciting.. or useful! Let's add a couple more
routes so that the application starts to shape up like this: 

|URL Pattern |Handler|Action 
|:----|:----|:-----|
|/ | home | Display the home page |
|/snippet/view| snippetView| Display a specific snippet|
|/snippet/create|snippetCreate|Create a new snippet|

## Customizing HTTP headers 

Let's now update our application so that the /snippet/create route only responds to HTTP requests which
use the POST method, like so

|Method|URL Pattern |Handler|Action 
|:----|:----|:----|:-----|
|ANY|/ | home | Display the home page |
|ANY|/snippet/view| snippetView| Display a specific snippet|
|POST|/snippet/create|snippetCreate|Create a new snippet|

Making this change is important because - later in our application build requests to the
/snippet/create route will result in a new snippet beging created in a database. Creating a new snippet
in a database is a non-idempotent action that changes the state of our server, so we should follow HTTP
good practice and restrict this route to act on POST request only. 

## The http.Error shortcut 

If you want to send a non-200 status code and a plain-text response body (like we are in the code
above) then it's a good opportunity to use the **http.Error()** shortcut. This is a lightweight helper
function which take a given message and status code, then calls the **w.WriteHeader()** and
**w.Write()** methods behind-the-scenes for us. 

## The net/http constants 

One final tweak we can make is to use constants from the net/http package for HTTP methods and status
codes, instead of writing the strings and integers ourselves. 

Specifically, we can use the constant http.MethodPost instead of the string "POST", and the constant
http.StatusMethodNotAllowed instead of the integer 405. Like so: 

# URL query strings 

While we're on the subject of routing, let's update the snippetView handler so that it accepts an id
query string parameter from the user like so: 

|Mehtod|Pattern|Handler|Action|
|:---|:--|:--|:--|
|ANY|/|home|Display the home page|
|ANY|/snippet/view?id=1|snippetView|Display a specific snippet|
|POST|/snippet/create|snippetCreate|Create a new snippet|

Later we'll use this id parameter to select a specific snippet from a database and show it to the user.
But for now, we'll just read the value of the id parameter and interpolate it with a placeholder
response. 

## Additional information 

### The internal directory 

It's important to point out that the directory name **internal** carries a special meaning and behavior
in Go: any packages which live under this directory can only be imported by code inside the parent of
the **internal** directory. In our cse, this means that any packages which live in **internal** can
only be imported by code inside our snippetbox project directory. 

Or, looking at it the other way, this means that any packages under **internal** cannot be imported by
code outside of our project. 

This is useful because it prevents other codebases from importing and relying on the packages in our
**internal** directory - even if the project code is publicly available somewhere like GitHub. 

## HTML tmeplating and inheritance 

Let's inject a bit of life into the project and develop a proper home page for our Snippetbox web
application. Over the next couple of chapters we'll work towards creating a page which looks like this: 

## The http.Fileserver handler 

Go's net/http package ships with a built-in http.FileServer handler which you can use to serve files
over HTTP from a specific directory. Let's add a new route to our application so that all requests
which begin with "/static/" are handled using this, like so:

|Mehtod|Pattern|Handler|Action|
|:---|:--|:--|:--|
|ANY|/|home|Display the home page|
|ANY|/snippet/view?id=1|snippetView|Display a specific snippet|
|POST|/snippet/create|snippetCreate|Create a new snippet|
|ANY|/static|http.FileServer|Serve a specific static file|

## The http.Handler interface 

Before we go any further there's a little theory that we should cover. It's a bit complicated, so if
you find this chapter hard-going don't worry. Carry on with the application build and circle back to it
later once you're more familiar with Go. 

In the previous chapters I've thrown around the term handler without explaining what it truly means.
Strictly speaking, what we mean by handler is an object which satifies the http.Handler interface. 
    type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
    }
In simple terms, this basically means that to be a handler an object must have a ServeHTTP() method
with the exact signature
    ServeHTTP(http.ResponseWriter, *http.Request)
So in its simplest form a handler might look something like this: 
    type home struct { }

    func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("This is m home page"))
    }
Here we have an object (in this case it's a home struct, but it could equally be a string or function
or anything else), and we've implemented a method with the signature ServeHTTP(http.ResponseWriter,
\*http.Request) on it. That's alll we need to make a handler. 

You could then register this with a servremux using the Handle method like so: 
    mux := http.NewServeMux() 
    mux.Handle("/", &home{})
When this servemux receives a HTTP request for '/', it will then call the ServeHTTP() method of the
home struct - which in turn writes the HTTP response. 

## Handler functions 

Now, creating an object just so we can implement a ServeHTTP() method on it is long-winded and a bit
confusing. Which is why in practice it's far more common to write your handlers as a normal function
(like we have been so far). For example: 

# Chapter 3

## Configuration and error handling 

In this section of the book we're going to do some housekeeping. We won't add much new functionality,
but instead focus on making improvements that'll make it easier to manage our application as it grows. 

You'll learn how to: 
1. Set configuration settings for you application at runtime in an easy and idiomatic way using
   command-line flags. 

2. Improve your application log messages to include more information, and manage them differently
   depending on the type (or level) of log messages. 

3. Make dependencies available to your handlers in a way that's extensible, type-safe, and doesn't get
   in a say when it comes to writing tests. 

4. Centralize error handling so that you don't need to repeat yourself when writing code. 

## Managing configuration settings 

Our web application's main.go file currently contains a couple of hard-coded configuration settings: 

+ The network address for the server to listen on (currently ":4000")
+ The file path for the static files directory (currently "./ui/static") 

Having these hard-coded isn't ideal. There's no separation between our configuration setttings and
code, and we can't change the settings at runtime (which is important if you need different settings
for development, testing and productioin environments).

In this chapter we'll start to improve that, and make the network address for our server configuration
at runtime. 

## Leveled logging 

At the moment in our **main.go** file we're outputting log messages using the **log.Printf()** and
**log.Fatal()** functions. 

Both these functions output messages via Go's standard logger, which - by default - prefixes messages
with the local date and time and writes them to the standard error stream. The log.Fatal() function
will also call os.Exit(1) after writing the message, causing the application to immediately exit. 

In our application, we can break apart our log messages into two distinct types - or levels. The first
type is **informational messages** 

## Dependency injection 

There's one more problem with our logging that we need to address. If you open up your **handlers.go**
file you'll notice that the **home** handler function is still writting error messages using Go's
standard logger, not the errorLog logger that we want to be using. 

There are a few different ways to do this, the simplest being to just put the dependencies in global
variables. But in general, it is good practice to inject dependencies into your handlers. It makes your
code more explicit, less error-prone and easier to unit test than if you use global variables. 

For applications where all your handlers are in the same package like ours, a neat way to inject
dependencies is to put them into a custom **application** struct, and then define your handler
functions as methods against **application**. 

## Additonal information 

### Closures for dependency injection 

The pattern that we're using to inject dependencies won't work if your handlers are spread across
multiple packages. In that case, an alternative approach is to create a config package exporting an
Application struct and have your handler functions close over this to form a closure. Very roughly:

## Isolating the application routes 

Whiel we're refactoring our code there's one more change worth making 

Our **main()** functioin is beginning to get a bit crowded, so to keep it clear and focused I'd like to
move the route declarations for the application into a standalone **routes.go** file. 

# Database-driven responses 

For our Snippetbox web application to become truly useful we need somewhere to store (or persist) the
data entered by users, and the ability to query this data store dynamically at runtime. 

There are many different data stores we could use for our application - each with different pros and
cons - but we'll opt for the popular relational database MySQL. 

## Setting up MySQL 

If you're following along, you'll need to install MySQL on your computer at this point. The offical
MySQL documentation contains comprehensive installation instructions for all types of operating
systems, but if you're using Mac OS you should be able to install it with: 
    brew install mysql 
Or if you're using a Linux distribution which supports apt (like Debian and Ubuntu) you can install it
with: 
    sudo apt install mysql-server 

    /* for a given DSN */
    func openDB(dsn string) (*sql.DB, error) {
        db, err := sql.Open("mysql", dsn)
        if err != nil {
            return nil, err 
        }
        if err = db.Ping(); err != nill {
            return nil, err 
        }
        return db, nil
    }

There're a few things about this code which are interesting: 
1. Notice how the import path for our drive is prefixed with an underscoe? This is because our main.go
   file doesn't actually use anything in the mysql package. So if we try to import it normally the Go
   compiler will raise an error. However, we need the driver's init() function to run so that it can
   registr itself with the database/sql package. The trick to getting around this is to alias the
   package name to the blank identifier. This is standard practice for most of Go's SQL drivers. 

2. The sql.Open()

For now, we'll create a skeleton database model and have it return a bit of dummy data. It won't do
much, but I'd like to explain the pattern before we get into the nitty-gritty of SQL queries. 

## Designing a database model 

In this chapter we're going to sketch out a database model for our project. 

If you don't like the term model, you might want to think it as a service layer or data access layer
instead. Whatever you prefer to call it, the idea is tha we will encapsulate the code for working with
MySQL in a separate package to the rest of our application. 

## Template actions and functions 

In this section we're going to look at the template actions and functions that Go provides. 

We've already talked about some of the actions - {{define}}, {{template}} and {{block}} - but there are
three more which you can use to control the display of dynamic data - {{if}}, {{with}} and {{range}}. 

|Action <img width=300/>|Description <img width=500/>|
|:----|:---|
|{{if .Foo}} C1 {{else}} C2 {{end}} | If **.Foo** is not empty then render the content C1, otherwise render the content C2.|
|{{with .Foo}} C1 {{else}} C2 {{end}} | If **.Foo** is not empty, then set dot to the value of **.Foo** and render the content C1, otherwise render the content C2. |
|{{range .Foo}} C1 {{else}} C2 {{end}} | If the length of **.Foo** is greater than zero then loop over each element, setting dot to the value of each element and rendering the content C1. If the length of **.Foo** is zero then renader the content C2. The underlying type of **.Foo** must be an array, slice, map, or channel. |

There are a few things about these actions to point out: 
1. For all three actions the {{else}} clause is optional. For instance, you can write {{if .Foo}} C1
   {{end}} if there's no C2 content that you want to render. 

2. The empty values are false, 0, any nil pointer or interface value, and any array, slice, map, or
   string of length zero. 

3. It's important to grasp that the **with** and **range** actions change the value of dot. Once you
   start using them, what dot represent can be different depending on where you are in the template and
   what you're doing. 

The **html/template** package also provides some template functions which you can use to add extra
logic to your templates and control what is rendered at runtime. You can find a complete listing of
functions here, but the most important ones are: 

|Function<img width=300/>| Description <img width=500/>|
|:----|:----|
|{{eq .Foo .Bar}} | Yields true if **.Foo** is equal to **.Bar** |
|{{ne .Foo .Bar}} | Yields true if **.Foo** is not equal to **.Bar** |
|{{not .Foo}} | Yields the boolean negation of **.Foo** | 
|{{or .Foo .Bar}} | Yields **.Foo** if **.Foo** is not empty; otherwise yields **.Bar** |
|{{index .Foo i}} | Yields the value of **.Foo** at index **i**. The underlying type of **.Foo** must be a map, slice or array, and **i** must be an integer value. |
|{{printf "%s-%s" .Foo .Bar}} | Yields a formatted string containing the **.Foo** and **.Bar** alues. Works in the same way as fmt.Sprintf(). |
|{{len .Foo}} | Yields the length of **.Foo** as an integer. |
|{{$bar := len .Foo}} | Assign the length of **.Foo** to the template variable **$bar** |

The final row is an example of declaring a template variable. Template variables are particularly
useful if you want to store the result from a function and use it in multiple places in your template.
Variable names must be prefixed by a dollar sign and can contain alphanumeric characters only.  

## Using the with action 

A good opportunity to use the **{{with}}** action is the **view.tmpl** file that we created in the
previous chapter. Go ahead and update it like so: 

## Caching templates
Before we add any more functionality to our HTML templates, it's good time to make some optimizations
to our codebase. There are two main at the moment: 

## Catching runtime errors 
As soon as we begin adding dynamic behavior to our HTML templates there's a risk of encountering
runtime errors. 

Let's add a deliberate error to the **view.tmpl** template and see wht happens.

## Common dynamic data 

In some web applications there may be common dynamic data that you want to include on more than one -
or even every - webpage. For example, you might want to include the name and profile picture of the
current user, or a CSRF token in all pages with forms. 

In our case let's begin with something simple, and say that we want to include the current year in the
footer on every page. 

To do this we'll begin by adding a new CurrentYear field to the templateData struct, like so.

The next step is to add a **newTemplateData()** helper method to our application, which will return a
**templateData** struct initialized with the current year. 

## Custom template functions 

In the last part of this section about templating and dynamic data, I'd like to explain how to create
your own custom functions to use in Go templtes. 

To illustrate this, let's create a custom **humanDate()** function which outputs datatimes in a nice
'humanized' format like 02 Jan 2022 at 15:04, instead of outputing dates in the default format of
2022-01-02 15:04:00 +0000 UTC like we are currently 

There are two steps to do this: 
1. We need to create a template.FuncMap object containing the custom **humanDate()** function 
2. We need to use the **template.Funcs()** method to register this before parsing the templates. 

# Middleware 

When you're building a web application there's probably some shared functioinality that you want to use
for many (or even all) HTTP requests. For example, you might want to log every request, compress every
response, or check a cache before passing the request to your handlers. 

A common way of organizing this shared functionality is to set it up as middleware. This is essentially
some self-contained code which independency acts on a request before or after your normal application
handlers. 

In this section of the book you'll learn: 
1. An idiomatic pattern for **building and using custom middleware** which is compatible with
   **net/http** and many third-party packages. 

2. How to create middleware which **sets useful security headers** on every HTTP response. 

3. How to create middleware which **logs the request** received by your application 

4. How to create middleware which **recovers panics** so that they are gracefully handled by your
   application 

5. How to create and use composable **middleware chains** to help manage and organize your middleware. 

## How middleware works 

## Setting security headers 

Let's put the pattern we learned in the previous chapter to use, and make our own middleware which
automatically adds the following HTTP security headers to every response, inline with current OWASP
guidance. 

## Request Logging 

Let's continue in the same vein and add some middleware to log HTTP requests. Specifically, we're going
to use the information logger that we created earlier to record the IP address of the user, and which
URL and method are being requested. 

# CHAPTER 7 

# Advanced routing 

In the next section of this book we're going to add a HTML form to our application so that users can create
new snippets. 

To make this work smoothly we'll first need to update our application routes so that requests to
/snippet/create are handled differently based on the request method. Specifically: 
1. For GET /snippet/create requests we want to show the user the HTML form for adding a new snippet. 
2. For POST /snippet/create requests we want to process this form data and then insert a new snippet record
   into our database. 

Essentially, we want to rejig our application routes and handles so that they end up looking like this: 
| Method | Pattern | Handler | Action |
|:-------|:--------|:--------|:-------|
| GET | / | home | Display the home page |
| GET | /snippet/view/:id | snippetView | Display a specific snippet |
| GET | /snippet/create | snippetCreate | Display a HTML form for creating a new snippet |
| POST | /snippet/create | snippetCreatePost | Create a new snippet |

## Choosing a router 
Ther a literally hundreds of third-party routers for Go to pick from. And (fortunately for unfortunately,
depending on your perspective) they all work a bit differently. They have different APIs, different logic for
matching routes, and different behavioral quirks. 

Out of all the third-party routers I've tried there are three that I recommend as a starting point:

## User authentication 

In this section of the book we're going to add some user authenticatioin functionality to our
application, so that only registered, logged-in users can create new snippets. Non-logged-in users will
still be able to view the snippets, and will also be able to sign up for an account.
For our application, the process will work like this: 

1. A user will register by visiting a form at /user/signup and entering their name, email address and
   password. We'll store this information in a new users database table (which we'll create in a
   moment).

2. A user will log in by visiting a form at /user/login and entering their email address and password. 

3. We will then check the database to see if the email and password they entered match one of the users
   in the users table. If there's a match, the user has authenticated successfully and we add the
   relevant id value for the user to their session data, using the key "authenticatedUserID".

4. When we recieve any subsequent requests, we can check the user's session data for a
   "authenticatedUserID" value. If it exists, we know that the user has already successfully logged in.
   We can keep checking this until the session expires, when the user will need to log in again. 

## Routes Setup 

Let's begin this section by adding five new routes to our application, so that it looks like this: 

| Method | Pattern | Handler | Action |
|:---|:----|:----|:---|:----|
| GET | / | home | Display the home page |
| GET | /snippet/view/:id | snippetView | Display a specific snippet | 
| GET | /snippet/create | snippetCreate | Display a HTML form for creating a new snippet |                

