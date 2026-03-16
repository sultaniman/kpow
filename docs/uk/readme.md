## [![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Пригости мене кавою" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Türkçe](../tr/readme.md) | [Qyrgyz](../qy/readme.md) | [Français](../fr/readme.md) | [Українська](readme.md) | [Русский](../ru/readme.md)

KPow — це самостійно розгортаний, орієнтований на конфіденційність сервер контактних форм, створений для безпечного спілкування без залежності від сторонніх сервісів.
Він підтримує сучасні стандарти шифрування — PGP, Age та RSA — щоб забезпечити шифрування повідомлень перед доставкою.
Ідеально підходить для розробників, які цінують конфіденційність, проєктів з відкритим кодом, незалежних вебсайтів, платформ для інформаторів та внутрішніх інструментів, що потребують безпечної, аудитованої та автономної обробки повідомлень.

## Запуск сервера

### За допомогою аргументів CLI

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

### За допомогою файлу конфігурації

> [!note]
> Аргументи CLI завжди перевизначають змінні оточення та файли конфігурації.

Порядок вирішення конфігурації:

1. Конфігурація — спочатку завантажується з файлу конфігурації, якщо він наданий,
2. Змінні оточення — перевизначають значення з файлу конфігурації,
3. Аргументи CLI — перевизначають змінні оточення та значення файлу конфігурації

```mermaid
flowchart TD
    A[Запуск] --> B{Файл конфігурації наданий?}
    B -- Так --> C[Завантажити файл конфігурації]
    B -- Ні --> D[Використати значення за замовчуванням]
    C --> D
    D --> E[Завантажити змінні оточення]
    E --> F[Застосувати аргументи CLI]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Перевірка файлу конфігурації

Виконайте команду `verify`, щоб завантажити конфігурацію та повідомити про будь-які
проблеми валідації без запуску сервера:

```sh
$ kpow verify --config=path-to-config.toml
```

### Змінні оточення

| Назва змінної           | Опис                                      | Тип    | За замовчуванням |
| ----------------------- | ----------------------------------------- | ------ | ---------------- |
| `KPOW_TITLE`            | Заголовок сервера                         | string | ""               |
| `KPOW_PORT`             | Порт сервера                              | int    | 8080             |
| `KPOW_HOST`             | Адреса хоста сервера                      | string | localhost        |
| `KPOW_LOG_LEVEL`        | Рівень логування                          | string | INFO             |
| `KPOW_MESSAGE_SIZE`     | Максимальний розмір повідомлення сервера  | int    | 240              |
| `KPOW_HIDE_LOGO`        | Чи приховувати логотип                    | bool   | false            |
| `KPOW_CUSTOM_BANNER`    | Файл користувацького банера               | string | ""               |
| `KPOW_LIMITER_RPM`      | Обмежувач: запитів на хвилину             | int    | 0                |
| `KPOW_LIMITER_BURST`    | Обмежувач: розмір сплеску                 | int    | -1               |
| `KPOW_LIMITER_COOLDOWN` | Обмежувач: час охолодження в секундах     | int    | -1               |
| `KPOW_MAILER_FROM`      | Email відправника                         | string | ""               |
| `KPOW_MAILER_TO`        | Email отримувача                          | string | ""               |
| `KPOW_MAILER_DSN`       | DSN поштового сервісу (рядок підключення) | string | ""               |
| `KPOW_WEBHOOK_URL`      | URL webhook                               | string | ""               |
| `KPOW_MAX_RETRIES`      | Максимальна кількість спроб надсилання    | int    | 2                |
| `KPOW_KEY_KIND`         | Тип ключа: `age`, `pgp` або `rsa`         | string | ""               |
| `KPOW_ADVERTISE`        | Чи оголошувати ключ                       | bool   | false            |
| `KPOW_KEY_PATH`         | Шлях до файлу ключа                       | string | ""               |
| `KPOW_INBOX_PATH`       | Шлях до вхідних                           | string | ""               |
| `KPOW_INBOX_CRON`       | Розклад cron для обробки вхідних          | string | `*/5 * * * *`    |

> [!note]
> Необхідно вказати `KPOW_MAILER_DSN` або `KPOW_WEBHOOK_URL`, щоб KPow міг доставляти повідомлення.

## Шифрування

KPow підтримує публічні ключі Age, PGP та RSA для шифрування повідомлень.
Вкажіть тип ключа за допомогою `--key-kind` (або `KPOW_KEY_KIND`) та
шлях до вашого публічного ключа за допомогою `--pubkey` (або `KPOW_KEY_PATH`).
Доступні варіанти `--key-kind`: `age`, `pgp` або `rsa`.

### Генерація ключів

Використовуйте стандартні інструменти командного рядка для створення сумісних публічних ключів:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

Використовуйте `age.pub` як значення для `--pubkey` (або `KPOW_KEY_PATH`).

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

Передайте ASCII-армований файл `pgp.pub` до `--pubkey`.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

Вкажіть `rsa_public.pem` як `--pubkey`. Публічний ключ повинен бути у форматі PKIX
PEM-кодованого RSA ключа (2048 біт або більше).

### Приклад файлу конфігурації

Замість прапорців CLI вкажіть ключ у файлі конфігурації TOML:

```toml
[key]
kind = "age"           # or "pgp" or "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### Примітка щодо RSA шифрування

