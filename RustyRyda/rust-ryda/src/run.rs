use std::process::{Command, ExitStatus};
use std::io::{self, Write};

fn main() -> io::Result<()> {
    // List of Go executable files for gRPC services
    let go_executables = vec![
        "path/to/service1",
        "path/to/service2",
        // Add more services as needed
    ];

    let mut handles = Vec::new();

    // Start each Go service
    for exe in go_executables {
        match Command::new(exe)
            .spawn() {
                Ok(child) => {
                    println!("Started: {}", exe);
                    handles.push(child);
                },
                Err(e) => {
                    eprintln!("Failed to start {}: {}", exe, e);
                }
            }
    }

    // Wait for all processes to complete
    for handle in handles {
        match handle.wait() {
            Ok(status) => {
                println!("Process exited with: {}", status);
            },
            Err(e) => {
                eprintln!("Failed to wait on child process: {}", e);
            }
        }
    }

    println!("All Go services completed.");
    Ok(())
}
