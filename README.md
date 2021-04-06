# sub2notion

Subscribe to a public Notion page for changes and run a command on every detected change.

## Use Case

To publish my blog I'm using [loconotion](https://github.com/leoncvlt/loconotion).
It turns Notion pages into static websites. Normally you would need to manually run the command every time you change something on that Notion page.

But with sub2notion you can automatically call this command (or any other command) whenever the provided public page changes.

Currently only the provided page id gets checked for changes, not the subpages. I wanted to keep it simple and it's enough for my use case because even on subpage changes you can just add a empty line at the end of the root page and it will change the modified time of that page (even if you undo the change) and therefore trigger the chosen command.

## How To Use

Set the following env variables

- NOTION_PAGEID (the id of the public page you want to get notified about)
- NOTION_CMD (the command to run)
  - I set the full path to a bash script that starts a loconotion Docker container, creates a new git commit and pushes that commit
- NOTION_CMD_ARGS (optional, arguments for the command)
- NOTION_CMD_INTERVAL_SECONDS (How many seconds between the checks, defaults to 60)

### Note

Currently it does not work on Windows (except if you run inside WSL).

## Building

Clone this repository and then run `go build -o sub2notion -ldflags="-s -w" app.go`

## Special Thanks

- kjk for [notionapi](https://github.com/kjk/notionapi)
- leoncvlt for [loconotion](https://github.com/leoncvlt/loconotion)
