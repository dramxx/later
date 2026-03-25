mod commands;
mod config;
mod gist;

fn main() {
    let args: Vec<String> = std::env::args().collect();

    if args.len() < 2 {
        print_usage();
        std::process::exit(1);
    }

    if args[1] == "-h" || args[1] == "--help" {
        print_usage();
        return;
    }

    match args[1].as_str() {
        "send" => commands::send::send(),
        "inbox" => commands::inbox::inbox(),
        "config" => commands::config::config(),
        _ => {
            eprintln!("Unknown command: {}", args[1]);
            print_usage();
            std::process::exit(1);
        }
    }
}

fn print_usage() {
    println!("Usage: later <command>");
    println!("Commands: send, inbox, config");
    println!();
    println!("Examples:");
    println!("  later config --init");
    println!("  later send https://example.com");
    println!("  later inbox");
    println!("  later inbox --clear");
    println!("  later inbox --pop 1");
    println!("  later config");
}
