package server

// App object.
type App struct {
    config *Config
}

// Init method used to initialize our application.
func Init(config *Config) (app *App) {
    return &App{config: config}
}