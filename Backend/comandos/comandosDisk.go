package comandos

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int32
	MBR_dsk_fit        [1]byte
	Mbr_partitions     [4]partition
}

type partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_s           int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          [4]byte
}

type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_s     int32
	Part_next  int32
	Part_name  [16]byte
}

type superBloque struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [19]byte
	S_umtime            [19]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_firts_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type Inodo struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

type b_content struct {
	B_name  [12]byte
	B_inodo int32
}

type bloqueCarpeta struct {
	B_content [4]b_content
}

type bloqueArchivos struct {
	B_content [64]byte
}

type bloqueApuntadores struct {
	B_pointers [16]int32
}

func EjecutarMkdisk(banderas []string) {

	//Leer Banderas
	path := "./disco1.dsk"
	unit := "m"
	fit := "f"
	size := -1

	//fmt.Println(fit)
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-size" {

			size, _ = strconv.Atoi(dupla[1])

		} else if dupla[0] == "-unit" {
			unit = dupla[1]

		} else if dupla[0] == "-path" {
			path = dupla[1]
			if strings.Contains(path, "\"") {
				path = path[1 : len(path)-1]
			}

		} else if dupla[0] == "-fit" {
			if dupla[1] == "bf" {
				fit = "b"
			} else if dupla[1] == "ff" {
				fit = "f"
			} else if dupla[1] == "wf" {
				fit = "w"
			}
		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
			return
		}
	}

	//Crear MBR
	var newMbr MBR

	if size < 0 {
		fmt.Println("Valor size invalido")
		Salida_comando += "Valor size invalido\n"
		return
	}

	if unit == "k" {
		newMbr.Mbr_tamano = int32(size) * 1024
	} else if unit == "m" {
		newMbr.Mbr_tamano = int32(size) * 1024 * 1024
	} else {
		fmt.Println("El valor del parametro -unit no es valido")
		Salida_comando += "El valor del parametro -unit no es valido\n"
		return
	}

	fechaActual := time.Now()
	fecha := fechaActual.Format("2006-01-02 15:04:05")
	copy(newMbr.Mbr_fecha_creacion[:], fecha)
	randomNum := rand.Intn(99) + 1
	newMbr.Mbr_dsk_signature = int32(randomNum)

	copy(newMbr.MBR_dsk_fit[:], fit)

	newMbr.Mbr_partitions[0].Part_status = [1]byte{'0'}
	newMbr.Mbr_partitions[1].Part_status = [1]byte{'0'}
	newMbr.Mbr_partitions[2].Part_status = [1]byte{'0'}
	newMbr.Mbr_partitions[3].Part_status = [1]byte{'0'}

	newMbr.Mbr_partitions[0].Part_type = [1]byte{'0'}
	newMbr.Mbr_partitions[1].Part_type = [1]byte{'0'}
	newMbr.Mbr_partitions[2].Part_type = [1]byte{'0'}
	newMbr.Mbr_partitions[3].Part_type = [1]byte{'0'}

	newMbr.Mbr_partitions[0].Part_fit = [1]byte{'0'}
	newMbr.Mbr_partitions[1].Part_fit = [1]byte{'0'}
	newMbr.Mbr_partitions[2].Part_fit = [1]byte{'0'}
	newMbr.Mbr_partitions[3].Part_fit = [1]byte{'0'}

	newMbr.Mbr_partitions[0].Part_start = 0
	newMbr.Mbr_partitions[1].Part_start = 0
	newMbr.Mbr_partitions[2].Part_start = 0
	newMbr.Mbr_partitions[3].Part_start = 0

	newMbr.Mbr_partitions[0].Part_s = 0
	newMbr.Mbr_partitions[1].Part_s = 0
	newMbr.Mbr_partitions[2].Part_s = 0
	newMbr.Mbr_partitions[3].Part_s = 0

	newMbr.Mbr_partitions[0].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	newMbr.Mbr_partitions[1].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	newMbr.Mbr_partitions[2].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	newMbr.Mbr_partitions[3].Part_name = [16]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}

	//Crear Carpetas y archivo

	dir, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error path\n"
	}

	//Crea todos los directorios
	err = os.MkdirAll(filepath.Dir(dir), 0777)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error crear directorios\n"
	}

	//Crea archivo
	archivo, err := os.Create(dir)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error al crear archivo disco\n"
	}

	defer archivo.Close()

	//Buffer 1024 bytes
	bufer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		bufer.WriteByte(0)
	}

	var totalBytes int = 0

	//Escribe 0s en archivo
	for totalBytes < int(newMbr.Mbr_tamano) {
		c, err := archivo.Write(bufer.Bytes())
		if err != nil {
			fmt.Println("Error al escribir en el archivo: ", err)
			Salida_comando += "Error al escribir en el archivo \n"
			return
		}
		totalBytes += c
	}

	archivo.Seek(0, 0)
	err = binary.Write(archivo, binary.LittleEndian, &newMbr)
	if err != nil {
		fmt.Println("Error al escribir el MBR en el disco: ", err)
		Salida_comando += "Error al escribir el MBR en el disco:\n"
		return
	}
	fmt.Println("Disco creado con exito")
	Salida_comando += "Disco creado con exito\n"

	archivo.Close()

}

func EjecutarRmdisk(banderas []string) {

	path := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			path = dupla[1]
			if strings.Contains(path, "\"") {
				path = path[1 : len(path)-1]
			}

		}
	}

	if path == "" {

		//mensaje falta path
		return

	}

	err := os.Remove(path)

	if err != nil {
		//no se pudo eliminar el disco
		fmt.Println("No se pudo eliminar el disco")
		Salida_comando += "No se pudo eliminar el disco\n"
		return
	}

	//disco eliminado con exito
	fmt.Println("Disco eliminado con exito")
	Salida_comando += "Disco eliminado con exito\n"
}

func EjecutarFdisk(banderas []string) {
	//Leer Banderas
	path := "./disco1.dsk"
	unit := "k"
	fit := "w"
	size := -1
	name := ""
	typePartition := "p"

	//fmt.Println(fit)
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-size" {

			size, _ = strconv.Atoi(dupla[1])

		} else if dupla[0] == "-unit" {
			unit = dupla[1]

		} else if dupla[0] == "-path" {
			path = dupla[1]
			if strings.Contains(path, "\"") {
				path = path[1 : len(path)-1]
			}

		} else if dupla[0] == "-fit" {
			if dupla[1] == "bf" {
				fit = "b"
			} else if dupla[1] == "ff" {
				fit = "f"
			} else if dupla[1] == "wf" {
				fit = "w"
			}
		} else if dupla[0] == "-type" {
			typePartition = dupla[1]
		} else if dupla[0] == "-name" {
			name = dupla[1]
			if strings.Contains(path, "\"") {
				name = name[1 : len(name)-1]
			}
		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
			return
		}
	}

	//comprobar banderas
	//fmt.Println(unit)

	//fmt.Println(typePartition)
	if size <= 0 {
		//mal dato de size
		return
	}

	if unit == "k" {
		size = size * 1024
	} else if unit == "m" {
		size = size * 1024 * 1024
	}

	if name == "" {
		fmt.Println("nombre invalido")
		Salida_comando += "nombre invalido\n"
		return
	}

	//abrir archivo
	archivo, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		Salida_comando += "Error al leer el MBR del disco\n"
		return
	}

	numPart := -1
	despTemp := binary.Size(MBR{}) + 1

	//encuentra particion disponible y puntero a esta
	for i, particion := range disk.Mbr_partitions {

		if particion.Part_s == int32(0) {
			numPart = i
			break

		} else {
			despTemp += int(particion.Part_s) + 1
		}

	}

	//obtiene particion extendida
	var partExtend partition
	for _, part := range disk.Mbr_partitions {
		if part.Part_type == [1]byte{'e'} {
			if typePartition == "e" {
				fmt.Println("Ya existe una particion extendida")
				Salida_comando += "Ya existe una particion extendida\n"
				return
			}
			partExtend = part

		}
	}

	if typePartition == "l" {

		//Creacion particion logica

		if partExtend.Part_type != [1]byte{'e'} {
			fmt.Println("No existe particion extendida")
			Salida_comando += "no existe una particion extendida\n"
			return
		}

		var ebr EBR
		despTemp = int(partExtend.Part_start)
		var prev_ebr EBR

		for {

			archivo.Seek(int64(despTemp), 0)
			binary.Read(archivo, binary.LittleEndian, &ebr)
			if ebr.Part_s != 0 {
				if strings.Contains(string(ebr.Part_name[:]), name) {
					fmt.Println("Error: El nombre de la particion ya existe")
					Salida_comando += "Error: el nombre de la particion ya existe\n"
					return
				}
				despTemp += int(ebr.Part_s) + 1 + binary.Size(EBR{})

				prev_ebr = ebr

			} else {
				break

			}

		}

		if int32(despTemp)+int32(binary.Size(EBR{}))+int32(size)+1 > partExtend.Part_start+partExtend.Part_s {
			fmt.Println("Error: No hay espacio para crear la particion")
			Salida_comando += "Error: no hay espacio para crear la particion\n"
			return
		}

		//Crear Mbr

		ebr.Part_mount = [1]byte{'0'}
		ebr.Part_fit = [1]byte{fit[0]}
		ebr.Part_start = int32(despTemp) + 1 + int32(binary.Size(EBR{}))
		ebr.Part_s = int32(size)
		ebr.Part_next = int32(-1)
		copy(ebr.Part_name[:], name)

		if prev_ebr.Part_s != 0 {
			//modificar anterior ebr si hay
			prev_ebr.Part_next = int32(despTemp)
			despTemp_viejoEbr := int64(despTemp) - int64(prev_ebr.Part_s) - int64(binary.Size(EBR{})+1)
			archivo.Seek(despTemp_viejoEbr, 0)
			binary.Write(archivo, binary.LittleEndian, &prev_ebr)
		}

		//modificar ebr anterior
		//ebr.Part_next = int32(despTemp)
		//Escribir el nuevo  y viejo EBR
		//archivo.Seek(int64(despTemp)-int64(ebr.Part_s), 0)
		//binary.Write(archivo, binary.LittleEndian, &ebr)
		archivo.Seek(int64(despTemp), 0)
		binary.Write(archivo, binary.LittleEndian, &ebr)

		archivo.Close()
		fmt.Println("Particion logica creada con exito")
		Salida_comando += "Particion logica creada con exito\n"

		return

	} else {

		var nuevaPar partition

		nuevaPar.Part_status = [1]byte{'0'}

		if typePartition == "p" || typePartition == "e" {

			nuevaPar.Part_type = [1]byte{typePartition[0]}
		} else {
			fmt.Println("Tipo de particion no valida")
			Salida_comando += "Tipo de particion no valida\n"
			return
		}
		nuevaPar.Part_fit = [1]byte{fit[0]}

		if numPart < 0 {
			fmt.Println("No hay particiones disponibles")
			Salida_comando += "No hay particiones disponibles\n"
			return
		}

		nuevaPar.Part_start = int32(despTemp)

		nuevaPar.Part_s = int32(size)

		nuevaPar.Part_name = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		copy(nuevaPar.Part_name[:], name)

		if despTemp+int(nuevaPar.Part_s)+1 > int(disk.Mbr_tamano) {
			fmt.Println("tamano insuficiente para la particion")
			Salida_comando += "tamano insuficiente para la particion\n"
			return
		}

		nuevaPar.Part_correlative = 0

		disk.Mbr_partitions[numPart] = nuevaPar

		archivo.Seek(0, 0)
		binary.Write(archivo, binary.LittleEndian, &disk)
		archivo.Close()

		fmt.Println("particion creada con exito")
		Salida_comando += "particion creada con exito\n"

	}

}

