## [![Test](https://github.com/sultaniman/kpow/actions/workflows/test.yml/badge.svg)](https://github.com/sultaniman/kpow/actions/workflows/test.yml)

<a href="https://coff.ee/sultaniman" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Offrez-moi un café" height="36"></a>

# KPow 💥

[English](../../readme.md) | [Deutsch](../de/readme.md) | [Türkçe](../tr/readme.md) | [Qyrgyz](../qy/readme.md) | [Français](readme.md) | [Українська](../uk/readme.md) | [Русский](../ru/readme.md)

KPow est un formulaire de contact auto-hébergé, axé sur la confidentialité, conçu pour une communication sécurisée sans dépendre de services tiers.
Il prend en charge les standards de chiffrement modernes — PGP, Age et RSA — pour garantir que les messages sont chiffrés avant leur livraison.
Idéal pour les développeurs soucieux de la vie privée, les projets open source, les sites web indépendants, les plateformes de lancement d'alerte et les outils internes nécessitant un traitement de messages sécurisé, auditable et autonome.

## Démarrer le serveur

### Avec des arguments CLI

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

### Avec un fichier de configuration

> [!note]
> Les arguments CLI ont toujours priorité sur les variables d'environnement et les fichiers de configuration.

Ordre de résolution de la configuration :

1. Configuration — charger d'abord depuis le fichier de configuration s'il est fourni,
2. Variables d'environnement — écrasent les valeurs du fichier de configuration,
3. Arguments CLI — écrasent les variables d'environnement et les valeurs du fichier de configuration

```mermaid
flowchart TD
    A[Démarrage] --> B{Fichier de config fourni ?}
    B -- Oui --> C[Charger le fichier de config]
    B -- Non --> D[Utiliser les valeurs par défaut]
    C --> D
    D --> E[Charger les variables d'environnement]
    E --> F[Appliquer les arguments CLI]
```

```sh
$ kpow start --config=path-to-config.toml
```

### Vérifier un fichier de configuration

Exécutez la commande `verify` pour charger une configuration et signaler tout
problème de validation sans démarrer le serveur :

```sh
$ kpow verify --config=path-to-config.toml
```

### Variables d'environnement

| Nom de la variable      | Description                                     | Type   | Par défaut    |
| ----------------------- | ----------------------------------------------- | ------ | ------------- |
| `KPOW_TITLE`            | Titre du serveur                                | string | ""            |
| `KPOW_PORT`             | Port du serveur                                 | int    | 8080          |
| `KPOW_HOST`             | Adresse hôte du serveur                         | string | localhost     |
| `KPOW_LOG_LEVEL`        | Niveau de journalisation                        | string | INFO          |
| `KPOW_MESSAGE_SIZE`     | Taille maximale des messages du serveur         | int    | 240           |
| `KPOW_HIDE_LOGO`        | Masquer le logo                                 | bool   | false         |
| `KPOW_CUSTOM_BANNER`    | Fichier de bannière personnalisée               | string | ""            |
| `KPOW_LIMITER_RPM`      | Limiteur de débit : requêtes par minute         | int    | 0             |
| `KPOW_LIMITER_BURST`    | Limiteur de débit : taille de rafale            | int    | -1            |
| `KPOW_LIMITER_COOLDOWN` | Limiteur de débit : temps de récupération (sec) | int    | -1            |
| `KPOW_MAILER_FROM`      | Adresse e-mail de l'expéditeur                  | string | ""            |
| `KPOW_MAILER_TO`        | Adresse e-mail du destinataire                  | string | ""            |
| `KPOW_MAILER_DSN`       | DSN du mailer (chaîne de connexion)             | string | ""            |
| `KPOW_WEBHOOK_URL`      | URL du webhook                                  | string | ""            |
| `KPOW_MAX_RETRIES`      | Nombre maximal de tentatives d'envoi            | int    | 2             |
| `KPOW_KEY_KIND`         | Type de clé : `age`, `pgp` ou `rsa`             | string | ""            |
| `KPOW_ADVERTISE`        | Publier la clé                                  | bool   | false         |
| `KPOW_KEY_PATH`         | Chemin vers le fichier de clé                   | string | ""            |
| `KPOW_INBOX_PATH`       | Chemin vers la boîte de réception               | string | ""            |
| `KPOW_INBOX_CRON`       | Planification cron pour la boîte de réception   | string | `*/5 * * * *` |

> [!note]
> `KPOW_MAILER_DSN` ou `KPOW_WEBHOOK_URL` doit être fourni pour que KPow puisse livrer les messages.

## Chiffrement

KPow prend en charge les clés publiques Age, PGP et RSA pour le chiffrement des messages.
Indiquez le type de clé avec `--key-kind` (ou `KPOW_KEY_KIND`) et le
chemin vers votre clé publique avec `--pubkey` (ou `KPOW_KEY_PATH`).
Options disponibles pour `--key-kind` : `age`, `pgp` ou `rsa`.

### Génération de clés

Utilisez les outils en ligne de commande courants pour créer des clés publiques compatibles :

#### Age

```sh
age-keygen -o age.key
grep "^# public key:" age.key | cut -d' ' -f3 > age.pub
```

Utilisez `age.pub` comme valeur pour `--pubkey` (ou `KPOW_KEY_PATH`).

#### PGP

```sh
gpg --quick-generate-key "Votre Nom <vous@example.com>"
gpg --armor --export vous@example.com > pgp.pub
```

Passez le fichier `pgp.pub` en format ASCII-armored à `--pubkey`.

#### RSA

