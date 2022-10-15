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


