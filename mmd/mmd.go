package mmd

import (
	"os"
	"path/filepath"
	"io"
	"ldd"
	"fmt"
)

type Mmd struct {
	Definitions []Definition
}

func New() ( *Mmd ) {
	return & Mmd {
		Definitions: make([]Definition,0),
	}
}

func ( mmd * Mmd ) AddDefinition( definition * Definition ) {
	mmd.Definitions = append( mmd.Definitions, *definition )
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

func ( mmd * Mmd ) Export( directory string ) error {

	files   := make(map[string]bool)
	tocheck := make(map[string]bool)
	checked := make(map[string]bool)

	for _, definition := range mmd.Definitions {
		for _, file := range definition.Files {
			files[ file ]   = true
			tocheck[ file ] = true
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
