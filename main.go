package main

import (
	"./irc" // https://github.com/thoj/go-ircevent
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Start configuration

var server = "irc.esper.net:6667"
var channels = []string{"#example", "#example2"}
var bot = irc.IRC("Mandatch", "Mibbit")
var default_topic = "default"
var web_port = 80

// End configuration

var myIP = "0.0.0.0" // do not touch - automatically obtained.
var lastcmd = int64(0)

func NLSplit(str string) []string {
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

func listDir(dir string) []string {
	return []string{}
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
	if strings.HasPrefix(event.Message, "! ") {
		if isAllowed(string(event.Nick)) {
			channel := event.Arguments[0]
			msg := strings.Split(event.Message, " ")
			Target := ""
			Topic := ""

			if msg[0] != "!" {
				return
			}
			ct, _ := strconv.ParseInt(time.Now().UTC().Format("20060102150405"), 10, 64)
			if ct-5 < lastcmd {
				return
			}
			lastcmd = ct
			if len(msg) < 2 {
				return
			}
			if msg[1] != "" {
				Target = msg[1]
			} else {
				return
			}
			if msg[1] == "*" {
				Target = "Everyone"
			}
			if len(msg) == 2 {
				Topic = default_topic
			}
			if len(msg) == 3 {
				Topic = msg[2]
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

func WebServer() {
	head := "<html><body><h1><a href='/''>HOME</a></h1>"
	foot := "</body></html>"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		derp := ""
		derp = derp + head
		derp = derp + "<h2>Topics</h2><h3><a href='/alist'>View Access List</a></h3><ul>"
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
		derp = derp + "</ul>"
		derp = derp + foot
		fmt.Fprintf(w, derp)
	})
	http.HandleFunc("/topics/", func(w http.ResponseWriter, r *http.Request) {
		channell := r.URL.Path[len("/topics/"):]
		derp := ""
		derp = derp + head
		if strings.HasPrefix(channell, ",") {
			channel := strings.Replace(channell, ",", "#", -1)
			derp = derp + "<h2>Topics (" + channel + ")</h2><ul>"
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
			derp = derp + "</ul>"
		}
		derp = derp + foot
		fmt.Fprintf(w, derp)
	})
	http.HandleFunc("/topic/", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/topic/"):]
		channel := ""
		if strings.HasPrefix(title, ",") {
			ddd := strings.Split(title, "/")
			channel = ddd[0]
			title = ddd[1]
			channel = strings.Replace(channel, ",", "#", -1)
		}
		if isTopic(string(title), channel) {
			derp := ""
			derp = derp + head
			derp = derp + channel + "/" + title
			derp = derp + "<pre>" + getTopic(title, channel) + "</pre>"
			derp = derp + foot
			fmt.Fprintf(w, derp)
		}
	})
	http.HandleFunc("/alist", func(w http.ResponseWriter, r *http.Request) {
		derp := ""
		derp = derp + head
		derp = derp + "<h2>Access List</h2>"
		derp = derp + "<ul>"
		al := getAList()
		for i := range al {
			derp = derp + "<li>" + al[i] + "</li>"
		}
		derp = derp + foot
		fmt.Fprintf(w, derp)
	})
	if web_port > 0 {
		wp := strconv.Itoa(web_port)
		fmt.Println("Serving web requests on port " + string(wp))
		http.ListenAndServe(":"+string(wp), nil)
	}
}

func main() {
	bot.Connect(server)
	bot.AddCallback("001", onConnect)
	bot.AddCallback("PRIVMSG", onMsg)
	bot.AddCallback("KICK", onKick)
	WebServer()
	bot.Loop()
}
