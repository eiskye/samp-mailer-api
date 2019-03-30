package server

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "net/smtp"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

// A Mailer represents an SMTP server.
type Mailer struct {
    Server		string
    Auth   		smtp.Auth
}

// Message storing mail message.
type Message struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// Map is a map[string]interface{} with additional helpful functionality.
type Map map[string]interface{}

// Run fires up the HTTP server and blocks until failure.
func (app *App) Run() {
    router := mux.NewRouter()

    router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
        // server respond with a json
        w.Header().Add("Content-Type", "application/json")

        // read body.
        bodyBytes, err := ioutil.ReadAll(r.Body)
        defer r.Body.Close()

        if err != nil {
            panic(err)
        }

        // we unmarshal our byteArray which contains our
        // jsonFile's content into 'msg' which we defined above.
        var msg Message
        json.Unmarshal(bodyBytes, &msg)

        // body will hold the entire email data we need send.
        var buf bytes.Buffer

        body := fmt.Sprintf("From: %s <%s>\r\n", msg.From, app.config.SMTP.Email)
        body += fmt.Sprintf("To: %s\r\n", msg.To)
        body += fmt.Sprintf("Subject: %s\r\n", msg.Subject)

        // need to set the mime version and context type so that
        // the html can be remdered properly in the email sent.
        body += fmt.Sprint("MIME-version: 1.0\r\n")
        body += fmt.Sprint("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
        body += msg.Body
        body += fmt.Sprint("\r\n")

        // write the body.
        buf.Write([]byte(body))

        mailer := makeAuth(app.config.SMTP.Server, app.config.SMTP.Email, app.config.SMTP.Password, app.config.SMTP.Port)
        fmt.Printf("Sending mail to %s from %s to %s\n", mailer.Server, app.config.SMTP.Email, msg.To)

        err = smtp.SendMail(mailer.Server, mailer.Auth, msg.From, []string{msg.To}, buf.Bytes())
        if err != nil {
            fmt.Printf("Error delivering email: %s", err)

            json.NewEncoder(w).Encode(Map{
                "message": "Sorry, there was a problem sending email.",
                "success": false,
            })

            return
        }

        json.NewEncoder(w).Encode(Map{
            "message": "Sucessfull - email has been sent.",
            "success": true,
        })

    }).Methods("POST")

    allowedOrigins := handlers.AllowedOrigins([]string{"*"})
    allowedMethods := handlers.AllowedMethods([]string{"POST"})

    log.Print("Application started. Press CTRL+C to shut down.")
    http.ListenAndServe(":" + app.config.Server.BindAddr, handlers.CORS(allowedOrigins, allowedMethods)(router))
}

// returns a mailer and make auth with SMTP server.
func makeAuth(host string, username string, password string, port int) *Mailer {
    return &Mailer{
        Auth:   smtp.PlainAuth("", username, password, host),
        Server: fmt.Sprintf("%s:%d", host, port),
    }
}