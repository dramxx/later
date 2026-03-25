use crate::config;
use std::env;
use std::io::{self, Write};

pub fn config() {
    let args: Vec<String> = env::args().collect();

    if args.len() >= 3 {
        match args[2].as_str() {
            "--init" => {
                init_config();
                return;
            }
            "--path" => match config::ensure_file() {
                Ok(path) => {
                    println!("{}", path.display());
                    return;
                }
                Err(e) => {
                    eprintln!("Error: {}", e);
                    std::process::exit(1);
                }
            },
            _ => {
                eprintln!("Usage: later config [--init|--path]");
                std::process::exit(1);
            }
        }
    }

    let path = match config::ensure_file() {
        Ok(p) => p,
        Err(e) => {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
    };

    match open_editor(&path) {
        Ok(_) => {}
        Err(e) => {
            eprintln!("Error: {}", e);
            println!("Config file: {}", path.display());
            println!("Run 'later config --init' to set it up without an editor.");
            std::process::exit(1);
        }
    }
}

fn init_config() {
    let path = match config::ensure_file() {
        Ok(p) => p,
        Err(e) => {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
    };

    print!("GitHub token (gist scope only): ");
    io::stdout().flush().unwrap();
    let token = read_line();

    print!("Private gist ID: ");
    io::stdout().flush().unwrap();
    let gist_id = read_line();

    if token.is_empty() || gist_id.is_empty() {
        eprintln!("Error: values cannot be empty");
        std::process::exit(1);
    }

    if let Err(e) = config::save(&token, &gist_id) {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }

    println!("✓ config saved to {}", path.display());
}

fn read_line() -> String {
    let mut input = String::new();
    io::stdin().read_line(&mut input).unwrap();
    input.trim().to_string()
}

#[allow(unreachable_code)]
fn open_editor(path: &std::path::Path) -> Result<(), Box<dyn std::error::Error>> {
    if let Some(editor) = env::var("VISUAL").ok().filter(|s| !s.is_empty()) {
        return exec_editor(&editor, path);
    }
    if let Some(editor) = env::var("EDITOR").ok().filter(|s| !s.is_empty()) {
        return exec_editor(&editor, path);
    }

    #[cfg(target_os = "windows")]
    {
        return exec_editor("notepad.exe", path);
    }

    #[cfg(target_os = "macos")]
    {
        return exec_editor("open", path);
    }

    #[cfg(target_os = "linux")]
    {
        for editor in &["xdg-open", "nano", "vim", "vi"] {
            if std::path::Path::new(&format!("/usr/bin/{}", editor)).exists() {
                return exec_editor(editor, path);
            }
        }
    }

    Err("no editor found".into())
}

fn exec_editor(editor: &str, path: &std::path::Path) -> Result<(), Box<dyn std::error::Error>> {
    let output = std::process::Command::new(editor).arg(path).output()?;
    if !output.status.success() {
        return Err(format!("editor exited with code: {}", output.status).into());
    }
    Ok(())
}
