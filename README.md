# later

Share notes with yourself across machines, via cli, using a private GitHub Gist as storage.

```
later send look at https://example.com

later inbox
```

---

## Setup

### 1. Create your private Gist

1. Go to [gist.github.com](https://gist.github.com)
2. Click **+** (New gist) in the top right
3. Set filename to `inbox.txt`
4. Put any placeholder text (GitHub requires at least one character)
5. Click **Create secret gist**
6. Copy the gist ID from the URL:
   `gist.github.com/yourusername/` **`1d726ef02757ca62c48defa1ab646bdc`**

---

### 2. Create a GitHub token

1. Go to GitHub → **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
2. Click **Generate new token (classic)**
3. Give it a name like `later`
4. Tick only the **`gist`** scope — nothing else needed
5. Click **Generate token**
6. **Copy it immediately** — GitHub shows it only once

---

### 3. Create the config file

```toml
[gist]
token = "ghp_xxxx..."
gist_id = "1d726ef02757ca62c48defa1ab646bdc"
```

**Windows:** `%APPDATA%\later\config.toml`

**Linux / Mac:** `~/.config/later/config.toml`

Repeat this on every machine you want to use `later` on. The config is the only setup required per machine.

---

### 4. Install the binary

**Windows** (PowerShell):

```powershell
powershell -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/dramxx/later/main/install.ps1 | iex"
```

**Linux / Mac:**

```bash
curl -sSL https://raw.githubusercontent.com/dramxx/later/main/install.sh | sh
```

The installer downloads the binary for your platform and adds it to your PATH automatically. Open a new terminal after installing.

---

## Usage

### Edit your config

```
later config
```

Opens `config.toml` in a text editor (Notepad on Windows, gedit on Linux). Creates the file with an empty template if it doesn't exist yet.

---

### Send something

```
later send https://magazine.sebastianraschka.com/
later send "read the KV cache section carefully"
later send https://arxiv.org/abs/2501.12345
```

---

### Read your inbox

```
later inbox
```

Output:

```
1  [2026-03-24 19:42]  https://magazine.sebastianraschka.com/
2  [2026-03-24 19:44]  read the KV cache section carefully
3  [2026-03-24 19:45]  https://arxiv.org/abs/2501.12345
```

---

### Clear everything

```
later inbox --clear
```

---

### Remove specific entries

```
later inbox --pop 2          # remove entry 2
later inbox --pop 1 3        # remove entries 1 and 3
```
