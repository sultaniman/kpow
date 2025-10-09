[![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)
---
<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Kofe alyp beriñiz" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Qyrgyz](readme.md)

KPow — öz aldınça orunatıla turgan, kupujaaluuluqqa bağyttalğan bağlanyş formasy, üçünçü tarap qyzmattaryna tayanbastan koopsuz kabarlaşuunu camdayt.
Al PGP, Age cana RSA sıyaktuu zamanbap şifrlöö standarttaryn qoldop, qabarlardı cetkirilgençe murun şifrlööt.
Kupujaaluuluktu bağalagan iştep çıguuçular, open source doolborlor, bağımsız sayttar, "whistleblower" platformalary cana içki quraldar ücün koopsuz, esep berüüçü cana öz aldınça qabar iştetüü çeçimi.

## Serverdi baştoo

### CLI parametrlerin qoldonuu

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

### Konfiguratsija fajlyn qoldonuu

> [!note]
> CLI parametrleri är daima ortamalyk özgörtülördön cana konfiguratsija fajlınan üstün çıgat.

Konfiguratsijany cañırlöö tartibi:

1. Alğaç konfiguratsija fajly (eger berilse) cüktölöt,
2. Anan ortamalyk özgörtülör (ENV) maantilerdi almastırat,
3. Akırı CLI parametrleri barın basıp ötöt.

```mermaid
flowchart TD
    A[Start] --> B{Config fajly barby?}
    B -- Ooba --> C[Konfiguratsija fajlyn cüktöö]
    B -- Joq --> D[Kalıptagı maantilerdi qoldonuu]
    C --> D
    D --> E[Ortamalyk özgörtülördü cüktöö]
    E --> F[CLI parametrlerin qoldonuu]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Konfiguratsija fajlyn tekşerüü

Serverdi iske qozğoodon muрун `verify` buyruğun qoldonuñuz, konfiguratsijany cüktöp, qatalar tuuraluu cabarlardy bildiret:

```sh
$ kpow verify --config=path-to-config.toml
```

### Ortamalyk özgörtülör

| Özgörtmönün aty         | Tasvirleme                           | Tibi   | Alğaçky maani |
| ----------------------- | ------------------------------------ | ------ | ------------- |
| `KPOW_TITLE`            | Serverdin atalyşy                    | string | ""            |
| `KPOW_PORT`             | Serverdin portu                      | int    | 8080          |
| `KPOW_HOST`             | Serverdin host dariegi               | string | localhost     |
| `KPOW_LOG_LEVEL`        | Logdun deñgeeli                      | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Server qabılday alğan maksimalduu qabar ölçömi | int    | 240           |
| `KPOW_HIDE_LOGO`        | Logo cäşirilsinbi                   | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Banner fajly                         | string | ""            |
| `KPOW_LIMITER_RPM`      | Rate limiter: minutuna suranuular    | int    | 0             |
| `KPOW_LIMITER_BURST`    | Rate limiter: burst ölçömi           | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Rate limiter: muzdatqıç sekunddarı   | int    | -1            |
| `KPOW_MAILER_FROM`      | Jönötüüçünün email dariegi           | string | ""            |
| `KPOW_MAILER_TO`        | Qabyldooçunun email dariegi          | string | ""            |
| `KPOW_MAILER_DSN`       | SMTP DSN (bağlanyş saptygy)          | string | ""            |
| `KPOW_WEBHOOK_URL`      | Webhook URL                          | string | ""            |
| `KPOW_MAX_RETRIES`      | Email ciberüü ücün qayra araaketter  | int    | 2             |
| `KPOW_KEY_KIND`         | Açkyç türü: `age`, `pgp` ce `rsa`    | string | ""            |
| `KPOW_ADVERTISE`        | Açkyç caryjalansynby                 | bool   | false         |
| `KPOW_KEY_PATH`         | Açkyç fajlynyn colu                  | string | ""            |
| `KPOW_INBOX_PATH`       | Inbox papkasy                        | string | ""            |
| `KPOW_INBOX_CRON`       | Inbox ücün cron-cadwal               | string | `*/5 * * * *` |

> [!note]
> KPow qabarlardy cetkizüü ücün coq degende `KPOW_MAILER_DSN` ce `KPOW_WEBHOOK_URL` berilişi şart.

## Şifrlöö

KPow Age, PGP cana RSA açyk açkyçtaryn qoldonup qabarlardy şifrlööt.
`--key-kind` (ce `KPOW_KEY_KIND`) menen açkyç türün,
`--pubkey` (ce `KPOW_KEY_PATH`) menen açyk açkyç fajlynyn colun körsötüüñüz.
Mümkün varianttar: `age`, `pgp`, `rsa`.

### Açkyçtardy casoo

Tömöndögü köb qoldonulgan CLI quraldardy paydalanyñız:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

`age.pub` fajlyn `--pubkey` (ce `KPOW_KEY_PATH`) sifatında qoldonuñuz.

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

ASCII formatındagı `pgp.pub` fajly `--pubkey` ücün dayar bolot.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` fajlyn `--pubkey` qa qoşyñuz.
Açyk açkyç PKIX PEM formatında boluşu cana keminde 2048 bittik boluşu zarıl.

### Konfig fajlynyñ mısaly

CLI flagtardyn ornuna açkyçtı TOML konfiguratsijasına belgileñiz:

```toml
[key]
kind = "age"           # ce "pgp" ce "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA şifrlöö tuuraluu eske saluu

Bul tuzum RSA şifrlööün OAEP padding cana SHA-256 hash menen qoldonot.
RSA açkyçtardy ce `message_size` parametrin caylaştırğanda tömönkü nuktalardy eske alıñız:

✅ **Açkyç cana algoritm talaptary**

- **RSA açkyç ündöştügü:** OAEP paddingti qoldoo şart (keminde 2048-bit önerilet).
- **Hash algoritmi:** Şifrlöö SHA-256 menen atqarıladı — deşifrlöö da uşul hash menen boluşu tiiş.

**OAEP paddingdin kölömü**

- Padding = 2 × Hash kölömü + 2 bayt
- SHA-256 ücün (Hash kölömü 32 bayt) padding 66 bayttı tüzöt

**Maksimalduu qabar ölçömlör**

| RSA açkyç uzundugu | Hash algoritmi | Hash kölömü | Padding kölömü | Maksimalduu qabar ölçömi |
| ------------------ | -------------- | ----------- | -------------- | ------------------------ |
| 2048 bit           | SHA-256        | 32 bayt     | 66 bayt        | 190 bayt                 |
| 4096 bit           | SHA-256        | 32 bayt     | 66 bayt        | 446 bayt                 |

⚠️ Çektöödön aşqan qabarlardı şifrlöögö dayar kılış ücün qısqartylat.

**Konfiguratsija boyunça keñeş**

TOML konfiguratsijasynda (`message_size`) maaniñizdi RSA açkyç uzunduguna laayıqtap belgileñiz. Mısalı:

```toml
[server]
message_size = 180  # 2048-bit RSA (SHA-256) ücün
```

## Mailer agymy

```mermaid
flowchart TD
    A[Cañı qabar ciberildi] --> B{Daroo ciberüü işke aştyby?}
    B -- Ooba --> C[Qabar ciberildi]
    B -- Joq --> D[Inbox papkaga saqtaluu]
    D --> E[Scheduler (cron) iske cügöt]
    E --> F[Qabarlardı oqoo]
    F --> G{Qayradan ciberüü araketi?}
    G -- Ooba --> H[Qabar ciberildi]
    G -- Joq --> E
```

## Webhook

`--webhook-url` (ce `KPOW_WEBHOOK_URL`) körsötülgön bolso, KPow şifrlengen forma maalymattaryn
körsötülgön endpointke JSON formatında POST menen ciberet:

```json
{
  "subject": "<form subject>",
  "content": "<encrypted message>",
  "hash": "<sha256-hash>"
}
```

Webhook URL `localhost` emes bolso, macburduu türdö HTTPS boluşu kerek.
HTTP status kody < 400 bolso — ijgilik dep eseptelet.

## Öndürüü

### Formany cañıltoo

Bun cana Tailwind CSS stilderdi casap çıgaruu ücün qoldonulat.
Stil bu laqtary `styles` papkasynda.
`just styles` — formanyñ stillerin cañıltoo cana casap çıkaruu,
`just error-styles` — qata betterdin stillerin casoo ücün.
Bul komandalar işteşi ücün `bun` cana `bunx` ornatylyşy şart.

### Bannerdi caalaştıruu

Formanı ıñgayılaştıruu ücün `--banner=/path/to/banner.html` ce `KPOW_CUSTOM_BANNER=/path/to/banner.html`
menen cañı banner qoşoolo bolot.
Berilgen banner HTML sanitizatsijadan ötöt; tömönkö tagtardy qoldonuuga bolot.

> [!note]
> Bannerdin içindegi elementterge `style` atributun qoşup stillöö alasyñız.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Litsenzija

KPow **Business Source License 1.1** astında litsenzijalanğan.

Üçünçü tarapka kommersijalyk hosted ce managed qyzmat korsetüü ücün qoşumça litsenzijasız paydalana albaysız.

**2028-12-04** küni bul doolbor **Apache License 2.0** astına ötöt.

- 📄 [`LICENSE`](../../LICENSE)
- 📄 [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Skrinşottor

![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)
---
![rate limited](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)
---
![csrf error](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
