# Oscar R Acuña — portfolio

A monochrome, manual-style portfolio served by Go. No JavaScript,
frontend framework, external fonts, or asset build step.

## Run locally

Requires Go 1.24 or later. From the repository root:

```sh
go run ./cmd/web
```

Open http://localhost:8080. The default address is `:8080`; override it with
`PORT=8081`, `PORT=:8081`, or a full address such as `PORT=127.0.0.1:8081`.
An optional `.env` file is supported (see `.env.example`). Restart after Go
changes; template and CSS edits only need a browser refresh.

## Edit

- `templates/home.page.tmpl`: biography, project descriptions, and links.
- `templates/base.layout.tmpl`: document layout, title, and search description.
- `static/style.css`: monochrome palette, typography, index, and responsive layout.
- `static/favicon.svg`: terminal prompt favicon.

The layout borrows from Unix manual pages: a contents index, concise sections,
and monospace type. All navigation uses ordinary links. There are no animations.
Keyboard focus styles and a skip link are included.

The copy draws on Oscar's stated IT operations experience, Linux/Proxmox lab,
Go projects, and RHCSA study game. It does not claim a completed certification.
Review biographical wording and public links before publishing.

## Check

```sh
go test ./...
go vet ./...
go build -o bin/portfolio ./cmd/web
```

Serve the binary from the repository root so it can find `templates/` and
`static/`. Deployment is separate; these commands do not publish or push.
