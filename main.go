package main

import (
	"./irc" // https://github.com/thoj/go-ircevent
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Channels     []string
	Server       string
	DefaultTopic string
	WebPort      int
	Prefix       string
	Nick         string
	Realname     string
	Cooldown     int64
	WebDesign    string
	Identify     bool
	Nickserv     string
	Command      string
}

// Start configuration
func config() bool {
	c, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = json.Unmarshal([]byte(c), &f)
	return true
}

//Get the config before more damage is done!
var f Config

var channels = []string{"#example", "#example2"}
var server = "irc.esper.net"
var bot = irc.IRC("Mandatch", "Mibbit")
var botnick = "Mandatch"
var default_topic = "default"
var web_port = 80
var prefix = "!"
var cooldown = int64(5)
var webdesign = "default.html"
var identify = false
var nickserv = "Nickserv"
var ncmd = "IDENTIFY mypassword"

// End configuration

var myIP = "0.0.0.0" // do not touch - automatically obtained.
var lastcmd = int64(0)

func NLSplit(str string) []string { // Utility to split a string by newline. Cross-platform.
	rstr := strings.Replace(string(str), "\r", "\n", -1)
	rstr = strings.Replace(string(rstr), "\n\n", "\n", -1)
	st := strings.Split(string(rstr), "\n")
	return st
}

func prepForIRC(str string, target string, channel string) string {
	dorp := strings.Replace(string(str), "{n}", "\x02"+target+"\x02", -1)
	dorp = strings.Replace(string(dorp), "{i}", myIP, -1)
	dorp = strings.Replace(string(dorp), "{c}", channel, -1)
	return dorp
}

func getAList() []string {
	list, err := ioutil.ReadFile("data/access.list")
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	return NLSplit(string(list))
}

func isAllowed(user string) bool {
	ulist := getAList()
	for i := range ulist {
		if strings.ToLower(user) == strings.ToLower(string(ulist[i])) {
			return true
		}
	}
	return false
}

func isTopic(topic string, channel string) bool {
	dirl, err := ioutil.ReadDir("data/topics")
	if err != nil {
		return false
	}
	for i := range dirl {
		if topic+".txt" == string(dirl[i].Name()) {
			return true
		}
		if channel == string(dirl[i].Name()) {
			dirll, err := ioutil.ReadDir("data/topics/" + channel)
			if err != nil {
				return false
			}
			for l := range dirll {
				if topic+".txt" == string(dirll[l].Name()) {
					return true
				}
			}
		}
	}
	return false
}

func getTopic(topic string, channel string) string {
	td, err := ioutil.ReadFile("data/topics/" + topic + ".txt")
	if err != nil {
		tdd, err := ioutil.ReadFile("data/topics/" + channel + "/" + topic + ".txt")
		if err != nil {
			return ""
		}
		return string(tdd)
	}
	return string(td)
}

func sendMessage(channel string, text string, target string) {
	if len(text) == 0 {
		return
	}
	msgs := NLSplit(text)

	if len(msgs) >= 3 {
		msgs = msgs[0:2]
	}
	for i := range msgs {
		bot.Privmsg(channel, prepForIRC(msgs[i], target, channel))
	}
	return
}

func onMsg(event *irc.Event) {
	fmt.Println("<["+event.Arguments[0]+"]", event.Nick+">", event.Message)
	if strings.HasPrefix(event.Message, prefix) {
		if isAllowed(string(event.Nick)) {
			channel := event.Arguments[0]
			msgg := event.Message[len(prefix):]
			msg := strings.Split(msgg, " ")
			Target := ""
			Topic := ""
			ct, _ := strconv.ParseInt(time.Now().UTC().Format("20060102150405"), 10, 64)
			if ct-cooldown < lastcmd {
				return
			}
			lastcmd = ct
			if len(msg) < 1 {
				return
			}
			if msg[0] != "" {
				Target = msg[0]
			} else {
				return
			}
			if msg[0] == "*" {
				Target = "Everyone"
			}
			if len(msg) == 1 {
				Topic = default_topic
			}
			if len(msg) == 2 {
				Topic = msg[1]
			}
			if Topic == "" {
				Topic = default_topic
			}
			Topic = strings.ToLower(Topic)
			if isTopic(string(Topic), channel) {
				Topicstr := getTopic(Topic, channel)
				sendMessage(channel, Topicstr, Target)
			}
		}
	}
	return
}

func onConnect(event *irc.Event) {
	//http://res.public-craft.com/myip.php - gets bot's host's IP.
	// TODO - Do this without involving my own services.
	resp, err := http.Get("http://res.public-craft.com/myip.php")
	if err != nil {
		myIP = "0.0.0.0"
	} else {
		dorp, _ := ioutil.ReadAll(resp.Body)
		myIP = string(dorp)
	}
	resp.Body.Close()
	if identify {
		bot.Privmsg(nickserv, ncmd)
	}
	for i := range channels {
		bot.Join(channels[i])
	}
}

