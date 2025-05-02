# KPow

Simple loopback form.

## Starting server

> [!NOTE] CLI arguments always override environment variables and configuration files.

### Using cli arguments

```sh
$ kpow start --pubkey=public_key_path \
             --password=password \
             --port=8080 \
             --host=localhost \
             --mailer-dsn=smtp://user:password@smtp.example.com:587
```

### Using configuration file

```sh
$ kpow start --config=path-to-config.toml
```

### Environment variables

| Name             | Purpose            | Value |
| ---------------- | ------------------ | ----- |
| KPOW_PUBKEY      | Path to public key | null  |
| KPOW_PASSWORD    | Password           | null  |
| KPOW_PORT        | Port               | null  |
| KPOW_HOST        | Host               | null  |
| KPOW_MAILER_DSN  | Mailer DSN         | null  |
| KPOW_MAILER_FROM | Mailer From        | null  |
| KPOW_MAILER_TO   | Mailer To          | null  |
