use std::env;
use std::error::Error;
use ocrs::{OcrEngine, OcrEngineParams};
use rten_tensor::prelude::*;
use rten_tensor::NdTensor;

fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: {} <image_path>", args[0]);
        std::process::exit(1);
    }

    let img_path = &args[1];
    
    // Use OcrEngineParams::default() or similar? 
    // Wait, let's see if we can use the simplest initialization.
    // I will write a script to download the standard models if needed, or see if it has a default builder.
    // For now, let's try cargo check to see if we can just initialize it.
    println!("Init...");
    Ok(())
}
