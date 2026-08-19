***REMOVED*** gak_telegram_bot
__gak_telegram_bot__ is Telegram Bot that help to conclude an insurance contract for employees by an employer. It has the ability to expand functionality.

__gak_telegram_bot__ work in the multithread mode. Config will let you agile setting bot. Abailable the next settings in config file: allow or deny user access to the bot; the users in admin group get access for a special administrative commands; access word give access for any command if user not in admin group (for example for demo); set log file. You can modify config file and setting will update online without terminat worling. Command-line options have priority over config files and environment variables.

__gak_telegram_bot__ has a commands:
1. /calc - An user input same data about business in interative mode. The bot get some information fron goverment site. After collect all info the bot colculate sum of the contract.
2. /send - An user can send a documents or a scans that include data for colculate sum of contract to the bot.
3. /get - An user in admin group can get it in any time.
4. /help - The short help for bot commands.
5. /about - About insurance company.

Using:
---
__gak_telegram_bot [-t security_token] [-v] [-d] [-c config_file] [-h]__

The Telegram token must be supplied with `-t` or the `TELEGRAM_BOT_TOKEN` environment variable. Do not store tokens in config files or commit them to Git.

***REMOVED******REMOVED******REMOVED******REMOVED*** Flags:
    -t, --token:    The Telegram Bot API security token.
    -v, --verbose:  Output log to StdOut (Default: off)
    -d, --debug:    Enable Telegram debug mode (Default: off)
    -c, --config:   The JSON config file (Default: bot.cfg)
    -h, --help:     Show help.

Any option have default value.

The config file is JSON and is loaded by Kong. The local `.env` file is also loaded by Kong's configuration loader and is ignored by Git.
---
```json
{
    "deny": "user1,user2",
    "allow": "all",
    "acc_word": "oQkkj1!_k",
    "adm_group": "Firecat1",
    "company": "Страховая Компания в г. Петропавловск (Казахстан)"
}
```

`.env` example:

```dotenv
TELEGRAM_BOT_TOKEN=your_token_here
```

