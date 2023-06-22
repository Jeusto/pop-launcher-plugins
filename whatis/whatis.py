#!/usr/bin/env python3
import json
import sys
import os
import subprocess


def run_whatis(arg):
    try:
        output = os.popen(f"whatis {arg}").read()
        output = subprocess.check_output(["whatis", arg], universal_newlines=True)

        return [
            (section_id, man_page.strip(), description.strip())
            for entry in output.splitlines()
            if (man_page := entry.split("(", 1)[0])
            and (section_id := entry.split("(", 1)[1].split(")", 1)[0])
            and (description := entry.split("-", 1)[1])
        ]
    except subprocess.CalledProcessError as e:
        print(e)
        return []


def open_man_page(arg):
    terminal, targ = detect_terminal()
    section = arg.section_number
    name = arg.name

    try:
        subprocess.Popen([terminal, targ, "man", section, name])

    except OSError as e:
        print(e)


def detect_terminal():
    symlink = "/usr/bin/x-terminal-emulator"
    fallback = "/usr/bin/gnome-terminal"

    if os.path.islink(symlink):
        terminal = os.readlink(symlink)
        return terminal, "-e"

    return fallback, "--"


# The entry for the search results
class ResultEntry:
    def __init__(self, section_number, name, description):
        self.section_number = section_number
        self.name = name
        self.description = description


# The plugin app
class App(object):
    def __init__(self):
        self.entries: list[ResultEntry] = []

    def append_entries(self, entries):
        for entry in entries:
            self.entries.append(ResultEntry(entry[0], entry[1], entry[2]))

    # When the user activates an entry
    def activate(self, index):
        if not self.entries:
            return

        open_man_page(self.entries[index])
        sys.stdout.write('"Close"\n')
        sys.stdout.flush()

    # When the user types something in the search bar
    def search(self, query):
        self.entries = []

        if query:
            query = query.split(" ", 1)[1]
            whatis_results = run_whatis(query)
            self.append_entries(whatis_results)

        for index, entry in enumerate(self.entries):
            sys.stdout.write(
                json.dumps(
                    {
                        "Append": {
                            "id": index,
                            "name": entry.name + " (" + entry.section_number + ")",
                            "description": entry.description,
                            "keywords": None,
                            "icon": None,
                            "exec": None,
                            "window": None,
                        }
                    }
                )
            )
            sys.stdout.write("\n")

        sys.stdout.write('"Finished"\n')
        sys.stdout.flush()


# Main execution
def main():
    app = App()

    for line in sys.stdin:
        try:
            request = json.loads(line)
            if "Search" in request:
                app.search(request["Search"])
            elif "Activate" in request:
                app.activate(request["Activate"])
        except json.decoder.JSONDecodeError:
            pass


if __name__ == "__main__":
    main()
