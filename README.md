# Ts3-Support-Bot
This repository provides you with a simple Teamspeak support bot.

## How to build
1. First things first, you need to install [Go](https://go.dev/)
2. Copy the repository
3. Install the dependency
```shell
go get github.com/multiplay/go-ts3
```
4. Then build your executable with
```shell
go build .
```
If you have any problems on Linux, try the following for building it:
```shell
CGO_ENABLED=0 go build .
```

### Removing depandency:
If you want to remove the depandency, after you build your teamspeakbot just do:
```shell
go get github.com/multiplay/go-ts3@none
```

## Configuration
The configuration is really straight straightforward. All configuration can be done in the JSON file ([config.json](https://github.com/loliebop/Ts3-Support-Bot/blob/main/config.json)).

```JSON
{
    "serverip": "localhost:10011", // IP:QueryPort 
    "User": "serveradmin",
    "Password": "",
    
    "ServerID": 1, // The teamspeak instance (normally its 1, if you dont have multiple teamspeak server instances)
    "SupportChannel": "3", // The support channel ID, where the client gets the option to create a support ticket.
    "TS3DefaultChannel": "1", // The channel in which you join after connecting to your teamspeak server.

    "Teams": { // "Teamname": [servergroup IDs which will be messaged at ticket creation]
        "scp": [9,7],
        "ts3": [11],
        "forum": [10]
    },

    "Messages": {
        "ticketCreated": "%s created a Ticket for %s!", // First %s is the clients name. Second %s is the teamname.
        "channelJoinMessage": "Hi %v!\nWelcome to the support waiting room! Use one of the following commands to create a ticket:" // %s is the clients name
    }
}
```