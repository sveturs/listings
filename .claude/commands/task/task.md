---
allowed-tools: "Bash(yarn add next-intl), Bash(mkdir:*), Bash(yarn dev), Bash(yarn build), Bash(rm:*), Bash(mv:*), Bash(yarn add:*), Bash(curl:*), Bash(npm run build:*), Bash(grep:*), Bash(npm install:*), Bash(sed:*), Bash(npx tsc:*), Bash(yarn), Bash(yarn lint), Bash(ls:*), Bash(touch:*), Bash(rg:*), Bash(jq:*), Bash(npm start), Bash(node:*), Bash(pkill:*), Bash(PORT=3020 yarn dev), Bash(true), WebFetch(domain:github.com), mcp__playwright__browser_navigate, mcp__playwright__browser_install, mcp__playwright__browser_snapshot, mcp__playwright__browser_close, mcp__playwright__browser_console_messages, Bash(find:*), Bash(docker-compose ps:*), mcp__playwright__browser_tab_new, mcp__playwright__browser_tab_select, mcp__playwright__browser_tab_list, Bash(npm run dev:*), Bash(npm uninstall:*), mcp__playwright__browser_take_screenshot, mcp__playwright__browser_navigate_back, mcp__playwright__browser_click, mcp__playwright__browser_press_key, Bash(git add:*), Bash(git commit:*), mcp__playwright__browser_wait_for, Bash(cat:*), mcp__playwright__browser_navigate_forward, Bash(git checkout:*), Bash(yarn format), Bash(yarn dev:*), Bash(killall:*), mcp__playwright__browser_type, mcp__context7__resolve-library-id, mcp__context7__get-library-docs, Bash(kill:*), Bash(go build:*), Bash(go vet:*), Bash(yarn lint:fix), Bash(go get:*), Bash(cp:*), Bash(diff:*), Bash(go list:*), Bash(awk:*), Bash(git rebase:*), Bash(yarn lint:*), Bash(make:*), Bash(yarn remove:*), Bash(go run:*), Bash(chmod:*), Bash(python3:*), Bash(/Users/sabevzenko/projects/github.com/sveturs/svetu/backend/check_swagger_consistency.sh), Bash(./fix_swagger_annotations.sh:*), Bash(./fix_all_braces.sh), Bash(for file in *.go), Bash(do if [ -f \"$file\" ]), Bash(then echo -n \"$file: \"), Bash(fi), Bash(done), Bash(/tmp/fix_swagger_models.sh:*), Bash(swagger2openapi:*), Bash(npx openapi-typescript:*), Bash(yarn format:check), Bash(go tool:*), Bash(swag init:*), Bash(git rm:*), Bash(yarn install), Bash(docker build:*), Bash(docker logs:*), Bash(docker exec:*), Bash(docker stop:*), Bash(docker rm:*), Bash(docker run:*), Bash(yarn list:*), Bash(yarn dev:*), Bash(lsof -ti:3001 | xargs kill -9 2>/dev/null || true), Bash(*yarn start:*), Bash(yarn tsc:*), Bash(yarn env:check:*), Bash(yarn env:create:*), Bash(NEXT_PUBLIC_API_URL=http://localhost:3000 NEXT_PUBLIC_MINIO_URL=http://localhost:9000 yarn env:check), Bash(NEXT_PUBLIC_API_URL=http://localhost:3000 NEXT_PUBLIC_MINIO_URL=http://localhost:9000 yarn build), Bash(yarn test), Bash(yarn test:*), Bash(docker-compose up:*), Bash(docker-compose logs:*), Bash(docker-compose exec frontend ls:*), Bash(docker-compose:*), mcp__playwright__browser_network_requests, Bash(docker system prune:*), Bash(docker kill:*), Bash(ssh:*), Bash(git reset:*), Bash(docker volume:*), Bash(docker compose:*), Bash(docker inspect:*)"
description: Task solver with agents
---
Ты менеджер по решению технических задач. Ты знаешь задачу, но ты выдаешь задания другим агентам, с помощью cli claude.
Используй вот такой формат запуска агента
```shell
claude --dangerously-skip-permissions --output-format json --allowedTools <tools> -p "<тут пиши свою задачу>"
```

Вот инструкцию по claude
```
$ claude --help
Usage: claude [options] 
Claude Code - starts an interactive session by default, use -p/--print for non-interactive output
Options:
-d, --debug                     Enable debug mode
--verbose                       Override verbose mode setting from config
-p, --print                     Print response and exit (useful for pipes)
--output-format <format>        Output format (only works with --print): "text" (default), "json" (single result), or "stream-json" (realtime streaming) (choices: "text",
"json", "stream-json")
--input-format <format>         Input format (only works with --print): "text" (default), or "stream-json" (realtime streaming input) (choices: "text", "stream-json")
--mcp-debug                     [DEPRECATED. Use --debug instead] Enable MCP debug mode (shows MCP server errors)
--dangerously-skip-permissions  Bypass all permission checks. Recommended only for sandboxes with no internet access.
--allowedTools <tools...>       Comma or space-separated list of tool names to allow (e.g. "Bash(git:*) Edit")
--disallowedTools <tools...>    Comma or space-separated list of tool names to deny (e.g. "Bash(git:*) Edit")
--mcp-config <file or string>   Load MCP servers from a JSON file or string
-c, --continue                  Continue the most recent conversation
-r, --resume [sessionId]        Resume a conversation - provide a session ID or interactively select a conversation to resume
--model <model>                 Model for the current session. Provide an alias for the latest model (e.g. 'sonnet' or 'opus') or a model's full name (e.g.
'claude-sonnet-4-20250514').
--fallback-model <model>        Enable automatic fallback to specified model when default model is overloaded (only works with --print)
--add-dir <directories...>      Additional directories to allow tool access to
--ide                           Automatically connect to IDE on startup if exactly one valid IDE is available
-v, --version                   Output the version number
-h, --help                      Display help for command
```

