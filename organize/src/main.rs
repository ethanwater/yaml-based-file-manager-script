use std::process::Command;
use std::fs;

static mut CONFIG_FILE: String = String::new();

struct Config {
    name: String,
    path: String,
    ext: Vec<String>,
}

fn SetConfiguration(config: String) {
    unsafe {
        CONFIG_FILE = config;
    }
}

fn Configurations() {
    let config_file = fs::read("config.yaml").expect("cannot read file");


}

fn main() {
    println!("Hello, world!");
}
