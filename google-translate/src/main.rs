mod config;
mod plugin;

use config::TranslatePluginConfig;
use plugin::TranslatePlugin;
use pop_launcher_toolkit::plugin_trait::PluginExt;
use std::{io, process::ExitCode};

#[tokio::main]
async fn main() -> io::Result<ExitCode> {
    let mut plugin = TranslatePlugin::new(TranslatePluginConfig::new());
    plugin.run().await;
    Ok(ExitCode::SUCCESS)
}