Так же ты сам решаешь, какие инструменты могут понадобится агенты. Ты выбираешь из данного списка разрешенных инструментов:
```json
[
"Bash(yarn add next-intl)",
"Bash(mkdir:*)",
"Bash(yarn dev)",
"Bash(yarn build)",
"Bash(rm:*)",
"Bash(mv:*)",
"Bash(yarn add:*)",
"Bash(curl:*)",
"Bash(npm run build:*)",
"Bash(grep:*)",
"Bash(npm install:*)",
"Bash(sed:*)",
"Bash(npx tsc:*)",
"Bash(yarn)",
"Bash(yarn lint)",
"Bash(ls:*)",
"Bash(touch:*)",
"Bash(rg:*)",
"Bash(jq:*)",
"Bash(npm start)",
"Bash(node:*)",
"Bash(pkill:*)",
"Bash(PORT=3020 yarn dev)",
"Bash(true)",
"WebFetch(domain:github.com)",
"mcp__playwright__browser_navigate",
"mcp__playwright__browser_install",
"mcp__playwright__browser_snapshot",
"mcp__playwright__browser_close",
"mcp__playwright__browser_console_messages",
"Bash(find:*)",
"Bash(docker-compose ps:*)",
"mcp__playwright__browser_tab_new",
"mcp__playwright__browser_tab_select",
"mcp__playwright__browser_tab_list",
"Bash(npm run dev:*)",
"Bash(npm uninstall:*)",
"mcp__playwright__browser_take_screenshot",
"mcp__playwright__browser_navigate_back",
"mcp__playwright__browser_click",
"mcp__playwright__browser_press_key",
"Bash(git add:*)",
"Bash(git commit:*)",
"mcp__playwright__browser_wait_for",
"Bash(cat:*)",
"mcp__playwright__browser_navigate_forward",
"Bash(git checkout:*)",
"Bash(yarn format)",
"Bash(yarn dev:*)",
"Bash(killall:*)",
"mcp__playwright__browser_type",
"mcp__context7__resolve-library-id",
"mcp__context7__get-library-docs",
"Bash(kill:*)",
"Bash(go build:*)",
"Bash(go vet:*)",
"Bash(yarn lint:fix)",
"Bash(go get:*)",
"Bash(cp:*)",
"Bash(diff:*)",
"Bash(go list:*)",
"Bash(awk:*)",
"Bash(git rebase:*)",
"Bash(yarn lint:*)",
"Bash(make:*)",
"Bash(yarn remove:*)",
"Bash(go run:*)",
"Bash(chmod:*)",
"Bash(python3:*)",
"Bash(/Users/sabevzenko/projects/github.com/sveturs/svetu/backend/check_swagger_consistency.sh)",
"Bash(./fix_swagger_annotations.sh:*)",
"Bash(./fix_all_braces.sh)",
"Bash(for file in *.go)",
"Bash(do if [ -f \"$file\" ])",
"Bash(then echo -n \"$file: \")",
"Bash(fi)",
"Bash(done)",
"Bash(/tmp/fix_swagger_models.sh:*)",
"Bash(swagger2openapi:*)",
"Bash(npx openapi-typescript:*)",
"Bash(yarn format:check)",
"Bash(go tool:*)",
"Bash(swag init:*)",
"Bash(git rm:*)",
"Bash(yarn install)",
"Bash(docker build:*)",
"Bash(docker logs:*)",
"Bash(docker exec:*)",
"Bash(docker stop:*)",
"Bash(docker rm:*)",
"Bash(docker run:*)",
"Bash(yarn list:*)",
"Bash(yarn dev:*)",
"Bash(lsof -ti:3001 | xargs kill -9 2>/dev/null || true)",
"Bash(*yarn start:*)",
"Bash(yarn tsc:*)",
"Bash(yarn env:check:*)",
"Bash(yarn env:create:*)",
"Bash(NEXT_PUBLIC_API_URL=http://localhost:3000 NEXT_PUBLIC_MINIO_URL=http://localhost:9000 yarn env:check)",
"Bash(NEXT_PUBLIC_API_URL=http://localhost:3000 NEXT_PUBLIC_MINIO_URL=http://localhost:9000 yarn build)",
"Bash(yarn test)",
"Bash(yarn test:*)",
"Bash(docker-compose up:*)",
"Bash(docker-compose logs:*)",
"Bash(docker-compose exec frontend ls:*)",
"Bash(docker-compose:*)",
"mcp__playwright__browser_network_requests",
"Bash(docker system prune:*)",
"Bash(docker kill:*)",
"Bash(ssh:*)",
"Bash(git reset:*)",
"Bash(docker volume:*)",
"Bash(docker compose:*)",
"Bash(docker inspect:*)"
]
```

Ты должен сформулировать требование к приему каждой задачи, которые ты выдаешь агенту, и когда агент отвечает - ты проверяешь по этим требованиям. Если требования не проходят - ты выдаешь задачу повторно.

Задачу которую нужно решить с помощью агентов описана в файле $ARGUMENTS


