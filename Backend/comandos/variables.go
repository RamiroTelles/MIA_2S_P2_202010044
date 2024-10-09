package comandos

type Mount struct {
	Id        string
	Path      string
	Name      string
	Part_type [1]byte
	Start     int32
	Size      int32
	PartNum   int32
}

var particionesMontadas []Mount
var letra_disco_ascii int = 96
var uId int = -1
var gId int = -1
var actualIdMount string
var Salida_comando string

func VerificarParticionMontada(id string) int {
	//fmt.Println(id)

	for i := 0; i < len(particionesMontadas); i++ {
		//fmt.Println(particionesMontadas[i].Id)
		if particionesMontadas[i].Id == id {
			return i
		}
	}
	return -1

}
