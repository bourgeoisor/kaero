<img src="https://github.com/user-attachments/assets/f34c127b-4199-4123-bb03-1142b102c348" width="200"/>

## What is Kaero?

Kaero is a simple terminal-based IRC client made in Golang. In Japanese, kaero 帰る means "to return (home)" which I find poetic for an IRC client. Its name is also close in sounds to kaeru かえる which means "frog" (thus the logo).

![screely-1736785076382](https://github.com/user-attachments/assets/a29c2621-36b5-4879-8f9c-fe0a5f89b00f)

## Features

These features are currently supported.

- 125 server replies (out of ~150 well known numerics) implemented
- 31 user commands implemented
- Can join up to 9 channels in one server
- TLS support

## Not supported

These features are not supported yet but may be implemented in the future.

- User prefixes
- Channel prefixes
- Config file
- Customizable colors
- Customizable keybinds
- Highlight mentions
- Direct messages
- CTCP
- DCC
- IRCv3
- CAP negotiations
- Identity registration
- User list as a side panel
- Multi-server support

## Keybinds

- `Ctrl + C`: quit the application
- `Enter`: toggle the input bar, or send message is non-empty
- `Alt + [0-9]`: swap between channels
- `PgUp` / `PgDown`: scroll through the chat log

## Commands

- `/ping`
- `/quit [<reason>]`
- `/nick <nickname>`
- `/oper <name> <password>`
- `/motd [<target>]`
- `/version [<target>]`
- `/admin [<target>]`
- `/connect <target server> [<port> [<remote server>]]`
- `/lusers`
- `/time [<server>]`
- `/stats <query> [<server>]`
- `/help [<subject>]`
- `/info`
- `/join <channel>{,<channel>} [<key>{,<key>}]`
- `/part <channel>{,<channel>} [<reason>]` (alias: `/leave ...`)
- `/topic <channel> [<topic>]`
- `/names <channel>{,<channel>}`
- `/list [<channel>{,<channel>}] [<elistcond>{,<elistcond>}]`
- `/invite <nickname> <channel>`
- `/kick <channel> <user> *( "," <user> ) [<comment>]`
- `/who <mask>`
- `/whois [<target>] <nick>`
- `/whowas <nick> [<count>]`
- `/kill <nickname> <comment>`
- `/rehash`
- `/restart`
- `/squit <server> <comment>`
- `/away [<text>]`
- `/links`
- `/userhost <nickname>{ <nickname>}`
- `/wallops <text>`

## How to run?

I am not currently publishing binaries, so they need to be built manually.

1. Clone the repository.

   ```sh
   git clone https://github.com/bourgeoisor/kaero
   cd kaero/
   ```

3. Fetch the required dependencies and build the binary.

   ```sh
   go get
   go build
   ```

5. Run the binary.

   ```sh
   ./kaero
   ```
