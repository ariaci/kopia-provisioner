# kopia-provisioner

A lightweight, idiomatic Go tool for declarative user management on Kopia servers. The *kopia-provisioner* reads a YAML configuration, normalizes it per host, and executes deterministic Add/Update/Delete operations against a Kopia server. Each action is reported using a Git‑style status notation (`-`, `A`, `U`, `D`, `!A`, `!U`, `!D`) for immediate clarity.

The design emphasizes modularity, predictable behavior, and clean separation of concerns.

## ✨ Features

- Declarative user provisioning via YAML
- Host-specific overrides for passwords and identity settings
- Supply passwords using chainable providers (`file`, `env`, `inline`) with support for provider pipelines (`>`)
- Typed password backends (`nil`, `plain`, `kopia-hash`)
- Deterministic normalization of all identities
- Clean Kopia CLI integration through a structured argument builder
- Git-style action reporting (`-`, `A`, `U`, `D`, `!A`, `!U`, `!D`)
- Context-aware validation (passwords required only where Kopia needs them)
- Minimal, robust Go architecture
- Simulation mode (--dry-run, -d) to preview all actions without executing Kopia commands
- Selective update (--update, -u) for users present in both YAML and the Kopia server
- Compute password hashes through the Kopia CLI to ensure full compatibility with Kopia’s internal hashing format

## 📦 Installation

### Using Go

```bash
go install github.com/ariaci/kopia-provisioner@latest
```

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

The provisioner is driven by a YAML file describing users and optional host-specific overrides.

Example:
```yaml
users:
  bob:
    default:
      password: "kopia-hash:<kopia password hash for bob@*>"
    hosts:
      home:
        password: "plain:<plain password for bob@home>"
      work:
  jeff:
    default:
      password: "plain:<plain default password for jeff@*>"
    hosts:
      - home
  sarah:
    default:
      password: "kopia-hash:<kopia password hash for sarah@*>"
    hosts:
      - work
```

## Normalization rules

- Defaults apply unless overridden by a host entry.
- Host overrides replace fields atomically.
- Passwords are required only for actions that need them (add, update).
- Delete operations do not require a password.
- The normalized configuration always reflects the final, effective state for a given host.

## 🚀 Running the provisioner

### Set up an initial yaml configuration

```bash
kopia-provisioner config init [-s|--scope password:{user|identity}] [-c|--config-file <kopia-repository>.config]
```

1) Determine all identities and users from the given repository
2) Generate an initial YAML configuration containing these entries
3) Add a ```nil:``` placeholder password for each identity/user based on configured password-scope (defaults to ```identity``` if no password scope is given)
4) Print the resulting configuration to stdout

Example:
```bash
kopia-provisioner config init --scope password:user > users.yaml
```

### Synchronize users with Kopia

```bash
kopia-provisioner users sync <your-yaml-config> [-u|--update] [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute Add, Delete (or Update) operations
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
kopia-provisioner users remove <your-yaml-config> [-u|--update] [-c|--config-file <kopia-repository>.config]
```

The tool will:
1) Load the YAML configuration
2) Normalize all identities for the selected host
3) Execute **only Delete** (or Update) operations
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

The provisioner uses a compact, Git‑inspired status notation:
| Symbol | Meaning |
|--------|---------|
| `-` | Nothing done |
| `A` | User added |
| `U` | User updated |
| `D` | User deleted |
| `!A` | Add failed |
| `!U` | Update failed |
| `!D` | Remove failed |

Example:

```bash
 A bob@home
 U bob@work
 D jeff@home
!D sarah@work
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
- Status renderer — prints -/A/U/D/!A/!U/!D

This separation keeps the codebase predictable, testable, and easy to extend.

## 📄 License

MIT — free to use, modify, and integrate.