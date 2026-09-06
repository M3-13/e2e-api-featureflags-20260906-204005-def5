VERDICT: CHANGES_REQUESTED

**Bewertungsrahmen**: `project_type: go-backend`, daher entfallen Pflichttexte/UI/Cookie-Banner/Barrierefreiheit. Der EU AI Act ist nicht anwendbar, weil keine KI-Funktion enthalten ist. Relevant sind DSGVO/GDPR und der Cyber Resilience Act (CRA).

## 1. DSGVO/GDPR

**Verarbeitete personenbezogene Daten**: Der Query-Parameter `user` in `GET /flags/{key}/evaluate?user={id}` ist potenziell personenbezogen. Die Verarbeitung ist transients; der Wert wird weder im In-Memory-Store gespeichert (`internal/store/store.go`) noch geloggt (`internal/httpapi/middleware.go`).

**Findings:**

- **Low — Rechtsgrundlage nicht im Code dokumentiert.** `internal/httpapi/evaluate.go` verarbeitet `r.URL.Query().Get("user")`, ohne dass im Projekt eine Rechtsgrundlage nach Art. 6 DSGVO erkennbar ist.  
  **Abhilfe:** In `README.md` oder `AGENTS.md` ein Datenschutzkapitel ergänzen: Zweck (deterministische Feature-Bewertung), Rechtsgrundlage (z. B. Art. 6 Abs. 1 lit. b/f DSGVO), Speicherdauer (keine dauerhafte Speicherung), Logging-Garantie (kein `user`-Logging). Der konkrete Betreiber muss die Rechtsgrundlage für seinen Einsatz prüfen.

- **Kein High/Critical-Befund:** AC-14 und AC-15 sind erfüllt. Es werden nur Methode, Pfad ohne Query-String und Statuscode geloggt. Der In-Memory-Store speichert ausschließlich Flag-Definitionen und keine User-IDs.

## 2. EU Cyber Resilience Act (CRA)

**Findings:**

- **High — Verwaltungs-API unverschlüsselt und ohne Authentifizierung.**  
  `main.go` ruft `http.ListenAndServe(":8080", newHandler())` auf. Der Dienst bindet damit an alle Netzwerk-Interfaces, ohne TLS und ohne Authentifizierung. Jeder mit Netzzugriff kann Feature-Flags anlegen, ändern oder löschen.  
  **Abhilfe:**  
  - Bind-Adresse konfigurierbar machen, Default `127.0.0.1:8080` statt `:8080`.  
  - `http.Server` mit TLS verwenden (`ServeTLS` oder vorgeschalteter TLS-Terminierung).  
  - Auth-Middleware für mindestens `POST /flags`, `PUT /flags/{key}` und `DELETE /flags/{key}` ergänzen (z. B. Bearer-Token/API-Key).  
  - Sichere Standardkonfiguration in `README.md` dokumentieren.

- **Medium — Keine HTTP-Server-Timeouts.**  
  `main.go` nutzt `http.ListenAndServe` und damit den Default-Server ohne `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout` und `IdleTimeout`.  
  **Abhilfe:** Eigenen `http.Server` in `main.go` erzeugen, z. B. mit `ReadHeaderTimeout: 5s`, `ReadTimeout: 10s`, `WriteTimeout: 10s`, `IdleTimeout: 60s`.

- **Low — SBOM/Sicherheitsdokumentation nicht sichtbar.**  
  Im sichtbaren Code ist keine CRA-Dokumentation (SBOM, Update-/Patch-Prozess, Sicherheitsannahmen) enthalten. Da offenbar nur die Standardbibliothek genutzt wird, ist das Lieferkettenrisiko gering, aber dokumentatorisch unvollständig.  
  **Abhilfe:** In `README.md`/`AGENTS.md` einen CRA-/Security-Abschnitt ergänzen: sichere Konfiguration, Support- und Update-Prozess, Meldestelle für Schwachstellen sowie ein generierbares SBOM.

## 3. EU AI Act

Nicht anwendbar. Der Service enthält kein KI-System, sondern ausschließlich deterministisches FNV-1a-Hashing (`internal/rollout/rollout.go`).

## 4. Pflichttexte & UI

Entfällt. Es handelt sich um eine reine REST-API ohne Endnutzer-UI. Rechtstexte, Datenschutzerklärung für Webseiten, Cookie-Banner oder Impressumspflichten entstehen erst bei einem zugehörigen Frontend bzw. beim Betreiber des Gesamtdienstes.

## 5. Barrierefreiheit

Entfällt. Keine öffentliche Web-UI vorhanden.

**Fazit**: Die datenschutzrelevanten Anforderungen im Code sind grundsätzlich gut umgesetzt. Vor einer Marktfreigabe müssen jedoch die CRA-Sicherheitsanforderungen erfüllt werden, insbesondere TLS/Authentifizierung und Server-Timeouts. Daher: CHANGES_REQUESTED.