```sh
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in rsa_private.pem -out rsa_public.pem
```

Fournissez `rsa_public.pem` comme `--pubkey`. La clé publique doit être une
clé RSA encodée en PEM au format PKIX (2048 bits ou plus).

### Exemple de fichier de configuration

Au lieu des options CLI, spécifiez la clé dans un fichier de configuration TOML :

```toml
[key]
kind = "age"           # ou "pgp" ou "rsa"
path = "/etc/kpow/key.pub"
advertise = false
```

### Note sur le chiffrement RSA

Ce système utilise le chiffrement RSA avec le rembourrage OAEP et l'algorithme de hachage SHA-256.
Veuillez suivre ces directives lors de l'utilisation de clés RSA et de la configuration des paramètres de message :

**Exigences de clé et d'algorithme**

- **Compatibilité des clés RSA :** Doit prendre en charge le rembourrage OAEP (taille recommandée : 2048 bits ou plus).
- **Algorithme de hachage :** Le chiffrement utilise SHA-256 — le déchiffrement doit utiliser le même.

**Surcharge du rembourrage OAEP**

- Taille du rembourrage = 2 x TailleHash + 2 octets
- Pour SHA-256 (TailleHash = 32 octets), la surcharge totale est de 66 octets

**Tailles maximales des messages**

| Taille de clé RSA | Algorithme de hachage | Taille du hash | Taille du rembourrage | Taille max. du message |
| ----------------- | --------------------- | -------------- | --------------------- | ---------------------- |
| 2048 bits         | SHA-256               | 32 octets      | 66 octets             | 190 octets             |
| 4096 bits         | SHA-256               | 32 octets      | 66 octets             | 446 octets             |

Les messages dépassant la taille maximale pour la clé seront tronqués avant le chiffrement.

**Conseil de configuration**

Dans votre configuration TOML (`message_size`), définissez la valeur en dessous de la taille maximale du message en fonction de la longueur de votre clé RSA. Par exemple :

```toml
[server]
message_size = 180  # pour RSA 2048 bits avec SHA-256
```

## Logique du mailer

```mermaid
flowchart TD
    A[Nouveau message soumis] --> B{Essayer d'envoyer immédiatement ?}
    B -- Succès --> C[Message envoyé]
    B -- Erreur --> D[Sauvegarder dans la boîte de réception]
    D --> E[Exécution du planificateur]
    E --> F[Lire les messages]
    F --> G{Essayer de renvoyer}
    G -- Succès --> H[Message envoyé]
    G -- Erreur --> E
```

## Webhook

Lorsque `--webhook-url` (ou `KPOW_WEBHOOK_URL`) est fourni, KPow envoie les
données du formulaire chiffrées au point de terminaison spécifié en format JSON via POST :

```json
{
    "subject": "<form subject>",
    "content": "<encrypted message>",
    "hash": "<sha256-hash>"
}
```

L'URL du webhook doit utiliser HTTPS sauf si elle pointe vers `localhost`. Tout code de statut HTTP
< 400 est considéré comme un succès.

## Docker

KPow est livré avec un Dockerfile et peut être facilement déployé dans des conteneurs :

```sh
docker build -t kpow .
docker run -p 8080:8080 \
  -v /path/to/key.pub:/app/key.pub \
  -e KPOW_KEY_KIND=age \
  -e KPOW_KEY_PATH=/app/key.pub \
  -e KPOW_WEBHOOK_URL=https://hooks.example.com/notify \
  kpow
```

## Vérification de santé

KPow fournit un point de terminaison `/health` pour l'orchestration de conteneurs et les répartiteurs de charge :

```sh
curl http://localhost:8080/health
# {"status":"ok"}
```

## Développement

### Formulaire personnalisé

Bun et Tailwind CSS sont utilisés pour construire les styles.
Les sources des styles se trouvent dans le dossier `styles`.
Utilisez `just styles` pour personnaliser et construire les styles du formulaire, et
`just error-styles` pour les pages d'erreur.
Les deux commandes nécessitent l'installation de `bun` et `bunx`.

### Bannière personnalisée

Il est possible de personnaliser le formulaire et d'ajouter une bannière personnalisée en utilisant `--banner=/path/to/banner.html` ou en définissant `KPOW_CUSTOM_BANNER=/path/to/banner.html`.
Le HTML de la bannière fournie sera assaini ; vous trouverez ci-dessous la liste des balises autorisées.

**Balises autorisées**

> [!note]
> Vous pouvez utiliser l'attribut `style` pour styliser votre bannière.

- `a`
- `p`
- `span`
- `img`
- `div`
- `ul,ol,li`
- `h1-h6`

## Licence

KPow est sous licence **Business Source License 1.1**.

Vous **ne pouvez pas utiliser** ce logiciel pour offrir un service commercial hébergé ou géré à des tiers sans acheter une licence commerciale séparée.

Le **04/12/2028**, ce projet sera re-licencié sous la **Apache License 2.0**.

- 📄 Voir [`LICENSE`](../../LICENSE)
- 📄 Voir [`LICENSE-BUSL`](../../LICENSE-BUSL)
- 📄 Voir [`LICENSE-APACHE`](../../LICENSE-APACHE)

## Captures d'écran

## ![form](https://github.com/sultaniman/kpow/blob/main/screenshots/form.png?raw=true)

## ![rate limited](https://github.com/sultaniman/kpow/blob/main/screenshots/rate-limited.png?raw=true)

![csrf error](https://github.com/sultaniman/kpow/blob/main/screenshots/csrf-error.png?raw=true)

<p align="center">✨ 🚀 ✨</p>
