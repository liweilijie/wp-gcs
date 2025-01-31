package main

import (
	"context"
	logging "github.com/ipfs/go-log/v2"
	"wp-gcs/config"
	"wp-gcs/database"
	"wp-gcs/gcs"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

var log = logging.Logger("main")

func main() {
	logging.SetDebugLogging()
	log.Info("gcs start...")
	app := config.InitConfig()

	db := database.InitDB(*app)
	database.InitialMigration(db)

	wpUploadsHandle := database.New(db)

	ctx := context.Background()
	gcs := gcs.NewGCS(ctx, app.GOOGLE_APPLICATION_PROJECT_ID, app.GOOGLE_APPLICATION_BUCKET)

	run(app, wpUploadsHandle, gcs)
}
