package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"log"
)

func failOnError(err error, msg string) {
	if (err != nil) {
	  log.Fatalf("%s: %s", msg, err)
	}
}

func localizacionChunks(nombreL string) string{

	f, err := os.Open("testlog.txt")
	failOnError(err, "Error en abrir log")
	defer f.Close()

	// hace Splits por cada linea por defecto.
	scanner := bufio.NewScanner(f)

	//var cantchunks int
	var listachunks string
	var t string
	var info []string
	var init int
	var n int
	var mark bool
	mark = false

	for scanner.Scan() {
		
		t = scanner.Text()
		
		if mark{
			
			info = strings.Fields(t)
			listachunks += info[1] + " "
			init++

			if (init == n) {
				return listachunks
			}
		}
		
		if strings.Contains(t, nombreL) {
			words := strings.Fields(t) //es como split por blankspace
			n, _ = strconv.Atoi(words[1])
			//cantchunks = n
			init = 0
			mark = true
			continue
		}
	}

	if err := scanner.Err(); err != nil {
    	log.Fatalf("No se pudo localizar chunks correctamente: %v", err)
	}
	
	return ""
}

func main(){
	
	fmt.Println(localizacionChunks("Estudio_en_Rosa_2"))
}
