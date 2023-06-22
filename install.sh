#!/bin/bash
LAUNCHER_PATH=$HOME/.local/share/pop-launcher
PLUGINS_DIR=$LAUNCHER_PATH/plugins
TEMP_DIR=$(mktemp -d)

install-all() {
     # Prep
     mkdir -p $PLUGINS_DIR
     git clone git@github.com:Jeusto/pop-launcher-plugins.git $TEMP_DIR
     mv $TEMP_DIR/* $PLUGINS_DIR

     # Install
     cd $PLUGINS_DIR
     for dir in */; do
         cd "$dir"
         if [ -f "package.json" ]; then
             npm install && npm run build
         elif [ -f "go.mod" ]; then
             go build
         elif [ -f "Cargo.toml" ]; then
             cargo build --release
         fi
         cd ..
     done

     # Cleanup
     rm -rf $TEMP_DIR
}

install-all