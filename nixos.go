package main

import "path/filepath"

func mkNixOSImage(configFname, cacheDir, vmID string) (string, error) {
	outputFname := filepath.Join(cacheDir, "nixos", vmID)
	err := run("nix-shell", "-p", "nixos-generators", "--run", "nixos-generate -f openstack -o "+outputFname+" -c "+configFname)
	if err != nil {
		return "", err
	}

	return outputFname, nil
}
