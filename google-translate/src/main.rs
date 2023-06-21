mod config;
mod plugin;

use config::PluginConfig;
use plugin::Plugin;
use pop_launcher_toolkit::plugin_trait::PluginExt;
use std::{io, process::ExitCode};

#[tokio::main]
async fn main() -> io::Result<ExitCode> {
    let plugin_config = PluginConfig::new();
    let mut plugin = Plugin::new(plugin_config);

    plugin.run().await;
    Ok(ExitCode::SUCCESS)
}
