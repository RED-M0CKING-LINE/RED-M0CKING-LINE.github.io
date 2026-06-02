# EA Website

---

## Quickstart

Make environment variables based on template

Set environment variables
```bash
set -a ; source .env ; set +a ;
```

### Native

```bash
# build wasm tooling
make wasm

# run dev server on :8080
make run
```

### Podman (or Docker)

```bash
# build and run the container
docker compose -f ./deploy/compose/compose.yaml up --build
```

### Kube (maybe this will work idc yet)

```bash
kubectl apply -f deployment.yaml
```

---

## Configuration

Use .env. See env.example

`COOKIE_SECRET` is required for production. You can generate it with:

```bash
make cookie-secret
# OR
openssl rand -hex 32
```

---

### Certificates

#### Real Certs

acme.sh, cron, have nginx reference token directory
#TODO finished docs

#### Dev Certs
```
openssl req -x509 -nodes -newkey rsa:4096 -days 3650 -keyout ./www.ethanashley.net.key -out ./www.ethanashley.net.crt -subj "/CN=www.ethanashley.net"
```

#### Dummy Certs (for dummy server)
```
openssl req -x509 -nodes -newkey rsa:2048 -days 3650 -keyout ./dummy.key -out ./dummy.crt -subj "/CN=_"
```


---

## OIDC setup

Set these four envs to enable authentication, otherwise the site will run in dev mode.

```
OIDC_ISSUER=https://idp.example/
OIDC_CLIENT_ID=...
OIDC_CLIENT_SECRET=...
OIDC_REDIRECT_URL=https://app.example/auth/callback
OIDC_SCOPES=openid,profile,email
```

---

## Blog
The container is designed to be stateless so all content is baked into the container. This requires a rebuild of the container to obtain new content.
Additionally, all runtime states are managed in cookies so that the backend can remain stateless

### Markdown authoring

Drop a file into `content/blog/<slug>.md`:

```markdown
---
title: "Title"
date: 2026-05-21T09:00:00Z
updated: 2026-06-1T10:00:00-04:00
summary: "One line summary"
author: "Me"
tags: ["topic", "meta"]
draft: false
---

Markdown body

GFM (tables, task lists, strikethrough), footnotes, typographer punctuation, and autolinks are enabled.

Output is sanitized!!!
```

Posts are loaded on startup
Drafts (`draft: true`) are hidden in `prod`, visible in `dev`

---

### Atom feed

Its Atom, its a feed, and Im here to feed it to you!

Generated from `content/blog/*.md` at every request, served at: `/feed.xml`

---

## WASM

The `/tools` page loads tools with `WebAssembly.instantiateStreaming`. 

To add another tool:
1. Add `wasm/<name>/main.go` exporting JS functions via `js.Global().Set`.
2. Add a Make rule and a section to `web/templates/pages/tools.html`.
3. Wire the form/result handlers in `web/static/js/tools.js`.

---

## CI/CD

`.github/workflows/ci.yaml` runs on every push/PR:

- `gofmt -l .` (fails on any drift)
- `go vet ./...`
- `go test -race ./...`
- builds the server binary and the WASM tool

On push to `main` or a `v*` tag, the image is built and (if `REGISTRY` secret
is set) pushed. A subsequent deploy job SSHes to `DEPLOY_HOST` and runs
`podman pull` + `podman run` to roll the container.

Secrets - none baked in, all placeholders:

| Secret              | Purpose                                       |
| ------------------- | --------------------------------------------- |
| `REGISTRY`          | Hostname of container registry (optional).    |
| `REGISTRY_USERNAME` | Registry login.                               |
| `REGISTRY_PASSWORD` | Registry token.                               |
| `DEPLOY_HOST`       | SSH target.                                   |
| `DEPLOY_USER`       | SSH user (needs `podman` perms).              |
| `DEPLOY_SSH_KEY`    | Private key.                                  |
| `DEPLOY_KNOWN_HOSTS`| Output of `ssh-keyscan -H <host>`.            |

Without `REGISTRY`, the publish job builds locally and the deploy job is
skipped - lint/test/build still gate every PR.

---

## Notes
- Inline CSS styles are stripped, they must go in the CSS file
- Both nginx and the app apply security headers, so changes must happen in both places


