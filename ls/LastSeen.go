package p

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

var (
	users map[string]ChannelChatters
)

// LastSeen gives everyone in chat points
func LastSeen(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		if ok := must(w, fmt.Errorf("no user param")); !ok {
			return
		}
	}
	user = strings.ToLower(user)

	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	if ok := must(w, err); !ok {
		return
	}

	anc := datastore.NameKey("Channels", "modesttim", nil)
	anc.Namespace = "FossaBot"
	uKey := datastore.NameKey("Points", user, anc)
	uKey.Namespace = "FossaBot"

	var (
		uChatter ChannelChatters
		uChannel Channel
	)
	err = client.Get(ctx, anc, &uChannel)
	if ok := must(w, err); !ok {
		return
	}
	err = client.Get(ctx, uKey, &uChatter)
	if ok := must(w, err); !ok {
		return
	}

	w.Write([]byte(fmt.Sprintf("%s was last seen %s", user, uChatter.LastSeen.Format(time.RFC1123))))
}

func must(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(500)
		w.Write([]byte("failed: " + err.Error()))
		fmt.Println("failed: " + err.Error())
		return false
	}
	return true
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
