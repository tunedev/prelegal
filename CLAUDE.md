# Prelegal Proect
---

## Overview

> This is a SaaS product to allow users to draft legal agreements based on templates in the templates directory. The user uses an AI chat in order to establish what document they want and how to fill in the fields. The available documents are covered in catalog.json file in the project root, included here: 

@catalog.json

The AI chat is not yet implemented — see Implementation Status below for what currently exists.

## Development process
---

When instructed to build a feature:

1. Use your Atlassian tools to read the feature instructions from Jira
2. Develop the feature - do not skip any step from the feature-dev step process
3. Thoroughly test the feature with unit and integration tests and fix any issues
4. Submit a PR using your github tools

## AI design
---

When writing code to make calls to LLMs, use your streaming ensure calls to OpenRouter is relatively looking fast to the user the model to use are `openai/gpt-oss-120b:free`, `qwen/qwen3-next-80b-a3b-instruct:free`, `openrouter/free`. you should use Structured outputs so that you can interpret the results and populate fields in the legal document.

There is an `OPENROUTER_API_KEY` in the .env file in the root.

## Technical design
---

The entire project should be packaged into a Docker container.   
The backend should be in backend/ using golang and echo.   
The frontend should still be in frontend/ 
The database should use sqlite and should be created form scratch each time the Docker container is brought up, allowing for a users table with signup and signin. 
Consider statically building the frontend and serving it via golang, if that will work.   
There should be scripts in scripts/ for:  
```bash
# Mac
scripts/start-mac.sh # Start
scripts/stop-mac.sh  # Stop

# Linux
scripts/start-linux.sh # Start
scripts/stop-linux.sh # Stop

# Windows
scripts/start-windows.ps1 # Start
scripts/stop-windows.ps1 # Stop
```
Backend available at http://localhost:8000

## Color Scheme
- Accent yellow: `#ecad0a`
- Blue Primary: `#209dd7`
- Purple Secondary: `#753991` (submit buttons)
- Dark Navy: `#032147` (headings)
- Gray Text `#888888`

## Implementation Status

- **PL-1 (done)**: Mutual NDA creator frontend (SvelteKit), form-driven, PDF export via browser print.
- **PL-2 (done)**: CommonPaper legal document template dataset added to `templates/`.
- **PL-3 (done)**: V1 technical foundation — Go/Echo backend in `backend/`, SQLite recreated from scratch on each startup with a `users` table (schema only, no signup/signin endpoints yet), frontend statically built and embedded into the Go binary via `go:embed` with SPA fallback, single multi-stage `Dockerfile`, and start/stop scripts in `scripts/` for Mac/Linux/Windows. Login is currently fake: it sets a client-side flag and does not call the backend.
- **Not yet built**: AI chat, OpenRouter integration, real signup/signin auth, and support for any document type other than the Mutual NDA.
