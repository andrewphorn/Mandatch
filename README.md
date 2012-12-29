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

Configuring Mandatch
--------------------

Configuring your brand spanking new copy of Mandatch is dead simple.

* Open up `config.json` in your favorite text editor (Preferably something like Notepad++ or Sublime Text. DO NOT USE NOTEPAD.)
* Change the various settings to whatever you want.
* Save and close `config.json`
* Run Mandatch.exe (or whatever the executable might be called on your system) and enjoy!


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