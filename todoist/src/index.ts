#!/usr/bin/env node
import { PopPlugin } from "pop-launcher-toolkit";
import { cleanSearchQuery } from "./utils";
import { Todoist } from "./todoist";

// Main execution
async function main() {
  const plugin = await TodoistPlugin.init();
  plugin.run();
}

const pages = [
  {
    id: 1,
    name: "Quick add task",
    icon: "",
  },
  // {
  //   id: 2,
  //   name: "Add task with options",
  //   icon: "",
  // },
  // {
  //   id: 3,
  //   name: "Search task",
  //   icon: "",
  // },
  // {
  //   id: 4,
  //   name: "Get random task to do",
  //   icon: "",
  // },
];

class TodoistPlugin extends PopPlugin {
  private query: string;
  private current_page = 0;
  private todoist: Todoist;

  private constructor(todoist: Todoist) {
    super();
    this.query = "";
    this.todoist = todoist;
  }

  static async init(): Promise<TodoistPlugin> {
    return new TodoistPlugin(await Todoist.init());
  }

  name(): string {
    return "Todoist";
  }

  show_options() {
    this.respond_with("Clear");
    pages.forEach((option) => {
      this.respond_with({
        Append: {
          id: option.id,
          name: option.name,
          description: "",
          keywords: null,
          icon: null,
          exec: null,
          window: null,
        },
      });
    });
    this.respond_with("Finished");
  }

  show_add_task() {
    this.respond_with("Clear");
    this.respond_with({
      Append: {
        id: 1,
        name: "Quick add task",
        description: "Task name: " + this.query,
        keywords: null,
        icon: null,
        exec: null,
        window: null,
      },
    });
    this.respond_with("Finished");
  }

  search(query: string) {
    this.query = cleanSearchQuery(query, "todo");

    switch (this.current_page) {
      case 0:
        this.show_options();
        break;
      case 1:
        this.show_add_task();
      default:
        break;
    }
  }

  async activate(index: number) {
    switch (this.current_page) {
      case 0:
        this.current_page = index;
        this.reset_searchbar();
        this.show_add_task();
        break;
      case 1:
        await this.todoist.add_task(this.query);
        this.show_notification(
          "Task added successfully",
          "Task: " + this.query,
          "/home/asaday/.local/share/pop-launcher/plugins/todoist/assets/plugin.svg",
          1000
        );
        this.exit();
        break;
      default:
        break;
    }
  }

  reset_searchbar(): void {
    this.respond_with({ Fill: "todo " });
  }

  exit(): void {
    this.respond_with("Close");
    this.respond_with("Finished");
    process.exit(0);
  }
}

main();
