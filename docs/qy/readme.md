## [![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Kofe alyp beriñiz" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Türkçe](../tr/readme.md) | [Qyrgyz](readme.md) | [Français](../fr/readme.md) | [Українська](../uk/readme.md) | [Русский](../ru/readme.md)

KPow — öz aldynça ornotula turğan, qupujaluuluqqa bağyttalğan bajlanyş formasy, üçünçü tarap qyzmattaryna tayanbastan qoopsuz qabarlaşuunu qamsyzdajt.
Al PGP, Age cana RSA syjaqtuu zamanbap şifrlöö standarttaryn qoldonup, qabarlardy cetkirilgençe murun şifrlööt.
Qupujaluuluktu bağalağan iştep çyğuuçular, open source dolboorlor, bağymsyz sajttar, "whistleblower" platformalary cana içki quraldar üçün qoopsuz, esep berüüçü cana öz aldynça qabar iştetüü çeçimi.

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
> CLI parametrleri ar dajyma ortomoluq özgörmölördön cana konfiguratsija fajlynan üstün çyğat.

Konfiguratsijany cañyrlöö tartibi:

1. Alğaç konfiguratsija fajly (eger berilse) cüktölöt,
2. Anan ortomoluq özgörmölör (ENV) maanilerdi almaştyrat,
3. Aqyry CLI parametrleri baryn basyp ötöt.

```mermaid
flowchart TD
    A[Start] --> B{Konfig fajly barby?}
    B -- Ooba --> C[Konfiguratsija fajlyn cüktöö]
    B -- Coq --> D[Qalyptağy maanilerdi qoldonuu]
    C --> D
    D --> E[Ortomoluq özgörmölördü cüktöö]
    E --> F[CLI parametrlerin qoldonuu]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Konfiguratsija fajlyn tekşerüü

Serverdi işke qozğoodon murun `verify` buyruğun qoldonuñuz, konfiguratsijany cüktöp, qatalar tuuraluu cabarlardy bildiret:

```sh
$ kpow verify --config=path-to-config.toml
```

### Ortomoluq özgörmölör

| Özgörmönün aty          | Tasvirleme                                     | Tibi   | Alğaçqy maani |
| ----------------------- | ---------------------------------------------- | ------ | ------------- |
| `KPOW_TITLE`            | Serverdin atalyşy                              | string | ""            |
| `KPOW_PORT`             | Serverdin portu                                | int    | 8080          |
| `KPOW_HOST`             | Serverdin host daregi                          | string | localhost     |
| `KPOW_LOG_LEVEL`        | Logdun deñgeeli                                | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Server qabylday alğan maksimalduu qabar ölçömü | int    | 240           |
| `KPOW_HIDE_LOGO`        | Logo caşyrylsynby                              | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Banner fajly                                   | string | ""            |
| `KPOW_LIMITER_RPM`      | Rate limiter: minutuna suranuular              | int    | 0             |
| `KPOW_LIMITER_BURST`    | Rate limiter: burst ölçömü                     | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Rate limiter: muzdatqyç sekunddary             | int    | -1            |
| `KPOW_MAILER_FROM`      | Cönötüüçünün email daregi                      | string | ""            |
| `KPOW_MAILER_TO`        | Qabylda alğan email daregi                     | string | ""            |
| `KPOW_MAILER_DSN`       | SMTP DSN (bajlanyş saptyğy)                    | string | ""            |
| `KPOW_WEBHOOK_URL`      | Webhook URL                                    | string | ""            |
| `KPOW_MAX_RETRIES`      | Email ciberüü üçün qajra araketter             | int    | 2             |
| `KPOW_KEY_KIND`         | Açqyç türü: `age`, `pgp` ce `rsa`              | string | ""            |
| `KPOW_ADVERTISE`        | Açqyç caryjalansyn                             | bool   | false         |
| `KPOW_KEY_PATH`         | Açqyç fajlynyn colu                            | string | ""            |
| `KPOW_INBOX_PATH`       | Inbox papkasy                                  | string | ""            |
| `KPOW_INBOX_CRON`       | Inbox üçün cron-cadybal                        | string | `*/5 * * * *` |

> [!note]
> KPow qabarlardy cetkizüü üçün coq degende `KPOW_MAILER_DSN` ce `KPOW_WEBHOOK_URL` beriliişi şart.

## Şifrlöö

KPow Age, PGP cana RSA açyq açqyçtaryn qoldonup qabarlardy şifrlööt.
`--key-kind` (ce `KPOW_KEY_KIND`) menen açqyç türün,
`--pubkey` (ce `KPOW_KEY_PATH`) menen açyq açqyç fajlynyn colun körsötüñüz.
Mümkün varianttar: `age`, `pgp`, `rsa`.

### Açqyçtardy casoo

Tömöndögü köb qoldonulğan CLI quraldardy pajdalanyñyz:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

`age.pub` fajlyn `--pubkey` (ce `KPOW_KEY_PATH`) sypatynda qoldonuñuz.

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

ASCII formatyndağy `pgp.pub` fajly `--pubkey` üçün dajyr bolot.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` fajlyn `--pubkey` qa qoşuñuz.
Açyq açqyç PKIX PEM formatynda boluşu cana keminde 2048 bittik boluşu zaryl.

