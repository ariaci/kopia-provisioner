# kopia-provisioner

A lightweight, idiomatic Go tool for declarative user management on Kopia servers. The *kopia-provisioner* reads a YAML configuration, normalizes it per host, and executes deterministic Add/Update/Delete operations against a Kopia server. Each action is reported using a Git‑style status notation (`AI`, `UI`, `DI`, `!AI`, `!UI`, `!DI`, `-DI`, ...) for immediate clarity.

The design emphasizes modularity, predictable behavior, and clean separation of concerns.

## ✨ Features

- Declarative user provisioning via YAML
- Host-specific overrides for passwords and identity settings
- Supply passwords using chainable providers (`file`, `env`, `inline`) with support for provider pipelines (`>`)
- Typed password backends (`nil`, `plain`, `kopia-hash`)
- Deterministic normalization of all identities
- Clean Kopia CLI integration through a structured argument builder
- Delete operations are snapshot-aware - override this behaviour by adding `--force`
- Git-style action reporting using a compact and forward-compatible 3‑character DSL (`AI`, `UI`, `DI`, `!AI`, `!UI`, `!DI`, `-DI`, ...)
- Context-aware validation (passwords required only where Kopia needs them)
- Minimal, robust Go architecture
- Simulation mode (`--dry-run`, `-d`) to preview all actions without executing Kopia commands
- Selective update (`--update`, `-u`) for users present in both YAML and the Kopia server
- Compute password hashes through the Kopia CLI to ensure full compatibility with Kopia’s internal hashing format

## 📦 Installation

### Using Go

```bash
go install github.com/ariaci/kopia-provisioner@latest
```

#### ⚠️ Note

Installing via `go install` works correctly, but the resulting binary will **not** contain:

* Windows Version Metadata (file properties such as Product Version, File Version, etc.)
* the embedded Git Commit Hash

This is expected behavior, because `go install` builds from the module cache without access to Git metadata.

For a fully featured binary with all resources included, download the pre‑compiled executables from the **Releases** page or build the project manually.

### 🐳 Using Docker

#### Running kopia-provisioner via Docker

The `kopia-provisioner` can also be executed directly from a Docker container.
This is useful when you prefer not to install Go locally, want reproducible environments, or integrate the provisioner into CI/CD pipelines.

#### Pull the image

To pull the latest version:
```bash
docker pull ghcr.io/ariaci/kopia-provisioner:latest
```

To pull a specific version:
```bash
docker pull ghcr.io/ariaci/kopia-provisioner:v1.2.3_kopia.0.23.1
```

#### Running the provisioner

Mount your configuration files into the container and execute the desired command:
```bash
docker run --rm \
  -v $(pwd)/users.yaml:/config/users.yaml \
  -v $(pwd)/repo.config:/config/repo.config \
  ghcr.io/ariaci/kopia-provisioner:latest \
  users sync /config/users.yaml --config-file /config/repo.config
```

#### Example: Add users

```bash
docker run --rm \
  -v $(pwd)/users.yaml:/config/users.yaml \
  -v $(pwd)/repo.config:/config/repo.config \
  ghcr.io/ariaci/kopia-provisioner:latest \
  users add /config/users.yaml --config-file /config/repo.config
```

#### Example: Compute Kopia password hash

```bash
docker run --rm \
  -v $(pwd)/repo.config:/config/repo.config \
  ghcr.io/ariaci/kopia-provisioner:latest \
  hash "my-password" --config-file /config/repo.config
```

#### Notes

* the Docker-Image is based on the official Kopia image and adds the `kopia-provisioner` binary
* the binary is located at `/usr/local/bin/kopia-provisioner` inside the container
* license files are available under `/usr/share/licenses/kopia-provisioner`
* you must mount your Kopia repository configuration (`*.config`) into the container
* additional files (e.g., password files) can be mounted as needed

### Build manually

#### Using Linux

```bash
git clone https://github.com/ariaci/kopia-provisioner 
cd kopia-provisioner 
go generate ./...
go build -o .build/kopia-provisioner
```

#### Using Windows

```bat
git clone https://github.com/ariaci/kopia-provisioner 
cd kopia-provisioner 
go generate ./...
go build -o .build/kopia-provisioner.exe
```

## 🧩 Configuration

The provisioner is driven by a YAML file describing global and optional user- and host-specific overrides.

### Abstract

```yaml
default:
  password: "plain:<plain default password for all identities>"
users:
  bob:
    default:
      password: "kopia-hash:<kopia password hash for all hosts of bob>"
    hosts:
      home:
        password: "plain:<plain password for identity bob@home>"
      work:
  jeff:
    default:
      password: "kopia-hash:<kopia password hash for all hosts of jeff>"
    hosts:
      - home
      - work
  johnny:
    hosts:
      - home
      - work
```

### Example

```yaml
users:
  rudi:
    hosts:
      home:
        password: "plain:rudi_pwd"
  brigitte:
    default:
      password: "plain:brigitte_pwd"
    hosts:
      - home
      - work
```

## Normalization rules

- Defaults apply unless overridden by a specific user default entry or host entry.
- User default overrides and host overrides replace fields atomically.
- Passwords are required only for actions that need them (add, update).
- Delete operations do not require a password.
- The normalized configuration always reflects the final, effective state for a given identity.
- The global default section supports the same fields and password providers as `users.<user>.default`, and participates in the same override hierarchy.

## 🚀 Running the provisioner

### Set up an initial yaml configuration

```bash
kopia-provisioner config init [-s|--scope password:{global|user|identity}] [-c|--config-file <kopia-repository>.config]
```

