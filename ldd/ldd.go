package ldd

import (
	"os/exec"
	"log"
	"text/scanner"
	"debug/elf"
)

func isDynamicExecutableOrLibrary( filename string ) bool {
	elf_file, err := elf.Open( filename )

	if err != nil {
		return false
	}

	elf_file.Close()

	return true
}

//func List_dynamic_dependencies_recursive( filename string ) []string {
//	files := make(map[string]bool)
//
//	for _, dep := List_dynamic_dependencies( filename ) {
//		files[ dep ] = true
//
//		for _
//	}
//}

func List_dynamic_dependencies( filename string ) []string {
	if ! isDynamicExecutableOrLibrary( filename ) {
		return nil
	}

	cmd := exec.Command( "ldd", filename )
	stdout,err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal( err )
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var s scanner.Scanner
	s.Init( stdout )

	//var p *Parser
	var items chan item
	
	_, items = NewParser( s )

	var list []string = make([]string,0)

	for item := range items {
		if item.typ == itemError {
			panic( item.val )
		}
		list = append( list, item.val )
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return list
}