### Konfig fajlynyn mysaly

CLI flagtardyn ornuna açqyçty TOML konfiguratsijasyna belgileñiz:

```toml
[key]
kind = "age"           # ce "pgp" ce "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA şifrlöö tuuraluu eske saluu

Bul tuzum RSA şifrlöönü OAEP padding cana SHA-256 hash menen qoldonot.
RSA açqyçtardy ce `message_size` parametrin cajlaştyrğanda tömönkü nuqtalardy eske alyñyz:

✅ **Açqyç cana algoritm talaptary**

- **RSA açqyç ündöştügü:** OAEP paddingti qoldoo şart (keminde 2048-bit önerilet).
- **Hash algoritmi:** Şifrlöö SHA-256 menen atqaryladı — deşifrlöö da uşul hash menen boluşu tiiş.

**OAEP paddingdin kölömü**

- Padding = 2 × Hash kölömü + 2 bajt
- SHA-256 üçün (Hash kölömü 32 bajt) padding 66 bajtty tüzöt

**Maksimalduu qabar ölçömlör**

| RSA açqyç uzunduğu | Hash algoritmi | Hash kölömü | Padding kölömü | Maksimalduu qabar ölçömü |
| ------------------ | -------------- | ----------- | -------------- | ------------------------ |
| 2048 bit           | SHA-256        | 32 bajt     | 66 bajt        | 190 bajt                 |
| 4096 bit           | SHA-256        | 32 bajt     | 66 bajt        | 446 bajt                 |

⚠️ Çektöödön aşqan qabarlardy şifrlöögö dajar qylyş üçün qysqartylat.

**Konfiguratsija bojunça keñeş**

TOML konfiguratsijasynda (`message_size`) maaniñizdi RSA açqyç uzunduğuna laajyqtap belgileñiz. Mysaly:

```toml
[server]
message_size = 180  # 2048-bit RSA (SHA-256) üçün
```

## Mailer ağymy

```mermaid
flowchart TD
    A[Cañy qabar ciberildi] --> B{Daroo ciberüü işke aştyby?}
    B -- Ooba --> C[Qabar ciberildi]
    B -- Coq --> D[Inbox papqağa saqtaluu]
    D --> E[Scheduler (cron) işke çügöt]
    E --> F[Qabarlardy oquu]
    F --> G{Qajradan ciberüü araketi?}
    G -- Ooba --> H[Qabar ciberildi]
    G -- Coq --> E
```

## Webhook

`--webhook-url` (ce `KPOW_WEBHOOK_URL`) körsötülgön bolso, KPow şifrlengen forma maalymattaryn
körsötülgön endpointke JSON formatynda POST menen ciberet:

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

Webhook URL `localhost` emes bolso, macburduu türdö HTTPS boluşu kerek.
HTTP status kodu < 400 bolso — ijgilik dep eseptelet.

## Docker

KPow Dockerfile menen kelip, kontejerlerge oñoj ornotulat:

```sh
docker build -t kpow .
docker run -p 8080:8080 \
  -v /path/to/key.pub:/app/key.pub \
  -e KPOW_KEY_KIND=age \
  -e KPOW_KEY_PATH=/app/key.pub \
  -e KPOW_WEBHOOK_URL=https://hooks.example.com/notify \
  kpow
```

## Salamatyn tekşerüü

KPow `/health` endpoint qamtyjt, kontejer orkestratsijasy cana load balancerler üçün:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Öndürüü

### Formany özgörtüü

Bun cana Tailwind CSS stilderdi casap çyğaruu üçün qoldonulat.
Stil bulaqdary `styles` papqasynda.
`just styles` — formanyñ stillerin cañyltoo cana casap çyğaruu,
`just error-styles` — qata betterdin stillerin casoo üçün.
Bul komandalar işteşi üçün `bun` cana `bunx` ornatylyşy şart.

### Bannerdi caalaştyruu

Formany yñğajylaştyruu üçün `--banner=/path/to/banner.html` ce `KPOW_CUSTOM_BANNER=/path/to/banner.html`
menen cañy banner qoşoolo bolot.
Berilgen banner HTML sanitizatsijadan ötöt; tömönkü tegdardy qoldonuuğa bolot.

> [!note]
> Bannerdin içindegi elementterge `style` atributun qoşup stildep alasyñyz.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Litsenzija

KPow **Business Source License 1.1** astynda litsenzijalanğan.

Üçünçü tarapqa kommersijalyq hosted ce managed qyzmat körsetüü üçün qoşumça litsenzijasyz pajdalana albajsyz.

**2028-12-04** küni bul dolboor **Apache License 2.0** astyna ötöt.

- 📄 [`LICENSE`](../../LICENSE)
- 📄 [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Skrinşottor

## ![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)

## ![rate limited](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)

![csrf error](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
