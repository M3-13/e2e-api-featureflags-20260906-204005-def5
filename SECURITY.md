VERDICT: CHANGES_REQUESTED

## Sicherheitsprüfung Feature-Flag-Service

Hinweis: Es wurde kein automatisierter Security-Scanner ausgeführt (`no applicable security scanners for this project type`). Die Bewertung basiert auf manueller Codeprüfung des zusammengeführten Stands.

### Erfüllte Sicherheitsanforderungen
- Body-Limit für POST/PUT korrekt über `http.MaxBytesReader` und `maxBodyBytes` = 1 MiB umgesetzt; Überschreitung liefert 413.
- Fehlerantworten enthalten ausschließlich das definierte JSON-Fehlerobjekt ohne interne Details.
- Das eigene Zugriffslog protokolliert nur Methode, Pfad ohne Query-String und Status; `user`-Werte werden nicht geloggt.
- Der In-Memory-Store speichert nur Flag-Definitionen, keine User-IDs aus Evaluate-Requests.
- Keine offensichtlichen Injection-Schwachstellen (SQL/Command/Pfad) oder unsichere Deserialisierung.

### Festgestellte Schwachstellen / Härtungslücken

#### 1. Fehlende Authentifizierung / Autorisierung
- **Schweregrad:** Mittel
- **Betroffene Stelle:** `main.go`, `newHandler()`/`main()`; alle `/flags*`-Routen
- **Risiko:** Jeder, der den Server erreicht, kann Feature Flags ansehen, anlegen, ändern und löschen. Feature Flags steuern häufig sichtbares Produktverhalten. Der Server lauscht auf `:8080` (alle Netzwerkinterfaces), was die Angriffsfläche vergrößert.
- **Konkreter Fix:** Vor Verwaltungsendpunkten eine Authentifizierung vorschalten, z. B. eine Middleware, die einen konfigurierbaren API-Key oder ein JWT im `Authorization`-Header prüft. Falls eine vollständige Auth-Lösung aktuell nicht im Umfang liegt, den Server mindestens an Loopback/private IP binden oder ausschließlich hinter einem authentifizierenden Reverse-Proxy betreiben. Die Tests und Dokumentation müssen dann entsprechend angepasst werden, sodass legitime Clients weiterhin zugreifen können.

#### 2. HTTP-Server ohne Timeouts
- **Schweregrad:** Mittel
- **Betroffene Stelle:** `main.go`, `main()`: `http.ListenAndServe(":8080", newHandler())`
- **Risiko:** Der verwendete Default-Server hat keine `ReadTimeout`-, `ReadHeaderTimeout`-, `WriteTimeout`- oder `IdleTimeout`-Werte. Dadurch sind Slowloris-/Slow-Body-Angriffe möglich, die Sockets und Goroutinen lange blockieren und den Dienst erschöpfen können.
- **Konkreter Fix:** Einen expliziten `http.Server` konfigurieren, z. B.:
  ```go
  srv := &http.Server{
      Addr:              ":8080",
      Handler:           newHandler(),
      ReadHeaderTimeout: 5 * time.Second,
      ReadTimeout:       10 * time.Second,
      WriteTimeout:      10 * time.Second,
      IdleTimeout:       60 * time.Second,
      MaxHeaderBytes:    1 << 20,
  }
  log.Fatal(srv.ListenAndServe())
  ```

#### 3. Log-Injection über nicht sanitisierte URL-Pfade
- **Schweregrad:** Niedrig (bis Mittel, abhängig vom Log-Konsumenten)
- **Betroffene Stelle:** `internal/httpapi/middleware.go`, `LoggingMiddleware`
- **Risiko:** `r.URL.Path` wird unmittelbar in `log.Printf("%s %s %d", ...)` verwendet. Ein Angreifer kann z. B. `%0A` in den Pfad einbauen, was nach dem URL-Decoding einen Zeilenumbruch erzeugt und gefälschte Logzeilen ermöglicht.
- **Konkreter Fix:** Pfad mit `%q` loggen oder explizit Steuerzeichen entfernen/escapen:
  ```go
  log.Printf("%s %q %d", r.Method, r.URL.Path, rec.status)
  ```

#### 4. Unverschlüsselter Transport und User-ID im Query-String
- **Schweregrad:** Niedrig
- **Betroffene Stelle:** `main.go`, `internal/httpapi/evaluate.go` (`r.URL.Query().Get("user")`)
- **Risiko:** User-IDs werden als Query-Parameter übertragen und können auf dem Transportweg oder in vorgelagerten Proxy-Logs sichtbar sein. Der eigene Logger ist korrekt, aber die Übertragung selbst ist unverschlüsselt.
- **Konkreter Fix:** TLS-Terminierung (z. B. `ListenAndServeTLS` mit gültigen Zertifikaten oder einen TLS-terminierenden Reverse-Proxy) einsetzen. Falls das API-Format weiterentwickelt werden kann, User-ID alternativ über einen nicht geloggten Header transportieren.

#### 5. JSON-Decoder akzeptiert zusätzliche JSON-Werte im Body
- **Schweregrad:** Niedrig
- **Betroffene Stelle:** `internal/httpapi/respond.go`, `DecodeJSON`
- **Risiko:** `json.NewDecoder(r.Body).Decode(v)` decodiert nur das erste JSON-Objekt. Ein Body wie `{"key":"k","enabled":true} {"key":"x"}` oder `{...}garbage` wird nicht vollständig als ungültig erkannt.
- **Konkreter Fix:** Nach der ersten `Decode` prüfen, dass kein weiterer JSON-Wert folgt, z. B.:
  ```go
  dec := json.NewDecoder(r.Body)
  if err := dec.Decode(v); err != nil { ... }
  if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
      WriteError(w, http.StatusBadRequest, "invalid JSON")
      return err
  }
  ```

Insgesamt liegt kein kritischer oder sofort ausnutzbarer Exploit im Code vor. Die genannten Punkte sollten vor einem produktiven Einsatz auf einem nicht isolierten Netz jedoch behoben werden.