func EjecutarMount(banderas []string) {

	path := ""
	name := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			path = dupla[1]
			if strings.Contains(path, "\"") {
				path = path[1 : len(path)-1]
			}

		} else if dupla[0] == "-name" {
			name = dupla[1]
			if strings.Contains(name, "\"") {
				name = name[1 : len(name)-1]
			}

		}
	}

	if path == "" {
		fmt.Println("Ruta no valida")
		Salida_comando += "Ruta no valida\n"
		return
	}

	if name == "" {
		fmt.Println("nombre no valido")
		Salida_comando += "nombre no valido\n"
		return
	}

	archivo, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		Salida_comando += "Error al leer el MBR del disco\n"
		return
	}
	numPart := -1
	correlativo := 0

	for i, part := range disk.Mbr_partitions {

		if part.Part_status == [1]byte{'1'} {
			correlativo++
			continue
		}

		if strings.Contains(string(part.Part_name[:]), name) {
			if string(part.Part_type[:]) == "p" {
				numPart = i
				continue
			} else {
				fmt.Println("Solo se pueden montar particiones Primarias")
				Salida_comando += "Solo se pueden montar particiones Primarias\n"
				return
			}

		}
	}

	if numPart == -1 {
		fmt.Println("particion no encontrada")
		Salida_comando += "particion no encontrada\n"
		return
	}
	if correlativo == 0 {
		letra_disco_ascii++
	}
	letra := string(rune(letra_disco_ascii))

	id := "44" + strconv.Itoa(correlativo+1) + letra

	var newMount Mount

	newMount.Id = id
	newMount.Path = path
	newMount.Name = name
	newMount.Part_type = [1]byte{'p'}
	newMount.Start = disk.Mbr_partitions[numPart].Part_start
	newMount.Size = disk.Mbr_partitions[numPart].Part_s
	newMount.PartNum = int32(numPart)

	disk.Mbr_partitions[numPart].Part_status = [1]byte{'1'}
	copy(disk.Mbr_partitions[numPart].Part_id[:], id)

	particionesMontadas = append(particionesMontadas, newMount)

	archivo.Seek(int64(0), 0)
	binary.Write(archivo, binary.LittleEndian, &disk)
	archivo.Close()

	fmt.Println("particion " + id + " montada")
	Salida_comando += "particion " + id + " montada\n"
	EjecutarLMount()

}

func EjecutarLMount() {

	if len(particionesMontadas) == 0 {
		fmt.Println("no hay particiones montadas")
		Salida_comando += "no hay particiones montadas\n"
		return
	}

	for _, mounts := range particionesMontadas {
		fmt.Println("--------------------------------------------------------")
		Salida_comando += "-----------------------------------------------------------\n"
		fmt.Print("Id: ")
		Salida_comando += "Id: "
		fmt.Println(mounts.Id)
		Salida_comando += mounts.Id + "\n"
		fmt.Print("Disco: ")
		Salida_comando += "Disco: "
		fmt.Println(mounts.Path)
		Salida_comando += mounts.Path + "\n"
		fmt.Print("Nombre Particion: ")
		Salida_comando += "Nombre Particion : "
		fmt.Println(mounts.Name)
		Salida_comando += mounts.Name + "\n"
		fmt.Print("Tipo: ")
		Salida_comando += "Tipo: "
		fmt.Println(string(mounts.Part_type[:]))
		Salida_comando += string(mounts.Part_type[:]) + "\n"
		fmt.Print("Inicio: ")
		Salida_comando += "Inicio: "
		fmt.Println(mounts.Start)
		Salida_comando += strconv.Itoa(int(mounts.Start)) + "\n"
		fmt.Print("Tamano: ")
		Salida_comando += "Tamano: "
		fmt.Println(mounts.Size)
		Salida_comando += strconv.Itoa(int(mounts.Size)) + "\n"
	}
	fmt.Println("--------------------------------------------------------")
	Salida_comando += "-----------------------------------------------------------\n"
}

func EjecutarLogin(banderas []string) {
	user := ""
	pass := ""
	id := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			user = dupla[1]

		} else if dupla[0] == "-pass" {
			pass = dupla[1]

		} else if dupla[0] == "-id" {
			id = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	if user == "" {
		fmt.Println("Ingrese el campo -user")
		Salida_comando += "Ingrese el campo -user\n"
		return
	}

	if pass == "" {
		fmt.Println("Ingrese el campo -pass")
		Salida_comando += "Ingrese el campo -pass\n"
		return
	}

	if id == "" {
		fmt.Println("Ingrese el campo -id")
		Salida_comando += "Ingrese el campo -id\n"
		return
	}

	if uId != -1 {
		fmt.Println("Ya hay una sesion iniciada")
		Salida_comando += "Ya hay una sesion iniciada\n"
		return

	}

	index := VerificarParticionMontada(id)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	//var inodoTemp Inodo
	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	if !(sblock.S_filesystem_type == 2) {
		println("El sistema de archivos no está formateado")
		Salida_comando += "El sistema de archivos no está formateado\n"
		return
	}

	result := leerArchivo("/users.txt", archivo, &sblock)

	lineas := strings.Split(result, "\n")
	lineas = lineas[:len(lineas)-1]
	nombreGrupo := ""
	for _, linea := range lineas {

		if linea[2] == 'u' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[3] == user && campos[4] == pass {
				uId, _ = strconv.Atoi(campos[0])
				nombreGrupo = campos[2]
				break

			}
		}

	}

	if nombreGrupo == "" {
		fmt.Println("Usuario no encontrado")
		return
	}

	for _, linea := range lineas {

		if linea[2] == 'g' {
			campos := strings.Split(linea, ",")
			if campos[2] == nombreGrupo {
				gId, _ = strconv.Atoi(campos[0])
				if linea[0] == '0' {
					fmt.Println("Grupo no encontrado")
					uId = -1
					gId = -1
					return
				}
				break

			}
		}

	}

	fmt.Println("Usuario Logueado con exito")
	Salida_comando += "Usuario Logueado con exito\n"
	fmt.Println("Usuario: " + strconv.Itoa(uId))
	Salida_comando += "Usuario: " + strconv.Itoa(uId) + "\n"
	fmt.Println("Grupo: " + strconv.Itoa(gId))
	Salida_comando += "Grupo: " + strconv.Itoa(gId) + "\n"
	actualIdMount = id

}

func EjecutarLogout() {
	if uId == -1 {
		fmt.Println("Aun no hay sesion iniciada")
		return
	}
	uId = -1
	gId = -1
	actualIdMount = ""
	fmt.Println("Sesion cerrada con Exito")
}

func EjecutarMkfs(banderas []string) {
	id := ""
	typeVar := "full"

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-id" {
			id = dupla[1]

		} else if dupla[0] == "-type" {
			typeVar = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	index := VerificarParticionMontada(id)
	if index == -1 {
		fmt.Println("Particion no montada")
		Salida_comando += "Particion no montada\n"
		return
	}
	//fmt.Println(typeVar)
	if typeVar != "full" {
		fmt.Println("Tipo de formateo invalido")
		Salida_comando += "Tipo de formateo invalido\n"
		return
	}

	var n int

	n = int(math.Floor(float64(int(particionesMontadas[index].Size)-int(binary.Size(superBloque{}))) / float64(4+int(binary.Size(Inodo{}))+3*int(binary.Size(bloqueArchivos{})))))

	crearExt2(index, n)
}

func EjecutarCat(banderas []string) {
	var filen []string

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if strings.Contains(dupla[0], "-file") {
			filen = append(filen, dupla[1])

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}
	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	//var inodoTemp Inodo
	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	if !(sblock.S_filesystem_type == 2) {
		println("El sistema de archivos no está formateado")
		Salida_comando += "El sistema de archivos no está formateado\n"
		return
	}

	for i, filePath := range filen {
		result := leerArchivo(filePath, archivo, &sblock)

		fmt.Println("---------------Archivo" + strconv.Itoa(i) + "---------------")
		Salida_comando += "---------------Archivo" + strconv.Itoa(i) + "---------------\n"
		fmt.Println(result)
		Salida_comando += result
	}

	archivo.Close()

}

func EjecutarMkdir(banderas []string) {
	ruta := ""
	p := false
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			ruta = dupla[1]

		} else if dupla[0] == "-p" {
			p = true
		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}
	if p {
		fmt.Println(p)
	}

	index := VerificarParticionMontada(actualIdMount)
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	crearDirectorio(ruta, index, archivo, &sblock)

}

