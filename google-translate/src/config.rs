use serde::{Deserialize, Serialize};
use std::{default::Default, fs, path::PathBuf};

#[derive(Debug, Serialize, Deserialize)]
pub struct PluginConfig {
    pub source_language: String,
    pub target_language: String,
    pub api_url: String,
}

impl PluginConfig {
    pub fn new() -> Self {
        let mut config_path = PathBuf::new();
        config_path.push(std::env::var("HOME").unwrap());
        config_path.push(".local/share/pop-launcher/plugins/google-translate/plugin.ron");

        let config_file = fs::read_to_string(config_path).unwrap_or_default();
        let config: PluginConfig = ron::from_str(&config_file).unwrap_or_default();

        PluginConfig {
            source_language: config.source_language,
            target_language: config.target_language,
            api_url: config.api_url,
        }
    }
}

impl Default for PluginConfig {
    fn default() -> Self {
        PluginConfig {
            source_language: "en".to_string(),
            target_language: "fr".to_string(),
            api_url: "".to_string(),
        }
    }
}
