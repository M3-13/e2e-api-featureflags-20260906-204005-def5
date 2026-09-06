# Feature-Flag-Service

Ein HTTP-basierter Feature-Flag-Service in Go, der ausschließlich die
Standardbibliothek (`net/http`) verwendet. Er verwaltet Feature-Flags in einem
thread-sicheren In-Memory-Store und bewertet Flags für einzelne Nutzer
deterministisch anhand eines stabilen Hashs und eines `rollout_percent`.

## Tech Stack

- **Sprache**: Go (1.22+)
- **Framework**: `net/http` (Standardbibliothek, `ServeMux` mit Methoden- und Wildcard-Mustern)
- **Speicher**: In-Memory mit `sync.Mutex`

## Installation

```bash
go mod download
```

## Ausführen

```bash
go run .
```

Der Server lauscht auf Port `8080`. Es sind keine Umgebungsvariablen
erforderlich.

## Endpunkte

| Methode | Pfad                     | Beschreibung                                        |
|---------|--------------------------|-----------------------------------------------------|
| POST    | `/flags`                 | Legt ein Flag an (Body: `{key, enabled, description?, rollout_percent?}`) → 201 |
| GET     | `/flags`                 | Listet alle Flags → 200 (leer: `[]`)                |
| GET     | `/flags/{key}`           | Liefert ein Flag → 200 / 404                        |
| PUT     | `/flags/{key}`           | Ändert `enabled`, `description`, `rollout_percent` → 200 / 404 |
| DELETE  | `/flags/{key}`           | Entfernt ein Flag → 204 / 404                       |
| GET     | `/flags/{key}/evaluate?user={id}` | Bewertet ein Flag für einen Nutzer → 200 `{"result":bool}` / 404 |
| GET     | `/healthz`               | Health-Check → 200 `{"status":"ok"}`                |

Fehlerantworten (400, 404, 500) enthalten ausschließlich das JSON-Fehlerobjekt
`{"error":"<msg>"}`.

## Features

- Thread-sicherer In-Memory-Store für Flag-Definitionen (`key`, `enabled`,
  `description`, `rollout_percent`).
- Deterministische Rollout-Bewertung pro Nutzer über einen stabilen Hash.
- Health-Endpoint `/healthz`.
- Begrenzung der Request-Body-Größe auf 1 MiB.