func EjecutarMkfile(banderas []string) {

	ruta := ""
	r := false
	size := 0
	cont := ""

	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-path" {
			ruta = dupla[1]

		} else if dupla[0] == "-r" {
			r = true
		} else if dupla[0] == "-size" {
			size1, err := strconv.Atoi(dupla[1])
			size = size1
			if err != nil {
				fmt.Println("No se pudo leer ell -size")
				return
			}
		} else if dupla[0] == "-cont" {
			cont = dupla[1]
		} else {
			fmt.Println("Parametro invalido")
		}
	}
	if r {
		fmt.Println(".")
	}

	if cont != "" {
		cont = obtenerArchivo_externo(cont)
	} else if size > 0 {
		for i := 0; i < size; i++ {
			cont += "0"
		}
	}

	if ruta == "" {
		fmt.Println("No se ingreso el -path")
		return
	}

	index := VerificarParticionMontada(actualIdMount)

	crearArchivo(ruta, cont, r, index)

	fmt.Println("Archivo creado con exito")
	Salida_comando += "Archivo creado con exito\n"

}

func crearDirectorio(ruta string, index int, archivo *os.File, sblock *superBloque) {
	lRuta := strings.Split(ruta[1:], "/")

	var nombreArchivo []string
	var numInodo int
	if len(lRuta) == 1 {
		numInodo = 0
		nombreArchivo = lRuta
	} else {
		nombreArchivo = lRuta[len(lRuta)-1:]
		lRuta = lRuta[:len(lRuta)-1]

		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}

	if numInodo == -1 {
		fmt.Println("No encontro la ruta para crear el directorio")
		Salida_comando += "No encontro la ruta para crear el directorio\n"
		return
	}

	err12 := rellenarBloque(numInodo, nombreArchivo[0], archivo, sblock, 0)

	if err12 == -1 {
		fmt.Println("Error al escribir el archivo")
		Salida_comando += "Error al escribir el archivo\n"
	}

	var newInodo Inodo

	newInodo.I_uid = int32(uId)
	newInodo.I_gid = int32(gId)
	newInodo.I_s = int32(0)

	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")

	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_type = [1]byte{'0'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	var newBloqueCarpeta bloqueCarpeta

	copy(newBloqueCarpeta.B_content[0].B_name[:], ".")
	newBloqueCarpeta.B_content[0].B_inodo = sblock.S_firts_ino

	copy(newBloqueCarpeta.B_content[1].B_name[:], "..")
	newBloqueCarpeta.B_content[1].B_inodo = int32(numInodo)

	newBloqueCarpeta.B_content[2].B_inodo = -1
	newBloqueCarpeta.B_content[3].B_inodo = -1

	newInodo.I_block[0] = sblock.S_first_blo

	escribir_nuevoBloque_c(archivo, sblock, &newBloqueCarpeta)
	escribir_nuevoInodo(archivo, sblock, &newInodo)

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err := binary.Write(archivo, binary.LittleEndian, sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		Salida_comando += "Error al escribir el superbloque\n"
		return
	}

	archivo.Close()

	fmt.Println("Directorio creado con exito")
	Salida_comando += "Directorio creado con exito\n"
}

func crearExt2(index int, n int) {
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco \n"
		return
	}
	defer archivo.Close()
	archivo.Seek(int64(particionesMontadas[index].Start), 0)

	// for i := int32(0); i < particionesMontadas[index].Size; i++ {

	// 	err := binary.Write(archivo, binary.LittleEndian, [1]byte{0})
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 	}
	// }

	//Buffer 1024 bytes
	bufer := new(bytes.Buffer)
	for i := 0; i < 1024; i++ {
		bufer.WriteByte(0)
	}

	var totalBytes int = 0

	//Escribe 0s en archivo
	for totalBytes < int(particionesMontadas[index].Size) {
		c, err := archivo.Write(bufer.Bytes())
		if err != nil {
			fmt.Println("Error al escribir en el archivo: ", err)
			Salida_comando += "Error al escribir en el archivo\n"
			return
		}
		totalBytes += c
	}

	if totalBytes != int(particionesMontadas[index].Size) {
		for i := totalBytes; i < int(particionesMontadas[index].Size); i++ {

			err := binary.Write(archivo, binary.LittleEndian, [1]byte{0})
			if err != nil {
				fmt.Println("Error: ", err)
				Salida_comando += "Error al formatear archivo\n"
			}
		}
	}

	var sbloque superBloque

	sbloque.S_filesystem_type = 2
	sbloque.S_bm_inode_start = int32(particionesMontadas[index].Size) + int32(binary.Size(superBloque{}))
	sbloque.S_bm_block_start = sbloque.S_bm_inode_start + int32(n)
	sbloque.S_inode_start = sbloque.S_bm_block_start + int32(3*n)
	sbloque.S_block_start = sbloque.S_inode_start + int32(n*int(binary.Size(Inodo{})))

	sbloque.S_inodes_count = int32(n)
	sbloque.S_blocks_count = int32(3 * n)

	sbloque.S_free_inodes_count = int32(n)
	sbloque.S_free_blocks_count = int32(3 * n)
	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")
	copy(sbloque.S_mtime[:], []byte(fechaF))
	copy(sbloque.S_umtime[:], []byte(fechaF))
	sbloque.S_mnt_count = 1
	sbloque.S_magic = 61267
	sbloque.S_inode_s = int32(binary.Size(Inodo{}))
	sbloque.S_block_s = int32(binary.Size(bloqueArchivos{}))
	sbloque.S_firts_ino = 0
	sbloque.S_first_blo = 0

	var newInodo Inodo

	var newblock bloqueCarpeta

	newInodo.I_uid = 1
	newInodo.I_gid = 1
	newInodo.I_s = 0
	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_block[0] = 0
	newInodo.I_type = [1]byte{'0'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	copy(newblock.B_content[0].B_name[:], ".")
	newblock.B_content[0].B_inodo = 0
	copy(newblock.B_content[1].B_name[:], "..")
	newblock.B_content[1].B_inodo = 0
	newblock.B_content[2].B_inodo = -1
	newblock.B_content[3].B_inodo = -1

	archivo.Seek(int64(sbloque.S_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newInodo)
	if err != nil {
		fmt.Println("Error al escribir el inodo 0: ", err)
		Salida_comando += "Error al escribir inodo 0\n"
		return
	}
	archivo.Seek(int64(sbloque.S_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &newblock)
	if err != nil {
		fmt.Println("Error al escribir el bloque 0: ", err)
		Salida_comando += "Error al escribir el bloque 0\n"
		return
	}

	sbloque.S_free_blocks_count--
	sbloque.S_free_inodes_count--
	sbloque.S_firts_ino++
	sbloque.S_first_blo++

	archivo.Seek(int64(sbloque.S_bm_block_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap inodo 0: ", err)
		Salida_comando += "Error al escribir el bitmap inodo 0\n"
		return
	}

	archivo.Seek(int64(sbloque.S_bm_inode_start), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})
	if err != nil {
		fmt.Println("Error al escribir el bitmap bloque 0: ", err)
		Salida_comando += "Error al escribir el bitmap bloque 0\n"
		return
	}

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sbloque)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		Salida_comando += "Error al escribir el superbloque\n"
		return
	}

	archivo.Close()

	uId = 1
	gId = 1

	//creo archivo users.txt
	crearArchivo("/users.txt", "1,g,root\n1,u,root,root,123\n", false, index)

	uId = -1
	gId = -1

	fmt.Println("EXT2 creado con exito")
	Salida_comando += "EXT2 creado con exito\n"
}

func crearArchivo(ruta string, cont string, r bool, index int) {

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}
	lRuta := strings.Split(ruta[1:], "/")

	var nombreArchivo []string
	var numInodo int
	if len(lRuta) == 1 {
		numInodo = 0
		nombreArchivo = lRuta
	} else {
		nombreArchivo = lRuta[len(lRuta)-1:]
		lRuta = lRuta[:len(lRuta)-1]

		numInodo = obtenerNumInodo(lRuta, archivo, &sblock)
	}

	if numInodo == -1 {
		fmt.Println("No encontro la ruta para crear el archivo")
		Salida_comando += "No encontro la ruta para crear el archivo\n"
		return
	}
	err12 := rellenarBloque(numInodo, nombreArchivo[0], archivo, &sblock, len(cont))

	if err12 == -1 {
		fmt.Println("Error al escribir el archivo")
		Salida_comando += "Error al escribir el archivo\n"
	}

	var newInodo Inodo

	newInodo.I_uid = int32(uId)
	newInodo.I_gid = int32(gId)
	newInodo.I_s = int32(len(cont))

	fechaActual := time.Now()
	fechaF := fechaActual.Format("2006-01-02 15:04:05")

	copy(newInodo.I_atime[:], []byte(fechaF))
	copy(newInodo.I_ctime[:], []byte(fechaF))
	copy(newInodo.I_mtime[:], []byte(fechaF))

	for i := int32(0); i < 15; i++ {
		newInodo.I_block[i] = -1
	}

	newInodo.I_type = [1]byte{'1'}
	newInodo.I_perm = [3]byte{'6', '6', '4'}

	escribir_archivo(archivo, &sblock, &newInodo, cont)

	escribir_nuevoInodo(archivo, &sblock, &newInodo)

	//sblock.S_firts_ino = get_siguienteInodoLibre(archivo, &sblock)
	//sblock.S_free_inodes_count--

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		Salida_comando += "Error al escribir el superbloque\n"
		return
	}

	archivo.Close()
}

func editarArchivo(index int, ruta string, cont string) {
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el archivo\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	var numInodo int
	if ruta == "/" {
		numInodo = 0
	} else {
		lRuta := strings.Split(ruta[1:], "/")

		numInodo = obtenerNumInodo(lRuta, archivo, &sblock)
	}

	eliminarBloquesInodo(numInodo, &sblock, archivo)

	var inodoTemp Inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(numInodo)), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	//escribo el nuevo contenido
	escribir_archivo(archivo, &sblock, &inodoTemp, cont)
	inodoTemp.I_s = int32(len(cont))

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(numInodo)), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo: ", err)
		Salida_comando += "Error al escribir el inodo\n"
		return
	}

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Write(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al escribir el superbloque: ", err)
		Salida_comando += "Error al escribir el superbloque\n"
		return
	}

	archivo.Close()

}

