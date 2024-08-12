package bootstrap

import "spectator.main/internals/mongo"

type Application struct {
	Config   *Config
	Mongo mongo.Client
}

func App() Application {
	app := &Application{}
	app.Config = InitConfig()
	app.Mongo = NewMongoDatabase(app.Config)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}