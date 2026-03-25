use crate::config;
use crate::gist;
use chrono::Local;
use std::env;

pub fn send() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 3 {
        eprintln!("Usage: later send <text>");
        std::process::exit(1);
    }

    let text = args[2..].join(" ");

    let cfg = match config::load() {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
    };

    let content = match gist::get_inbox(&cfg) {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
    };

    let timestamp = Local::now().format("%Y-%m-%d %H:%M").to_string();
    let new_line = format!("[{}]  {}", timestamp, text);

    let new_content = if content.is_empty() {
        format!("LATER\n\n{}\n", new_line)
    } else {
        if !content.ends_with('\n') {
            format!("{}\n{}\n", content, new_line)
        } else {
            format!("{}{}\n", content, new_line)
        }
    };

    if let Err(e) = gist::update_inbox(&cfg, &new_content) {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }

    println!("✓ saved");
}
