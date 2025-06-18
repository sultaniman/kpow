# KPow 💥

A simple loopback form.

## Starting the server

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
  --key-kind=rsa \
  --advertise-key \
  --inbox-path=/data/inbox \
  --inbox-cron="*/5 * * * *" \
  --batch-size=10 \
  --log-level=INFO \
  --banner=/etc/kpow/banner.html \
  --hide-logo \
  --message-size=512
```

### Using a configuration file

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
| `KPOW_KEY_KIND`         | Key kind: `age`, `pgp`, or `rsa`      | string | ""            |
| `KPOW_ADVERTISE`        | Whether to advertise the key          | bool   | false         |
| `KPOW_KEY_PATH`         | Path to the key file                  | string | ""            |
| `KPOW_INBOX_PATH`       | Path to inbox                         | string | ""            |
| `KPOW_INBOX_CRON`       | Cron schedule for inbox processing    | string | `*/5 * * * *` |
| `KPOW_INBOX_BATCH_SIZE` | Inbox batch size                      | int    | 5             |

## Encryption

KPow supports Age, PGP, and RSA public keys for encrypting messages.
Provide the key kind with `--key-kind` (or `KPOW_KEY_KIND`) and the
path to your public key with `--pubkey` (or `KPOW_KEY_PATH`).
Available `--key-kind` options: `age`, `pgp`, or `rsa`.

### RSA Encryption Note

This system uses RSA encryption with OAEP padding and the SHA-256 hashing algorithm.
Please follow these guidelines when using RSA keys and configuring message parameters:

✅ **Key and Algorithm Requirements**

- **RSA Key Compatibility:** Must support OAEP padding (recommended size is 2048 bits or greater).
- **Hashing Algorithm:** Encryption uses SHA-256 — decryption must use the same.

**OAEP Padding Overhead**

- Padding size = 2 × HashSize + 2 bytes
- For SHA-256 (HashSize = 32 bytes), total padding is 66 bytes

**Maximum Message Sizes**

| RSA Key Size | Hash Algorithm | Hash Size | Padding Size | Max Message Size |
| ------------ | -------------- | --------- | ------------ | ---------------- |
| 2048 bits    | SHA-256        | 32 bytes  | 66 bytes     | 190 bytes        |
| 4096 bits    | SHA-256        | 32 bytes  | 66 bytes     | 446 bytes        |

⚠️ Messages must not exceed the maximum size, or encryption will fail with an error.

**Configuration Hint**

In your TOML config (`message_size`), set the value below the maximum message size based on your RSA key length. For example:

```toml
[server]
message_size = 180  # for 2048-bit RSA with SHA-256
```

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

Bun and Tailwind CSS are used to build the styles.
The style sources are in the `styles` folder.
Use `just styles` to customize and build the form styles, and
`just error-styles` for the error pages.
Both commands require `bun` and `bunx` to be installed.
