use serde::Deserialize;
use std::fs;
use std::io;
use std::path::PathBuf;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub gist: GistConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct GistConfig {
    pub token: String,
    pub gist_id: String,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            gist: GistConfig {
                token: String::new(),
                gist_id: String::new(),
            },
        }
    }
}

fn config_path() -> PathBuf {
    if cfg!(target_os = "windows") {
        let base = std::env::var("APPDATA").unwrap_or_else(|_| {
            format!(
                "{}\\AppData\\Roaming",
                std::env::var("USERPROFILE").unwrap()
            )
        });
        PathBuf::from(base).join("later").join("config.toml")
    } else {
        directories::BaseDirs::new()
            .map(|d| d.config_dir().join("later").join("config.toml"))
            .unwrap_or_else(|| PathBuf::from(".later/config.toml"))
    }
}

pub fn ensure_file() -> io::Result<PathBuf> {
    let path = config_path();
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)?;
    }
    if !path.exists() {
        let template = r#"# later config
[gist]
token = ""
gist_id = ""
"#;
        fs::write(&path, template)?;
    }
    Ok(path)
}

pub fn save(token: &str, gist_id: &str) -> io::Result<PathBuf> {
    let path = ensure_file()?;
    let content = format!("[gist]\ntoken = \"{}\"\ngist_id = \"{}\"\n", token, gist_id);
    fs::write(&path, content)?;
    Ok(path)
}

pub fn load() -> Result<Config, Box<dyn std::error::Error>> {
    let path = config_path();
    if !path.exists() {
        return Err(format!(
            "config not found at {} — run 'later config --init'",
            path.display()
        )
        .into());
    }
    let content = fs::read_to_string(&path)?;
    let cfg: Config =
        toml::from_str(&content).map_err(|e| format!("failed to parse config: {}", e))?;
    if cfg.gist.token.is_empty() {
        return Err("missing 'token' in config — run 'later config --init'".into());
    }
    if cfg.gist.gist_id.is_empty() {
        return Err("missing 'gist_id' in config — run 'later config --init'".into());
    }
    Ok(cfg)
}

#[allow(dead_code)]
pub fn path() -> PathBuf {
    config_path()
}
