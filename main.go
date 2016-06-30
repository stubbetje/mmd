package main

import (
	"fmt"
	"flag"
	"os"
	"mmd"
)

func main() {

	tempFolder := fmt.Sprintf( "%s%c%s-%d-tmp", os.TempDir(), os.PathSeparator, "mmd", os.Getpid() )
	outputDir  := flag.String( "output-dir", tempFolder, "Output directory" )

	flag.Parse()

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


	fmt.Println( "folder = ", *outputDir )

	if err := mmd.Export( *outputDir ); err != nil {
		panic( err )
	}
}
