package analizador

import (
	"Proy1/Backend/comandos"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/rs/cors"
)

func Analizar(comandoEntero string) {

	//fmt.Println(comandoEntero)
	comandoEntero = strings.ToLower(comandoEntero)

	analComando := regexp.MustCompile("^[A-Za-z]+")
	comando := analComando.FindAllString(comandoEntero, 1)
	analBanderas := regexp.MustCompile("-[A-Za-z0-9]+(=(([a-zA-Z0-9/._]+)|(\"[^\n]*\")))?")
	banderas := analBanderas.FindAllString(comandoEntero, -1)

	if comando != nil {
		ejecutarComando(comando, banderas)
	}

}

func ejecutarComando(comando []string, banderas []string) {

	switch strings.ToLower(comando[0]) {

	case "execute":
		//ejecutar execute
		EjecExecute(banderas)
		//fmt.Println("ejecutar execute")
		break

	case "mkdisk":

		comandos.EjecutarMkdisk(banderas)
		break

	case "rmdisk":

		comandos.EjecutarRmdisk(banderas)
		break
	case "fdisk":
		comandos.EjecutarFdisk(banderas)
		break
	case "mount":
		comandos.EjecutarMount(banderas)
		break
	case "lmount":
		comandos.EjecutarLMount()
		break
	case "mkfs":
		comandos.EjecutarMkfs(banderas)
		break

	case "login":
		comandos.EjecutarLogin(banderas)
		break

	case "logout":
		comandos.EjecutarLogout()
		break
	case "cat":
		comandos.EjecutarCat(banderas)
		break
	case "mkdir":
		comandos.EjecutarMkdir(banderas)
		break
	case "mkfile":
		comandos.EjecutarMkfile(banderas)
		break

	case "rep":
		//fmt.Println("ejecutar rep")
		EjecRep(banderas)
		break

	case "exit":
		fmt.Println("cerrando aplicacion")
		os.Exit(0)

	}

}

func EjecRep(banderas []string) {

	name := ""
	path := ""
	id := ""
	path_file_ls := ""
	fmt.Println(path_file_ls)
	//ruta := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {

			name = dupla[1]
			if strings.Contains(name, "\"") {
				name = name[1 : len(name)-1]
			}

		} else if dupla[0] == "-path" {
			path = dupla[1]
			if strings.Contains(path, "\"") {
				path = path[1 : len(path)-1]
			}

		} else if dupla[0] == "-id" {
			id = dupla[1]
		} else {
			fmt.Println("Parametro invalido")
			comandos.Salida_comando += "Parametro invalido\n"
		}
	}

	switch name {
	case "mbr":
		//reporte mbr
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepMbr(index, path)
		break
	case "disk":
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepDisk(index, path)
		break

	case "inode":
		//reporte inodo
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepInodos(index, path)
		break

	case "block":
		//reporte block
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepBloques(index, path)
		break
	case "bm_inode":
		//reporte bitmap inodo
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepBmInodos(index, path)
		break
	case "bm_bloc":
		//reporte bitmap block
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepBmBloques(index, path)

		break
	case "sb":
		//reporte sb
		index := comandos.VerificarParticionMontada(id)
		if index == -1 {
			fmt.Println("Id no encontrada")
			comandos.Salida_comando += "Id no encontrada\n"
			return
		}
		comandos.RepSb(index, path)
		break
	case "file":
		//reporte file

		break
	case "ls":
		//reporte ls
		break

	default:
		fmt.Println("nombre no valido")
		comandos.Salida_comando += "nombre no valido\n"
		return
	}

}

func EjecExecute(banderas []string) {

	dupla := strings.Split(banderas[0], "=")
	//fmt.Println(banderas)
	if dupla[0] == "-path" {
		//fmt.Println(dupla[1])
		archivo, err := os.Open(dupla[1])

		if err != nil {
			fmt.Println("Error al abrir el archivo: ", err)
			comandos.Salida_comando += "Error al abrir el archivo \n"
			return
		}
		defer archivo.Close()

		scanner := bufio.NewScanner(archivo)

		for scanner.Scan() {
			linea := scanner.Text()
			//fmt.Println(linea)
			if len(linea) == 0 {
				continue
			}
			if linea[0] == '#' {
				//fmt.Println(linea)
				continue
			}
			Analizar(linea)
			//fmt.Println(linea)
		}
	}
}

type Cmd_API struct {
	Cmd string `json:"cmd"`
}

func InitApp() {
	fmt.Println("API Backend Proyecto 1 Archivos")
	comandos.Salida_comando = ""

	mux := http.NewServeMux()

	/* Ejemplo 7 */
	// Endpoint tipo POST
	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		// Configuracion de la cabecera
		w.Header().Set("Content-Type", "application/json")
		var Content Cmd_API
		body, _ := io.ReadAll(r.Body)
		// Arreglo  de bytes a Json
		//fmt.Println(string(body))
		json.Unmarshal(body, &Content)
		// Ejecuta el comando

		split_cmd(Content.Cmd)
		// Respuesta del servidor

		comandosUsados := comandos.Salida_comando
		//fmt.Println("----------------")
		//fmt.Println(comandosUsados)
		response := map[string]string{
			"result": comandosUsados,
		}

		// Convertir el mapa a JSON y enviarlo como respuesta
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error al crear respuesta JSON", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		w.Write(jsonResponse)

		// Limpio la salida de comandos
		comandos.Salida_comando = ""
	})

	fmt.Println("Servidor en el puerto 5000")
	// Configuracion de cors
	// handler := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:5173/"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders:   []string{"Content-Type", "application/json"},
	// 	AllowCredentials: true,
	// }).Handler(mux)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func split_cmd(cont string) {
	//cont = cont[1 : len(cont)-1]

	arr_com := strings.Split(cont, "\n")
	//fmt.Println(len(arr_com))
	//fmt.Println(arr_com)
	for _, linea := range arr_com {
		if len(linea) == 0 {
			continue
		}
		if linea[0] == '#' {
			//fmt.Println(linea)
			comandos.Salida_comando += linea + "\n"
			continue
		}
		//fmt.Println(linea)
		Analizar(linea)
		//fmt.Println(linea)
	}

}
