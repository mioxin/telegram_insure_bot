# gak_telegram_bot
__gak_telegram_bot__ is Telegram Bot that help to conclude an insurance contract for employees by an employer. It has the ability to expand functionality.

__gak_telegram_bot__ work in the multithread mode. Config will let you agile setting bot. Abailable the next settings in config file: allow or deny user access to the bot; the users in admin group get access for a special administrative commands; access word give access for any command if user not in admin group (for example for demo); set log file. You can modify config file and setting will update online without terminat worling. A command line options is more priotity then a config file settings.

__gak_telegram_bot__ has a commands:
1. /calc - An user input same data about business in interative mode. The bot get some information fron goverment site. After collect all info the bot colculate sum of the contract.
2. /send - An user can send a documents or a scans that include data for colculate sum of contract to the bot. 
3. /get - An user in admin group can get it in any time.
4. /help - The short help for bot commands.
5. /about - About insurance company.

Using:
---
__gak_telegram_bot <-t security_token> <-v> <-d> <-c config_file> <-h>__

#### Flags:
    -t:             The Telegram Bot API security token. The value override secure token in config file.
    -v, -verbouse:  Output log to StdOut (Default: off)
    -d:             The debug mode output more serv info to log (Default: off)
    -c:             The config file (Default: bot.cfg)
    -h, -help:      Show help (Default: false)

Any option have default value.

The config file will contains addition options.
---
#### List of user (comma separated) for deny access, (available key "all")

deny = user1, user2

#### List of user (comma separated) for allow access, (available key "all")

allow = all

#### Access word for denyed user for get access to any command. It need input after the command.
#### For example: /get <acc_word>

acc_word = oQkkj1!_k

#### List of users (comma separated) for access to admin command

adm_group = Firecat1