func EjecMkUsr(banderas []string) {
	name := ""
	pass := ""
	group := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			name = dupla[1]

		} else if dupla[0] == "-pass" {
			pass = dupla[1]

		} else if dupla[0] == "-grp" {
			group = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		Salida_comando += "No se ingreso el campo -name\n"
		return
	}

	if pass == "" {
		fmt.Println("Ingrese el campo -pass")
		Salida_comando += "No se ingreso el campo -pass\n"
		return
	}

	if group == "" {
		fmt.Println("Ingrese el campo -grp")
		Salida_comando += "No se ingreso el campo -grp\n"
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede crear usuarios")
		Salida_comando += "Solo el usuario root puede crear usuarios\n"
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	var inodoTemp Inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	txt := leerArchivo("/users.txt", archivo, &sblock)
	nexId := nextIdUser(txt)

	txt += strconv.Itoa(nexId) + ",U," + group + "," + name + "," + pass + "\n"

	archivo.Close()
	editarArchivo(index, "/users.txt", txt)

	fmt.Println("Usuario creado con exito")
	Salida_comando += "Usuario creado con exito\n"

}

func EjecRmUsr(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-user" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		Salida_comando += "No se ingreso el campo -name\n"
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede eliminar usuarios")
		Salida_comando += "Solo el usuario root puede eliminar usuarios\n"
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	var inodoTemp Inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	txt := leerArchivo("/users.txt", archivo, &sblock)

	lineas := strings.Split(txt, "\n")
	lineas = lineas[:len(lineas)-1]

	band := false
	for i, linea := range lineas {

		if linea[2] == 'U' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[3] == name {
				nuevaLinea := []byte(linea)
				nuevaLinea[0] = '0'
				lineas[i] = string(nuevaLinea)
				band = true
				break

			}
		}

	}

	newTxt := ""
	for _, linea := range lineas {

		newTxt += linea + "\n"

	}

	if !band {
		fmt.Println("No se encontro el Usuario")
		Salida_comando += "No se encontro el Usuario\n"
		return
	}
	archivo.Close()
	editarArchivo(index, "/users.txt", newTxt)

	fmt.Println("Usuario eliminado con exito")
	Salida_comando += "Usuario eliminado con exito \n"

}

func EjecMkGrp(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		Salida_comando += "No se ingreso el campo -name\n"
		return
	}
	if uId != 1 {
		fmt.Println("Solo el usuario root puede crear grupos")
		Salida_comando += "Solo el usuario root puede crear grupos\n"
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	var inodoTemp Inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	txt := leerArchivo("/users.txt", archivo, &sblock)
	nexId := nextIdGroup(txt)

	txt += strconv.Itoa(nexId) + ",G," + name + "\n"

	archivo.Close()
	editarArchivo(index, "/users.txt", txt)

	fmt.Println("Grupo creado con exito")
	Salida_comando += "Grupo creado con exito\n"

}

func EjecRmGrp(banderas []string) {
	name := ""
	for _, valor := range banderas {
		dupla := strings.Split(valor, "=")

		if dupla[0] == "-name" {
			name = dupla[1]

		} else {
			fmt.Println("Parametro invalido")
			Salida_comando += "Parametro invalido\n"
		}
	}

	if name == "" {
		fmt.Println("Ingrese el campo -name")
		Salida_comando += "No se ingreso el campo -name\n"
		return
	}

	if uId != 1 {
		fmt.Println("Solo el usuario root puede eliminar grupos")
		Salida_comando += "Solo el usuario root puede eliminar grupos\n"
	}

	index := VerificarParticionMontada(actualIdMount)

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	var inodoTemp Inodo

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	txt := leerArchivo("/users.txt", archivo, &sblock)

	lineas := strings.Split(txt, "\n")
	lineas = lineas[:len(lineas)-1]

	band := false
	for i, linea := range lineas {

		if linea[2] == 'G' && linea[0] != '0' {
			campos := strings.Split(linea, ",")
			if campos[2] == name {
				nuevaLinea := []byte(linea)
				nuevaLinea[0] = '0'
				lineas[i] = string(nuevaLinea)
				band = true
				break

			}
		}

	}
	newTxt := ""
	for _, linea := range lineas {

		newTxt += linea + "\n"

	}

	if !band {
		fmt.Println("No se encontro el grupo")
		Salida_comando += "No se encontro el grupo\n"
		return
	}
	archivo.Close()
	editarArchivo(index, "/users.txt", newTxt)

	fmt.Println("Grupo eliminado con exito")
	Salida_comando += "Grupo eliminado con exito\n"

}

func nextIdUser(txt string) int {

	lineas := strings.Split(txt, "\n")
	id := 1

	lineas = lineas[:len(lineas)-1]

	var numTemp int
	for _, linea := range lineas {
		if linea[2] == 'U' {
			numTemp, _ = strconv.Atoi(string(linea[0]))
			if numTemp > id {
				id = numTemp
			}

		}
	}

	return id + 1

}

func nextIdGroup(txt string) int {

	lineas := strings.Split(txt, "\n")
	id := 1

	lineas = lineas[:len(lineas)-1]

	var numTemp int
	for _, linea := range lineas {
		if linea[2] == 'G' {
			numTemp, _ = strconv.Atoi(string(linea[0]))
			if numTemp > id {
				id = numTemp
			}

		}
	}

	return id + 1

}

func eliminarBloquesInodo(numInodo int, sblock *superBloque, archivo *os.File) {
	var inodoTemp Inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(numInodo)), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	for i, ptr := range inodoTemp.I_block {
		if ptr != -1 {
			eliminarBloque(int(ptr), sblock, archivo)
			inodoTemp.I_block[i] = -1
		}
	}

	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(numInodo)), 0)
	err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al escribir el inodo: ", err)
		Salida_comando += "Error al escribir el inodo\n"
		return
	}
}

func eliminarBloque(ptr int, sblock *superBloque, archivo *os.File) {
	bufer := make([]byte, 64)

	archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)
	err := binary.Write(archivo, binary.LittleEndian, &bufer)
	if err != nil {
		fmt.Println("Error al eliminar el bloque: ", err)
		Salida_comando += "Error al eliminar el bloque\n"
		return
	}

	archivo.Seek(int64(sblock.S_bm_block_start+int32(ptr)), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{0})
	if err != nil {
		fmt.Println("Error al eliminar el bitmap Bloque: ", err)
		Salida_comando += "Error al eliminar el bitmap Bloque\n"
		return
	}
	sblock.S_free_blocks_count++

	sblock.S_first_blo = get_siguienteBloqueLibre(archivo, sblock)
}

func leerArchivo(ruta string, archivo *os.File, sblock *superBloque) string {

	lRuta := strings.Split(ruta[1:], "/")

	//var nombreArchivo []string
	var numInodo int
	if len(lRuta) == 0 {
		fmt.Println("Ruta Archivo / no valido")
		Salida_comando += "Ruta Archivo / no valido\n"
		//nombreArchivo = lRuta
	} else {
		//lRuta = lRuta[:len(lRuta)-1]
		//nombreArchivo = lRuta[len(lRuta)-1:]
		numInodo = obtenerNumInodo(lRuta, archivo, sblock)
	}
	if numInodo == -1 {
		fmt.Println("Ruta no encontrada\n")
		Salida_comando += "Ruta no encontrada\n"
		return ""
	}

	var inodoTemp Inodo
	archivo.Seek(int64(sblock.S_inode_start+(int32(numInodo)*int32(binary.Size(Inodo{})))), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return ""
	}
	result := leerBloquesInodo(&inodoTemp, archivo, sblock)
	return result
}

func leerBloquesInodo(inodoTemp *Inodo, archivo *os.File, sblock *superBloque) string {
	if inodoTemp.I_type != [1]byte{'1'} {
		fmt.Println("El inodo no es de tipo archivo\n")
		Salida_comando += "El inodo no es de tipo archivo\n"
		return ""
	}

	result := ""

	for i, ptr := range inodoTemp.I_block {
		if ptr != -1 {
			if i == 12 {
				//apuntador simple
				return result
			} else if i == 13 {
				//apuntador doble
				return result
			} else if i == 14 {
				//apuntador triple
				return result
			} else {
				//Directo
				result += leerBloqueArchivos_directo(int(ptr), archivo, sblock)
			}

		}
	}

	return result

}

func leerBloqueArchivos_directo(numBloque int, archivo *os.File, sblock *superBloque) string {

	result := ""
	var bloqueTemp bloqueArchivos

	despTemp := sblock.S_block_start + (int32(numBloque) * int32(binary.Size(bloqueArchivos{})))
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &bloqueTemp)

	if err != nil {
		fmt.Println("Error al leer el bloque de archivo")
		Salida_comando += "Error al leer el bloque de archivo\n"
		return ""
	}

	result = strings.TrimRight(string(bloqueTemp.B_content[:]), string(rune(0)))
	return result
}

func obtenerNumInodo(ruta []string, archivo *os.File, sblock *superBloque) int {

	numInodo := 0

	for _, nombreCarpeta := range ruta {
		numInodo = buscarEnBloques(nombreCarpeta, numInodo, archivo, sblock)
		if numInodo == -1 {
			//no lo encontro, verificar el R para crear la nueva carpeta e numInodo= ptr nueva carpeta
			return numInodo
		}
	}

	return numInodo

}

func buscarEnBloques(nombre string, numInodo int, archivo *os.File, sblock *superBloque) int {

	despTemp := int(sblock.S_inode_start) + numInodo*(binary.Size(Inodo{}))

	var inodoTemp Inodo
	//var bloqueCarpetaTemp bloqueCarpeta
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return -1
	}

	for i, ptr := range inodoTemp.I_block {

		if ptr != -1 {

			if i == 12 {
				//apuntadir directo
				temp := buscar_inodo_bloqueIndirecto(int(ptr), nombre, 0, archivo, sblock)
				if temp != -1 {
					return temp
				}
			} else if i == 13 {
				//doble directo
				temp := buscar_inodo_bloqueIndirecto(int(ptr), nombre, 1, archivo, sblock)
				if temp != -1 {
					return temp
				}
			} else if i == 14 {
				//triple indirecto
				temp := buscar_inodo_bloqueIndirecto(int(ptr), nombre, 2, archivo, sblock)
				if temp != -1 {
					return temp
				}
			} else {
				//directo
				temp := buscar_inodo_bloqueDirecto(int(ptr), nombre, archivo, sblock)
				if temp != -1 {
					return temp
				}
			}

		}
	}

	return -1
}

