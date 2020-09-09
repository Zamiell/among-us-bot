# among-us-bot

A Discord bot that helps to keep track of the players in an ongoing game of <a href="https://store.steampowered.com/app/945360/Among_Us/">Among Us</a> as well as provide a waiting list.

<br />

## Installation

* `cd /root/`
* `git clone git@github.com:Zamiell/among-us-bot.git`
* `sqlite3 database.sqlite3 < ./install/database_schema.sql`
* `cp ./install/supervisor/among-us-bot.conf /etc/supervisor/conf.d/among-us-bot.conf`
* `supervisorctl reread`
* `supervisorctl start among-us-bot`
