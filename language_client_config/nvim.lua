--[[
Require this file to enable

Simplest way is to require it in your root `init.lua`

Let's say you save this file as `./custom/plugins/dbwf-ls.lua`:
	
	require("custom.plugins.dbwf-ls")
]]
-- BufRead* alone will not recognise newly created file from netrw buf
vim.api.nvim_create_autocmd({ "BufRead", "BufNewfile" }, {
	pattern = "*.flow.yaml",
	callback = function()
		local client_dbwf = vim.lsp.get_active_clients({ filter = { name = "dbwf-ls" } })[1]
		if not client_dbwf then
			client_dbwf = vim.lsp.start_client {
				name = "dbwf-ls",
				cmd = { vim.fn.expand("$HOME/.config/dbwf-ls/main") },
				on_attach = require "custom.plugins.init"["on_attach"]
			}

			if not client_dbwf then
				vim.notify "Cannot start client"
				return
			end
		else
			client_dbwf = client_dbwf.id
		end
		local buf = vim.api.nvim_get_current_buf()
		vim.lsp.buf_attach_client(buf, client_dbwf)
	end,
})
