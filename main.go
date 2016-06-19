package main

import (
	"fmt"
	"os"
	"mmd"
)

func main() {

	bash := & mmd.Definition {
		Files: []string {
			"/usr/bin/bash",
			"/usr/bin/ls",
			"/usr/bin/cat",
			"/usr/bin/less",
			"/usr/share/terminfo/x/xterm",
		},
	}

	strace := & mmd.Definition {
		Files: []string { "/usr/bin/strace" },
	}

	mmd := mmd.New()

	mmd.AddDefinition( bash )
	mmd.AddDefinition( strace )

	folder := fmt.Sprintf( "%s%c%s-%d-tmp", os.TempDir(), os.PathSeparator, "mmd", os.Getpid() )

	fmt.Println( "folder = ", folder )

	if err := mmd.Export( folder ); err != nil {
		panic( err )
	}
}
