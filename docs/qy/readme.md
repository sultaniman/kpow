[![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)
---
<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Kofe alyp beriñiz" height="36"></a>

# KPow 💥

KPow bul öz aldynça tejlengen, kupujaluuluk türdö bagyttalgan bajlanyş formasy,
üçünçü tarap qyzmattarğa tayansyz, koopsuz bajlanyşka mümkündük beret.
Bul Age, PGP cana RSA syjaktuu zamanbap şifrlöö standartyn qoldonup,
cönötülgön qatty şifrlep bekitet.

Bul kupujaluulukka önöktör, Open Source projektter, köz karandysyz sajttar üçün ideal.

## Serverdi candyruu

### CLI parametirleri menen

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

### Konfiguratsija fajlyn paydalanuu

> [!note]
> CLI parametirleri ar daima orto moynunda turat: alardy ortamende coldoo pyson.

Konfiguratsijanyn tartibi:

1. Konfigurasija fajlyn cüktöö;
2. Sistem cöjrösünön çüktöö (ENV);
3. Sonunda CLI parametirleri prioritet alat.

```mermaid
flowchart TD
    A[Start] --> B{Config Fajly Barby?}
    B -- ooba --> C[Konfig Fajlyn Jüktöö]
    B -- cok --> D[Konfig defaulttaryn qoldonuu]
    C --> D
    D --> E[Environment Ozgörtüülördü Jüktöö]
    E --> F[CLI Parametirlerin Qoldonuu]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Konfiguratsija fajlyn tekşerüü

Serverdi candyruu aldynda konfiguratsijany syndoo:

```sh
$ kpow verify --config=path-to-config.toml
```

### Aylanma Özgörtülör (Environment variables)

| Özgörmö Aty             | Deskripciya                     | Tipi   | Algaçky maanisi |
| ----------------------- | ------------------------------- | ------ | --------------- |
| `KPOW_TITLE`            | Server aty                      | string | ""              |
| `KPOW_PORT`             | Server portu                    | int    | 8080            |
| `KPOW_HOST`             | Server host daregi              | string | localhost       |
| `KPOW_LOG_LEVEL`        | Log daracasy                    | string | INFO            |
| `KPOW_MESSAGE_SIZE`     | Qabardyn ölçömü                 | int    | 240             |
| `KPOW_HIDE_LOGO`        | Logo casyrylsynby               | bool   | false           |
| `KPOW_CUSTOM_BANNER`    | Bannerdin fajly                 | string | ""              |
| `KPOW_LIMITER_RPM`      | Bir request/m sany              | int    | 0               |
| `KPOW_LIMITER_BURST`    | Burst requestterdin sany        | int    | -1              |
| `KPOW_LIMITER_COOLDOWN` | Rate limit kuldaun ubakyty      | int    | -1              |
| `KPOW_MAILER_FROM`      | Jönötüüçü emaily                | string | ""              |
| `KPOW_MAILER_TO`        | Kabyldooçu emaily               | string | ""              |
| `KPOW_MAILER_DSN`       | SMTP DSN                        | string | ""              |
| `KPOW_WEBHOOK_URL`      | Webhook URL                     | string | ""              |
| `KPOW_MAX_RETRIES`      | Retry sany                      | int    | 2               |
| `KPOW_KEY_KIND`         | Açkyç türü: `age`, `pgp`, `rsa` | string | ""              |
| `KPOW_ADVERTISE`        | Açkyç caryjalansynby            | bool   | false           |
| `KPOW_KEY_PATH`         | Açkyç fajlynyn coly             | string | ""              |
| `KPOW_INBOX_PATH`       | Inbox folderin coly             | string | ""              |
| `KPOW_INBOX_CRON`       | Inboxtu iştetüü cron tartibi    | string | `*/5 * * * *`   |

> [!note]
> KPowğo Messec cetkirüü üçün `KPOW_MAILER_DSN` ce `KPOW_WEBHOOK_URL` boluşu kerek.

## Şifrlöö

KPow Age, PGP, cana RSA açyk açkyç menen qatty şifrlöö üçün koldonulat.
`--key-kind` (ce `KPOW_KEY_KIND`) parametri menen açkyç türün, `--pubkey` (ce `KPOW_KEY_PATH`) menen açkyç fajlynyn colun körsötüñüz.
Mümkün varianttar: `age`, `pgp`, `rsa`.

### Açkyç casoo

CLI tirkeme arkyluu:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

`age.pub` fajlyn `--pubkey` arqyluu qoldonuñuz.

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

`--pubkey` üçün `pgp.pub` fajlyn beriniz.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` fajly `--pubkey` sifatynda qoldonulut. Açyk açkyç PKIX PEM formatynda boluşu kerek.

### Konfig misaly

CLI flagtardyn ornuna açqyçty TOML fajl menen körsötüñüz:

```toml
[key]
kind = "age"           # ce "pgp" ce "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA Şifrlöö belgesi

Bul sistemas RSA OAEP padding cana SHA-256 heşetüü menen iştöö.
Açkyçtyn uzunduğuna caraşa qabar da çektelet.
Misal üçün, 2048-bittik RSA menen `message_size = 180`.

## Maler agymy

```mermaid
flowchart TD
    A[Cañy qabar ciberüü] --> B{Daroo ciberüü işke aştyby?}
    B -- Ooba --> C[Qabar ciberildi]
    B -- Cok --> D[Inbox folderge saktoo]
    D --> E[Cron cügürüü]
    E --> F[Qabarlardy oquu]
    F --> G{Qajra ciberüügö macburbu?}
    G -- Ooba --> H[Qabar ciberildi]
    G -- Ooba --> E
```

## Webhook

`--webhook-url` (ce `KPOW_WEBHOOK_URL`) berseñiz, KPow şifrlengen maalymaty JSON formatynda korsötülgön endpointke POST arqyluu ciberilet:

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

Webhook URL HTTPS boluşu şart, `localhost` bolboso. HTTP code < 400 bolsa, ijgilik qatary eseptelet.

## Öndürüü

### Formany özğörtüü

Bun cana Tailwind CSS stil casoo üçün qoldonulat.

- Stil fajldary `styles` folderinde.
- `just styles` bujruğu stilderdi casoo üçün.
- `just error-styles` - qata betlerin stilleri.

Bul komandalar üçün `bun` cana `bunx` kerek.

### Bannerdi özğörtüü

`--banner=/path/to/banner.html` ce `KPOW_CUSTOM_BANNER=/path/to/banner.html` menen biriktirip, öz bannerdi qoşo alasyz.
Bannerdin HTMLi tazalanuu cana filtrden ötöt, tuura tagtardyn tizmesi tömönküdöj:

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Litsenziya

KPow **Business Source License 1.1** menen litsenziyalangan.
Siz programmany kommersijalyk hosttoo uchun üçünçü tarapka qyzmat körsötüü üçün lisenzijasyz pajdalana albajsyz.
**2028-12-04** ta projekt **Apache License 2.0** menen daroo litsenziyalanat.

- 📄 [`LICENSE`](./LICENSE)
- 📄 [`LICENSE-BUSL`](./LICENSE-BUSL)
- 📄 [`LICENSE-APACHE`](./LICENSE-APACHE)

## Skrinşottor

![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)
---
![rate limited](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)
---
![csrf error](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
