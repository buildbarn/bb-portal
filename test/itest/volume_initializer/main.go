package main

import (
	"fmt"
	"os"
)

type directory struct {
	path string
	mode os.FileMode
	uid  int
	gid  int
}

var volumeDirectories = []directory{
	{path: "/volumes/storage-ac/persistent_state", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/storage-cas/persistent_state", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/storage-fsac/persistent_state", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/jaeger/data", mode: 0o700, uid: 10001, gid: 10001},
	{path: "/volumes/postgres", mode: 0o700, uid: 65532, gid: 65532},
	{path: "/volumes/postgres/data", mode: 0o700, uid: 65532, gid: 65532},
	{path: "/volumes/worker-hardlinking/build", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/worker-hardlinking/cache", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/worker-fuse", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/worker-fuse/build", mode: 0o777, uid: -1, gid: -1},
	{path: "/volumes/worker-fuse/cas/persistent_state", mode: 0o777, uid: -1, gid: -1},
}

func main() {
	if err := initializeDirectories(volumeDirectories); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initializeDirectories(directories []directory) error {
	for _, directory := range directories {
		if err := os.MkdirAll(directory.path, directory.mode); err != nil {
			return fmt.Errorf("create %s: %w", directory.path, err)
		}
		if err := os.Chmod(directory.path, directory.mode); err != nil {
			return fmt.Errorf("set permissions on %s: %w", directory.path, err)
		}
		if directory.uid >= 0 || directory.gid >= 0 {
			if err := os.Chown(directory.path, directory.uid, directory.gid); err != nil {
				return fmt.Errorf("set ownership on %s: %w", directory.path, err)
			}
		}
	}
	return nil
}
