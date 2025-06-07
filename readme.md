# KPow 💥

Simple loopback form.

## Starting server

### Using cli arguments

```sh
$ kpow start --pubkey=public_key_path \
             --port=8080 \
             --host=localhost \
             --mailer=smtp://user:password@smtp.example.com:587 \
             --mailer-from=sender@example.com \
             --mailer-to=recipient@example.com
```

### Using configuration file

> [!NOTE]
> CLI arguments always override environment variables and configuration files.

Configuration resolution order:

1. Configuration - first load from config file if provided,
2. Environment variables - Environment variables override values from configuration file,
3. CLI arguments - CLI arguments override environment variables and configuration file values

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

| Name             | Purpose                | Default Value |
| ---------------- | ---------------------- | ------------- |
| KPOW_KEY_KIND    | Key Kind               | null          |
| KPOW_KEY_PATH    | Path to public key     | null          |
| KPOW_ADVERTISE   | Show pubkey on website | false         |
| KPOW_MAILER_DSN  | Mailer DSN             | null          |
| KPOW_MAILER_FROM | Mailer From            | null          |
| KPOW_MAILER_TO   | Mailer To              | null          |
| KPOW_TITLE       | Page title             | ""            |
| KPOW_PORT        | Port                   | 8080          |
| KPOW_HOST        | Host                   | localhost     |
| KPOW_LOG_LEVEL   | Log Level              | info          |

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
