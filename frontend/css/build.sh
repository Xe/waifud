#!/usr/bin/env nix-shell
#! nix-shell -p nodePackages.clean-css-cli -i bash

cleancss -o ./xess.css ./src/xess.css ./src/admin.css