func buscar_inodo_bloqueDirecto(ptr int, nombre string, archivo *os.File, sblock *superBloque) int {
	var bloqueCarpetaTemp bloqueCarpeta
	despTemp := int(sblock.S_block_start) + int(ptr*binary.Size(bloqueCarpeta{}))
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &bloqueCarpetaTemp)
	if err != nil {
		fmt.Println("Error al leer el bloqueCarpeta: ", err)
		Salida_comando += "Error al leer el bloqueCarpeta\n"
		return -1
	}
	for _, cont := range bloqueCarpetaTemp.B_content {
		if strings.Contains(string(cont.B_name[:]), nombre) {
			return int(cont.B_inodo)
		}
	}
	return -1
}

func buscar_inodo_bloqueIndirecto(ptr int, nombre string, profundidad int, archivo *os.File, sblock *superBloque) int {
	var bloqueApuntadoresTemp bloqueApuntadores
	despTemp := int(sblock.S_block_start) + int(ptr*binary.Size(bloqueApuntadores{}))
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &bloqueApuntadoresTemp)
	if err != nil {
		fmt.Println("Error al leer el bloqueCarpeta: ", err)
		Salida_comando += "Error al leer el bloqueCarpeta\n"
		return -1
	}
	for _, cont := range bloqueApuntadoresTemp.B_pointers {
		if cont != -1 {
			if profundidad == 0 {
				temp := buscar_inodo_bloqueDirecto(int(cont), nombre, archivo, sblock)
				if temp != -1 {
					return temp
				}
			} else {
				temp := buscar_inodo_bloqueIndirecto(int(cont), nombre, profundidad-1, archivo, sblock)
				if temp != -1 {
					return temp
				}
			}
		}
	}
	return -1
}

func rellenarBloque(numInodo int, nombre string, archivo *os.File, sblock *superBloque, size int) int {

	archivo.Seek(int64(sblock.S_inode_start+int32(numInodo)*int32(binary.Size(Inodo{}))), 0)
	var inodoTemp Inodo
	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return -1
	}

	for i, ptr := range inodoTemp.I_block {
		if ptr != -1 {
			//Rellena Bloque
			if i == 12 {
				//simple indirecto
				temp := rellenarBloque_indirecto(int(ptr), nombre, 0, archivo, sblock)
				if temp != -1 {
					return temp
				}
				continue
			} else if i == 13 {
				//doble indirecto
				temp := rellenarBloque_indirecto(int(ptr), nombre, 1, archivo, sblock)
				if temp != -1 {
					return temp
				}
				continue
			} else if i == 14 {
				//triple indirecto
				temp := rellenarBloque_indirecto(int(ptr), nombre, 2, archivo, sblock)
				if temp != -1 {
					return temp
				}
				continue
			} else {
				//directo
				temp := rellenarBloque_directo(int(ptr), nombre, archivo, sblock)
				if temp != -1 {
					return temp
				}
				continue
			}
		} else {

			//Crear Bloque
			if i == 12 {
				//simple indirecto
				var newBloqueCarpeta bloqueCarpeta
				copy(newBloqueCarpeta.B_content[0].B_name[:], nombre)
				newBloqueCarpeta.B_content[0].B_inodo = sblock.S_firts_ino
				temp := sblock.S_firts_ino
				newBloqueCarpeta.B_content[1].B_inodo = -1
				newBloqueCarpeta.B_content[2].B_inodo = -1
				newBloqueCarpeta.B_content[3].B_inodo = -1
				escribir_nuevoBloque_c(archivo, sblock, &newBloqueCarpeta)

				var newBloquePuntero1 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero1.B_pointers); i++ {
					newBloquePuntero1.B_pointers[i] = -1
				}

				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_firts_ino
				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero1)

				inodoTemp.I_block[i] = temp
				inodoTemp.I_s += int32(size)

				archivo.Seek(int64(int(sblock.S_inode_start)+(numInodo*binary.Size(Inodo{}))), 0)
				err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
				if err != nil {
					fmt.Println("Error al escribir el inodo")
					Salida_comando += "Error al escribir el inodo\n"
					return -1
				}
				return int(temp)

			} else if i == 13 {
				//doble indirecto
				var newBloqueCarpeta bloqueCarpeta
				copy(newBloqueCarpeta.B_content[0].B_name[:], nombre)
				newBloqueCarpeta.B_content[0].B_inodo = sblock.S_firts_ino
				temp := sblock.S_firts_ino
				newBloqueCarpeta.B_content[1].B_inodo = -1
				newBloqueCarpeta.B_content[2].B_inodo = -1
				newBloqueCarpeta.B_content[3].B_inodo = -1

				escribir_nuevoBloque_c(archivo, sblock, &newBloqueCarpeta)

				var newBloquePuntero1 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero1.B_pointers); i++ {
					newBloquePuntero1.B_pointers[i] = -1
				}

				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_firts_ino
				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero1)

				var newBloquePuntero2 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero2.B_pointers); i++ {
					newBloquePuntero2.B_pointers[i] = -1
				}
				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_firts_ino

				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero2)

				inodoTemp.I_block[i] = temp
				inodoTemp.I_s += int32(size)

				archivo.Seek(int64(int(sblock.S_inode_start)+(numInodo*binary.Size(Inodo{}))), 0)
				err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
				if err != nil {
					fmt.Println("Error al escribir el inodo")
					Salida_comando += "Error al escribir el inodo\n"
					return -1
				}
				return int(temp)

			} else if i == 14 {
				//triple indirecto
				var newBloqueCarpeta bloqueCarpeta
				copy(newBloqueCarpeta.B_content[0].B_name[:], nombre)
				newBloqueCarpeta.B_content[0].B_inodo = sblock.S_firts_ino
				temp := sblock.S_first_blo
				newBloqueCarpeta.B_content[1].B_inodo = -1
				newBloqueCarpeta.B_content[2].B_inodo = -1
				newBloqueCarpeta.B_content[3].B_inodo = -1

				escribir_nuevoBloque_c(archivo, sblock, &newBloqueCarpeta)

				var newBloquePuntero1 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero1.B_pointers); i++ {
					newBloquePuntero1.B_pointers[i] = -1
				}

				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_first_blo
				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero1)

				var newBloquePuntero2 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero2.B_pointers); i++ {
					newBloquePuntero2.B_pointers[i] = -1
				}
				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_first_blo

				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero2)

				var newBloquePuntero3 bloqueApuntadores

				for i := 0; i < len(newBloquePuntero3.B_pointers); i++ {
					newBloquePuntero3.B_pointers[i] = -1
				}
				newBloquePuntero1.B_pointers[0] = temp
				temp = sblock.S_first_blo

				escribir_nuevoBloque_a(archivo, sblock, &newBloquePuntero3)

				inodoTemp.I_block[i] = temp
				inodoTemp.I_s += int32(size)

				archivo.Seek(int64(int(sblock.S_inode_start)+(numInodo*binary.Size(Inodo{}))), 0)
				err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
				if err != nil {
					fmt.Println("Error al escribir el inodo")
					Salida_comando += "Error al escribir el inodo\n"
					return -1
				}

				return int(temp)

			} else {
				//directo
				var newBloqueCarpeta bloqueCarpeta
				copy(newBloqueCarpeta.B_content[0].B_name[:], nombre)
				newBloqueCarpeta.B_content[0].B_inodo = sblock.S_firts_ino
				temp := sblock.S_first_blo

				newBloqueCarpeta.B_content[1].B_inodo = -1
				newBloqueCarpeta.B_content[2].B_inodo = -1
				newBloqueCarpeta.B_content[3].B_inodo = -1
				escribir_nuevoBloque_c(archivo, sblock, &newBloqueCarpeta)

				inodoTemp.I_block[i] = temp
				inodoTemp.I_s += int32(size)

				archivo.Seek(int64(int(sblock.S_inode_start)+(numInodo*binary.Size(Inodo{}))), 0)
				err = binary.Write(archivo, binary.LittleEndian, &inodoTemp)
				if err != nil {
					fmt.Println("Error al escribir el inodo")
					Salida_comando += "Error al escribir el inodo\n"
					return -1
				}

				return int(temp)
			}
		}
	}
	return -1

}

func rellenarBloque_directo(ptr int, nombre string, archivo *os.File, sblock *superBloque) int {
	var bloqueCarpetaTemp bloqueCarpeta
	despTemp := int(sblock.S_block_start) + int(ptr*binary.Size(bloqueCarpeta{}))
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &bloqueCarpetaTemp)
	if err != nil {
		fmt.Println("Error al leer el bloqueCarpeta: ", err)
		Salida_comando += "Error al leer el bloqueCarpeta\n"
		return -1
	}

	for i := 0; i < 4; i++ {
		if bloqueCarpetaTemp.B_content[i].B_inodo == -1 {
			//rellenar bloque
			bloqueCarpetaTemp.B_content[i].B_inodo = sblock.S_firts_ino
			copy(bloqueCarpetaTemp.B_content[i].B_name[:], nombre)

			archivo.Seek(int64(despTemp), 0)
			binary.Write(archivo, binary.LittleEndian, &bloqueCarpetaTemp)

			return ptr
		}
	}
	return -1

}

