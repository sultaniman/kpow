## [![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Угости меня кофе" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Türkçe](../tr/readme.md) | [Qyrgyz](../qy/readme.md) | [Français](../fr/readme.md) | [Українська](../uk/readme.md) | [Русский](readme.md)

KPow — это самостоятельно размещаемая контактная форма с упором на конфиденциальность, предназначенная для безопасного общения без использования сторонних сервисов.
Поддерживаются современные стандарты шифрования — PGP, Age и RSA — для обеспечения шифрования сообщений перед доставкой.
Идеально подходит для разработчиков, заботящихся о конфиденциальности, проектов с открытым исходным кодом, независимых веб-сайтов, платформ для информаторов и внутренних инструментов, требующих безопасной, аудируемой и автономной обработки сообщений.

## Запуск сервера

### С помощью аргументов CLI

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

### С помощью конфигурационного файла

> [!note]
> Аргументы CLI всегда переопределяют переменные окружения и конфигурационные файлы.

Порядок разрешения конфигурации:

1. Конфигурация — сначала загружается из конфигурационного файла, если он указан,
2. Переменные окружения — переопределяют значения из конфигурационного файла,
3. Аргументы CLI — переопределяют переменные окружения и значения из конфигурационного файла

```mermaid
flowchart TD
    A[Старт] --> B{Указан конфигурационный файл?}
    B -- Да --> C[Загрузить конфигурационный файл]
    B -- Нет --> D[Использовать значения по умолчанию]
    C --> D
    D --> E[Загрузить переменные окружения]
    E --> F[Применить аргументы CLI]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Проверка конфигурационного файла

Выполните команду `verify`, чтобы загрузить конфигурацию и сообщить о любых
проблемах валидации без запуска сервера:

```sh
$ kpow verify --config=path-to-config.toml
```

### Переменные окружения

| Имя переменной          | Описание                                 | Тип    | По умолчанию  |
| ----------------------- | ---------------------------------------- | ------ | ------------- |
| `KPOW_TITLE`            | Заголовок сервера                        | string | ""            |
| `KPOW_PORT`             | Порт сервера                             | int    | 8080          |
| `KPOW_HOST`             | Адрес хоста сервера                      | string | localhost     |
| `KPOW_LOG_LEVEL`        | Уровень логирования                      | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Максимальный размер сообщения            | int    | 240           |
| `KPOW_HIDE_LOGO`        | Скрывать ли логотип                      | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Файл пользовательского баннера           | string | ""            |
| `KPOW_LIMITER_RPM`      | Ограничитель: запросов в минуту          | int    | 0             |
| `KPOW_LIMITER_BURST`    | Ограничитель: размер всплеска            | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Ограничитель: время ожидания в секундах  | int    | -1            |
| `KPOW_MAILER_FROM`      | Email отправителя                        | string | ""            |
| `KPOW_MAILER_TO`        | Email получателя                         | string | ""            |
| `KPOW_MAILER_DSN`       | DSN почтовика (строка подключения)       | string | ""            |
| `KPOW_WEBHOOK_URL`      | URL вебхука                              | string | ""            |
| `KPOW_MAX_RETRIES`      | Максимальное количество попыток отправки | int    | 2             |
| `KPOW_KEY_KIND`         | Тип ключа: `age`, `pgp` или `rsa`        | string | ""            |
| `KPOW_ADVERTISE`        | Публиковать ли ключ                      | bool   | false         |
| `KPOW_KEY_PATH`         | Путь к файлу ключа                       | string | ""            |
| `KPOW_INBOX_PATH`       | Путь к папке входящих                    | string | ""            |
| `KPOW_INBOX_CRON`       | Расписание cron для обработки входящих   | string | `*/5 * * * *` |

> [!note]
> Необходимо указать `KPOW_MAILER_DSN` или `KPOW_WEBHOOK_URL`, чтобы KPow мог доставлять сообщения.

## Шифрование

KPow поддерживает публичные ключи Age, PGP и RSA для шифрования сообщений.
Укажите тип ключа с помощью `--key-kind` (или `KPOW_KEY_KIND`) и
путь к публичному ключу с помощью `--pubkey` (или `KPOW_KEY_PATH`).
Доступные значения `--key-kind`: `age`, `pgp` или `rsa`.

### Генерация ключей

Используйте стандартные инструменты командной строки для создания совместимых публичных ключей:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

Используйте `age.pub` в качестве значения для `--pubkey` (или `KPOW_KEY_PATH`).

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

Передайте ASCII-armored файл `pgp.pub` в `--pubkey`.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

Укажите `rsa_public.pem` в качестве `--pubkey`. Публичный ключ должен быть в формате PKIX
PEM-закодированного RSA-ключа (2048 бит или более).

### Пример конфигурационного файла

Вместо флагов CLI укажите ключ в конфигурационном файле TOML:

```toml
[key]
kind = "age"           # or "pgp" or "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### Примечание о шифровании RSA

