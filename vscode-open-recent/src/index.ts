#!/usr/bin/env node
import { RecentlyOpened, Vscode } from "./vscode.js";
import { PluginExt, PluginResponse, respondWith, runPlugin } from "./pop.js";
import { cleanSearchQuery, getIcon, log } from "./utils.js";

// Main execution
function main() {
  Plugin.init().then((plugin) => {
    log(JSON.stringify(plugin, null, 2));
    runPlugin(plugin);
  });
}

class Plugin implements PluginExt {
  vscode: Vscode | null;
  private query: string;
  private results: RecentlyOpened;

  private constructor(vscode: Vscode) {
    this.vscode = vscode;
    this.query = "";
    this.results = [];
  }

  static async init(): Promise<Plugin> {
    return new Plugin(await Vscode.init());
  }

  name(): string {
    return "vscode-open-recent";
  }

  search(query: string) {
    this.query = cleanSearchQuery(query, "vs");

    if (!this.query) {
      this.respond("Clear");
      this.vscode?.recentlyOpened.forEach((entry, idx) => {
        this.respond({
          Append: {
            id: idx,
            name: entry.name,
            description: entry.path,
            keywords: null,
            icon: { Name: getIcon(entry.type) },
            exec: null,
            window: null,
          },
        });
      });
      this.respond("Finished");
      return;
    }

    this.results = this.vscode?.search(this.query) || [];

    if (this.results.length === 0) {
      this.respond("Clear");
      this.respond({
        Append: {
          id: 0,
          name: "No results found.",
          description: "Try a different query.",
          keywords: null,
          icon: { Name: getIcon("info") },
          exec: null,
          window: null,
        },
      });

      this.respond("Finished");
      return;
    }

    this.respond("Clear");
    this.results.forEach((entry, idx) => {
      this.respond({
        Append: {
          id: idx,
          name: entry.name,
          description: entry.path,
          keywords: null,
          icon: { Name: getIcon(entry.type) },
          exec: null,
          window: null,
        },
      });
    });
    this.respond("Finished");
  }

  activate(index: number) {
    const entry = this.results[index];
    this.vscode?.open(entry);
    this.respond("Close");
  }

  exit(): void {
    process.exit(0);
  }

  run() {
    runPlugin(this);
  }

  respond(response: PluginResponse) {
    respondWith(response);
  }
}

main();
