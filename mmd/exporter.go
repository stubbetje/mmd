package mmd

import (
	"os"
	"path/filepath"
	"io"
	"ldd"
	"fmt"
)

type Exporter struct {
	Definitions []Definition
}

func NewExporter() ( *Exporter ) {
	return & Exporter {
		Definitions: make([]Definition,0),
	}
}

func ( exporter * Exporter ) AddDefinition( definition * Definition ) {
	exporter.Definitions = append( exporter.Definitions, *definition )
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func ( exporter * Exporter ) Export( directory string ) error {

	files   := make(map[string]bool)
	tocheck := make(map[string]bool)
	checked := make(map[string]bool)

	for _, definition := range exporter.Definitions {
		for _, file := range definition.Content {
			files[ file.File ]   = true
			tocheck[ file.File ] = true
		}
	}

	for {
		var finished = true

		for file := range tocheck {

			for _, dep := range ldd.List_dynamic_dependencies( file ) {
				files[ dep ] = true

				if ! checked[ dep ] {
					finished       = false
					tocheck[ dep ] = true
				}
			}

			checked [ file ] = true
			delete( tocheck, file )
		}

		if finished {
			break
		}
	}

	for filename := range files {
		// skip linux-vdso.so.1
		if filename == "linux-vdso.so.1" {
			continue
		}

		fmt.Println( "COPY ", filename )

		destination := fmt.Sprintf( "%s%s", directory, filename )
		dir         := filepath.Dir( destination )

		fmt.Println( "  TO ", destination )

		if err := os.MkdirAll( dir, 0700 ); err != nil {
			//panic( fmt.Sprintf( "error creating directory %s: %q", dir, err ) )
			return err
		}

		if err := copyFileContents( filename, destination ); err != nil {
			return err
		}

		info, err := os.Stat( filename )
		if err != nil {
			return err
		}

		if err := os.Chmod( destination, info.Mode() ); err != nil {
			return err
		}
	}

	dockerFile, err := os.Create( filepath.Join( directory, "Dockerfile" ) )
	if err != nil {
		return err
	}

	contents := "FROM scratch\nADD . /"

	_, err = dockerFile.WriteString( contents )

	if err != nil {
		return err
	}

	if err := dockerFile.Close(); err != nil {
		return err
	}

	return nil
}