1) Determine all identities and users from the given repository
2) Generate an initial YAML configuration containing these entries
3) Add a `nil:` placeholder password for each identity/user/global based on configured password-scope (defaults to `identity` if no password scope is given)
4) Print the resulting configuration to stdout

Example:
```bash
kopia-provisioner config init --scope password:user > users.yaml
```

### Synchronize users with Kopia

```bash
kopia-provisioner users sync <your-yaml-config> [-f|--force] [-u|--update] [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute Add, Delete (or Update) operations
   - Skip Delete in case of associated snapshots (status: `-DI`)
   - Use `--force` to force delete of identities with associated snapshots (status: `DI`)
   - Delete operations follow the status DSL:
     - `DI` — identity deleted
     - `-DI` — delete skipped due to associated snapshots
     - `!DI` — delete failed (Kopia error)
4) Print a status line for each action

### Add missing users to Kopia

```bash
kopia-provisioner users add <your-yaml-config> [-u|--update] [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute **only Add** (or Update) operations
4) Print a status line for each action

### Remove users from Kopia

```bash
kopia-provisioner users remove <your-yaml-config> [-f|--force] [-u|--update] [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute **only Delete** (or Update) operations
   - Skip Delete in case of associated snapshots (status: `-DI`)
   - Use `--force` to force delete of identities with associated snapshots (status: `DI`)
   - Delete operations follow the status DSL:
     - `DI` — identity deleted
     - `-DI` — delete skipped due to associated snapshots
     - `!DI` — delete failed (Kopia error)
4) Print a status line for each action

### Update users that exist in both YAML and Kopia

```bash
kopia-provisioner users update <your-yaml-config> [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute **only Update** operations
4) Print a status line for each action

### Generate a Kopia password hash

```bash
kopia-provisioner hash <text> [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Use the Kopia CLI to compute a password hash in the context of the given repository
2) Print the resulting hash, equivalent to running:
   ```bash
   kopia --config-file <kopia-repository>.config server users hash-password --user-password <text>
   ```

## 📊 Status output

The provisioner uses a compact, Git‑inspired status notation.

The formal grammar of the status notation DSL is defined below (EBNF):
```ebnf
status      = state , operation , object ;

state       = " " | "!" | "-" ;
operation   = "A" | "U" | "D" | "-" ;
object      = "I" ;
```

The grammar is designed to be forward‑compatible with additional object types (e.g. snapshots) and operations (e.g. move).

| State     | Meaning              |
|-----------|----------------------|
| `<space>` | Operation succeeded  |
| `!`       | Operation failed     |
| `-`       | Operation suppressed |

| Operation | Meaning |
|-----------|---------|
| `A`       | Add     |
| `U`       | Update  |
| `D`       | Delete  |
| `-`       | No-op   |

| Object | Meaning  |
|--------|----------|
| `I`    | Identity |

### Example

```bash
 AI bob@home      # added bob@home successfully
 UI bob@work      # updated bob@work successfully
 DI jeff@home     # deleted jeff@home successfully
-DI jeff@work     # deleting jeff@work suppressed
!DI sarah@work    # deleting sarah@work failed
```

## 🔐 Password providers

Passwords in the configuration are defined using a small provider‑based DSL.
This mechanism allows you to supply plain passwords or password hashes from different sources such as environment variables, files, or inline YAML values.

Providers can be chained using the `>` character.
Each provider receives the output of the previous one, allowing you to build transformation pipelines.

If no provider is specified, __kopia‑provisioner automatically uses the__ `inline` provider.

The formal grammar of the password DSL is defined below (EBNF):
```ebnf
password        = provider-chain type ":" value ;
provider-chain  = [ backend { ">" backend } ">" ] ;

type            = "plain" | "kopia-hash" | "nil" ;
backend         = "file" | "env" | "inline" ;

value           = { ? any character ? } ;
```

Examples:
```
env>file>plain:DB_PWD_FILE
file>kopia-hash:secrets/.bob-pwdhash
inline>plain:supersecretpassword
plain:supersecretpassword
nil:
```

### 📁 `file` - Load password from a file

Reads the password from a file on disk.
Relative paths are resolved against the configuration’s base directory.

Example:

```yaml
password: "file>plain:secrets/db-password.txt"
```

Behavior:
- value is a file path
- relative paths → resolved against the base directory
- file is read lazily during pipeline execution
- surrounding whitespace is trimmed
- errors if the file does not exist or cannot be read

### 🔧 `env` - Load password from an environment variable

Reads the password from an environment variable.

Example:

```yaml
password: "env>plain:DB_PASSWORD"
```

Behavior:
- value is the environment variable name
- errors if the variable is not set
- whitespace is trimmed

### 🔑 `inline` - Inline password value
A static password defined directly in the configuration.

Example:

```yaml
password: "inline>plain:supersecret"
```

Behavior:
- value is used as‑is
- no additional processing

## 🔐 Password backends

Passwords are provided through a typed interface. Each backend implements:

```Go
KopiaArguments() ([]string, error)
```

### Empty password

```yaml
password: "nil:"
```

### Plain password

```yaml
password: "plain:<password>"
```

### Kopia password hash

```yaml
password: "kopia-hash:<password hash>"
```

## 🛠 Architecture overview

The provisioner is structured into clear layers:

- YAML parsing — loads raw user definitions
- Normalization — merges defaults and host overrides into a final identity
- Password interface — encapsulates password retrieval and argument generation
- Argument builder — constructs server users <action> <identity> ...
- Action executor — runs Kopia commands
- Status renderer — prints status of executed operations

This separation keeps the codebase predictable, testable, and easy to extend.

## 📄 License

MIT — free to use, modify, and integrate.
