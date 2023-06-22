import readline from "readline";
import { log } from "./utils";

export interface PluginExt {
  name(): string;
  run(): void;
  search(query: string): void;
  activate(id: Index): void;
  activate_context?(id: Index, context: Index): void;
  complete?(id: Index): void;
  context?(id: Index): void;
  exit?(): void;
  interrupt?(): void;
  quit?(id: Index): void;
  respond_with?(response: PluginResponse): void;
  init_logging?(): void;
}

export class PopPlugin implements PluginExt {
  name(): string {
    return "pop";
  }

  respond_with(response: PluginResponse) {
    process.stdout.write(`${JSON.stringify(response)}\n`);
  }

  run() {
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
      terminal: false,
    });

    rl.on("line", (line) => {
      try {
        const request: Request = JSON.parse(line) as Request;
        log(request);

        switch (request) {
          case "Exit":
            this.exit?.();
          case "Interrupt":
            this.interrupt?.();

          default:
            if ((request as ActivateEvent).Activate !== undefined) {
              this.activate((request as ActivateEvent).Activate);
            }
            if (
              (request as ActivateContextEvent).ActivateContext !== undefined
            ) {
              const { id, context } = (request as ActivateContextEvent)
                .ActivateContext;
              this.activate_context?.(id, context);
            }
            if ((request as CompleteEvent).Complete !== undefined) {
              this.complete?.((request as CompleteEvent).Complete);
            }
            if ((request as ContextEvent).Context !== undefined) {
              this.context?.((request as ContextEvent).Context);
            }
            if ((request as QuitEvent).Quit !== undefined) {
              this.quit?.((request as QuitEvent).Quit);
            }
            if ((request as SearchEvent).Search !== undefined) {
              this.search((request as SearchEvent).Search);
            }
        }
      } catch (err) {
        log(err);
      }
    });
  }

  search(query: string) {}
  activate(id: Index) {}
  complete(id: Index) {}

  context(id: Index) {}
  activate_context?(id: Index, context: Index): {};

  exit() {}
  interrupt() {}
  quit(id: Index) {}

  init_logging() {}
}

// PluginResponse Types
type Index = number;
type GpuPreference = "Default" | "NonDefault";
type IconSource = { Name: string } | { Mime: string };

type ContextOption = {
  id: number;
  name: string;
};

export type PluginSearchResult = {
  id: number;
  name: string;
  description: string;
  keywords: Array<string> | null;
  icon: IconSource | null;
  exec: string | null;
  window: [number, number] | null;
};

export type PluginResponse =
  | { Append: PluginSearchResult }
  | { Context: { id: Index; options: Array<ContextOption> } }
  | { DesktopEntry: { path: string; gpu_preference: GpuPreference } }
  | { Fill: string }
  | "Clear"
  | "Close"
  | "Finished";

// Event types
type ExitEvent = "Exit";
type InterruptEvent = "Interrupt";

export interface ActivateEvent {
  Activate: number;
}
interface CompleteEvent {
  Complete: number;
}
interface ContextEvent {
  Context: number;
}
export interface QuitEvent {
  Quit: number;
}
export interface SearchEvent {
  Search: string;
}

export interface ActivateContextEvent {
  ActivateContext: {
    id: number;
    context: number;
  };
}

export type Request =
  | ActivateEvent
  | ActivateContextEvent
  | CompleteEvent
  | ContextEvent
  | ExitEvent
  | InterruptEvent
  | QuitEvent
  | SearchEvent;
