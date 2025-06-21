[![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

# KPow 💥


[English](readme.md) | [Qyrgyz](docs/qy/readme.md)

KPow is a self-hosted, privacy-focused contact form designed for secure communication without relying on third-party services.
It supports modern encryption standards — PGP, Age, and RSA — to ensure that messages are encrypted before delivery.
Ideal for privacy-conscious developers, open source projects, independent websites, whistleblower platforms,
and internal tools that require secure, auditable, and self-contained message handling.



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
  --log-level=INFO \
  --banner=/etc/kpow/banner.html \
  --hide-logo \
  --message-size=512
```

### Using a configuration file

> [!note]
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
    C --> D
    D --> E[Load Environment Variables]
    E --> F[Apply CLI Arguments]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Verifying a configuration file

Run the `verify` command to load a configuration and report any
validation problems without starting the server:

```sh
$ kpow verify --config=path-to-config.toml
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
| `KPOW_CUSTOM_BANNER`    | Custom banner file                    | string | ""            |
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

## Encryption

KPow supports Age, PGP, and RSA public keys for encrypting messages.
Provide the key kind with `--key-kind` (or `KPOW_KEY_KIND`) and the
path to your public key with `--pubkey` (or `KPOW_KEY_PATH`).
Available `--key-kind` options: `age`, `pgp`, or `rsa`.

### Generating Keys

Use common command‑line tools to create compatible public keys:

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

Use `age.pub` as the value for `--pubkey` (or `KPOW_KEY_PATH`).

#### PGP

```sh
gpg --quick-generate-key "Your Name <you@example.com>"
gpg --armor --export you@example.com > pgp.pub
```

Pass the ASCII-armored `pgp.pub` file to `--pubkey`.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

Provide `rsa_public.pem` as `--pubkey`. The public key must be a PKIX
PEM‑encoded RSA key (2048 bits or greater).

### Config file example

Instead of CLI flags, specify the key in a TOML config file:

```toml
[key]
kind = "age"           # or "pgp" or "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

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

⚠️ Messages exceeding the maximum size for the key will be trimmed before encryption.

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
    E --> F[Read Messages]
    F --> G{Try to Re-Send}
    G -- Success --> H[Message Sent]
    G -- Error --> E
```

## Webhook

When `--webhook-url` (or `KPOW_WEBHOOK_URL`) is provided, KPow will POST the
encrypted form data to the specified endpoint in JSON format:

```json
{
  "subject": "<form subject>",
  "content": "<encrypted message>",
  "hash": "<sha256-hash>"
}
```

The webhook URL must use HTTPS unless it points to `localhost`. Any HTTP status
code < 400 is considered a success.

## Development

### Custom form

Bun and Tailwind CSS are used to build the styles.
The style sources are in the `styles` folder.
Use `just styles` to customize and build the form styles, and
`just error-styles` for the error pages.
Both commands require `bun` and `bunx` to be installed.

### Custom banner

It is possible to customize the form and add a custom banner using `--banner=/path/to/banner.html` or by setting `KPOW_CUSTOM_BANNER=/path/to/banner.html`.
HTML in the provided banner will be sanitized, below you can see the list of allowed tags.

**Allowed tags**

> [!note]
> You can use `style` attribute to style your banner.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## License

KPow is licensed under the **Business Source License 1.1**.

You **may not use** this software to offer a commercial hosted or managed service to third parties without purchasing a separate commercial license.

On **2028-12-04**, this project will be re-licensed under the **Apache License 2.0**.

- 📄 See [`LICENSE`](./LICENSE)
- 📄 See [`LICENSE-BUSL`](./LICENSE-BUSL)
- 📄 See [`LICENSE-APACHE`](./LICENSE-APACHE)
