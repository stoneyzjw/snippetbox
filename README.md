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
