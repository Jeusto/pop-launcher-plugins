import path from "path";
import os from "os";
import fs from "fs";

export function log(obj: any) {
  const timestamp = new Date().toISOString();
  const logFile = path.join(os.homedir(), ".vscode-recent.log");
  const line = `${timestamp} ${JSON.stringify(obj)}\n`;

  fs.appendFileSync(logFile, line, { encoding: "utf-8", flag: "a" });
}

export function cleanSearchQuery(query: string, plugin_id: string): string {
  return query.trim().replace(new RegExp(`${plugin_id}(\\s|$)`), "");
}

export function getIcon(type: string) {
  switch (type) {
    case "file":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/file.svg";
    case "workspace":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/workspace.svg";
    case "folder":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/folder.svg";
    case "info":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/info.svg";
    case "error":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/error.svg";
    case "warning":
      return "/home/asaday/.local/share/pop-launcher/plugins/vscode-open-recent/assets/warning.svg";
    default:
      return "";
  }
}
