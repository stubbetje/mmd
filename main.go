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

//	bash,err := mmd.LoadFromYamlFile( "definitions/bash.mmd.yaml" )
//	if( err != nil ) {
//		panic( err )
//	}
//
//	strace,err := mmd.LoadFromYamlFile( "definitions/strace.mmd.yaml" )
//	if( err != nil ) {
//		panic( err )
//	}
	fmt.Printf( "%v\n", flag.Args() )

	exporter := mmd.NewExporter()

	for _, packageFile := range flag.Args() {
		def, err := mmd.LoadFromYamlFile( packageFile )
		if( err != nil ) {
			panic( err )
		}
		exporter.AddDefinition( def )
	}

//	mmd.AddDefinition( bash )
//	mmd.AddDefinition( strace )

	fmt.Println( "folder = ", *outputDir )

	fmt.Printf( "%v\n", exporter )

	if err := exporter.Export( *outputDir ); err != nil {
		panic( err )
	}
}
