

<p align="center"><img src="./logo.png" alt="PartyPin" width="120" /></p>

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)
![PWA](https://img.shields.io/badge/PWA-ready-5A0FC8)
![Uploads](https://img.shields.io/badge/Uploads-Local%20FS-informational)
![License](https://img.shields.io/badge/License-MIT-green)

Code‑locked event gallery. Create an **organizer PIN** (event ID), share it with guests, and let them upload photos under simple rules (gallery on/off, allow any uploads, require *taken today*, etc.). Runs as a Go server with a PWA front‑end.

---

## What it does

- **Create event** → generates a short **eventId** and stores its config to `events/<eventId>.json`.
- **Guests upload photos** to `/upload` using the eventId.
- **Host controls rules** via the event config:
  - `AllowGallery`: enable a public gallery view.
  - `AllowAny`: allow any file (vs. stricter checks in your UI).
  - `RequireTakenToday`: enforce “taken today” rule (enforced client/UI + server-side checks you add).
- **Browse images** from `/images` or directly via `/uploads/...` (static).



---

## Quick start

```bash
go run .
# → Listening on http://localhost:6060
```

Project serves:
- **PWA UI** from `public/`
- **Event configs** from `events/` (created on demand)
- **Uploads** from `uploads/` (served at `/uploads/…`)

You can open `http://localhost:6060` to use the UI, or call the API directly.

---

## API

Base URL: `http://localhost:6060`

### 1) Create event

**POST** `/create-event`  
`Content-Type: application/json`

Request body (EventConfig):
```json
{
  "title": "Claire & Leo — Birthday",
  "allowGallery": true,
  "allowAny": false,
  "requireTakenToday": true
}
```

Response:
```json
{ "eventId": "ABC123" }
```

This writes `events/ABC123.json` with the provided config.

---

### 2) Get event config

**GET** `/event-config?eventId=ABC123`

Typical response:
```json
{
  "title": "Claire & Leo - Birthday",
  "allowGallery": true,
  "allowAny": false,
  "requireTakenToday": true
}
```

> The exact shape/params are handled in `internal.HandleEventConfig`. The UI calls this to hydrate pages.

---

### 3) Upload image

**POST** `/upload`  
`Content-Type: multipart/form-data`

Form fields:
- `eventId` — the event code you got from `/create-event`
- `image` — the image file to upload

Example:
```bash
curl -F eventId=ABC123 \
     -F image=@photo.jpg \
     http://localhost:6060/upload
```

Response: JSON describing the stored file (shape defined in `internal.HandleUpload`), or a 4xx/5xx error.

---

### 4) List images

**GET** `/images?eventId=ABC123`

Returns a JSON list of images for that event (shape defined in `internal.HandleImages`), e.g.,
```json
[
  { "url": "/uploads/ABC123/1730123456_aa12bb_photo.jpg", "name": "photo.jpg" }
]
```

---

### 5) Static uploads

All uploaded files are available under:
```
/uploads/<eventId>/<stored-filename>
```
They are served read‑only via `http.FileServer`.


## Notes on safety

- **PINs are short** by design for convenience. Add rate‑limiting and auth if you expose this publicly.
- Enforce **content rules** (size/type/exif “taken today”) in both client and server.
- Consider **moderation** (human or automated) if `AllowGallery` is enabled.

---

## Configuration

This minimal server runs with defaults:
- Port: **6060**
- Storage: local `events/` and `uploads/`

You can front it with a reverse proxy (Traefik, Nginx) and add TLS, caching, and auth.

