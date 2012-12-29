Mandatch
========
Usage
-----

* Install your OS's golang compiler
* `$ go build`
* ./Mandatch(.exe|what have you)
* Spam channels.

Editing
-------

Want to edit Mandatch? Great! Do the github twiddledance and do whatever you want. Sell it for a profit, if you can! Want to claim you made it? Go for it! I'd like it if you didn't, but it's your choice.

Setting up Mandatch
-------------------

Setting Mandatch up isn't too hard. Ignoring the bit where you have to install a golang compiler, it's just 3-5 simple modifications to the source, and then a compile.

If you'll notice here...
```
// Start configuration

var server = "irc.esper.net:6667"
var channels = []string{"#example", "#example2"}
var bot = irc.IRC("Mandatch", "Mibbit")
var default_topic = "default"
var web_port = 80

// End configuration
```

The default IRC server is at `irc.esper.net` on port `6667`. Chances are, you don't want this. Change it to whatever you want.

`channels` is a list of channels. At the moment, there are two channels in the list. To add more, simply add a comma after the last channel, and add another quote-surrounded channel.

`bot` is a struct containing a bunch of information about the bot. The important bit here is the first argument to irc.IRC ("Mandatch") and the second argument ("Mibbit"). The first argument is the nickname. The second argument is the username and realname.

`default_topic` is the topic to use when a topic isn't specified in the command. To disable this, I guess you could just set it to nothing (two quotes next to eachother, like thus: ""). I have NO idea what that would do, and I have no intention of trying it.

`web_port` is what port the webserver should listen on. To disable the webserver, set `web_port` to `0`

Topics
------

# Adding Topics #

Adding topics to Mandatch is as easy as creating files in the `data/topics` directory.

Follow these simple steps to add topics!
* Locate and enter the `data` directory.
* Locate and enter the `topics` directory.
* Make a text file (.txt) using the name of the topic you want to add. For example, to add the topic `test`, the file would be named `text.txt`
* Open up the text file and add your words to it!

If you want to add a topic for a specific channel, after locating and entering the `topics` directory, follow these steps.
* Create a folder named `#channel`, except replace the part that says `channel` with your channel's name.
* Follow the rest of the steps given in the previous guide.

# Variables #

* `{i}` - Is replaced with the bot's IP when sent to the channel. `Example: {i} -> 127.0.0.1`
* `{n}` - Is replaced with the #target#'s name when sent to the channel. `Example: {n} -> AndrewPH`
* `{c}` - Is replaced with the channel the topic was issued in.


Reporting Bugs
--------------

If you come across a bug in a copy of Mandatch that has had only the configuration modified, please open an issue here on github.

Requesting Features
-------------------

Open an issue here on github. Or, if you want to be REALLY helpful, make a pull request with a copy of Mandatch that has the features you want.

Disclaimer
----------

I made Mandatch to suit MY needs. If you're wondering why x feature hasn't already been implemented, it's because I haven't needed it yet. I will gladly accept feature additions, though.