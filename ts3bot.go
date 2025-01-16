package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/multiplay/go-ts3"
)

/*
	Chat:
	- TESTME: PingtTeamler (wenn online (multible clients))

*/

type Config struct {
	QueryIP  string `json:"serverip"`
	User     string
	Password string

	ServerID          int
	SupportChannel    string
	TS3DefaultChannel string

	Teams map[string]int
}

var cfg Config
var BotID string

func getNameFromUID(client *ts3.Client, clid string) string {
	type nameResponse struct {
		Name string `ms:"client_nickname"`
	}

	res := &nameResponse{}

	cmd := ts3.NewCmd("clientinfo")
	cmd.WithArgs(
		ts3.NewArg("clid", clid),
	).WithResponse(&res)

	_, err := client.ExecCmd(cmd)
	if err != nil {
		slog.Error("sendMsg threw an error", "err", err)
		return ""
	}

	return res.Name
}

func sendMsg(client *ts3.Client, clid string, msg string) {
	cmd := ts3.NewCmd("sendtextmessage")
	cmd.WithArgs(
		ts3.NewArg("targetmode", "1"),
		ts3.NewArg("target", clid),
		ts3.NewArg("msg", msg),
	)
	if _, err := client.ExecCmd(cmd); err != nil {
		slog.Error("sendMsg threw an error", "err", err)
	}
}

func sendMsgInt(client *ts3.Client, clid int, msg string) {
	cmd := ts3.NewCmd("sendtextmessage")
	cmd.WithArgs(
		ts3.NewArg("targetmode", "1"),
		ts3.NewArg("target", clid),
		ts3.NewArg("msg", msg),
	)
	if _, err := client.ExecCmd(cmd); err != nil {
		slog.Error("sendMsg threw an error", "err", err)
	}
}

func getUserByGroup(client *ts3.Client, gid int) {
	cl, err := client.Server.ClientList("-groups")
	if err != nil {
		slog.Error("getUserByGroup Error", "error", err)
		return
	}

	for _, user := range cl {
		groups := *user.ServerGroups

		for _, group := range groups {
			if group != gid {
				fmt.Println(group, gid)
				continue
			}

			sendMsgInt(client, user.ID, "Es wurde ein Ticket erstellt!")
		}
	}
}

func moveUser(client *ts3.Client, userid string, cid string) {
	cmd := ts3.NewCmd("clientmove")
	cmd.WithArgs(
		ts3.NewArg("cid", cid),
		ts3.NewArg("clid", userid),
	)
	_, err := client.ExecCmd(cmd)
	if err != nil {
		slog.Error("Error moving Client", "error", err)
	}
}

func createChannel(client *ts3.Client, invokerName string, invokerID string, team string) {
	channelName := strings.ToUpper(team) + " | \u200b" + invokerName

	type channelCreateRespone struct {
		ChannelID string `ms:"cid"`
	}
	i := &channelCreateRespone{}

	cmd := ts3.NewCmd("channelcreate").WithArgs(
		ts3.NewArg("channel_name", channelName),
		ts3.NewArg("cpid", 3),
	).WithResponse(&i)

	if _, err := client.ExecCmd(cmd); err != nil {
		slog.Error("", "err", err)
	}

	moveUser(client, invokerID, i.ChannelID)
	moveUser(client, BotID, cfg.TS3DefaultChannel)
}

func LoadConfig() {
	var config Config

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}

	cfg = config
}

func main() {
	LoadConfig()

	c, err := ts3.NewClient(cfg.QueryIP)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err := c.Login(cfg.User, cfg.Password); err != nil {
		log.Fatal(err)
	}

	if v, err := c.Version(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Supportbot is running:", *v)
	}

	if err := c.Use(cfg.ServerID); err != nil {
		log.Fatal(err)
	}

	coninfo, err := c.Whoami()
	if err != nil {
		log.Print("Coudn't retrive Bot ID:", err)
	}
	BotID = strconv.Itoa(coninfo.ClientID)

	if err := c.SetNick("Support Bot"); err != nil {
		log.Print("Bot couldn't set his Nickname:", err)
	}

	eventHandler(c)
}

func moveEvent(client *ts3.Client, data map[string]string) {
	clID := data["clid"]
	ctID := data["ctid"]

	if ctID == cfg.SupportChannel {
		msg := fmt.Sprintf("Hi %v!\nWillkommen im Support Warteraum! Benutze ein der folgenden Befehle für dein Support Ticket:\n", getNameFromUID(client, clID))
		sendMsg(client, clID, msg)

		for teamNames, _ := range cfg.Teams {
			sendMsg(client, clID, "!"+teamNames)
		}

		return
	}
}

func textEvent(client *ts3.Client, data map[string]string) {
	var invokerID string = data["invokerid"]
	var invokerName string = data["invokername"]
	var msg string = data["msg"]
	var targetmode string = data["targetmode"]

	if targetmode != "1" {
		return
	}
	if invokerID == BotID {
		return
	}
	if string(msg[0]) != "!" {
		return
	}

	teamID, exist := cfg.Teams[msg[1:]]
	if !exist {
		return
	}

	createChannel(client, invokerName, invokerID, msg[1:])
	getUserByGroup(client, teamID)
}

func eventHandler(client *ts3.Client) {
	if err := client.Register(ts3.ChannelEvents); err != nil {
		log.Fatal("Coudn't register ChannelEvents: ", err)
	}

	if err := client.Register(ts3.TextPrivateEvents); err != nil {
		log.Fatal("Coudn't register TextPrivateEvents: ", err)
	}

	for notification := range client.Notifications() {
		switch notification.Type {
		case "clientmoved":
			moveEvent(client, notification.Data)
		case "clientleftview":
			fmt.Println("disconnected")
		case "cliententerview":
			fmt.Println("cliententerview")
		case "textmessage":
			textEvent(client, notification.Data)
		}
	}
}
