# Developerler üçün maalymat

Bul dokument KPow projektke qoşuluu üçün cardam beret.

- **cmd/** – Cobra menen kurulgan CLI. `start` bujruğu uşul cerde.
- **config/** – Konfiguraciya strukturalary, kömekçiler. `GetConfig` fajldardy, sistem çöjrödön, CLI flagtardan biriktiret.
- **server/** – Negizgi qoşumça kody. HTTP serverdi, formany, şifrlöö, mailerler cana cron qyzmaty bar.
- **styles/** – Tailwind CSS stilderi. `just styles` kompiliasijasyn atkarat.
- **art/** – Dokumentasijada ce web interfejste pajdalanuuçu süröttör.
1. **Go'nu cüktöö** – Projekt Go modul sistemasyn pajdalanat. Go 1.21+ cüktöö kerek.
2. **Bun** – `just styles` üçün kerek.
3. **Serverdi candandyruu**

```sh
./kpow start
   CLI flagtary çöjrödön maanilerin cana konfiguraciya fajlyn bastary.
## Konfigurasija
Yrastoo parametrler TOML fajl menen, ortom özgörtüülör menen ce CLI flagtary menen berilet. `config/config.go` fajly bar parametrlerden habar beret. `config.toml` cana `example.env` dajarlap qoyulgän.
- **Server** – Port, host, log cana requestterdi çektöö.
- **Mailerler** – SMTP ce webhook arkyluu cönötüü. Ijgiliktüü cetkirilbese inbox folderine saktoo.
- **Şifrlöö** – `age`, `pgp`, `rsa` açyk açkyç kol'donulat.
- **Schedulér** – Cron job inboxdon cönötüügö araket qylat.

Konfigurasija fajlda açkyçty körsötüü misaly:

### Konfigurasija agymy

    A[Başy] --> B{Config Barby?}
    B -- Ooba --> C[Config Fajlyn Okuu]
    C --> E[Sistem Ozgörtülörü]
### Konfigurasijany tekşerüü
```ş
- **Cron qyzmaty** `server/cron/` içinde. Inbokstu ciberüü.
- **Şifrlöö logikasy** `server/enc/` içinde.
```ş

```ş
```ş
### Mailer agymy
```ş
- **Middleware** `server/server.go` – CSRF, rate limit, maalymatty çektöö.
    B -- Qata --> D[Inboxko Saluu]
    D --> E[Cron candatuu]
    E --> F[Qabarlardy oquu]
    F --> G{Qajra ciberüü}
```go
3. Caña feature üçün test qoşuu kerek.
    A[Bastaş] --> B{Config Barby?}
    B -- Ooba --> C[Config Faylyn Okuu]
    B -- Jok --> D[Default Maanilär]
    C --> E[Ortam Ozgörtülör]
    D --> E
    E --> F[CLI Flagtar]
```

### Konfiguraciyany tekşerüü

```sh
./kpow verify --config=config.toml
```

## Önök Tips

- **Şablondor** `server/templates/` da, forma cana qata better.
- **Middleware** `server/server.go` – CSRF, rate limit, body çeklö.
- **Cron jobtor** `server/cron/` içinde. Inbox-čisti ciberüü.
- **Şifrlöö kyroyçtary** `server/enc/` içinde.

### Klyuç casoo

Age:

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

PGP:

```sh
gpg --quick-generate-key "Sizdin Atynyz <siz@example.com>"
gpg --armor --export siz@example.com > pgp.pub
```

RSA:

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` PKIX PEM formatynda boluşu kerek.

### Mailer qaytuu agymy

```mermaid
flowchart TD
    A[Caña habar] --> B{Darhol cönötüü?}
    B -- Ooba --> C[Ciberildi]
    B -- Qata --> D[Inboxka Saluu]
    D --> E[Cron cügürüü]
    E --> F[Xabarlardy oqoo]
    F --> G{Qayra ciberüü}
    G -- Ooba --> H[Ciberildi]
    G -- Qata --> E
```

## Testterdi Jügürtüü

```sh
go test ./...
```

(Testter qoldonmolordun cükteelüüsü üçün internet kerek bolup qaluşu mümkin.)

## Salym Qosuu

1. Repo-to fork casap, feature branch açynyz.
2. `gofmt` menen standart format saktañyz.
3. Caña funzionalnost üçün testter qosuu kerek.
4. PR cönötüü.

Formanyn, şifrlöö cana retry logikany tuluk bilüü üçün `readme.md` cana `server` paketindegi komentariyalarni qarañyz.

## Çygaruu

1. `just test` atqaryñyz.
2. `just build` ce GoReleaser menen binardyq faildar.
3. Litsenziyalar tekşerüü.
4. Siriätti tekşerip, syrly malymat joqtuğuna ynynuu.
5. Git tag qoşu.

Proyekt azyrça Business Source License 1.1 astynda, README-da belgilengendey 2028-12-04 de Apache License 2.0 go ote tygonot.
