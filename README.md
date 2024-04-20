# dbwf-ls
Databricks workflow yaml language server.

An "intellisense" server to help with creating databricks workflow file.

## Why?

Databricks workflow file is a way to define a workflow on databricks using `json` file.
However, there are 2 problems with that:

- `json` file doesn't allow comments. Yes most the time it's not needed. But when it's needed, json can't do it.
- There's no code suggestions, no documentation close to source, no diagnostics, etc.

This language server aims to solve those problems by:

- Use `yaml` file instead of json, then it can be parsed to json using any language.
With a slight custom file extension, `.flow.yaml`.
This means it still get the `yaml` syntax highlighting while the language server only handles the juicy parts.
- Well this is the juicy part

## What?

This server comes with the following capabilities:

```golang
TextDocumentSync
HoverProvider
DefinitionProvider
CodeActionProvider
CompletionProvider
DocumentFormattingProvider
```

## Demo

Will be here, at some point

## ⚙️  Setup

Note: Before going further, you should know that this is currently a setup for `Neovim`, simply because neovim works with custom binary language server seamlessly while other editors don't.
I don't know about `Emacs` but VSCode sure is a pain to make this work so, for now, it doesn't. Just `Neovim`

First clone and navigate to this repository with your terminal (I'm assuming you're a fantastic developer).

Then build and save this binary to config path. If it errors, create the `~/.config/dbfw-ls` in advance is a good try but it really shouldn't:

```bash
go build main.go && mv ./main ~/.config/dbwf-ls
```

Then save `language_client_config/nvim.lua` somewhere in your nvim config.
Let's say you have a simple setup and you have a root `init.lua` and you save this file as `lua/custom/plugins/dbwf-ls.lua`.

Now you can import this into your root config:

```lua
require("custom.plugins.dbwf-ls")
```

If you have a more complex config, I think you know how to deal with this.

With this setup, you have the `dbwf-ls` server attaches to any `.flow.yaml` file you open or create. All log for the current session will be saved at `~/.config/dbfw-ls/log.txt`

Enjoy!

## Dependencies

None. Because I like to type. Also avoid dependencies hell.

## 👏 Contributing

Let's see if this is useful for more than just me.

