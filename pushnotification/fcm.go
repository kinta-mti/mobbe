package pushnotification

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func Init(ctx context.Context) (*messaging.Client, error) {
	log.Print("[pushnotification.init] called!!")
	opt := option.WithCredentialsFile("ypgmerchant-5bc86-firebase-adminsdk-tly3r-804a007e11.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	// Get the FCM object
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return fcmClient, nil
}
