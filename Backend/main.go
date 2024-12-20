package main

import (
	"Proy1/Backend/analizador"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//:D
	//var banderas []string
	//banderas = append(banderas, "-path=/home/archivos/user/docs/usac/")

	//analizador.Analizar("execute -path=./pruebas.smia")

	// var banderas []string
	// banderas = append(banderas, "-user=root")
	// banderas = append(banderas, "-pass=123")
	// banderas = append(banderas, "-id=441a")
	// fmt.Println(banderas)
	// comandos.EjecutarLogin(banderas)
	//comandos.EjecLs(banderas)
	analizador.InitApp()
	// for {
	// 	leerComando()
	// }

	// dir, err := filepath.Abs("./carpeta/disco2.dsk")

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = os.MkdirAll(filepath.Dir(dir), 0777)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// archivo, err := os.Create(dir)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// archivo.Close()

}

func leerComando() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese un comando: ")
	comando, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al ingresar el comando: ", err)
		return
	}

	comando = strings.TrimSpace(comando)
	if len(comando) == 0 {
		return
	}
	if comando[0] != '#' {
		analizador.Analizar(comando)
	}

}
