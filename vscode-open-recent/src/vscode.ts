import { join } from "path";
import { open } from "sqlite";
import { exec } from "child_process";
import Fuse from "fuse.js";
import sqlite3 from "sqlite3";

const PATH_DIRS = ["/usr/bin", "/bin", "/snap/bin"];
const EXECUTABLE = "code";
const CONFIG_PATH = ".config/Code/User/globalStorage/state.vscdb";
const FUSE_OPTIONS = {
  keys: ["type", "name", "path"],
};

type Db = DbRow[];
type DbRecentlyOpened = {
  entries: DbEntry[];
};

interface DbRow {
  key: string;
  value: any;
}

interface DbEntry {
  workspace?: DbWorkspace;
  folderUri?: string;
  fileUri?: string;
  label?: string;
  remoteAuthority?: string;
}

interface DbWorkspace {
  id: string;
  configPath: string;
}

export interface Entry {
  type: "file" | "workspace" | "folder";
  name: string;
  path: string;
  full_path: string;
}

export type RecentlyOpened = Entry[];

export class Vscode {
  recentlyOpened: RecentlyOpened;
  private vscodeIsInstalled: boolean;
  private configPath: string;

  private constructor() {
    this.recentlyOpened = [];
    this.vscodeIsInstalled = false;
    this.configPath = CONFIG_PATH;
  }

  static async init() {
    const vscode = new Vscode();
    vscode.recentlyOpened = await vscode.fetchRecentlyOpened();
    return vscode;
  }

  private async fetchRecentlyOpened(): Promise<RecentlyOpened> {
    try {
      const db = await open({
        filename: join(process.env.HOME || "", this.configPath),
        driver: sqlite3.cached.Database,
      });

      const rows = await db.all<Db>(
        "SELECT * FROM ItemTable WHERE key = 'history.recentlyOpenedPathsList'"
      );

      if (rows.length === 0) {
        throw new Error("Could not find recently opened paths");
      }

      const row = rows[0] as DbRow;
      const rowValue = JSON.parse(row.value) as DbRecentlyOpened;
      const recent = rowValue.entries;

      return recent.map((entry) => {
        return this.convertToEntry(entry);
      });
    } catch (error) {
      throw error;
    }
  }

  search(query: string) {
    const fuse = new Fuse(this.recentlyOpened, FUSE_OPTIONS);
    const results = fuse.search(query).slice(0, 8);

    return results.map((result) => {
      return result.item;
    }, []);
  }

  async open(entry: Entry) {
    exec(`${EXECUTABLE} ${entry.path}`, (err) => {
      if (err) {
        throw err;
      }
    });
  }

  private convertToEntry(entry: DbEntry): Entry {
    if (entry.workspace) {
      return {
        type: "workspace",
        name: getFileNameFromPath(entry.workspace.configPath),
        path: getPathWithoutFileName(entry.workspace.configPath),
        full_path: entry.workspace.configPath.replace("file://", ""),
      };
    } else if (entry.folderUri) {
      return {
        type: "folder",
        name: getFileNameFromPath(entry.folderUri),
        path: getPathWithoutFileName(entry.folderUri),
        full_path: entry.folderUri.replace("file://", ""),
      };
    } else if (entry.fileUri) {
      return {
        type: "file",
        name: getFileNameFromPath(entry.fileUri),
        path: getPathWithoutFileName(entry.fileUri),
        full_path: entry.fileUri.replace("file://", ""),
      };
    } else {
      throw new Error("Unknown entry type");
    }
  }
}

function getPathWithoutFileName(uri: string): string {
  return uri
    .replace("file://", "")
    .split("/")
    .slice(0, -1)
    .join("/")
    .replace(".code-workspace", "")
    .replace("/home/" + process.env.USER, "~");
}

function getFileNameFromPath(uri: string): string {
  return (
    uri
      .replace("file://", "")
      .split("/")
      .pop()
      ?.replace(".code-workspace", "") || ""
  );
}
