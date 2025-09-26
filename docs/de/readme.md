[![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)
---
<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Spendiere mir einen Kaffee" height="36"></a>

# KPow 💥

**Sprache:** [Deutsch](readme.md) · [English](../../readme.md) · [Qyrgyz](../qy/readme.md)

KPow ist ein selbst gehostetes, auf Privatsphäre ausgerichtetes Kontaktformular,
das sichere Kommunikation ohne Drittanbieter ermöglicht.
Es nutzt moderne Verschlüsselungsstandards wie Age, PGP und RSA,
um eingehende Nachrichten zu verschlüsseln und zu schützen.

Ideal für datenschutzbewusste Organisationen, Open-Source-Projekte und unabhängige Seiten.

## Server starten

### Über CLI-Parameter

```sh
$ kpow start \
  --config=/etc/kpow/config.toml \
  --port=8080 \
  --host=0.0.0.0 \
  --limiter-rpm=100 \
  --limiter-burst=20 \
  --limiter-cooldown=10 \
  --mailer-from=sender@example.com \
  --mailer-to=recipient@example.com \
  --mailer-dsn=smtp://user:password@smtp.example.com:587 \
  --max-retries=3 \
  --webhook-url=https://hooks.example.com/notify \
  --pubkey=/keys/key.pub \
  --key-kind=rsa \
  --advertise-key \
  --inbox-path=/data/inbox \
  --inbox-cron="*/5 * * * *" \
  --log-level=INFO \
  --banner=/etc/kpow/banner.html \
  --hide-logo \
  --message-size=512
```

### Nutzung einer Konfigurationsdatei

> [!note]
> CLI-Flags haben immer Vorrang; setze sie bewusst ein.

Die Reihenfolge der Konfiguration:

1. Konfigurationsdatei laden
2. Umgebungsvariablen anwenden (ENV)
3. CLI-Parameter überschreiben die vorherigen Werte

```mermaid
flowchart TD
    A[Start] --> B{Config-Datei vorhanden?}
    B -- ja --> C[Konfigurationsdatei laden]
    B -- nein --> D[Standardwerte verwenden]
    C --> D
    D --> E[Umgebungsvariablen laden]
    E --> F[CLI-Parameter anwenden]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Konfigurationsdatei prüfen

Überprüfe die Konfiguration vor dem Start des Servers:

```sh
$ kpow verify --config=path-to-config.toml
```

### Umgebungsvariablen

| Variable                | Beschreibung                      | Typ    | Standardwert    |
| ----------------------- | --------------------------------- | ------ | --------------- |
| `KPOW_TITLE`            | Anzeigename des Servers           | string | ""              |
| `KPOW_PORT`             | Server-Port                       | int    | 8080            |
| `KPOW_HOST`             | Host-Adresse                      | string | localhost       |
| `KPOW_LOG_LEVEL`        | Log-Level                         | string | INFO            |
| `KPOW_MESSAGE_SIZE`     | Maximale Nachrichtengröße         | int    | 240             |
| `KPOW_HIDE_LOGO`        | Logo ausblenden                   | bool   | false           |
| `KPOW_CUSTOM_BANNER`    | Pfad zur Banner-Datei             | string | ""              |
| `KPOW_LIMITER_RPM`      | Requests pro Minute               | int    | 0               |
| `KPOW_LIMITER_BURST`    | Burst-Anzahl                      | int    | -1              |
| `KPOW_LIMITER_COOLDOWN` | Cooldown für Rate-Limit           | int    | -1              |
| `KPOW_MAILER_FROM`      | Absenderadresse                   | string | ""              |
| `KPOW_MAILER_TO`        | Empfängeradresse                  | string | ""              |
| `KPOW_MAILER_DSN`       | SMTP-DSN                          | string | ""              |
| `KPOW_WEBHOOK_URL`      | Webhook-URL                       | string | ""              |
| `KPOW_MAX_RETRIES`      | Anzahl der Wiederholungsversuche  | int    | 2               |
| `KPOW_KEY_KIND`         | Schlüsseltyp: `age`, `pgp`, `rsa` | string | ""              |
| `KPOW_ADVERTISE`        | Schlüssel bekannt machen?         | bool   | false           |
| `KPOW_KEY_PATH`         | Pfad zur Schlüsseldatei           | string | ""              |
| `KPOW_INBOX_PATH`       | Pfad zum Inbox-Ordner             | string | ""              |
| `KPOW_INBOX_CRON`       | Cron-Zeitplan für die Inbox       | string | `*/5 * * * *`   |

> [!note]
> Für den Versand benötigt KPow entweder `KPOW_MAILER_DSN` oder `KPOW_WEBHOOK_URL`.

## Verschlüsselung

KPow nutzt öffentliche Schlüssel für Age, PGP und RSA,
um eingehende Nachrichten asymmetrisch zu verschlüsseln.
Gib den Typ mit `--key-kind` (bzw. `KPOW_KEY_KIND`) und den Pfad mit
`--pubkey` (bzw. `KPOW_KEY_PATH`) an.
Mögliche Werte: `age`, `pgp`, `rsa`.

### Schlüssel erzeugen

Per CLI:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

Verwende `age.pub` mit `--pubkey`.

#### PGP

```sh
gpg --quick-generate-key "Dein Name <du@example.com>"
gpg --armor --export du@example.com > pgp.pub
```

`pgp.pub` anschließend mit `--pubkey` angeben.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` bei `--pubkey` nutzen. Der Schlüssel muss im PKIX-PEM-Format vorliegen.

### Beispielkonfiguration

Statt CLI-Flags kannst du den Schlüssel in der TOML-Datei definieren:

```toml
[key]
kind = "age"           # oder "pgp" bzw. "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA-Verschlüsselung

KPow verwendet RSA OAEP mit SHA-256. Die maximale Nachrichtengröße hängt
dabei von der Schlüssellänge ab. Bei einem 2048-Bit-Schlüssel empfiehlt
sich `message_size = 180`.

## Mailer-Ablauf

```mermaid
flowchart TD
    A[Neue Nachricht] --> B{Direkt senden?}
    B -- Ja --> C[Nachricht gesendet]
    B -- Nein --> D[In Inbox speichern]
    D --> E[Cron-Job starten]
    E --> F[Nachrichten lesen]
    F --> G{Erneut senden?}
    G -- Ja --> H[Nachricht gesendet]
    G -- Ja --> E
```

## Webhook

Mit `--webhook-url` (oder `KPOW_WEBHOOK_URL`) sendet KPow die
verschlüsselte Nachricht als JSON per POST an das Ziel:

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

Die URL muss HTTPS verwenden, außer bei `localhost`.
Antworten mit Status < 400 gelten als erfolgreich.

## Entwicklung

### Formular anpassen

Zur Erstellung der Styles kommen Bun und Tailwind CSS zum Einsatz.

- Styles liegen im Ordner `styles`.
- `just styles` baut die regulären Styles.
- `just error-styles` erzeugt die Fehlerseiten-Styles.

Für diese Befehle werden `bun` und `bunx` benötigt.

### Banner anpassen

Mit `--banner=/path/to/banner.html` oder
`KPOW_CUSTOM_BANNER=/path/to/banner.html` kannst du dein eigenes Banner einbinden.
Der HTML-Inhalt wird gesäubert; folgende Tags sind erlaubt:

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Lizenz

KPow steht unter der **Business Source License 1.1**.
Das Hosting als Drittanbieterdienst ist ohne zusätzliche Lizenz nicht gestattet.
Am **04.12.2028** wechselt das Projekt automatisch zur **Apache License 2.0**.

- 📄 [`LICENSE`](./LICENSE)
- 📄 [`LICENSE-BUSL`](./LICENSE-BUSL)
- 📄 [`LICENSE-APACHE`](./LICENSE-APACHE)

## Screenshots

![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)
---
![rate limited](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)
---
![csrf error](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
