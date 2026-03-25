use crate::config::Config;
use serde::{Deserialize, Serialize};

const BASE_URL: &str = "https://api.github.com";

#[derive(Deserialize)]
struct GistResponse {
    files: std::collections::HashMap<String, GistFile>,
}

#[derive(Deserialize)]
struct GistFile {
    content: String,
}

#[derive(Serialize)]
struct PatchRequest {
    files: std::collections::HashMap<String, PatchFile>,
}

#[derive(Serialize)]
struct PatchFile {
    content: String,
}

pub fn get_inbox(cfg: &Config) -> Result<String, Box<dyn std::error::Error>> {
    let url = format!("{}/gists/{}", BASE_URL, cfg.gist.gist_id);
    let client = reqwest::blocking::Client::builder()
        .timeout(std::time::Duration::from_secs(10))
        .build()?;

    let response = client
        .get(&url)
        .header("User-Agent", "later-cli")
        .header("Authorization", format!("Bearer {}", cfg.gist.token))
        .header("Accept", "application/vnd.github+json")
        .send()?;

    let status = response.status();
    if !status.is_success() {
        let body = response.text()?;
        return Err(format!("HTTP {}: {}", status, body).into());
    }

    let gist: GistResponse = response.json()?;
    Ok(gist
        .files
        .get("inbox.txt")
        .map(|f| f.content.clone())
        .unwrap_or_default())
}

pub fn update_inbox(cfg: &Config, content: &str) -> Result<(), Box<dyn std::error::Error>> {
    let url = format!("{}/gists/{}", BASE_URL, cfg.gist.gist_id);

    let mut files = std::collections::HashMap::new();
    files.insert(
        "inbox.txt".to_string(),
        PatchFile {
            content: content.to_string(),
        },
    );

    let patch = PatchRequest { files };

    let client = reqwest::blocking::Client::builder()
        .timeout(std::time::Duration::from_secs(10))
        .build()?;

    let response = client
        .patch(&url)
        .header("User-Agent", "later-cli")
        .header("Authorization", format!("Bearer {}", cfg.gist.token))
        .header("Accept", "application/vnd.github+json")
        .header("Content-Type", "application/json")
        .json(&patch)
        .send()?;

    let status = response.status();
    if !status.is_success() {
        let body = response.text()?;
        return Err(format!("HTTP {}: {}", status, body).into());
    }

    Ok(())
}
