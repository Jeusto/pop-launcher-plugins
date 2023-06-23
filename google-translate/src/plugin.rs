use crate::config::TranslatePluginConfig;
use pop_launcher_toolkit::{
    launcher::{PluginResponse, PluginSearchResult},
    plugin_trait::{async_trait, PluginExt},
};
use reqwest::Client;

pub struct TranslatePlugin {
    query: String,
    translation: String,
    config: TranslatePluginConfig,
}

#[async_trait]
impl PluginExt for TranslatePlugin {
    fn name(&self) -> &str {
        "google-translate"
    }

    async fn search(&mut self, query: &str) {
        // Remove prefix and suffix from query
        self.query = query.strip_prefix("tr ").unwrap_or_default().to_string();
        self.query = self
            .query
            .split("--")
            .next()
            .unwrap_or_default()
            .to_string();

        // Detect languages from query
        let (source_language, target_language) = self.extract_languages(query);
        if let Some(source_language) = source_language {
            self.config.source_language = source_language;
        }
        if let Some(target_language) = target_language {
            self.config.target_language = target_language;
        }

        println!("Query: {}", self.query);
        println!("Source language: {}", self.config.source_language);
        println!("Target language: {}", self.config.target_language);

        if self.translation.is_empty() {
            self.show_single_result("Start typing and then press enter to search.")
                .await;
        } else {
            self.show_single_result(&self.translation).await;
        }
    }

    async fn activate(&mut self, _id: u32) {
        self.respond_with(PluginResponse::Fill("tr ".to_string()))
            .await;

        self.show_single_result("Translating...").await;

        self.translation = self
            .translate(&self.query)
            .await
            .unwrap_or_else(|_| "Error".to_string());

        self.show_single_result(&self.translation).await;
    }
}

impl TranslatePlugin {
    pub fn new(config: TranslatePluginConfig) -> Self {
        Self {
            query: String::new(),
            translation: String::new(),
            config,
        }
    }

    async fn show_single_result(&self, message: &str) {
        self.respond_with(PluginResponse::Clear).await;
        self.respond_with(PluginResponse::Append(PluginSearchResult {
            id: 0,
            name: message.to_string(),
            ..Default::default()
        }))
        .await;
        self.respond_with(PluginResponse::Finished).await;
    }

    async fn translate(&self, query: &str) -> Result<String, Box<dyn std::error::Error>> {
        let client = Client::new();

        let resp = client
            .get(self.config.api_url.as_str())
            .query(&[
                ("sl", self.config.source_language.as_str()),
                ("tl", self.config.target_language.as_str()),
                ("q", query),
            ])
            .send()
            .await?
            .json::<serde_json::Value>()
            .await?;

        let translated_word = resp
            .get("sentences")
            .and_then(|sentences| sentences.get(0))
            .and_then(|sentence| sentence.get("trans"))
            .and_then(|trans| trans.as_str())
            .ok_or("Unable to parse response")?;

        Ok(translated_word.to_string())
    }

    fn extract_languages(&self, query: &str) -> (Option<String>, Option<String>) {
        let mut languages = query.split_whitespace();
        let mut source_language = None;
        let mut target_language = None;

        while let Some(word) = languages.next() {
            if word == "--" {
                source_language = languages.next().map(|s| s.to_string());
                target_language = languages.next().map(|s| s.to_string());
                break;
            }
        }

        (source_language, target_language)
    }
}