Система использует RSA-шифрование с дополнением OAEP и алгоритмом хеширования SHA-256.
Пожалуйста, следуйте этим рекомендациям при использовании RSA-ключей и настройке параметров сообщений:

✅ **Требования к ключам и алгоритмам**

- **Совместимость RSA-ключа:** Должен поддерживать дополнение OAEP (рекомендуемый размер — 2048 бит или более).
- **Алгоритм хеширования:** Шифрование использует SHA-256 — дешифрование должно использовать тот же алгоритм.

**Накладные расходы дополнения OAEP**

- Размер дополнения = 2 × HashSize + 2 байта
- Для SHA-256 (HashSize = 32 байта) общий размер дополнения составляет 66 байт

**Максимальные размеры сообщений**

| Размер RSA-ключа | Алгоритм хеширования | Размер хеша | Размер дополнения | Макс. размер сообщения |
| ---------------- | -------------------- | ----------- | ----------------- | ---------------------- |
| 2048 бит         | SHA-256              | 32 байта    | 66 байт           | 190 байт               |
| 4096 бит         | SHA-256              | 32 байта    | 66 байт           | 446 байт               |

⚠️ Сообщения, превышающие максимальный размер для ключа, будут обрезаны перед шифрованием.

**Подсказка по настройке**

В конфигурационном файле TOML (`message_size`) установите значение ниже максимального размера сообщения в зависимости от длины вашего RSA-ключа. Например:

```toml
[server]
message_size = 180  # for 2048-bit RSA with SHA-256
```

## Логика отправки

```mermaid
flowchart TD
    A[Новое сообщение отправлено] --> B{Попытка отправить сразу?}
    B -- Успех --> C[Сообщение отправлено]
    B -- Ошибка --> D[Сохранить в папку входящих]
    D --> E[Запуск планировщика]
    E --> F[Чтение сообщений]
    F --> G{Попытка повторной отправки}
    G -- Успех --> H[Сообщение отправлено]
    G -- Ошибка --> E
```

## Webhook

Если указан `--webhook-url` (или `KPOW_WEBHOOK_URL`), KPow отправит
зашифрованные данные формы на указанный эндпоинт в формате JSON методом POST:

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

URL вебхука должен использовать HTTPS, за исключением `localhost`. Любой HTTP-код
статуса < 400 считается успешным.

## Docker

KPow поставляется с Dockerfile и легко разворачивается в контейнерах:

```sh
docker build -t kpow .
docker run -p 8080:8080 \
  -v /path/to/key.pub:/app/key.pub \
  -e KPOW_KEY_KIND=age \
  -e KPOW_KEY_PATH=/app/key.pub \
  -e KPOW_WEBHOOK_URL=https://hooks.example.com/notify \
  kpow
```

## Проверка состояния

KPow предоставляет эндпоинт `/health` для оркестрации контейнеров и балансировщиков нагрузки:

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Разработка

### Пользовательская форма

Bun и Tailwind CSS используются для сборки стилей.
Исходники стилей находятся в папке `styles`.
Используйте `just styles` для настройки и сборки стилей формы, а
`just error-styles` для страниц ошибок.
Обе команды требуют установленных `bun` и `bunx`.

### Пользовательский баннер

Можно настроить форму и добавить пользовательский баннер с помощью `--banner=/path/to/banner.html` или установив `KPOW_CUSTOM_BANNER=/path/to/banner.html`.
HTML в предоставленном баннере будет санитизирован, ниже приведён список разрешённых тегов.

**Разрешённые теги**

> [!note]
> Вы можете использовать атрибут `style` для стилизации вашего баннера.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Лицензия

KPow лицензирован под **Business Source License 1.1**.

Вы **не можете использовать** это программное обеспечение для предоставления коммерческого хостингового или управляемого сервиса третьим лицам без приобретения отдельной коммерческой лицензии.

**04.12.2028** этот проект будет перелицензирован под **Apache License 2.0**.

- 📄 См. [`LICENSE`](../../LICENSE)
- 📄 См. [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 См. [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Скриншоты

## ![форма](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)

## ![ограничение частоты](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)

![ошибка CSRF](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
