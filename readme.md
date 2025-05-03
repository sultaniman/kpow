# KPow 💥

Simple loopback form.

## Starting server

> [!NOTE] CLI arguments always override environment variables and configuration files.

### Using cli arguments

```sh
$ kpow start --pubkey=public_key_path \
             --password=password \
             --port=8080 \
             --host=localhost \
             --mailer=smtp://user:password@smtp.example.com:587 \
             --mailer-from=sender@example.com \
             --mailer-to=recipient@example.com
```

### Using configuration file

```sh
$ kpow start --config=path-to-config.toml
```

### Environment variables

| Name             | Purpose            | Value |
| ---------------- | ------------------ | ----- |
| KPOW_PUBKEY_PATH | Path to public key | null  |
| KPOW_PASSWORD    | Password           | null  |
| KPOW_PORT        | Port               | null  |
| KPOW_HOST        | Host               | null  |
| KPOW_MAILER_DSN  | Mailer DSN         | null  |
| KPOW_MAILER_FROM | Mailer From        | null  |
| KPOW_MAILER_TO   | Mailer To          | null  |
