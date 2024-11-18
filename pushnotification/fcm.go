package pushnotification

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FcmClient *messaging.Client

func Init(ctx context.Context) error {
	log.Print("[pushnotification.init] called!!")
	opt := option.WithCredentialsFile("ypgmerchant-5bc86-firebase-adminsdk-tly3r-804a007e11.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}
	// Get the FCM object
	_fcmClient, err := app.Messaging(ctx)
	FcmClient = _fcmClient
	if err != nil {
		return err
	}
	return nil
}

func SendNotification(
	token string,
	message string,
) error {
	//Send to One Token
	_, err := FcmClient.Send(context.Background(), &messaging.Message{
		Token: token,
		Data: map[string]string{
			message: message,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
