use crate::config::PluginConfig;
use pop_launcher_toolkit::{
    launcher::{PluginResponse, PluginSearchResult},
    plugin_trait::{async_trait, PluginExt},
};
use reqwest::Client;

pub struct Plugin {
    query: String,
    config: PluginConfig,
}

impl Plugin {
    pub fn new(config: PluginConfig) -> Self {
        Self {
            query: String::new(),
            config,
        }
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

        println!("{}", translated_word);

        Ok(translated_word.to_string())
    }
}

#[async_trait]
impl PluginExt for Plugin {
    fn name(&self) -> &str {
        "vscode-open-recent"
    }

    async fn search(&mut self, query: &str) {
        let query = query.strip_prefix("tr ");

        self.query = match query {
            Some(q) => q.to_string(),
            None => "".to_string(),
        };

        self.respond_with(PluginResponse::Append(PluginSearchResult {
            id: 0,
            name: "Start typing and then press enter to search.".to_string(),
            description: "".to_string(),
            keywords: None,
            icon: None,
            exec: None,
            window: None,
        }))
        .await;
        self.respond_with(PluginResponse::Finished).await;
    }

    async fn activate(&mut self, _id: u32) {
        self.respond_with(PluginResponse::Append(PluginSearchResult {
            id: 0,
            name: "Translating...".to_string(),
            description: "".to_string(),
            keywords: None,
            icon: None,
            exec: None,
            window: None,
        }))
        .await;
        self.respond_with(PluginResponse::Finished).await;

        let translation = self
            .translate(&self.query)
            .await
            .unwrap_or_else(|_| "Error".to_string());

        self.respond_with(PluginResponse::Clear).await;
        self.respond_with(PluginResponse::Append(PluginSearchResult {
            id: 0,
            name: translation,
            description: "".to_string(),
            keywords: None,
            icon: None,
            exec: None,
            window: None,
        }))
        .await;
        self.respond_with(PluginResponse::Finished).await;
    }
}
