package comandos

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Mount struct {
	Id        string
	Path      string
	Name      string
	Part_type [1]byte
	Start     int32
	Size      int32
	PartNum   int32
}

var PathDiscos []string
var ParticionesMontadas []Mount
var letra_disco_ascii int = 96
var UId int = -1
var GId int = -1
var ActualIdMount string
var Salida_comando string

func VerificarParticionMontada(id string) int {
	//fmt.Println(id)

	for i := 0; i < len(ParticionesMontadas); i++ {
		//fmt.Println(particionesMontadas[i].Id)
		if ParticionesMontadas[i].Id == id {
			return i
		}
	}
	return -1

}

func GetLsDiscos() []LsFile {

	var datos []LsFile

	for _, disco := range PathDiscos {
		datos = append(datos, LsFile{disco, 0})
	}

	return datos

}

func GetLsParticiones(path string) []LsFile {
	var datos []LsFile

	archivo, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return nil
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		Salida_comando += "Error al leer el MBR del disco\n"
		return nil
	}

	for _, particion := range disk.Mbr_partitions {

		if particion.Part_s != int32(0) {
			datos = append(datos, LsFile{strings.TrimRight(string(particion.Part_name[:]), string(rune(0))), 1})
			break

		}
	}
	archivo.Close()
	return datos
}

func EliminarPathDisco(path string) {
	var newDisks []string

	for _, disk := range PathDiscos {
		if disk != path {
			newDisks = append(newDisks, disk)
		}
	}

	// Reemplazar el slice original con el nuevo
	PathDiscos = newDisks

}

func VerificarNombreMontado(nombre string) int {

	for i := 0; i < len(ParticionesMontadas); i++ {
		//fmt.Println(particionesMontadas[i].Id)
		if ParticionesMontadas[i].Name == nombre {
			return i
		}
	}
	return -1

}

func RestartValues() {
	letra_disco_ascii = 96
	UId = -1
	GId = -1
	ActualIdMount = ""
	Salida_comando = ""
	ParticionesMontadas = ParticionesMontadas[:0]
}