func onKick(event *irc.Event) {
	if strings.ToLower(event.Message) != "leave" {
		if strings.ToLower(event.Arguments[1]) == strings.ToLower(bot.GetNick()) {
			bot.Join(event.Arguments[0])
			sendMessage(event.Arguments[0], "If you want me to stay out, kick me with the reason 'leave'.", "") // Yep
		}
	}
}
func renderPage(c string, p string) string {
	tpl, err := ioutil.ReadFile("data/html/" + webdesign)
	if err != nil {
		return "There was an issue loading the template data: File not Found."
	}
	ret := ""
	tpll := ""
	for i := range tpl {
		tpll = tpll + string(tpl[i])
	}
	// If somebody has a better solution for converting a []byte to a string, I'd like to hear it.
	// Tried, like, everything.
	ret = strings.Replace(tpll, "{b}", botnick, -1)
	ret = strings.Replace(ret, "{p}", p, -1)
	ret = strings.Replace(ret, "{c}", c, -1)
	return ret
}
func WebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := "Topics"
		derp := ""
		dirl, err := ioutil.ReadDir("data/topics")
		if err != nil {
			return
		}
		for i := range dirl {
			if strings.HasSuffix(dirl[i].Name(), ".txt") {
				n := dirl[i].Name()[:len(dirl[i].Name())-4]
				derp = derp + "<li><a href='/topic/" + n + "'>" + n + "</a></li>"
			}
			if strings.HasPrefix(dirl[i].Name(), "#") {
				ddd := strings.Replace(dirl[i].Name(), "#", ",", -1)
				derp = derp + "<li><a href='/topics/" + ddd + "'>" + dirl[i].Name() + "</a></li>"
			}
		}
		fmt.Fprintf(w, renderPage(derp, page))
	})
	http.HandleFunc("/topics/", func(w http.ResponseWriter, r *http.Request) {
		channell := r.URL.Path[len("/topics/"):]
		derp := ""
		page := ""
		if strings.HasPrefix(channell, ",") {
			channel := strings.Replace(channell, ",", "#", -1)
			derp = derp + "<h2>Topics (" + channel + ")</h2><ul>"
			page = "Topics (" + channel + ")"
			dirl, err := ioutil.ReadDir("data/topics/" + channel)
			if err != nil {
				return
			}
			for i := range dirl {
				if strings.HasSuffix(dirl[i].Name(), ".txt") {
					n := dirl[i].Name()[:len(dirl[i].Name())-4]
					derp = derp + "<li><a href='/topic/" + channell + "/" + n + "'>" + n + "</a></li>"
				}
			}
		}
		fmt.Fprintf(w, renderPage(derp, page))
	})
	http.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/topic/"):]
		channel := ""
		page := ""
		if strings.HasPrefix(title, ",") {
			ddd := strings.Split(title, "/")
			channel = ddd[0]
			title = ddd[1]
			channel = strings.Replace(channel, ",", "#", -1)
		}
		if isTopic(string(title), channel) {
			derp := ""
			page = "Viewing Topic (" + channel + "/" + title + ")"
			derp = derp + "<pre>" + getTopic(title, channel) + "</pre>"
			fmt.Fprintf(w, renderPage(derp, page))
		}
	})
	http.HandleFunc("/alist", func(w http.ResponseWriter, r *http.Request) {
		derp := ""
		page := "Access List"
		derp = derp + "<ul>"
		al := getAList()
		for i := range al {
			derp = derp + "<li>" + al[i] + "</li>"
		}
		fmt.Fprintf(w, renderPage(derp, page))
	})
	if web_port != 0 {
		wp := strconv.Itoa(web_port)
		fmt.Println("Serving web requests on port " + string(wp))
		http.ListenAndServe(":"+string(wp), nil)
	}
}

func main() {
	if !config() {
		fmt.Println("Error in loading config. Aborting.")
		return
	}
	channels = f.Channels
	server = f.Server
	bot = irc.IRC(f.Nick, f.Realname)
	default_topic = f.DefaultTopic
	web_port = f.WebPort
	prefix = f.Prefix
	cooldown = f.Cooldown
	webdesign = f.WebDesign
	botnick = f.Nick
	identify = f.Identify
	nickserv = f.Nickserv
	ncmd = f.Command

	bot.Connect(server)
	bot.AddCallback("001", onConnect)
	bot.AddCallback("PRIVMSG", onMsg)
	bot.AddCallback("KICK", onKick)
	WebServer()
	bot.Loop()
}
