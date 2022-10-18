package main 

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" 
    "github.com/go-playground/form/v4"
    "flag"
    "log"
    "net/http"
    "html/template"
    "os"
    /* 
     * Import the models package that we just created. You need to prefix this with 
     * whatever model path you set up back in chapter 02.01 (Project Setup and Creating
     * a Module) so that the import statement looks like this: 
     * "{your-module-path}/internal/models". If you can't remember what module path you 
     * used, you can find it at the top of the go.mod file. 
     */
    "github.com/stoneyzjw/snippetbox/internal/models"
    _ "github.com/go-sql-driver/mysql"
)

/* 
 * Define an application struct to hold the application-wide dependencies for the 
 * web application. For now we'll only include fields for the two custom loggers, but 
 * we'll add more to it as the build progress. 
 */
type application struct {
    errorLog    *log.Logger 
    infoLog     *log.Logger
    snippets    *models.SnippetModel
    templateCache map[string]*template.Template
    // Add a formDecorder field to hold a pointer to a form.Decorder instance
    formDecoder *form.Decoder
}

func main() {
    /* 
     * Define a new command-line flag with the name 'addr', a default value of ":4000"
     * and some short help text explaining what the flag controls. The value of the 
     * falg will be stored in the addr variable at runtime.
     */
    addr := flag.String("addr", ":4000", "HTTP network address")

    /* Define a new command-line flag for the MySQL DSN string. */ 
    dsn := flag.String("dsn", "stoney:Ly_2123320@/snippetbox?parseTime=true", "MySQL data source name")
    /* 
     * Importantly, we use the flag.Parse() function to parse the command-line flag. 
     * This reads in the command-line flag value and assigns it to the addr 
     * variable. You need to call this *before* you ue the addr variable 
     * otherwise it will always contain the default value of ":4000". If any errors are 
     * encoutered during parsing the application will be terminated. 
     */
    flag.Parse()

    /* 
     * Use log.New() to create a logger for writing information messages. This takes 
     * three parameters: the destination to write the logs to (os.Stdout), a string 
     * prefix for message (INFO followed by a tab), and flags to indicate what 
     * additional information to include (local date and time), Note that the flags 
     * are joined using the bitwise OR operator |.
     */
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)

    /* 
     * Create a logger for writing error messages in the same way, but use stderr ad 
     * the destination and use the log.Lshortfile flag to include the relevant
     * file name and line number
     */
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

    /* 
     * To keep the main() function tidy I've put the code for creating a connection 
     * pool into the separate openDB() function below. We pass openDB() the DSN 
     * from the command-line flag.
     */
    db, err := openDB(*dsn) 
    if err != nil {
        errorLog.Fatal(err)
    }

    /*
     * Wealso defer a call to db.Close(), so that the connection pool is closed 
     * before the main() function exits. 
     */
    defer db.Close()

    // Initialize a new template cache ... 
    templateCache, err := newTemplateCache() 
    if err != nil {
        errorLog.Fatal(err)
    }

    // Initialize a decoder instance ... 
    formDecoder := form.NewDecoder() 

    /* 
     * Initialize a new instance of our application struct, containing the 
     * dependencies. 
     */ 
    app := &application {
        errorLog:  errorLog, 
        infoLog:   infoLog,
        snippets:  &models.SnippetModel{DB: db},
        templateCache: templateCache, 
        formDecoder: formDecoder,
    }

    /*
     * Initialize a new http.Server struct. We set the Addr and Handler fields so 
     * that the server uses the same network address and routes as before, and set 
     * the ErrorLog field so that the server now uses the custom errorLog logger in 
     * the event of any problems. 
     */ 
     srv := &http.Server {
         Addr:          *addr, 
         ErrorLog:      errorLog, 
         Handler:       app.routes(),
     }
     infoLog.Printf("Starting server on %s", *addr)
     /* 
      * Because the err variable is now already declared in the code above, we need 
      * to ue the assignment operator = here, instead of the := 'declare and assign' 
      * operator. 
      */
     err = srv.ListenAndServe() 
     errorLog.Fatal(err)

}

/* for a given DSN */
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err 
    }
    if err = db.Ping(); err != nil {
        return nil, err 
    }
    return db, nil
}
