package main

import "path/filepath"

func mkNixOSImage(configFname, cacheDir, vmID string) (string, error) {
	outputFname := filepath.Join(cacheDir, "nixos", vmID+".qcow2")
	err := run("nix-shell", "-p", "nixos-generators", "--run", "nixos-generate -f qcow -o "+outputFname)
	if err != nil {
		return "", err
	}

	return outputFname, nil
}
