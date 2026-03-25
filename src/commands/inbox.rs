use crate::config;
use crate::gist;
use std::env;

pub fn inbox() {
    let args: Vec<String> = env::args().collect();

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

    let lines: Vec<&str> = content.lines().collect();
    let (non_entry_lines, entry_lines): (Vec<&str>, Vec<&str>) =
        lines.iter().partition(|l| !l.starts_with('['));

    if args.len() >= 3 && args[2] == "--clear" {
        let new_content: String = non_entry_lines
            .iter()
            .filter(|l| !l.is_empty())
            .map(|l| *l)
            .collect::<Vec<_>>()
            .join("\n");
        let final_content = if new_content.is_empty() {
            String::new()
        } else {
            new_content + "\n"
        };

        if let Err(e) = gist::update_inbox(&cfg, &final_content) {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
        println!("✓ cleared");
        return;
    }

    if args.len() >= 3 && args[2] == "--pop" {
        if args.len() < 4 {
            eprintln!("Usage: later inbox --pop <n> [n...]");
            std::process::exit(1);
        }

        let indices: Result<Vec<usize>, _> = args[3..].iter().map(|s| s.parse::<usize>()).collect();
        let indices = indices.unwrap_or_else(|_| {
            eprintln!("Error: invalid index");
            std::process::exit(1);
        });

        for &idx in &indices {
            if idx == 0 || idx > entry_lines.len() {
                eprintln!("Error: invalid index: {}", idx);
                std::process::exit(1);
            }
        }

        let new_entries: Vec<&str> = entry_lines
            .iter()
            .enumerate()
            .filter(|(i, _)| !indices.contains(&(i + 1)))
            .map(|(_, &l)| l)
            .collect();

        let new_content: String = non_entry_lines
            .iter()
            .filter(|l| !l.is_empty())
            .chain(new_entries.iter())
            .map(|&l| l)
            .collect::<Vec<_>>()
            .join("\n");
        let final_content = if new_content.is_empty() {
            String::new()
        } else {
            new_content + "\n"
        };

        if let Err(e) = gist::update_inbox(&cfg, &final_content) {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }

        let removed = indices.len();
        if removed == 1 {
            println!("✓ removed 1 entry");
        } else {
            println!("✓ removed {} entries", removed);
        }
        return;
    }

    if entry_lines.is_empty() {
        println!("inbox is empty");
        return;
    }

    for (i, line) in entry_lines.iter().enumerate() {
        println!("{}  {}", i + 1, line);
    }
}