func rellenarBloque_indirecto(ptr int, nombre string, profundidad int, archivo *os.File, sblock *superBloque) int {
	var bloqueApuntadoresTemp bloqueApuntadores
	despTemp := int(sblock.S_block_start) + int(ptr*binary.Size(bloqueApuntadores{}))
	archivo.Seek(int64(despTemp), 0)
	err := binary.Read(archivo, binary.LittleEndian, &bloqueApuntadoresTemp)
	if err != nil {
		fmt.Println("Error al leer el bloqueCarpeta: ", err)
		Salida_comando += "Error al leer el bloqueCarpeta\n"
		return -1
	}
	for _, cont := range bloqueApuntadoresTemp.B_pointers {
		if cont != -1 {
			if profundidad == 0 {
				temp := rellenarBloque_directo(int(cont), nombre, archivo, sblock)
				if temp != -1 {
					return temp
				}
			} else {
				temp := rellenarBloque_indirecto(int(cont), nombre, profundidad-1, archivo, sblock)
				if temp != -1 {
					return temp
				}
			}
		}
	}

	return -1

}

func escribir_archivo(archivo *os.File, sblock *superBloque, newInodo *Inodo, cont string) {
	cantBloques := len(cont) / 64

	if cantBloques > 4390 {
		fmt.Println("no se pudo escribir el archivo,archivo muy grande")
		Salida_comando += "Archivo muy grande para escribir\n"
		return
	}

	for i := 0; i <= cantBloques; i++ {

		if i < 12 {
			//bloques directos
			var newBloqueArchivos bloqueArchivos

			if len(cont) > 64 {
				copy(newBloqueArchivos.B_content[:], cont[:64])
				cont = cont[64:]
			} else {
				copy(newBloqueArchivos.B_content[:], cont)
			}

			newInodo.I_block[i] = sblock.S_first_blo

			escribir_nuevoBloque_d(archivo, sblock, &newBloqueArchivos)

		} else if i < 28 {
			//apuntadores simples
			return
		} else if i < 284 {
			//apuntadores dobles
			return
		} else {
			//Apuntadores triples
			return
		}
	}
}

func escribir_nuevoBloque_c(archivo *os.File, sblock *superBloque, bloque *bloqueCarpeta) {

	despTemp := int(sblock.S_block_start) + (int(sblock.S_first_blo) * binary.Size(bloqueCarpeta{}))

	archivo.Seek(int64(despTemp), 0)
	err := binary.Write(archivo, binary.LittleEndian, bloque)

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	despTemp = int(sblock.S_bm_block_start) + int(sblock.S_first_blo)

	archivo.Seek(int64(despTemp), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	sblock.S_first_blo = get_siguienteBloqueLibre(archivo, sblock)
	sblock.S_free_blocks_count--
}

func escribir_nuevoBloque_a(archivo *os.File, sblock *superBloque, bloque *bloqueApuntadores) {

	despTemp := int(sblock.S_block_start) + (int(sblock.S_first_blo) * binary.Size(bloqueApuntadores{}))

	archivo.Seek(int64(despTemp), 0)
	err := binary.Write(archivo, binary.LittleEndian, bloque)

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	despTemp = int(sblock.S_bm_block_start) + int(sblock.S_first_blo)

	archivo.Seek(int64(despTemp), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	sblock.S_first_blo = get_siguienteBloqueLibre(archivo, sblock)
	sblock.S_free_blocks_count--
}

func escribir_nuevoBloque_d(archivo *os.File, sblock *superBloque, bloque *bloqueArchivos) {

	despTemp := int(sblock.S_block_start) + (int(sblock.S_first_blo) * binary.Size(bloqueArchivos{}))

	archivo.Seek(int64(despTemp), 0)
	err := binary.Write(archivo, binary.LittleEndian, bloque)

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	despTemp = int(sblock.S_bm_block_start) + int(sblock.S_first_blo)

	archivo.Seek(int64(despTemp), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})

	if err != nil {
		fmt.Println("Error al escribir el bloque")
		Salida_comando += "Error al escribir el bloque\n"
		return
	}

	sblock.S_first_blo = get_siguienteBloqueLibre(archivo, sblock)
	sblock.S_free_blocks_count--
}

func escribir_nuevoInodo(archivo *os.File, sblock *superBloque, newInodo *Inodo) {

	despTemp := int(sblock.S_inode_start) + (int(sblock.S_firts_ino) * binary.Size(Inodo{}))

	archivo.Seek(int64(despTemp), 0)
	err := binary.Write(archivo, binary.LittleEndian, newInodo)

	if err != nil {
		fmt.Println("Error al escribir el inodo")
		Salida_comando += "Error al escribir el inodo\n"
		return
	}

	despTemp = int(sblock.S_bm_inode_start) + int(sblock.S_firts_ino)

	archivo.Seek(int64(despTemp), 0)
	err = binary.Write(archivo, binary.LittleEndian, [1]byte{1})

	if err != nil {
		fmt.Println("Error al escribir el inodo")
		Salida_comando += "Error al escribir el inodo\n"
		return
	}

	sblock.S_firts_ino = get_siguienteInodoLibre(archivo, sblock)
	sblock.S_free_inodes_count--
}

func get_siguienteInodoLibre(archivo *os.File, sblock *superBloque) int32 {

	var tempByte [1]byte
	archivo.Seek(int64(sblock.S_bm_inode_start), 0)
	for i := 0; i < int(sblock.S_inodes_count); i++ {

		err := binary.Read(archivo, binary.LittleEndian, &tempByte)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			Salida_comando += "Error al leer el bitmap inodos\n"
			return -1
		}
		if tempByte == [1]byte{0} {
			return int32(i)
		}
	}
	return -1

}

func get_siguienteBloqueLibre(archivo *os.File, sblock *superBloque) int32 {

	var tempByte [1]byte
	archivo.Seek(int64(sblock.S_bm_block_start), 0)
	for i := 0; i < int(sblock.S_blocks_count); i++ {

		err := binary.Read(archivo, binary.LittleEndian, &tempByte)
		if err != nil {
			fmt.Println("Error al leer el bitmap Bloques: ", err)
			Salida_comando += "Error al leer el bitmap Bloques\n"
			return -1
		}
		if tempByte == [1]byte{0} {
			return int32(i)
		}
	}
	return -1

}

