#!/usr/bin/env node
import { PopPlugin } from "./pop";
import { cleanSearchQuery, getIcon, log } from "./utils";
import { RecentlyOpened, Vscode } from "./vscode";

// Main execution
async function main() {
  const plugin = await VscodePlugin.init();
  plugin.run();
}

class VscodePlugin extends PopPlugin {
  private vscode: Vscode | null;
  private query: string;
  private results: RecentlyOpened;

  private constructor(vscode: Vscode) {
    super();
    this.query = "";
    this.vscode = vscode;
    this.results = vscode.recentlyOpened.slice(0, 8) || [];
  }

  static async init(): Promise<VscodePlugin> {
    return new VscodePlugin(await Vscode.init());
  }

  name(): string {
    return "vscode-open-recent";
  }

  search(query: string) {
    log("Search query: " + query);
    this.query = cleanSearchQuery(query, "vs");

    if (!this.query) {
      this.respond_with("Clear");
      this.results.forEach((entry, idx) => {
        this.respond_with({
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
      log("Finished");
      this.respond_with("Finished");
      return;
    }

    this.results = this.vscode?.search(this.query) || [];

    if (this.results.length === 0) {
      this.respond_with("Clear");
      this.respond_with({
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
      this.respond_with("Finished");
      return;
    }

    this.respond_with("Clear");
    this.results.forEach((entry, idx) => {
      this.respond_with({
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
    this.respond_with("Finished");
  }

  activate(index: number) {
    const entry = this.results[index];
    this.vscode?.open(entry);
    this.respond_with("Close");
  }

  exit(): void {
    process.exit(0);
  }
}

main();
