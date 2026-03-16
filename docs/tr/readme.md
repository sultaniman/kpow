## [![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Bana bir kahve ısmarla" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Türkçe](readme.md) | [Qyrgyz](../qy/readme.md) | [Français](../fr/readme.md) | [Українська](../uk/readme.md) | [Русский](../ru/readme.md)

KPow, üçüncü taraf hizmetlere bağımlı olmadan güvenli iletişim sağlamak için tasarlanmış, kendi sunucunuzda barındırılan, gizlilik odaklı bir iletişim formudur.
PGP, Age ve RSA gibi modern şifreleme standartlarını destekleyerek mesajların teslim edilmeden önce şifrelenmesini sağlar.
Gizliliğe önem veren geliştiriciler, açık kaynak projeler, bağımsız web siteleri, ihbar platformları ve güvenli, denetlenebilir, bağımsız mesaj işleme gerektiren dahili araçlar için idealdir.

## Sunucuyu Başlatma

### CLI argümanları ile

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

### Yapılandırma dosyası ile

> [!note]
> CLI argümanları her zaman ortam değişkenlerini ve yapılandırma dosyalarını geçersiz kılar.

Yapılandırma çözümleme sırası:

1. Yapılandırma - eğer belirtilmişse önce yapılandırma dosyasından yüklenir,
2. Ortam değişkenleri - yapılandırma dosyasındaki değerleri geçersiz kılar,
3. CLI argümanları - ortam değişkenlerini ve yapılandırma dosyası değerlerini geçersiz kılar

```mermaid
flowchart TD
    A[Başla] --> B{Yapılandırma Dosyası Belirtildi mi?}
    B -- Evet --> C[Yapılandırma Dosyasını Yükle]
    B -- Hayır --> D[Varsayılan Yapılandırmayı Kullan]
    C --> D
    D --> E[Ortam Değişkenlerini Yükle]
    E --> F[CLI Argümanlarını Uygula]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Yapılandırma dosyasını doğrulama

Sunucuyu başlatmadan yapılandırmayı yüklemek ve doğrulama sorunlarını raporlamak için `verify` komutunu çalıştırın:

```sh
$ kpow verify --config=path-to-config.toml
```

### Ortam değişkenleri

| Değişken Adı            | Açıklama                              | Tür    | Varsayılan    |
| ----------------------- | ------------------------------------- | ------ | ------------- |
| `KPOW_TITLE`            | Sunucu başlığı                        | string | ""            |
| `KPOW_PORT`             | Sunucu portu                          | int    | 8080          |
| `KPOW_HOST`             | Sunucu ana bilgisayar adresi          | string | localhost     |
| `KPOW_LOG_LEVEL`        | Günlük seviyesi                       | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Maksimum sunucu mesaj boyutu          | int    | 240           |
| `KPOW_HIDE_LOGO`        | Logonun gizlenip gizlenmeyeceği       | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Özel banner dosyası                   | string | ""            |
| `KPOW_LIMITER_RPM`      | Hız sınırlayıcı: dakika başına istek  | int    | 0             |
| `KPOW_LIMITER_BURST`    | Hız sınırlayıcı: anlık istek limiti   | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Hız sınırlayıcı: bekleme süresi (sn)  | int    | -1            |
| `KPOW_MAILER_FROM`      | Gönderici e-posta adresi              | string | ""            |
| `KPOW_MAILER_TO`        | Alıcı e-posta adresi                  | string | ""            |
| `KPOW_MAILER_DSN`       | Posta DSN (bağlantı dizesi)           | string | ""            |
| `KPOW_WEBHOOK_URL`      | Webhook URL                           | string | ""            |
| `KPOW_MAX_RETRIES`      | E-posta gönderimi için maks. deneme   | int    | 2             |
| `KPOW_KEY_KIND`         | Anahtar türü: `age`, `pgp` veya `rsa` | string | ""            |
| `KPOW_ADVERTISE`        | Anahtarın ilan edilip edilmeyeceği    | bool   | false         |
| `KPOW_KEY_PATH`         | Anahtar dosyasının yolu               | string | ""            |
| `KPOW_INBOX_PATH`       | Gelen kutusu yolu                     | string | ""            |
| `KPOW_INBOX_CRON`       | Gelen kutusu işleme cron zamanlaması  | string | `*/5 * * * *` |

> [!note]
> KPow'un mesajları iletebilmesi için `KPOW_MAILER_DSN` veya `KPOW_WEBHOOK_URL` değerlerinden en az biri sağlanmalıdır.

## Şifreleme

KPow, mesajları şifrelemek için Age, PGP ve RSA genel anahtarlarını destekler.
Anahtar türünü `--key-kind` (veya `KPOW_KEY_KIND`) ile, genel anahtarınızın
yolunu ise `--pubkey` (veya `KPOW_KEY_PATH`) ile belirtin.
Kullanılabilir `--key-kind` seçenekleri: `age`, `pgp` veya `rsa`.

### Anahtar Oluşturma

Uyumlu genel anahtarlar oluşturmak için yaygın komut satırı araçlarını kullanın:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

`--pubkey` (veya `KPOW_KEY_PATH`) değeri olarak `age.pub` dosyasını kullanın.

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

ASCII zırhlı `pgp.pub` dosyasını `--pubkey` parametresine verin.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

`rsa_public.pem` dosyasını `--pubkey` olarak belirtin. Genel anahtar, PKIX
PEM kodlu bir RSA anahtarı (2048 bit veya daha büyük) olmalıdır.

### Yapılandırma dosyası örneği

CLI bayrakları yerine anahtarı bir TOML yapılandırma dosyasında belirtebilirsiniz:

```toml
[key]
kind = "age"           # veya "pgp" veya "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### RSA Şifreleme Notu

