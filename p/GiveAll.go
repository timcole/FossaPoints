package p

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var (
	users map[string]ChannelChatters
)

// GiveAll gives everyone in chat points
func GiveAll(w http.ResponseWriter, r *http.Request) {
	tmi := getChatters(w)
	users = getPoints(w)
	var keys []*datastore.Key
	var keysData []*ChannelChatters

	_ = users
	for _, chatter := range tmi.Chatters.Staff {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "staff")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}
	for _, chatter := range tmi.Chatters.Admins {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "admin")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}
	for _, chatter := range tmi.Chatters.GlobalMods {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "global_mod")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}
	for _, chatter := range tmi.Chatters.Moderators {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "moderator")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}
	for _, chatter := range tmi.Chatters.Vips {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "vip")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}
	for _, chatter := range tmi.Chatters.Viewers {
		fmt.Println(chatter)
		key, data := manageCurrency(chatter, "viewer")
		keys = append(keys, key)
		keysData = append(keysData, &data)
	}

	putAllPoints(w, keys, keysData)
	w.Write([]byte("points added"))
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

func manageCurrency(user, role string) (*datastore.Key, ChannelChatters) {
	var chatter ChannelChatters
	if users[user].Chatter == user {
		chatter = users[user]
		chatter.LastSeen = time.Now()
		chatter.Currency += 5
	} else {
		chatter = ChannelChatters{
			Chatter:   user,
			Currency:  5,
			Role:      role,
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		}
	}

	anc := datastore.NameKey("Channels", "modesttim", nil)
	anc.Namespace = "FossaBot"
	key := datastore.NameKey("Points", chatter.Chatter, anc)
	key.Namespace = "FossaBot"

	return key, chatter
}

func getChatters(w http.ResponseWriter) *Chatters {
	req, err := http.NewRequest("GET", "https://tmi.twitch.tv/group/user/modesttim/chatters", nil)
	must(w, err)

	res, err := (&http.Client{Timeout: time.Minute}).Do(req)
	must(w, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	must(w, err)

	var tmi = &Chatters{}
	err = json.Unmarshal(body, &tmi)
	must(w, err)

	return tmi
}

func putAllPoints(w http.ResponseWriter, keys []*datastore.Key, data []*ChannelChatters) bool {
	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	must(w, err)

	_, err = client.PutMulti(ctx, keys, data)
	must(w, err)

	return true
}

func getPoints(w http.ResponseWriter) map[string]ChannelChatters {
	var ctx = context.Background()
	client, err := datastore.NewClient(ctx, "timcole-me")
	must(w, err)

	anc := datastore.NameKey("Channels", "modesttim", nil)
	anc.Namespace = "FossaBot"
	query := datastore.NewQuery("Points").Namespace(anc.Namespace).Ancestor(anc)

	var chanChatters []ChannelChatters
	var chanChattersMap = make(map[string]ChannelChatters)
	keys, err := client.GetAll(ctx, query, &chanChatters)
	must(w, err)

	for i, k := range keys {
		chanChatters[i].Chatter = k.Name
		chanChattersMap[k.Name] = chanChatters[i]
	}

	return chanChattersMap
}

type ChannelChatters struct {
	Chatter   string    `datastore:"-" json:"chatter"`
	Currency  int       `datastore:"currency" json:"currency"`
	Role      string    `datastore:"role" json:"role"`
	FirstSeen time.Time `datastore:"first_seen" json:"first_seen"`
	LastSeen  time.Time `datastore:"last_seen" json:"last_seen"`
}

type Chatters struct {
	ChatterCount int64         `json:"chatter_count"`
	Chatters     ChattersClass `json:"chatters"`
}

type ChattersClass struct {
	Vips       []string `json:"vips"`
	Moderators []string `json:"moderators"`
	Staff      []string `json:"staff"`
	Admins     []string `json:"admins"`
	GlobalMods []string `json:"global_mods"`
	Viewers    []string `json:"viewers"`
}
