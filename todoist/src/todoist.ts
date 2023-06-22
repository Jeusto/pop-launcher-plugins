import { TodoistApi } from "@doist/todoist-api-typescript";
import fs from "fs";
import path from "path";

const API_URL = "https://api.todoist.com/rest/v2";
const CONFIG_FILE = "config.json";
const PLUGIN_PATH = "/.local/share/pop-launcher/plugins/todoist/";

export class Todoist {
  api: TodoistApi;

  private constructor() {
    this.api = new TodoistApi(retrieveApiKey());
  }

  static async init() {
    return new Todoist();
  }

  async add_task(
    name: string,
    priority: number = 1,
    due_date: string = "today"
  ) {
    await this.api.addTask({
      content: name,
      priority: 1,
      dueString: due_date,
    });
  }
}

function retrieveApiKey(): string {
  const home = process.env.HOME;
  if (!home) throw new Error("HOME environment variable not set");

  const config_path = path.join(home, PLUGIN_PATH, CONFIG_FILE);
  const config = JSON.parse(fs.readFileSync(config_path, "utf8"));

  return config.TODOIST_API_KEY;
}