Bu sistem, OAEP dolgusu ve SHA-256 karma algoritması ile RSA şifreleme kullanır.
RSA anahtarlarını kullanırken ve mesaj parametrelerini yapılandırırken lütfen aşağıdaki yönergeleri izleyin:

✅ **Anahtar ve Algoritma Gereksinimleri**

- **RSA Anahtar Uyumluluğu:** OAEP dolgusunu desteklemelidir (önerilen boyut 2048 bit veya daha büyük).
- **Karma Algoritması:** Şifreleme SHA-256 kullanır — şifre çözme de aynısını kullanmalıdır.

**OAEP Dolgu Ek Yükü**

- Dolgu boyutu = 2 × KarmaBoyutu + 2 bayt
- SHA-256 için (KarmaBoyutu = 32 bayt), toplam dolgu 66 bayttır

**Maksimum Mesaj Boyutları**

| RSA Anahtar Boyutu | Karma Algoritması | Karma Boyutu | Dolgu Boyutu | Maks. Mesaj Boyutu |
| ------------------ | ----------------- | ------------ | ------------ | ------------------ |
| 2048 bit           | SHA-256           | 32 bayt      | 66 bayt      | 190 bayt           |
| 4096 bit           | SHA-256           | 32 bayt      | 66 bayt      | 446 bayt           |

⚠️ Anahtar için maksimum boyutu aşan mesajlar şifrelemeden önce kırpılacaktır.

**Yapılandırma İpucu**

TOML yapılandırmanızda (`message_size`), RSA anahtar uzunluğunuza göre maksimum mesaj boyutunun altında bir değer belirleyin. Örneğin:

```toml
[server]
message_size = 180  # 2048-bit RSA ve SHA-256 için
```

## Posta gönderim mantığı

```mermaid
flowchart TD
    A[Yeni Mesaj Gönderildi] --> B{Hemen Göndermeyi Dene?}
    B -- Başarılı --> C[Mesaj Gönderildi]
    B -- Hata --> D[Gelen Kutusu Klasörüne Kaydet]
    D --> E[Zamanlayıcı Çalışması]
    E --> F[Mesajları Oku]
    F --> G{Yeniden Göndermeyi Dene}
    G -- Başarılı --> H[Mesaj Gönderildi]
    G -- Hata --> E
```

## Webhook

`--webhook-url` (veya `KPOW_WEBHOOK_URL`) belirtildiğinde, KPow şifrelenmiş form verilerini belirtilen uç noktaya JSON formatında POST isteği ile gönderir:

```json
{
    "subject": "<form konusu>",
    "content": "<şifrelenmiş mesaj>",
    "hash": "<sha256-özet>"
}
```

Webhook URL'si, `localhost`'a yönlendirilmedikçe HTTPS kullanmalıdır. 400'den küçük herhangi bir HTTP durum kodu başarılı kabul edilir.

## Docker

KPow, Dockerfile ile birlikte gelir ve kolayca konteynerlere dağıtılabilir:

```sh
docker build -t kpow .
docker run -p 8080:8080 \
  -v /path/to/key.pub:/app/key.pub \
  -e KPOW_KEY_KIND=age \
  -e KPOW_KEY_PATH=/app/key.pub \
  -e KPOW_WEBHOOK_URL=https://hooks.example.com/notify \
  kpow
```

## Sağlık Kontrolü

KPow, konteyner orkestrasyonu ve yük dengeleyiciler için `/health` uç noktası sunar:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Geliştirme

### Özel form

Stilleri derlemek için Bun ve Tailwind CSS kullanılır.
Stil kaynakları `styles` klasöründedir.
Form stillerini özelleştirmek ve derlemek için `just styles`,
hata sayfaları için ise `just error-styles` komutunu kullanın.
Her iki komut da `bun` ve `bunx` kurulumunu gerektirir.

### Özel banner

Formu özelleştirmek ve özel bir banner eklemek için `--banner=/path/to/banner.html` kullanabilir veya `KPOW_CUSTOM_BANNER=/path/to/banner.html` ortam değişkenini ayarlayabilirsiniz.
Sağlanan banner içindeki HTML temizlenecektir; aşağıda izin verilen etiketlerin listesini görebilirsiniz.

**İzin verilen etiketler**

> [!note]
> Banner'ınızı biçimlendirmek için `style` niteliğini kullanabilirsiniz.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Lisans

KPow, **Business Source License 1.1** altında lisanslanmıştır.

Bu yazılımı, ayrı bir ticari lisans satın almadan üçüncü taraflara ticari barındırmalı veya yönetilen bir hizmet sunmak için **kullanamazsınız**.

**04.12.2028** tarihinde bu proje **Apache License 2.0** altında yeniden lisanslanacaktır.

- 📄 Bkz. [`LICENSE`](../../LICENSE)
- 📄 Bkz. [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 Bkz. [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Ekran Görüntüleri

## ![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)

## ![hız sınırlandı](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)

![csrf hatası](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
