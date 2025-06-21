# KPow 💥

KPow bul oz-ozun hosttoolgon, kuptuuluukka bagyttalgan baiylanyş formasy, üçünçü tarap xyzmattarğa tayansyz, koopsuz baylanyşka mümkindik beret.
Bul Age, PGP cana RSA syyaktuu zamanbap şifrlöö standartu tartyp koldoonu, cönötülgön xatty şifrlöp cazdyrat.
Bul kupuuluukka önöktör, açykk kaynak proyekttor, müüstakyl saytlar, tor aykyndoo platformalary, cana koopsuz, audit tuu, toluq önkörülgön xattar menen ishtöö kerek içki kuraldar üçün ideal.

## Serverdi baskiloo

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

### Konfiguraciya fajlyn paydalanuu

> [!note]
> CLI parametirleri ar daima orto moynunda turat: alardy ortamende coldoo pyson.

Konfiguraciyanyn tartibi:

1. Konfiguraciya fajlyn cükteü;
2. Aylandyñ čeyn geçkenn minezderd (ENV);
3. CLI parametirleri baryn basty alyştyrat.


```mermaid
flowchart TD
    A[Start] --> B{Config Fayly Barby?}
    B -- ooba --> C[Konfig Faylyn Jükteü]
    B -- cok --> D[Konfigdyn nökmö boyunça]
    C --> D
    D --> E[Environment Ozgörtülördü Jükteü]
    E --> F[CLI Parametirlerin Qoldonuu]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Konfiguraciya fajlyn tekşerüü

Serverdi baştağynda katali malymatlar bar-cogun teksherüü:

```sh
$ kpow verify --config=path-to-config.toml
```

### Aylanma Özgörtülör (Environment variables)

| Özgörtü Aty            | Deskripciya                       | Türü  | Default |
| --------------------- | --------------------------------- | ---- | ------- |
| `KPOW_TITLE`          | Server atkasy                     | string| ""      |
| `KPOW_PORT`           | Server portu                      | int   | 8080    |
| `KPOW_HOST`           | Server host addressi              | string| localhost|
| `KPOW_LOG_LEVEL`      | Log därecesi                     | string| INFO    |
| `KPOW_MESSAGE_SIZE`   | Maks xabardyn ölçömi             | int   | 240     |
| `KPOW_HIDE_LOGO`      | Logo casyrylsynby                | bool  | false   |
| `KPOW_CUSTOM_BANNER`  | Bannerdin fajly                  | string| ""      |
| `KPOW_LIMITER_RPM`    | Bir minuttaga süryö sany         | int   | 0       |
| `KPOW_LIMITER_BURST`  | Burst ölçömi                      | int   | -1      |
| `KPOW_LIMITER_COOLDOWN`| Söndürüü müdööti                | int   | -1      |
| `KPOW_MAILER_FROM`    | Joöntöçü email                   | string| ""      |
| `KPOW_MAILER_TO`      | Kabyldooçu email                 | string| ""      |
| `KPOW_MAILER_DSN`     | SMTP DSN                          | string| ""      |
| `KPOW_WEBHOOK_URL`    | Webhook URL                       | string| ""      |
| `KPOW_MAX_RETRIES`    | Qaytaruu sany                     | int   | 2       |
| `KPOW_KEY_KIND`       | Klyuç türü: `age`, `pgp`, `rsa`   | string| ""      |
| `KPOW_ADVERTISE`      | Klyuç caryyalansynby              | bool  | false   |
| `KPOW_KEY_PATH`       | Klyuç fajlynyn coly              | string| ""      |
| `KPOW_INBOX_PATH`     | Inbox folderin coly             | string| ""      |
| `KPOW_INBOX_CRON`     | Inboxti iştetüü cron cädvali     | string| `*/5 * * * *` |

## Şifrlöö

KPow Age, PGP, cana RSA publikalyk klyuhtar menen xatty şifrlööyü koldoyt.
`--key-kind` (ce `KPOW_KEY_KIND`) parametri menen klyuç türün, `--pubkey` (ce `KPOW_KEY_PATH`) menen klyuç fajlynyn colun körsötüñüz.
Mümkün varianttar: `age`, `pgp`, `rsa`.

### Klyuhtar casoo

Köb atalkan konzol kuraldary menen koldonuluuda:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

`age.pub` fajlyn `--pubkey` boluup qoldonulunuz.

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

`--pubkey` üçün `pgp.pub` fajlyn berriniz.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` fajly `--pubkey` sifatynda qoldonulut. Publikalyk klyuç PKIX PEM kod formatynda boluşu kerek.

### Konfig misaly

CLI flagtardyn ornuna açqyçty TOML fajl menen körsötüñüz:

```toml
[key]
kind = "age"           # ce "pgp" ce "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA Şifrlöö belgesi

Bul sistemas RSA OAEP padding cana SHA-256 xeshootsu menen ishtöö.
Klyuqtun uzunduguna karap maks xabar ölçömi tetkiklenet.
Misal üçün, 2048-bittik RSA menen message_size = 180.

## Maler logikasy

```mermaid
flowchart TD
    A[Caña habar ciberildi] --> B{Darhol ciberüüga araktylyshaby?}
    B -- Iygylyk --> C[Habar ciberildi]
    B -- Qata --> D[Inbox folderge saktoo]
    D --> E[Cron cügürüü]
    E --> F[Xabarlardy oqoo]
    F --> G{Qayra ciberüüga araktylyshaby?}
    G -- Iygylyk --> H[Habar ciberildi]
    G -- Qata --> E
```

## Webhook

`--webhook-url` (ce `KPOW_WEBHOOK_URL`) berseniz, KPow şifrlöngön maglymatty JSON formatynda koorsötülgön endpointke POST qylat:

```json
{
  "subject": "<form subject>",
  "content": "<encrypted message>",
  "hash": "<sha256-hash>"
}
```

Webhook URL HTTPS boluşu şart, `localhost` bolboso. HTTP code < 400 bolsa, iygylyktuu.

## Öndürüü

### Formdy özğörtüü

Bun cana Tailwind CSS stil casoo üçün paydalanyladi.
Stil fajldary `styles` folderinde.
`just styles` bujruğu stilderdi casoo üçün.
`just error-styles` - qata betlerin stilleri.
Bul komandalar `bun` cana `bunx` kuraldaryn talab kyladi.

### Bannerdi özğörtüü

`--banner=/path/to/banner.html` ce `KPOW_CUSTOM_BANNER=/path/to/banner.html` menen biriktirip, öz bannerdi qoşo alasyz.
Banner dinHTML sanitized bolot, ruqsat berilgen tagtardyn tizmesi tömönködöy:

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Litsenziya

KPow **Business Source License 1.1** menen litsenziyalangan.
Siz programmany kommersiyalyk hosttoolup içinçü tarapka xyzmat körsötöö üçün özünçö litsenziya satyp almasangiz, paydalana albaysyz.
**2028-12-04** ta projekt **Apache License 2.0** menen qayta litsenziyalanat.

- 📄 [`LICENSE`](./LICENSE)
- 📄 [`LICENSE-BUSL`](./LICENSE-BUSL)
- 📄 [`LICENSE-APACHE`](./LICENSE-APACHE)