Ця система використовує RSA шифрування з OAEP доповненням та алгоритмом хешування SHA-256.
Будь ласка, дотримуйтесь цих рекомендацій при використанні RSA ключів та налаштуванні параметрів повідомлень:

✅ **Вимоги до ключа та алгоритму**

- **Сумісність RSA ключа:** Повинен підтримувати OAEP доповнення (рекомендований розмір — 2048 біт або більше).
- **Алгоритм хешування:** Шифрування використовує SHA-256 — дешифрування повинно використовувати той самий.

**Накладні витрати OAEP доповнення**

- Розмір доповнення = 2 × HashSize + 2 байти
- Для SHA-256 (HashSize = 32 байти) загальне доповнення становить 66 байтів

**Максимальні розміри повідомлень**

| Розмір RSA ключа | Алгоритм хешування | Розмір хешу | Розмір доповнення | Макс. розмір повідомлення |
| ---------------- | ------------------ | ----------- | ----------------- | ------------------------- |
| 2048 біт         | SHA-256            | 32 байти    | 66 байтів         | 190 байтів                |
| 4096 біт         | SHA-256            | 32 байти    | 66 байтів         | 446 байтів                |

⚠️ Повідомлення, що перевищують максимальний розмір для ключа, будуть обрізані перед шифруванням.

**Підказка щодо конфігурації**

У вашому файлі конфігурації TOML (`message_size`) встановіть значення нижче максимального розміру повідомлення залежно від довжини вашого RSA ключа. Наприклад:

```toml
[server]
message_size = 180  # for 2048-bit RSA with SHA-256
```

## Логіка поштового сервісу

```mermaid
flowchart TD
    A[Нове повідомлення надіслано] --> B{Спробувати надіслати одразу?}
    B -- Успіх --> C[Повідомлення надіслано]
    B -- Помилка --> D[Зберегти у папку вхідних]
    D --> E[Запуск планувальника]
    E --> F[Читання повідомлень]
    F --> G{Спробувати надіслати повторно}
    G -- Успіх --> H[Повідомлення надіслано]
    G -- Помилка --> E
```

## Webhook

Коли вказано `--webhook-url` (або `KPOW_WEBHOOK_URL`), KPow надсилатиме POST-запит із
зашифрованими даними форми на вказаний ендпоінт у форматі JSON:

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

URL webhook повинен використовувати HTTPS, якщо він не вказує на `localhost`. Будь-який HTTP-код
стану < 400 вважається успішним.

## Docker

KPow постачається з Dockerfile і може бути легко розгорнутий у контейнерах:

```sh
docker build -t kpow .
docker run -p 8080:8080 \
  -v /path/to/key.pub:/app/key.pub \
  -e KPOW_KEY_KIND=age \
  -e KPOW_KEY_PATH=/app/key.pub \
  -e KPOW_WEBHOOK_URL=https://hooks.example.com/notify \
  kpow
```

## Перевірка стану

KPow надає ендпоінт `/health` для оркестрації контейнерів та балансувальників навантаження:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Розробка

### Користувацька форма

Bun та Tailwind CSS використовуються для збірки стилів.
Джерела стилів знаходяться у папці `styles`.
Використовуйте `just styles` для налаштування та збірки стилів форми, та
`just error-styles` для сторінок помилок.
Обидві команди потребують встановлених `bun` та `bunx`.

### Користувацький банер

Є можливість налаштувати форму та додати користувацький банер за допомогою `--banner=/path/to/banner.html` або встановивши `KPOW_CUSTOM_BANNER=/path/to/banner.html`.
HTML у наданому банері буде очищено; нижче наведено список дозволених тегів.

**Дозволені теги**

> [!note]
> Ви можете використовувати атрибут `style` для стилізації вашого банера.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Ліцензія

KPow ліцензований під **Business Source License 1.1**.

Ви **не маєте права** використовувати це програмне забезпечення для надання комерційного хостингового або керованого сервісу третім сторонам без придбання окремої комерційної ліцензії.

**04.12.2028** цей проєкт буде переліцензований під **Apache License 2.0**.

- 📄 Див. [`LICENSE`](../../LICENSE)
- 📄 Див. [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 Див. [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Знімки екрана

## ![форма](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)

## ![обмеження швидкості](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)

![помилка CSRF](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
