# KPow 💥

Simple loopback form.

## Starting server

### Using CLI arguments

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
  --key-kind=RSA \
  --advertise-key \
  --inbox-path=/data/inbox \
  --inbox-cron="*/5 * * * *" \
  --batch-size=10 \
  --log-level=INFO \
  --banner=/etc/kpow/banner.html \
  --hide-logo \
  --message-size=512
```

### Using configuration file

> [!NOTE]
> CLI arguments always override environment variables and configuration files.

Configuration resolution order:

1. Configuration - first load from config file if provided,
2. Environment variables - override values from configuration file,
3. CLI arguments - override environment variables and config file values

```mermaid
flowchart TD
    A[Start] --> B{Config File Provided?}
    B -- Yes --> C[Load Config File]
    B -- No --> D[Use Config Defaults]

    C --> E[Load Environment Variables]
    C --> D
    D --> E

    E --> F[Apply CLI Arguments]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Environment variables

| Variable Name           | Description                           | Type   | Default       |
| ----------------------- | ------------------------------------- | ------ | ------------- |
| `KPOW_TITLE`            | Server title                          | string | ""            |
| `KPOW_PORT`             | Server port                           | int    | 8080          |
| `KPOW_HOST`             | Server host address                   | string | localhost     |
| `KPOW_LOG_LEVEL`        | Logging level                         | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Maximum server message size           | int    | 240           |
| `KPOW_HIDE_LOGO`        | Whether to hide the logo              | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Custom banner text                    | string | ""            |
| `KPOW_LIMITER_RPM`      | Rate limiter: requests per minute     | int    | 0             |
| `KPOW_LIMITER_BURST`    | Rate limiter: burst size              | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Rate limiter: cooldown in seconds     | int    | -1            |
| `KPOW_MAILER_FROM`      | Mailer sender email                   | string | ""            |
| `KPOW_MAILER_TO`        | Mailer recipient email                | string | ""            |
| `KPOW_MAILER_DSN`       | Mailer DSN (connection string)        | string | ""            |
| `KPOW_WEBHOOK_URL`      | Webhook URL                           | string | ""            |
| `KPOW_MAX_RETRIES`      | Max retry attempts for sending emails | int    | 2             |
| `KPOW_KEY_KIND`         | Key kind (e.g., type of key)          | string | ""            |
| `KPOW_ADVERTISE`        | Whether to advertise the key          | bool   | false         |
| `KPOW_KEY_PATH`         | Path to the key file                  | string | ""            |
| `KPOW_INBOX_PATH`       | Path to inbox                         | string | ""            |
| `KPOW_INBOX_CRON`       | Cron schedule for inbox processing    | string | `*/5 * * * *` |
| `KPOW_INBOX_BATCH_SIZE` | Inbox batch size                      | int    | 5             |

## Mailer logic

```mermaid
flowchart TD
    A[New Message Submitted] --> B{Try to Send Immediately?}
    B -- Success --> C[Message Sent]
    B -- Error --> D[Save to Inbox Folder]
    D --> E[Scheduler Run]
    E --> F[Read Messages in Batches]
    F --> G{Try to Re-Send}
    G -- Success --> H[Message Sent]
    G -- Error --> E
```

## Development

### Custom form

Bun and Tailwindcss is used to build the styles.

TODO