func RepMbr(index int, path string) {
	archivo, err := os.Open(particionesMontadas[index].Path)

	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte MBR \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">MBR</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_tamano</TD><TD>" + strconv.Itoa(int(disk.Mbr_tamano)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_fecha_creacion</TD><TD>" + string(disk.Mbr_fecha_creacion[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">mbr_disk_signature</TD><TD>" + strconv.Itoa(int(disk.Mbr_dsk_signature)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">dsk_fit</TD><TD>" + string(disk.MBR_dsk_fit[:]) + "</TD></TR>\n"
	var ebrTemp EBR
	var despTemp int

	for _, part := range disk.Mbr_partitions {
		if part.Part_type == [1]byte{'e'} {
			name := strings.TrimRight(string(part.Part_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_status</TD><TD>" + string(part.Part_status[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_type</TD><TD>p</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(part.Part_fit[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(part.Part_start)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(part.Part_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"
			despTemp = int(part.Part_start)
			archivo.Seek(int64(despTemp), 0)
			err = binary.Read(archivo, binary.LittleEndian, &ebrTemp)
			if err != nil {
				fmt.Println("Error al leer el EBR: ", err)
				Salida_comando += "Error al leer el EBR\n"
				return
			}
			for ebrTemp.Part_s > 0 {
				name = strings.TrimRight(string(ebrTemp.Part_name[:]), string(rune(0)))
				Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion Logica</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_mount</TD><TD>" + string(ebrTemp.Part_mount[:]) + "</TD></TR>\n"

				Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(ebrTemp.Part_fit[:]) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_start)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_s)) + "</TD></TR>\n"

				Dot += "<TR><TD bgcolor=\"lightgrey\">part_next</TD><TD>" + strconv.Itoa(int(ebrTemp.Part_next)) + "</TD></TR>\n"
				Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"

				despTemp += int(ebrTemp.Part_s) + 1 + binary.Size(EBR{})
				archivo.Seek(int64(despTemp), 0)
				err = binary.Read(archivo, binary.LittleEndian, &ebrTemp)
				if err != nil {
					fmt.Println("Error al leer el EBR: ", err)
					Salida_comando += "Error al leer el EBR\n"
					return
				}

			}

		} else if part.Part_type == [1]byte{'p'} {
			name := strings.TrimRight(string(part.Part_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Particion</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_status</TD><TD>" + string(part.Part_status[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_type</TD><TD>p</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_fit</TD><TD>" + string(part.Part_fit[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_start</TD><TD>" + strconv.Itoa(int(part.Part_start)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_s</TD><TD>" + strconv.Itoa(int(part.Part_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">part_name</TD><TD>" + name + "</TD></TR>\n"

		}
	}

	Dot += "</TABLE>>];\n}"

	crearArchivoDot(path, Dot)
}

func RepDisk(index int, path string) {
	archivo, err := os.Open(particionesMontadas[index].Path)

	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var disk MBR
	archivo.Seek(int64(0), 0)
	err = binary.Read(archivo, binary.LittleEndian, &disk)
	if err != nil {
		fmt.Println("Error al leer el MBR del disco: ", err)
		Salida_comando += "Error al leer el MBR\n"
		return
	}

	sizeMBR := int(disk.Mbr_tamano)
	libre := int(disk.Mbr_tamano)

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte Disk \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "node[shape=record, color=lightgrey]a0[label=\"MBR"

	for _, part := range disk.Mbr_partitions {

		if part.Part_s != 0 {
			libre -= int(part.Part_s)
			Dot += "|"
			if part.Part_type == [1]byte{'e'} {

				Dot += "{Extendida"
				libreExtendida := part.Part_s

				var ebr EBR
				desp := int(part.Part_start)
				archivo.Seek(int64(desp), 0)
				err := binary.Read(archivo, binary.LittleEndian, &ebr)
				if err != nil {
					fmt.Println("Error al leer el EBR: ", err)
					Salida_comando += "Error al leer el EBR\n"
					return
				}

				if ebr.Part_s != 0 {
					libreExtendida -= ebr.Part_s
					Dot += "|{EBR"

					Dot += "|Logica"
					porcentaje := (float64(ebr.Part_s) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					//libre -= int(ebr.Part_s)

					desp += int(ebr.Part_s) + 1 + binary.Size(EBR{})
					archivo.Seek(int64(desp), 0)
					err = binary.Read(archivo, binary.LittleEndian, &ebr)
					if err != nil {
						fmt.Println("Error al leer el ebr: ", err)
						Salida_comando += "Error al leer el EBR\n"
						return
					}

					for {
						if ebr.Part_s == 0 {
							break
						}
						libreExtendida -= ebr.Part_s
						Dot += "|EBR"
						Dot += "|Logica"
						porcentaje := (float64(ebr.Part_s) * float64(100)) / float64(sizeMBR)
						Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
						//libre -= int(ebr.Part_s)

						desp += int(ebr.Part_s) + 1 + binary.Size(EBR{})
						archivo.Seek(int64(desp), 0)
						err = binary.Read(archivo, binary.LittleEndian, &ebr)
						if err != nil {
							fmt.Println("Error al leer el ebr: ", err)
							Salida_comando += "Error al leer el EBR\n"
							return
						}
					}

					if libreExtendida > 0 {
						Dot += "|Libre"
						porcentaje := (float64(libreExtendida) * float64(100)) / float64(sizeMBR)
						Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					}
					Dot += "}}"
				} else {
					Dot += "|Libre"
					porcentaje := (float64(part.Part_s) * float64(100)) / float64(sizeMBR)
					Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
					Dot += "}"
				}
			} else {
				Dot += "Primaria"
				porcentaje := (float64(part.Part_s) * float64(100)) / float64(sizeMBR)
				Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
			}
		}
	}

	if libre > 0 {
		Dot += "|Libre"
		porcentaje := (float64(libre) * float64(100)) / float64(sizeMBR)
		Dot += "\\n" + fmt.Sprintf("%.2f", porcentaje) + "%\\n"
	}
	Dot += "\"];\n}"

	crearArchivoDot(path, Dot)

}

func RepBmInodos(index int, path string) {

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	bitmap := ""
	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}
	size := int(sblock.S_bm_block_start - sblock.S_bm_inode_start)
	var temp [1]byte
	archivo.Seek(int64(sblock.S_bm_inode_start), 0)
	for i := 0; i < size; i++ {
		err = binary.Read(archivo, binary.LittleEndian, &temp)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			Salida_comando += "Error al leer el bitmap inodos\n"
			return
		}
		if temp == [1]byte{1} {
			bitmap += "1 "
		} else {
			bitmap += "0 "
		}
		if i%20 == 0 {
			bitmap += "\n"
		}

	}
	archivo.Close()

	crearTxt(path, bitmap)

}

func RepBmBloques(index int, path string) {

	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	bitmap := ""
	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}
	size := int(sblock.S_bm_block_start-sblock.S_bm_inode_start) * 3
	var temp [1]byte
	archivo.Seek(int64(sblock.S_bm_block_start), 0)
	for i := 0; i < size; i++ {
		err = binary.Read(archivo, binary.LittleEndian, &temp)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			Salida_comando += "Error al leer el bitmap inodos\n"
			return
		}
		if temp == [1]byte{1} {
			bitmap += "1 "
		} else {
			bitmap += "0 "
		}
		if i%20 == 0 {
			bitmap += "\n"
		}

	}
	archivo.Close()

	crearTxt(path, bitmap)

}

func RepSb(index int, path string) {
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	Dot := "digraph grid {bgcolor=\"slategrey\" label=\" Reporte SuperBlock \"layout=dot "
	Dot += "labelloc = \"t\"edge [weigth=1000 style=dashed color=red4 dir = \"both\" arrowtail=\"open\" arrowhead=\"open\"]"
	Dot += "a0[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">SuperBlock</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_filesystem_type</TD><TD>" + strconv.Itoa(int(sblock.S_filesystem_type)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inodes_count</TD><TD>" + strconv.Itoa(int(sblock.S_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_blocks_count</TD><TD>" + strconv.Itoa(int(sblock.S_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_blocks_count</TD><TD>" + strconv.Itoa(int(sblock.S_free_blocks_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_free_inodes_count</TD><TD>" + strconv.Itoa(int(sblock.S_free_inodes_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mtime</TD><TD>" + string(sblock.S_mtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_umtime</TD><TD>" + string(sblock.S_umtime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_mnt_count</TD><TD>" + strconv.Itoa(int(sblock.S_mnt_count)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_magic</TD><TD>" + strconv.Itoa(int(sblock.S_magic)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_size</TD><TD>" + strconv.Itoa(int(sblock.S_inode_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_size</TD><TD>" + strconv.Itoa(int(sblock.S_block_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_ino</TD><TD>" + strconv.Itoa(int(sblock.S_firts_ino)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_first_blo</TD><TD>" + strconv.Itoa(int(sblock.S_first_blo)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_inode_start</TD><TD>" + strconv.Itoa(int(sblock.S_bm_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_bm_block_start</TD><TD>" + strconv.Itoa(int(sblock.S_bm_block_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_inode_start</TD><TD>" + strconv.Itoa(int(sblock.S_inode_start)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">s_block_start</TD><TD>" + strconv.Itoa(int(sblock.S_block_start)) + "</TD></TR>\n"
	Dot += "</TABLE>>];\n}"

	crearArchivoDot(path, Dot)
}

func RepInodos(index int, path string) {
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	Dot := "digraph grid {\nbgcolor=\"slategrey\";\n label=\" Reporte Inodos \";\n layout=dot;\n "
	Dot += "labelloc = \"t\"; \n edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=open arrowhead=open];\n"
	var inodoTemp Inodo
	archivo.Seek(int64(sblock.S_inode_start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo: ", err)
		Salida_comando += "Error al leer el inodo\n"
		return
	}

	Dot += "inodo"
	Dot += strconv.Itoa(0)
	Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Inodo " + strconv.Itoa(0) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_uid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_gid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_s</TD><TD>" + strconv.Itoa(int(inodoTemp.I_s)) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_atime</TD><TD>" + string(inodoTemp.I_atime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_ctime</TD><TD>" + string(inodoTemp.I_ctime[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_mtime</TD><TD>" + string(inodoTemp.I_mtime[:]) + "</TD></TR>\n"
	for i, ptr := range inodoTemp.I_block {
		Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
	}

	Dot += "<TR><TD bgcolor=\"lightgrey\">I_type</TD><TD>" + string(inodoTemp.I_type[:]) + "</TD></TR>\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\">I_perm</TD><TD>" + string(inodoTemp.I_perm[:]) + "</TD></TR>\n"

	Dot += "</TABLE>>];\n"

	var byteTemp [1]byte
	for i := 1; i < int(sblock.S_inodes_count); i++ {
		archivo.Seek(int64(sblock.S_bm_inode_start+int32(i)), 0)
		err = binary.Read(archivo, binary.LittleEndian, &byteTemp)
		if err != nil {
			fmt.Println("Error al leer el bitmap inodos: ", err)
			Salida_comando += "Error al leer el bitmap inodos\n"
			return
		}
		if byteTemp == [1]byte{1} {
			archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(i)), 0)
			err = binary.Read(archivo, binary.LittleEndian, &inodoTemp)

			if err != nil {
				fmt.Println("Error al leer el inodo: ", err)
				Salida_comando += "Error al leer el bitmap inodos\n"
				return
			}
			Dot += "inodo"
			Dot += strconv.Itoa(i)
			Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">Inodo " + strconv.Itoa(i) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_uid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_gid</TD><TD>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_s</TD><TD>" + strconv.Itoa(int(inodoTemp.I_s)) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_atime</TD><TD>" + string(inodoTemp.I_atime[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_ctime</TD><TD>" + string(inodoTemp.I_ctime[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_mtime</TD><TD>" + string(inodoTemp.I_mtime[:]) + "</TD></TR>\n"
			for i, ptr := range inodoTemp.I_block {
				Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
			}

			Dot += "<TR><TD bgcolor=\"lightgrey\">I_type</TD><TD>" + string(inodoTemp.I_type[:]) + "</TD></TR>\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\">I_perm</TD><TD>" + string(inodoTemp.I_perm[:]) + "</TD></TR>\n"

			Dot += "</TABLE>>];\n"
			Dot += "inodo" + strconv.Itoa(i-1) + " -> inodo" + strconv.Itoa(i) + ";\n"
		}

	}
	Dot += "}"

	crearArchivoDot(path, Dot)
}

func RepBloques(index int, path string) {
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		Salida_comando += "Error al abrir el disco\n"
		return
	}
	defer archivo.Close()

	var sblock superBloque
	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &sblock)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		Salida_comando += "Error al leer el superbloque\n"
		return
	}

	Dot := "digraph grid {\n bgcolor=\"slategrey\";\n label=\" Reporte Bloques \";\n layout=dot;\n "
	Dot += "labelloc = \"t\";\n edge [weight=1000 style=dashed color=red4 dir = \"both\" arrowtail=open arrowhead=open];\n"
	var bCarpeta bloqueCarpeta
	archivo.Seek(int64(sblock.S_block_start), 0)
	err = binary.Read(archivo, binary.LittleEndian, &bCarpeta)
	if err != nil {
		fmt.Println("Error al leer el bloque carpetas: ", err)
		Salida_comando += "Error al leer el bloque carpetas\n"
		return
	}

	Dot += "bloque"
	Dot += strconv.Itoa(0)
	Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
	Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(0) + "</TD></TR>\n"

	Dot += "<TR><TD bgcolor=\"lightgrey\">b_name</TD><TD>b_inodo</TD></TR>\n"
	for _, cont := range bCarpeta.B_content {
		nam := strings.TrimRight(string(cont.B_name[:]), string(rune(0)))
		Dot += "<TR><TD bgcolor=\"lightgrey\">" + nam + "</TD><TD>" + strconv.Itoa(int(cont.B_inodo)) + "</TD></TR>\n"

	}

	Dot += "</TABLE>>];\n"
	var byteTemp [1]byte
	var bArchivo bloqueArchivos
	for i := 1; i < int(sblock.S_inodes_count); i++ {
		archivo.Seek(int64(sblock.S_bm_block_start+int32(i)), 0)
		err = binary.Read(archivo, binary.LittleEndian, &byteTemp)
		if err != nil {
			fmt.Println("Error al leer el bitmap bloques: ", err)
			Salida_comando += "Error al leer el bitmap bloques\n"
			return
		}
		if byteTemp == [1]byte{1} {
			archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(i)), 0)
			err = binary.Read(archivo, binary.LittleEndian, &bArchivo)

			if err != nil {
				fmt.Println("Error al leer el bloque archivos: ", err)
				Salida_comando += "Error al leer el bloque archivos\n"
				return
			}
			Dot += "bloque"
			Dot += strconv.Itoa(i)
			Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(i) + "</TD></TR>\n"
			cont := quitarCaracteresEspeciales(string(bArchivo.B_content[:]))
			Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">" + cont + "</TD></TR>\n"

			Dot += "</TABLE>>];\n"
			Dot += "bloque" + strconv.Itoa(i-1) + " -> bloque" + strconv.Itoa(i) + ";\n"
		}

	}

	Dot += "}"

	crearArchivoDot(path, Dot)
}

func RepTree(index int, path string) {

	//Abrir el disco
	archivo, err := os.OpenFile(particionesMontadas[index].Path, os.O_RDWR, 0664)
	if err != nil {
		fmt.Println("Error al abrir el disco: ", err)
		return
	}
	defer archivo.Close()

	archivo.Seek(int64(particionesMontadas[index].Start), 0)
	//Leer el superbloque
	var sb superBloque
	err = binary.Read(archivo, binary.LittleEndian, &sb)
	if err != nil {
		fmt.Println("Error al leer el superbloque: ", err)
		return
	}

	//Buscar el inodo raiz
	var raiz Inodo
	archivo.Seek(int64(sb.S_inode_start), 0)
	binary.Read(archivo, binary.LittleEndian, &raiz)
	Dot := "digraph H {\n"
	Dot += "node [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"1\"];\n"
	Dot += "node [shape=plaintext];\n"
	Dot += "graph [bb=\"0,0,352,154\"];\n"
	Dot += "rankdir=LR;\n"
	Dot += crearDotNodoTree(0, archivo, sb)
	Dot += "}"

	// archivoDot, err := os.Create("reporteTree.dot")
	// if err != nil {
	// 	fmt.Println("Error al crear el archivo .dot: ", err)
	// 	return
	// }
	// defer archivoDot.Close()

	// _, err = archivoDot.WriteString(Dot)
	// if err != nil {
	// 	fmt.Println("Error al escribir el archivo .dot: ", err)
	// 	return
	// }

	// cmd := exec.Command("dot", "-T", "png", "reporteTree.dot", "-o", "reporteTree.png")

	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println("Error al generar la imagen: ", err)
	// 	return
	// }

	// fmt.Println("Reporte generado con exito")

	crearArchivoDot(path, Dot)

}

func crearDotNodoTree(numInodo int, archivo *os.File, sblock superBloque) string {
	var inodoTemp Inodo
	archivo.Seek(int64(sblock.S_inode_start+int32(binary.Size(Inodo{}))*int32(numInodo)), 0)

	err := binary.Read(archivo, binary.LittleEndian, &inodoTemp)
	if err != nil {
		fmt.Println("Error al leer el inodo")
		return ""
	}

	Dot := "inodo" + strconv.Itoa(numInodo) + "[label = <\n"
	Dot += "<TABLE border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	Dot += "<tr><td bgcolor=\"lightgrey\" colspan=\"2\">Inodo" + strconv.Itoa(numInodo) + "</td></tr>\n"
	Dot += "<tr><td>i_uid</td><td>" + strconv.Itoa(int(inodoTemp.I_uid)) + "</td></tr>\n"
	Dot += "<tr><td>i_gid</td><td>" + strconv.Itoa(int(inodoTemp.I_gid)) + "</td></tr>\n"
	Dot += "<tr><td>i_size</td><td>" + strconv.Itoa(int(inodoTemp.I_s)) + "</td></tr>\n"
	Dot += "<tr><td>i_atime</td><td>" + string(inodoTemp.I_atime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_ctime</td><td>" + string(inodoTemp.I_ctime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_mtime</td><td>" + string(inodoTemp.I_mtime[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_type</td><td>" + string(inodoTemp.I_type[:]) + "</td></tr>\n"
	Dot += "<tr><td>i_perm</td><td>" + string(inodoTemp.I_perm[:]) + "</td></tr>\n"
	nodosDot := ""
	enlacesDot := ""
	for i, ptr := range inodoTemp.I_block {
		Dot += "<TR><TD bgcolor=\"lightgrey\">I_block[" + strconv.Itoa(i) + "]</TD><TD port='" + strconv.Itoa(i) + "'>" + strconv.Itoa(int(ptr)) + "</TD></TR>\n"
		if ptr != -1 {

			enlacesDot += "inodo" + strconv.Itoa(numInodo) + ":" + strconv.Itoa(i) + " -> bloque" + strconv.Itoa(int(ptr)) + ";\n"

			nodosDot += crearDotBloqueTree(int(ptr), string(inodoTemp.I_type[:]), archivo, sblock)
		}

	}
	Dot += "</TABLE>>];\n"
	Dot += nodosDot
	Dot += enlacesDot

	return Dot
}

func crearDotBloqueTree(ptr int, tipo string, archivo *os.File, sblock superBloque) string {
	Dot := ""
	if strings.Contains(tipo, "1") {
		var bloqueT bloqueArchivos
		archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)

		err := binary.Read(archivo, binary.LittleEndian, &bloqueT)

		if err != nil {
			fmt.Println("Error al leer el inodo")
			return ""
		}

		Dot += "bloque"
		Dot += strconv.Itoa(ptr)
		Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(ptr) + "</TD></TR>\n"
		cont := strings.TrimRight(string(bloqueT.B_content[:]), string(rune(0)))
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">" + cont + "</TD></TR>\n"

		Dot += "</TABLE>>];\n"

	} else {
		var bloqueT bloqueCarpeta
		archivo.Seek(int64(sblock.S_block_start+int32(binary.Size(bloqueArchivos{}))*int32(ptr)), 0)

		err := binary.Read(archivo, binary.LittleEndian, &bloqueT)

		if err != nil {
			fmt.Println("Error al leer el inodo")
			return ""
		}

		Dot += "bloque"
		Dot += strconv.Itoa(ptr)
		Dot += "[shape=none, color=lightgrey, label=<\n<TABLE cellspacing=\"3\" cellpadding=\"2\" style=\"rounded\" >\n"
		Dot += "<TR><TD bgcolor=\"lightgrey\" colspan=\"2\">bloque " + strconv.Itoa(ptr) + "</TD></TR>\n"

		Dot += "<TR><TD bgcolor=\"lightgrey\">b_name</TD><TD>b_inodo</TD></TR>\n"
		enlacesDot := ""
		nodosDot := ""
		for i, cont := range bloqueT.B_content {

			nam := strings.TrimRight(string(cont.B_name[:]), string(rune(0)))
			Dot += "<TR><TD bgcolor=\"lightgrey\">" + nam + "</TD><TD port= '" + strconv.Itoa(i) + "'>" + strconv.Itoa(int(cont.B_inodo)) + "</TD></TR>\n"

			if cont.B_inodo != -1 {
				if nam != "." && nam != ".." {
					nodosDot += crearDotNodoTree(int(cont.B_inodo), archivo, sblock)
					enlacesDot += "bloque" + strconv.Itoa(ptr) + ":" + strconv.Itoa(i) + " -> inodo" + strconv.Itoa(int(cont.B_inodo)) + ";\n"
				}
			}

		}
		Dot += "</TABLE>>];\n"
		Dot += nodosDot
		Dot += enlacesDot
	}

	return Dot
}

func crearArchivoDot(path string, Dot string) {

	dir, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error ruta archivo\n"
	}

	//Crea todos los directorios
	err = os.MkdirAll(filepath.Dir(dir), 0777)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error al crear los directorios fisicos\n"
	}

	//Crear el archivo .dot
	DotName := path[:len(path)-3]

	archivoDot, err := os.Create(DotName + "dot")
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		Salida_comando += "Error al crear el archivo .dot\n"
		return
	}

	defer archivoDot.Close()

	_, err = archivoDot.WriteString(Dot)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		Salida_comando += "Error al crear el archivo .dot\n"
		return
	}

	//Generar la imagen
	cmd := exec.Command("dot", "-T", "png", DotName+"dot", "-o", DotName+"png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar la imagen: ", err)
		Salida_comando += "Error al generar la imagen\n"
		return
	}

	fmt.Println("Reporte generado con exito")
	Salida_comando += "Reporte generado con exito\n"

}

func crearTxt(path string, cont string) {

	dir, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error ruta archivo\n"
	}

	//Crea todos los directorios
	err = os.MkdirAll(filepath.Dir(dir), 0777)

	if err != nil {
		fmt.Println(err)
		Salida_comando += "Error al crear directorios fisicos\n"
	}

	DotName := path[:len(path)-3]

	archivoDot, err := os.Create(DotName + "txt")
	if err != nil {
		fmt.Println("Error al crear el archivo .dot: ", err)
		Salida_comando += "Error al crear el archivo .dot\n"
		return
	}

	defer archivoDot.Close()

	_, err = archivoDot.WriteString(cont)
	if err != nil {
		fmt.Println("Error al escribir el archivo .dot: ", err)
		Salida_comando += "Error al escribir el archivo .dot\n"
		return
	}

	fmt.Println("Reporte generado con exito")
	Salida_comando += "Reporte generado con exito\n"
}

func obtenerArchivo_externo(path string) string {
	archivo, err := os.Open(path)

	if err != nil {
		fmt.Println("Error al abrir el archivo: ", err)
		return ""
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	lineas := ""
	for scanner.Scan() {
		lineas += scanner.Text()

	}

	return lineas
}

func quitarCaracteresEspeciales(cont string) string {

	result := cont

	for i := 0; i < 32; i++ {
		result = strings.ReplaceAll(result, string(rune(i)), "")
	}
	result = strings.ReplaceAll(result, string(rune(127)), "")
	result = strings.ReplaceAll(result, string(rune(255)), "")
	result = strings.ReplaceAll(result, "\xff", "")

	return result
}
