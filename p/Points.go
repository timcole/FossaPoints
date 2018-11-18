package p

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var (
	users map[string]ChannelChatters
)

// GetPoints gives everyone in chat points
func GetPoints(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		must(w, fmt.Errorf("no user param"))
	}

	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	must(w, err)

	anc := datastore.NameKey("Channels", "modesttim", nil)
	anc.Namespace = "FossaBot"
	uKey := datastore.NameKey("Points", user, anc)
	uKey.Namespace = "FossaBot"

	var (
		uChatter ChannelChatters
		uChannel Channel
	)
	err = client.Get(ctx, anc, &uChannel)
	must(w, err)
	err = client.Get(ctx, uKey, &uChatter)
	must(w, err)

	w.Write([]byte(fmt.Sprintf("%s currently has %d %s", user, uChatter.Currency, uChannel.CurrencyName)))
}

func must(w http.ResponseWriter, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(500)
		w.Write([]byte("failed: " + err.Error()))
		fmt.Println("failed: " + err.Error())
		os.Exit(1)
	}
}

type Channel struct {
	CreatedAt    time.Time `datastore:"created_at" json:"created_at"`
	CurrencyGain int       `datastore:"currency_gain" json:"currency_gain"`
	CurrencyName string    `datastore:"currency_name" json:"currency_name"`
}

type ChannelChatters struct {
	Currency  int       `datastore:"currency" json:"currency"`
	Role      string    `datastore:"role" json:"role"`
	FirstSeen time.Time `datastore:"first_seen" json:"first_seen"`
	LastSeen  time.Time `datastore:"last_seen" json:"last_seen"`
}
