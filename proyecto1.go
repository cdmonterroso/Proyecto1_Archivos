package main

import	(
	"io/ioutil"
	"os/exec"
	//"unsafe"
	"os"
	"fmt"
	"encoding/binary"
	"bytes"
	"bufio"
	"log"
	"strings"
	"strconv"
	//"time"
)

type Estructura struct{
	Tamanio int64 //8 bytes
	Estado byte //1 byte
	Nombre [8]byte //8 bytes
}

type Montaje struct{
	Id string
	Path string
	NombreParticion string
	Letra string
	Numero int
}

type Mbr struct{
	Tamanio int64
	Fecha [20]byte
	Disco int64
	Particion1 Particion
	Particion2 Particion
	Particion3 Particion
	Particion4 Particion
}

type Particion struct{
	Estado byte
	Tipo byte
	Ajuste byte
	Inicio int64
	Tamanio int64
	Nombre [16]byte
}

type Ebr struct{
	Estado byte
	Ajuste byte
	Inicio int64
	Tamanio int64
	Siguiente int64
	Nombre [16]byte
}

func crearArchivoBinario(tamanio int, ruta string, nombre string){
	//Creando carpeta si no existiera
	err := os.MkdirAll(ruta,0777)
	if err != nil{
		panic(err)
	}
	//CREANDO ARCHIVO BINARIO
	var direccion string
	direccion = (ruta + "/" + nombre)
	//Creando archivo binario
	file, er := os.Create(direccion)
	defer file.Close()
	if er != nil{
		fmt.Println("Error al crear archivo")
		log.Fatal(er)
	}
	//Creando archivo para reporte MBR
	fileR, e := os.Create("reporteMBR.dot")
	defer fileR.Close()
	if e != nil{
		fmt.Println("Error al crear reporte")
		log.Fatal(e)
	}

	//Creando archivo para reporte DISK
	fileD, eD := os.Create("reporteDISK.dot")
	defer fileD.Close()
	if eD != nil{
		fmt.Println("Error al crear reporte")
		log.Fatal(eD)
	}

	var caracter byte = 0;
	c := &caracter

	//Escribimos un 0 en el inicio del archivo.
	var binario1 bytes.Buffer
	binary.Write(&binario1, binary.BigEndian, c)
	file.Write(binario1.Bytes())
	//Nos posicionamos en el ultimo byte (primera posicion es 0)	
	file.Seek(int64(tamanio-1),0)

	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, c)
	file.Write(binario2.Bytes())

	//Escribir en un rango del archivo binario
	/*inicio := 20
	fin := 59
	for i:=inicio; i < fin; i++ {
		file.Seek(int64(i),0)
		var binTemp bytes.Buffer
		binary.Write(&binTemp,binary.BigEndian,byte(1))
		file.Write(binTemp.Bytes())
	}*/

	//RESERVANDO EL ESPACIO DEL EL MBR QUE CONTENDRA LAS 4 PARTICONES
	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.
	//AL INICIO EL MBR ESTARA EN BLANCO
	miMBR := Mbr{} 
	//tempParticion := Particion{}
	//fmt.Println("Tamaño del struct de una particion: ",binary.Size(tempParticion))
	//fmt.Println("Tamaño del mbr: ", binary.Size(miMBR))
	miMBR.Tamanio = int64(tamanio)
	miMBR.Disco = 15
	fechaCreacion := "15/09/2020_11:42:00"
	//fechaCreacion := time.Now()
	copy(miMBR.Fecha[0:],fechaCreacion)
	/*miMBR.Particion4 = Particion{Estado:1,Tipo:'p',Ajuste:'f',Inicio:176,Tamanio:20480}
	copy(miMBR.Particion4.Nombre[0:],"particion prueba")*/
	fmt.Println("MBR al inicio: ", miMBR)
	x := &miMBR
	var binario3 bytes.Buffer
	binary.Write(&binario3,binary.BigEndian, x)
	file.Write(binario3.Bytes())
	fmt.Println("Se creo el disco correctamente")
	
}

func escribirArchivoBinario(mbr Mbr, rutaArchivo string){
	//file,_ := os.Open(rutaArchivo)
	var file,_ = os.OpenFile(rutaArchivo, os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()

	//Escribir estructura en el archivo binario
	file.Seek(0, 0) // nos posicionamos en el inicio del archivo.
	
	x := &mbr
	var binario3 bytes.Buffer
	binary.Write(&binario3,binary.BigEndian, x)
	file.Write(binario3.Bytes())
	fmt.Println("Se actualizo el mbr... ", mbr)
	generarReporteMBR(mbr)
	generarReporteDISK(mbr)
}

func escribirEbr(ebr Ebr, rutaArchivo string){
	var file,_ = os.OpenFile(rutaArchivo,os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()

	file.Seek(ebr.Inicio,0)
	x := &ebr
	var binario4 bytes.Buffer
	binary.Write(&binario4,binary.BigEndian,x)
	file.Write(binario4.Bytes())
	fmt.Println("Se actualizo el ebr... ", ebr)
	//generarReporteEBR(ebr)
}

func generarReporteMBR(mbr Mbr){
	
	var nombreTemp string = ConvertirString(mbr.Particion1.Nombre[:])
	fmt.Println("Texto: ", nombreTemp)
	var strReporte string 
	fmt.Println("\n---REPORTE MBR---")
	fmt.Printf("Tamaño disco: %d\n",mbr.Tamanio)
	fmt.Printf("Numero de disco %d\n",mbr.Disco)
	fmt.Printf("Fecha de creacion %s\n",mbr.Fecha)
	fmt.Println("PARTICION 1")
	fmt.Printf("Ajuste: %d\n",mbr.Particion1.Ajuste)
	fmt.Printf("Estado: %d\n",mbr.Particion1.Estado)
	fmt.Printf("Inicio: %d\n",mbr.Particion1.Inicio)
	fmt.Printf("Nombre: %s\n",mbr.Particion1.Nombre)
	fmt.Printf("Tamaño: %d\n",mbr.Particion1.Tamanio)
	fmt.Printf("Tipo: %d\n",mbr.Particion1.Tipo)

	strReporte = "digraph t { tabla [ shape = plaintext\n color = black\n label=<\n"
	strReporte += "<table border='1' cellborder='1'>"
	strReporte += ("<tr><td colspan='2'> MBR </td></tr>")
	strReporte += ("<tr><td> Tamaño disco: </td> <td> " + strconv.Itoa(int(mbr.Tamanio)) + "</td></tr>" )
	strReporte += ("<tr><td> Numero disco: </td> <td> " + strconv.Itoa(int(mbr.Disco)) + "</td></tr>" )
	strReporte += ("<tr><td> Fecha disco: </td> <td> " + ConvertirString(mbr.Fecha[:]) + "</td></tr>" ) 
	strReporte += "<tr><td colspan='2'> PARTICION 1 </td></tr>"
	strReporte += ("<tr><td> Ajuste: </td> <td> " + strconv.Itoa(int(mbr.Particion1.Ajuste)) + "</td></tr>" )
	strReporte += ("<tr><td> Estado: </td> <td> " + strconv.Itoa(int(mbr.Particion1.Estado)) + "</td></tr>" )
	strReporte += ("<tr><td> Inicio: </td> <td> " + strconv.Itoa(int(mbr.Particion1.Inicio)) + "</td></tr>" ) 
	strReporte += ("<tr><td> Nombre: </td> <td> " + ConvertirString(mbr.Particion1.Nombre[0:]) + "</td></tr>" )
	strReporte += ("<tr><td> Tamaño: </td> <td> " + strconv.Itoa(int(mbr.Particion1.Tamanio)) + "</td></tr>" )
	strReporte += ("<tr><td> Tipo: </td> <td> " + strconv.Itoa(int(mbr.Particion1.Tipo)) + "</td></tr>" ) 
	strReporte += "<tr><td colspan='2'> PARTICION 2 </td></tr>"
	strReporte += ("<tr><td> Ajuste: </td> <td> " + strconv.Itoa(int(mbr.Particion2.Ajuste)) + "</td></tr>" )
	strReporte += ("<tr><td> Estado: </td> <td> " + strconv.Itoa(int(mbr.Particion2.Estado)) + "</td></tr>" )
	strReporte += ("<tr><td> Inicio: </td> <td> " + strconv.Itoa(int(mbr.Particion2.Inicio)) + "</td></tr>" ) 
	strReporte += ("<tr><td> Nombre: </td> <td> " + ConvertirString(mbr.Particion2.Nombre[0:]) + "</td></tr>" )
	strReporte += ("<tr><td> Tamaño: </td> <td> " + strconv.Itoa(int(mbr.Particion2.Tamanio)) + "</td></tr>" )
	strReporte += ("<tr><td> Tipo: </td> <td> " + strconv.Itoa(int(mbr.Particion2.Tipo)) + "</td></tr>" ) 
	strReporte += "<tr><td colspan='2'> PARTICION 3 </td></tr>"
	strReporte += ("<tr><td> Ajuste: </td> <td> " + strconv.Itoa(int(mbr.Particion3.Ajuste)) + "</td></tr>" )
	strReporte += ("<tr><td> Estado: </td> <td> " + strconv.Itoa(int(mbr.Particion3.Estado)) + "</td></tr>" )
	strReporte += ("<tr><td> Inicio: </td> <td> " + strconv.Itoa(int(mbr.Particion3.Inicio)) + "</td></tr>" ) 
	strReporte += ("<tr><td> Nombre: </td> <td> " + ConvertirString(mbr.Particion3.Nombre[0:]) + "</td></tr>" )
	strReporte += ("<tr><td> Tamaño: </td> <td> " + strconv.Itoa(int(mbr.Particion3.Tamanio)) + "</td></tr>" )
	strReporte += ("<tr><td> Tipo: </td> <td> " + strconv.Itoa(int(mbr.Particion3.Tipo)) + "</td></tr>" ) 
	strReporte += "<tr><td colspan='2'> PARTICION 4 </td></tr>"
	strReporte += ("<tr><td> Ajuste: </td> <td> " + strconv.Itoa(int(mbr.Particion4.Ajuste)) + "</td></tr>" )
	strReporte += ("<tr><td> Estado: </td> <td> " + strconv.Itoa(int(mbr.Particion4.Estado)) + "</td></tr>" )
	strReporte += ("<tr><td> Inicio: </td> <td> " + strconv.Itoa(int(mbr.Particion4.Inicio)) + "</td></tr>" ) 
	strReporte += ("<tr><td> Nombre: </td> <td> " + ConvertirString(mbr.Particion4.Nombre[0:]) + "</td></tr>" )
	strReporte += ("<tr><td> Tamaño: </td> <td> " + strconv.Itoa(int(mbr.Particion4.Tamanio)) + "</td></tr>" )
	strReporte += ("<tr><td> Tipo: </td> <td> " + strconv.Itoa(int(mbr.Particion4.Tipo)) + "</td></tr>" ) 
	
	strReporte += "</table> >]; }"
	
	byteReporte := make([]byte,len(strReporte))
	copy(byteReporte[0:],strReporte)
	fmt.Println("strReporte: ",strReporte)
	err := ioutil.WriteFile("reporteMBR.dot", byteReporte, 0644)
    if err != nil {
        panic(err)
	}
	cmd := exec.Command("dot", "-Tpng", "/home/daniel/Documentos/Archivos/Proyecto1/reporteMBR.dot","-o","reporteMBR.png")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	er := cmd.Run()
	if er != nil{
		os.Stderr.WriteString(err.Error())
	}
	//exec.Command("dot -Tpng reporteMBR.dot -o reporteMBR.png")
	/*comando := exec.Command("dot -Tpng reporteMBR.dot -o reporteMBR.png")
	if err := comando.Run(); err !=nil{
		fmt.Println("Se creo el reporte")
	}*/
	exec.Command("display reporteMBR.png &").Output()

}

func generarReporteDISK(mbr Mbr){
	var strReporte string
	strReporte = "digraph t { tabla [ shape = plaintext\n color = black\n label=<\n"
	strReporte += "<table border='1' cellborder='1'>"
	strReporte += ("<tr><td> MBR </td>")
	var nuevoNombre [16]byte
	if mbr.Particion1.Nombre == nuevoNombre { //No tiene nombre
		strReporte += "<td> LIBRE </td>"
	}else{ //Si tiene nombre
		var tipo string
		if mbr.Particion1.Tipo == 'e'{
			tipo = "e"
		}else if mbr.Particion1.Tipo == 'p'{
			tipo = "p"
		}
		strReporte += "<td> " + ConvertirString(mbr.Particion1.Nombre[0:]) + ", " + tipo + " </td>"
	}
	if mbr.Particion2.Nombre == nuevoNombre {
		strReporte += "<td> LIBRE </td>"
	}else{
		var tipo string
		if mbr.Particion2.Tipo == 'e'{
			tipo = "e"
		}else if mbr.Particion2.Tipo == 'p'{
			tipo = "p"
		}
		strReporte += "<td> " + ConvertirString(mbr.Particion2.Nombre[0:]) + ", " + tipo + " </td>"
	}
	if mbr.Particion3.Nombre == nuevoNombre {
		strReporte += "<td> LIBRE </td>"
	}else{
		var tipo string
		if mbr.Particion3.Tipo == 'e'{
			tipo = "e"
		}else if mbr.Particion3.Tipo == 'p'{
			tipo = "p"
		}
		strReporte += "<td> " + ConvertirString(mbr.Particion3.Nombre[0:]) + ", " + tipo + " </td>"
	}
	if mbr.Particion4.Nombre == nuevoNombre {
		strReporte += "<td> LIBRE </td>"
	}else{
		var tipo string
		if mbr.Particion4.Tipo == 'e'{
			tipo = "e"
		}else if mbr.Particion4.Tipo == 'p'{
			tipo = "p"
		}
		strReporte += "<td> " + ConvertirString(mbr.Particion4.Nombre[0:]) + ", " + tipo + " </td>"
	}
	strReporte += "</tr></table> >]; }"
	
	byteReporte := make([]byte,len(strReporte))
	copy(byteReporte[0:],strReporte)
	fmt.Println("strReporte: ",strReporte)
	err := ioutil.WriteFile("reporteDISK.dot", byteReporte, 0644)
    if err != nil {
        panic(err)
	}
	cmd := exec.Command("dot", "-Tpng", "/home/daniel/Documentos/Archivos/Proyecto1/reporteDISK.dot","-o","reporteDISK.png")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	er := cmd.Run()
	if er != nil{
		os.Stderr.WriteString(err.Error())
	}
}

func leerArchivoEntrada(){
	lectura := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese comando... ")
	comando,_ := lectura.ReadString('\n')
	var listaComando []string
	listaComando = strings.Split(comando, " ")
	//EJECUTA EL COMANDO INGRESADO (ABRIR EL ARCHIVO DE ENTRADA)
	var path string
	if listaComando[0] == "exec"{
		path = strings.ReplaceAll(listaComando[1],"-path->","")
		path = strings.ReplaceAll(path,"\n","")
	}else{
		fmt.Println("Comando incorrecto")
	}
	file, er := os.Open(path)
	//defer file.Close()
	if er != nil{
		log.Fatal(er)
		fmt.Println("Error al abrir archivo")
	}
	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan(){
		i++
		linea := scanner.Text()
		//fmt.Println("Comando: " + linea)
		ejecutarComando(linea)
	}
	file.Close()
}

func leerArchivoBinario(archivo string) Mbr{
	//fmt.Println("Iniciando lectura del archivo binario")
	file,err := os.Open(archivo)
	if err != nil {
		log.Fatal("Error al abrir archivo binario ", err)
	}
	defer file.Close()
	m := Mbr{} //Variable de tipo MBR
	size	 := binary.Size(m) //Obteniendo tamaño del mbr
	fmt.Println("Tamaño de la variable tipo mbr ", size)
	//Lee la cantidad de (size) bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en buffer, necesario para codificar binario
	buffer := bytes.NewBuffer(data)
	fmt.Println("buffer: ", buffer)
	//Decodificamos y guardamos en variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("Error en lectura ", err)
	}
	//fmt.Println("Se decodifico el archivo binario")

	//Se imprimen los valores guardados en el struct
	fmt.Println("Decodificado ",m)
	//fmt.Printf("Texto %s\n",m.Particion1.Nombre)
	/*fmt.Println("Tamaño del disco: ", m.Tamanio)
	fmt.Println("Fecha en que se creo el disco: ", m.Fecha)
	fmt.Println("Numero del disco: ", m.Disco)*/
	return m
}

func ConvertirString(c []byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func leerBytes(file *os.File, numero int) []byte {
	//fmt.Println("Iniciando lectura de archivo binario")
	bytes := make([]byte, numero) //array de bytes
	_, er := file.Read(bytes)
	if er != nil{
		log.Fatal("No se pudo leer el archivo binario ", er)
	}
	//fmt.Println("Se leyo el archivo binario")
	return bytes
}

func interpretar(lineaComando string){
	/*finalizar := 0
	//Leyendo el comando para abrir el archvio de entrada

	if comando == "x\n"{
		finalizar = 1
	}else{
		if comando != "" {
			var listaComando []string
			listaComando = strings.Split(comando, " ")
			//EJECUTA EL COMANDO INGRESADO (ABRIR EL ARCHIVO DE ENTRADA)
			ejecutarComando(listaComando)

		}
	}

	//Leyendo los comandos del archvio de entrada
	for finalizar != 1 {

	}*/
}

func ejecutarComando(lineaComando string) {
	fmt.Println()
	commandArray := strings.Split(lineaComando," ")
	data := strings.ToLower(commandArray[0])
	var tamanio, ruta, nombre, unidad string
	var tipo, ajuste, eliminar, agregar string
	//COMANDO PARA CREAR EL ARCHIVO BINARIO
	if data == "mkdisk" {
		fmt.Println(commandArray)
		indice := 0
		indice++
		//fmt.Println("Comando para crear el archivo binario")
		//Tamaño del disco
		data = strings.ToLower(commandArray[indice])
		if strings.Contains(data,"size"){
			tamanio = strings.ReplaceAll(data, "-size->", "")
			indice++
			fmt.Println("Tamaño del disco: " + tamanio)
		}
		//Ruta del disco
		data = strings.ToLower(commandArray[indice])
		if strings.Contains(data, "path"){
			ruta = strings.ReplaceAll(data, "-path->", "")
			indice++
			if strings.Contains(ruta, "\""){
				var rutaTemp string
				rutaTemp = strings.ToLower(commandArray[indice])
				for !(strings.Contains(rutaTemp,"\"")){
					ruta = strings.Join([]string{ruta, rutaTemp}, " ")
					//ruta = ruta + " " + rutaTemp
					indice ++
					rutaTemp = strings.ToLower(commandArray[indice])
				}
				//ruta = strings.Join([]string{ruta, rutaTemp}, " ")
				ruta = ruta + " " + rutaTemp
				indice++
			}
			fmt.Println("Ruta del disco: " + ruta)
		}
		//Nombre del disco
		data = strings.ToLower(commandArray[indice])
		if strings.Contains(data, "name"){
			nombre = strings.ReplaceAll(data, "-name->", "")
			indice++
			fmt.Println("Nombre del disco: " + nombre)
		}
		//Unidades del disco
		conversion,_ := strconv.Atoi(tamanio)
		//fmt.Println(len(commandArray))
		tamArray := len(commandArray)
		commandArray[(tamArray-1)] = strings.ToLower(commandArray[(tamArray-1)])
		if strings.Contains(commandArray[(tamArray-1)],"unit"){ //Tiene unidad
			//KB o MB
			data = strings.ToLower(commandArray[indice])
			unidad = strings.ReplaceAll(data, "-unit->", "")
			fmt.Println("Unidades del disco: " + unidad)
			if(unidad == "k"){
				crearArchivoBinario((conversion*1024),strings.ReplaceAll(ruta,"\"", ""),nombre)
			}else if (unidad == "m"){
				crearArchivoBinario((conversion*1024*1024),strings.ReplaceAll(ruta,"\"", ""),nombre)
			}
		}else{
			//MB
			crearArchivoBinario((conversion*1024*1024),strings.ReplaceAll(ruta,"\"", ""),nombre)
		}
	}else if data == "pause"{ //PAUSA
		lectura := bufio.NewReader(os.Stdin)
		fmt.Println("Enter para continuar... ")
		comando,_ := lectura.ReadString('\n')
		if comando == ""{
			//Continuar...
		}
	}else if data == "rmdisk" { //COMANDO PARA ELIMINAR UN DISCO
		fmt.Println(commandArray)
		indice := 0
		indice++
		fmt.Println("Comando para eliminar el archivo binario")
		data = strings.ToLower(commandArray[indice])
		if strings.Contains(data, "path"){
			ruta = strings.ReplaceAll(data, "-path->", "")
			indice++
			if strings.Contains(ruta, "\""){
				var rutaTemp string
				rutaTemp = strings.ToLower(commandArray[indice])
				for !(strings.Contains(rutaTemp,"\"")){
					ruta = strings.Join([]string{ruta, rutaTemp}, " ")
					//ruta = ruta + " " + rutaTemp
					indice ++
					rutaTemp = strings.ToLower(commandArray[indice])
				}
				//ruta = strings.Join([]string{ruta, rutaTemp}, " ")
				ruta = ruta + " " + rutaTemp
				indice++
			}
			fmt.Println("Ruta del disco a eliminar: " + ruta)
			err := os.Remove(strings.ReplaceAll(ruta,"\"", ""))
			if err != nil {
				fmt.Printf("Error eliminando archivo: %v\n", err)
			} else {
				fmt.Println("Eliminado correctamente")
			}
		}
	}else if(data == "fdisk"){ //COMANDO PARA CREAR O ELIMINAR UNA PARTICION
		fmt.Println(commandArray)
		//data = commandArray[1]
		for i := 0; i < len(commandArray); i++ {
			if strings.Contains(strings.ToLower(commandArray[i]),"delete"){ //Se eliminara una particion
				fmt.Println("Se eliminara una particion")
				for j := 1; j < len(commandArray); j++ {
					data = strings.ToLower(commandArray[j])
					if strings.Contains(data, "delete"){
						eliminar = strings.ReplaceAll(data, "-delete->", "")
						fmt.Println("Tipo de eliminacion: " + eliminar)
					}
					if strings.Contains(data, "path"){
						ruta = strings.ReplaceAll(data, "-path->", "")
						fmt.Println("Ruta de la particion a eliminar: " + ruta)
					}
					if strings.Contains(data, "path"){
						ruta = strings.ReplaceAll(data, "-path->", "")
						if strings.Contains(ruta, "\""){
							j++
							var rutaTemp string
							rutaTemp = strings.ToLower(commandArray[j])
							for !(strings.Contains(rutaTemp,"\"")){ //No viene comillas (true)
								ruta = strings.Join([]string{ruta, rutaTemp}, " ")
								//ruta = ruta + " " + rutaTemp
								j++
								rutaTemp = strings.ToLower(commandArray[j])
							}
							//ruta = strings.Join([]string{ruta, rutaTemp}, " ")
							ruta = ruta + " " + rutaTemp
						}
						fmt.Println("Ruta de la particion a eliminar: " + ruta)
					}
					if strings.Contains(data, "name"){
						nombre = strings.ReplaceAll(data, "-name->", "")
						fmt.Println("Nombre de la particion a eliminar: " + nombre)
					}
					
				}
				eliminarParticion(eliminar,ruta,nombre)
			}
			if strings.Contains(strings.ToLower(commandArray[i]),"size"){ //Se creara una particion
				fmt.Println("Se creara una particion")
				for k := 0; k < len(commandArray); k++ {
					data = strings.ToLower(commandArray[k])
					if strings.Contains(data,"size"){
						tamanio = strings.ReplaceAll(data, "-size->", "")
						fmt.Println("Tamanio del disco: " + tamanio)
					}
					if strings.Contains(data, "type"){
						tipo = strings.ReplaceAll(data, "-type->", "")
						fmt.Println("Tipo de particion: " + tipo)
					}
					if strings.Contains(data, "unit"){
						unidad = strings.ReplaceAll(data, "-unit->", "")
						fmt.Println("Unidades de la particion: " + unidad)
					}
					if strings.Contains(data, "fit"){
						ajuste = strings.ReplaceAll(data, "-fit->", "")
						fmt.Println("Ajuste de la particion: " + ajuste)
					}
					if strings.Contains(data, "path"){
						ruta = strings.ReplaceAll(data, "-path->", "")
						if strings.Contains(ruta, "\""){
							k++
							var rutaTemp string
							rutaTemp = strings.ToLower(commandArray[k])
							for !(strings.Contains(rutaTemp,"\"")){ //No viene comillas (true)
								ruta = strings.Join([]string{ruta, rutaTemp}, " ")
							
								k++
								rutaTemp = strings.ToLower(commandArray[k])
							}
							ruta = ruta + " " + rutaTemp
						}
						fmt.Println("Ruta de la particion a crear: " + ruta)
					}
					if strings.Contains(data, "name"){
						nombre = strings.ReplaceAll(data, "-name->", "")
						if strings.Contains(ruta, "\""){
							k++
							var nombreTemp string
							nombreTemp = commandArray[k]
							for !(strings.Contains(nombreTemp,"\"")){ //No viene comillas (true)
								nombre = strings.Join([]string{nombre, nombreTemp}, " ")
								k++
								nombreTemp = commandArray[k]
							}
							nombre = nombre + " " + nombreTemp
						}
						fmt.Println("Nombre de la particion a crear: " + nombre)
					}
				}
				crearParticion(nombre, ruta, tamanio, tipo, unidad, ajuste)
			}
			if strings.Contains(strings.ToLower(commandArray[i]), "add"){ //Se agregara espacio a la particion
				fmt.Println("Se añadira espacio a la particion")
				for k := 0; k < len(commandArray); k++ {
					data = strings.ToLower(commandArray[k])
					if strings.Contains(data,"add"){
						agregar = strings.ReplaceAll(data, "-add->", "")
						fmt.Println("Tamanio a agregar: " + agregar)
					}
					if strings.Contains(data, "unit"){
						unidad = strings.ReplaceAll(data, "-unit->", "")
						fmt.Println("Unidades de la particion: " + unidad)
					}
					if strings.Contains(data, "path"){
						ruta = strings.ReplaceAll(data, "-path->", "")
						if strings.Contains(ruta, "\""){
							k++
							var rutaTemp string
							rutaTemp = strings.ToLower(commandArray[k])
							for !(strings.Contains(rutaTemp,"\"")){ //No viene comillas (true)
								ruta = strings.Join([]string{ruta, rutaTemp}, " ")
							
								k++
								rutaTemp = strings.ToLower(commandArray[k])
							}
							ruta = ruta + " " + rutaTemp
						}
						fmt.Println("Ruta de la particion: " + ruta)
					}
					if strings.Contains(data, "name"){
						nombre = strings.ReplaceAll(data, "-name->", "")
						if strings.Contains(ruta, "\""){
							k++
							var nombreTemp string
							nombreTemp = commandArray[k]
							for !(strings.Contains(nombreTemp,"\"")){ //No viene comillas (true)
								nombre = strings.Join([]string{nombre, nombreTemp}, " ")
								k++
								nombreTemp = commandArray[k]
							}
							nombre = nombre + " " + nombreTemp
						}
						fmt.Println("Nombre de la particion: " + nombre)
					}
				}
				agregarEspacio(nombre, ruta, agregar, unidad)
			}
		}

	}
}

func crearParticion(nombre string, ruta string, tamanio string, tipo string, unidad string, ajuste string){
	fmt.Println()
	if tipo == "l"{
		fmt.Println("Particion logica...")
	}else{
		if nombre == "" || ruta == "" || tamanio == "" {
			fmt.Println("Faltan datos, no se creo la particion")
		}else{
			particionNueva := Particion{Estado:1}
	
			//Asignando tipo al struct
			//Validar que solo haya una particion extendida en un disco
			if tipo == "p" { //Primaria
				particionNueva.Tipo = 'p'
			}else if tipo == "e" { //Extendida
				//SOLO PUEDE HABER 1 EXTENDIDA
				particionNueva.Tipo = 'e'
			}else{
				//VALIDAR QUE NO ESTEN YA LAS 4 PARTICIONES POR DISCO
				particionNueva.Tipo = 'p'
			}
	
			//Asignando ajuste al struct
			if ajuste == "bf" {
				particionNueva.Ajuste = 'b'
			}else if ajuste == "ff" {
				particionNueva.Ajuste = 'f'
			}else if ajuste == "wf"{
				particionNueva.Ajuste = 'w'
			}else if ajuste == ""{
				particionNueva.Ajuste = 'w'
			}
	
			//Asignando tamaño al struct
			tamConv,_ := strconv.Atoi(tamanio)
			if unidad == "b" {
	
			}else if unidad == "k" {
				tamConv = (tamConv * 1024)
			}else if unidad == "m" {
				tamConv = (tamConv *1024 *1024)
			}else if unidad == "" {
				tamConv = (tamConv *1024)
			}else{
				fmt.Println("Unidad incorrecta")
			}
			particionNueva.Tamanio = int64(tamConv)
	
			//Asignando nombre al struct
			if nombre == "" {
				fmt.Println("Nombre incorrecto")
			}else{
				copy(particionNueva.Nombre[0:],nombre)
			}
			fmt.Println("Particion a crear ",particionNueva)
	
			//Asignando inicio al struct
			//Si el disco esta vacio se coloca en despues del mbr
			fmt.Println("...")
			mbrTemp := leerArchivoBinario(ruta)
			tamanioMBR := binary.Size(mbrTemp)
			fmt.Println("Tamaño del mbr leido: ", tamanioMBR)
			//IMPRIMIR LOS DATOS LEIDOS PARA COMPROBAR QUE SE LEYO CORRECTAMENTE
			fmt.Println("Imprimiendo mbr leido: ", mbrTemp)
	
			//1. Validar que solo haya una particion extendida en un disco
			//2. SI EL DISCO ESTA VACIO, SE INSERTA DESPUES DEL MBR
			//3. Si el disco tiene particiones, se coloca despues de la ulitma particion (sin pasar de 4 particiones)
			if (particionNueva.Tipo=='e') && (mbrTemp.Particion1.Tipo=='e' || mbrTemp.Particion2.Tipo=='e' || mbrTemp.Particion3.Tipo=='e' || mbrTemp.Particion4.Tipo=='e') {
				fmt.Println("Ya existe una particion extendida en el disco")
			}else{
				//NO DEBE REPETIRSE EL NOMBRE DE LAS PARTICIONES
				if(mbrTemp.Particion1.Nombre==particionNueva.Nombre || mbrTemp.Particion2.Nombre==particionNueva.Nombre || mbrTemp.Particion3.Nombre==particionNueva.Nombre || mbrTemp.Particion4.Nombre==particionNueva.Nombre){
					fmt.Printf("Ya existe una particion con este nombre: %s\n", particionNueva.Nombre)
				}else{
					//DEBE HABER ESPACIO DISPONIBLE EN EL DISCO
					// tam total del disco - tamMBR -tamP1 - tamP2 -tamP3 - tamP4 = espacio libre //bytes
					// espacio libre >= tam particion nueva //bytes
					//espacioLibre := mbrTemp.Tamanio - int64(binary.Size(mbrTemp) - mbrTemp.Particion1.Tamanio - mbrTemp.Particion2.Tamanio - mbrTemp.Particion3.Tamanio - mbrTemp.Particion4.Tamanio
					
					if(mbrTemp.Particion1.Estado == 0){ //La particion 1 esta vacia
						//Pero hay que preguntar si hay particiones llenas a la par ->, quiere decir que se elimino la particion 1 anteriormente
						if mbrTemp.Particion2.Estado != 0 || mbrTemp.Particion3.Estado != 0 || mbrTemp.Particion4.Estado != 0{
							//Obtener 
							inicioEspacioLibre := tamanioMBR+1
							var finEspacioLibre int64
							if mbrTemp.Particion2.Estado != 0 {
								//Espacio libre en particion 1
								finEspacioLibre = mbrTemp.Particion2.Inicio-1
								inicioParticion := inicioEspacioLibre
								tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
								//Preguntar si la nueva particion cabe en ese espacio libre
								if tamLibre >= particionNueva.Tamanio {
									//Se inserta la particion
									particionNueva.Inicio = int64(inicioParticion)
									mbrTemp.Particion1 = particionNueva
									escribirArchivoBinario(mbrTemp,ruta)
									//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
								}else{
									fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
								}
							}else if mbrTemp.Particion2.Estado == 0 && mbrTemp.Particion3.Estado != 0 {
								//Espacio libre en particion 1 y 2
								finEspacioLibre = mbrTemp.Particion3.Inicio-1
								inicioParticion := inicioEspacioLibre
								tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
								//Preguntar si la nueva particion cabe en ese espacio libre
								if tamLibre >= particionNueva.Tamanio {
									//Se inserta la particion
									particionNueva.Inicio = int64(inicioParticion)
									mbrTemp.Particion1 = particionNueva
									escribirArchivoBinario(mbrTemp,ruta)
									//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
								}else{
									fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
								}
							}else if mbrTemp.Particion2.Estado == 0 && mbrTemp.Particion3.Estado == 0 && mbrTemp.Particion4.Estado != 0 {
								//Espacio libre en particion 1, 2 y 3
								finEspacioLibre = mbrTemp.Particion4.Inicio-1
								inicioParticion := inicioEspacioLibre
								tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
								//Preguntar si la nueva particion cabe en ese espacio libre
								if tamLibre >= particionNueva.Tamanio {
									//Se inserta la particion
									particionNueva.Inicio = int64(inicioParticion)
									mbrTemp.Particion1 = particionNueva
									escribirArchivoBinario(mbrTemp,ruta)
									//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
								}else{
									fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
								}
							}
						}else{
							//EL DISCO ESTA VACIO, SE PUEDE INSERTAR EN PARTICION 1
							//VALIDAR QUE EL TAMAÑO DE LA PARTICION A CREAR NO SEA MAYOR AL ESPACIO RESTANTE DEL DISCO
							tamanioDisco := mbrTemp.Tamanio - int64(tamanioMBR)
							if tamanioDisco < particionNueva.Tamanio {
								fmt.Println("La particion no cabe en el disco")
							}else{
								inicioParticion := tamanioMBR+1
								particionNueva.Inicio = int64(inicioParticion)
								fmt.Println("La particion 1 esta vacia, inicia en: ", inicioParticion)
								//Coloco la nueva particion en la particion 1
								mbrTemp.Particion1 = particionNueva
								//Actualizo el mbr
								fmt.Println("Nuevo mbr con informacion actualizada: ", mbrTemp)
								//fmt.Println("Particion a guardar: ", mbrTemp.Particion1)
								escribirArchivoBinario(mbrTemp, ruta)
							}
						}
					}else if(mbrTemp.Particion2.Estado == 0){ //La particion 2 esta vacia y la 1 ocupada
						//Pero hay que preguntar si hay particiones llenas a la par ->, quiere decir que se elimino la particion 2 anteriormente
						if mbrTemp.Particion3.Estado != 0 || mbrTemp.Particion4.Estado != 0{
							//Obtener 
							inicioEspacioLibre := int64(tamanioMBR+1) + mbrTemp.Particion1.Tamanio
							var finEspacioLibre int64
							if mbrTemp.Particion3.Estado != 0 {
								//Espacio libre solamente en particion 2
								finEspacioLibre = mbrTemp.Particion3.Inicio-1
								inicioParticion := inicioEspacioLibre
								tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
								//Preguntar si la nueva particion cabe en ese espacio libre
								if tamLibre >= particionNueva.Tamanio {
									//Se inserta la particion
									particionNueva.Inicio = int64(inicioParticion)
									mbrTemp.Particion2 = particionNueva
									escribirArchivoBinario(mbrTemp,ruta)
									//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
								}else{
									fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
								}
							}else if mbrTemp.Particion3.Estado == 0 && mbrTemp.Particion4.Estado != 0 {
								//Espacio libre en particion 2 y 3
								finEspacioLibre = mbrTemp.Particion4.Inicio-1
								inicioParticion := inicioEspacioLibre
								tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
								//Preguntar si la nueva particion cabe en ese espacio libre
								if tamLibre >= particionNueva.Tamanio {
									//Se inserta la particion
									particionNueva.Inicio = int64(inicioParticion)
									mbrTemp.Particion2 = particionNueva
									escribirArchivoBinario(mbrTemp,ruta)
									//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
								}else{
									fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
								}
							}
						}else{
							//SE PUEDE INSERTAR EN PARTICION 2
							//VALIDAR QUE EL TAMAÑO DE LA PARTICION A CREAR NO SEA MAYOR AL ESPACIO RESTANTE DEL DISCO
							tamanioDisco := mbrTemp.Tamanio - int64(tamanioMBR) - mbrTemp.Particion1.Tamanio
							if tamanioDisco < particionNueva.Tamanio {
								fmt.Println("La particion no cabe en el disco")
							}else{
								inicioParticion := tamanioMBR+1 + int(mbrTemp.Particion1.Tamanio)
								particionNueva.Inicio = int64(inicioParticion)
								fmt.Println("La particion 2 esta vacia, inicia en: ", inicioParticion)
								//Coloco la nueva particion en la particion 1
								mbrTemp.Particion2 = particionNueva
								//Actualizo el mbr
								fmt.Println("Nuevo mbr con informacion actualizada: ", mbrTemp)
								//fmt.Println("Particion a guardar: ", mbrTemp.Particion2)
								escribirArchivoBinario(mbrTemp, ruta)
							}
						}
					}else if(mbrTemp.Particion3.Inicio == 0){ //La particion 3 esta vacia, la 1 y 2 ocupada
						//Pero hay que preguntar si hay particiones llenas a la par ->, quiere decir que se elimino la particion 2 anteriormente
						if mbrTemp.Particion4.Estado != 0{
							//Obtener 
							inicioEspacioLibre := int64(tamanioMBR+1) + mbrTemp.Particion1.Tamanio + mbrTemp.Particion2.Tamanio
							var finEspacioLibre int64
							//Espacio libre solamente en particion 3
							finEspacioLibre = mbrTemp.Particion4.Inicio-1
							inicioParticion := inicioEspacioLibre
							tamLibre := (finEspacioLibre-int64(inicioEspacioLibre))
							//Preguntar si la nueva particion cabe en ese espacio libre
							if tamLibre >= particionNueva.Tamanio {
								//Se inserta la particion
								particionNueva.Inicio = int64(inicioParticion)
								mbrTemp.Particion3 = particionNueva
								escribirArchivoBinario(mbrTemp,ruta)
								//EL ESPACIO QUE SOBRA SE CONVIERTE EN FRAGMENTACION
	
							}else{
								fmt.Println("El tamaño de la particion a crear es mayor al espacio libre") //3)
							}
						}else{
							//SE PUEDE INSERTAR EN PARTICION 3
							//VALIDAR QUE EL TAMAÑO DE LA PARTICION A CREAR NO SEA MAYOR AL ESPACIO RESTANTE DEL DISCO
							tamanioDisco := mbrTemp.Tamanio - int64(tamanioMBR) - mbrTemp.Particion1.Tamanio - mbrTemp.Particion2.Tamanio
							if tamanioDisco < particionNueva.Tamanio {
								fmt.Println("La particion no cabe en el disco")
							}else{
								inicioParticion := tamanioMBR+1 + int(mbrTemp.Particion1.Tamanio) + int(mbrTemp.Particion2.Tamanio)
								particionNueva.Inicio = int64(inicioParticion)
								fmt.Println("La particion 3 esta vacia, inicia en: ", inicioParticion)
								mbrTemp.Particion3 = particionNueva
								fmt.Println("Nuevo mbr con informacion actualizada: ", mbrTemp)
								//fmt.Println("Particion a guardar: ", mbrTemp.Particion3)
								escribirArchivoBinario(mbrTemp, ruta)
							}
						}
					}else if(mbrTemp.Particion4.Inicio == 0){ //La particion 4 es la unica disponible
						//SE PUEDE INSERTAR EN PARTICION 4
						//VALIDAR QUE EL TAMAÑO DE LA PARTICION A CREAR NO SEA MAYOR AL ESPACIO RESTANTE DEL DISCO
						tamanioDisco := mbrTemp.Tamanio - int64(tamanioMBR) - mbrTemp.Particion1.Tamanio - mbrTemp.Particion2.Tamanio - mbrTemp.Particion3.Tamanio
						if tamanioDisco < particionNueva.Tamanio {
							fmt.Println("La particion no cabe en el disco")
						}else{
							inicioParticion := tamanioMBR+1 + int(mbrTemp.Particion1.Tamanio) + int(mbrTemp.Particion2.Tamanio) + int(mbrTemp.Particion3.Tamanio)
							particionNueva.Inicio = int64(inicioParticion)
							fmt.Println("La particion 4 esta vacia, inicia en: ", inicioParticion)
							mbrTemp.Particion4 = particionNueva
							fmt.Println("Nuevo mbr con informacion actualizada: ", mbrTemp)
							//fmt.Println("Particion a guardar: ", mbrTemp.Particion3)
							escribirArchivoBinario(mbrTemp, ruta)
						}
					}else{
						fmt.Println("Disco lleno")
					}
				}
				
			}
		}
	}
		
}

func eliminarParticion(eliminar string,ruta string,nombre string){
	fmt.Println("...")
	mbrTemp := leerArchivoBinario(ruta)
	tamanioMBR := binary.Size(mbrTemp)
	fmt.Println("Tamaño del mbr leido: ", tamanioMBR)
	//IMPRIMIR LOS DATOS LEIDOS PARA COMPROBAR QUE SE LEYO CORRECTAMENTE
	fmt.Println("Imprimiendo mbr leido: ", mbrTemp)

	nombrePartcionTemp1 := ConvertirString(mbrTemp.Particion1.Nombre[0:])
	nombrePartcionTemp2 := ConvertirString(mbrTemp.Particion2.Nombre[0:])
	nombrePartcionTemp3 := ConvertirString(mbrTemp.Particion3.Nombre[0:])
	nombrePartcionTemp4 := ConvertirString(mbrTemp.Particion4.Nombre[0:])
	if nombrePartcionTemp1 == nombre {
		mbrTemp.Particion1.Ajuste = 0
		mbrTemp.Particion1.Estado = 0
		mbrTemp.Particion1.Inicio = 0
		mbrTemp.Particion1.Tamanio = 0
		mbrTemp.Particion1.Tipo = 0
		var nuevoNombre [16]byte
		copy(mbrTemp.Particion1.Nombre[0:],nuevoNombre[0:])
		if mbrTemp.Particion1.Tipo == 'p'{ //PRIMARIA

		}else{ //EXTENDIDA
			//IR A LA POSICION DE LA PARTICION EXTENDIDA EN EL ARCHIVO BINARIO Y ELIMINAR EL CONTENIDO
			if eliminar == "full" {
				
			}
		}
		escribirArchivoBinario(mbrTemp,ruta)
	}else if nombrePartcionTemp2 == nombre{
		mbrTemp.Particion2.Ajuste = 0
		mbrTemp.Particion2.Estado = 0
		mbrTemp.Particion2.Inicio = 0
		mbrTemp.Particion2.Tamanio = 0
		mbrTemp.Particion2.Tipo = 0
		var nuevoNombre [16]byte
		copy(mbrTemp.Particion2.Nombre[0:],nuevoNombre[0:])
		if mbrTemp.Particion1.Tipo == 'p'{ //PRIMARIA

		}else{ //EXTENDIDA
			//IR A LA POSICION DE LA PARTICION EXTENDIDA EN EL ARCHIVO BINARIO Y ELIMINAR EL CONTENIDO
			if eliminar == "full" {
				
			}
		}
		escribirArchivoBinario(mbrTemp,ruta)
	}else if nombrePartcionTemp3 == nombre{
		mbrTemp.Particion3.Ajuste = 0
		mbrTemp.Particion3.Estado = 0
		mbrTemp.Particion3.Inicio = 0
		mbrTemp.Particion3.Tamanio = 0
		mbrTemp.Particion3.Tipo = 0
		var nuevoNombre [16]byte
		copy(mbrTemp.Particion3.Nombre[0:],nuevoNombre[0:])
		if mbrTemp.Particion1.Tipo == 'p'{ //PRIMARIA

		}else{ //EXTENDIDA
			//IR A LA POSICION DE LA PARTICION EXTENDIDA EN EL ARCHIVO BINARIO Y ELIMINAR EL CONTENIDO
			if eliminar == "full" {
				
			}
		}
		escribirArchivoBinario(mbrTemp,ruta)
	}else if nombrePartcionTemp4 == nombre{
		mbrTemp.Particion4.Ajuste = 0
		mbrTemp.Particion4.Estado = 0
		mbrTemp.Particion4.Inicio = 0
		mbrTemp.Particion4.Tamanio = 0
		mbrTemp.Particion4.Tipo = 0
		var nuevoNombre [16]byte
		copy(mbrTemp.Particion4.Nombre[0:],nuevoNombre[0:])
		if mbrTemp.Particion1.Tipo == 'p'{ //PRIMARIA

		}else{ //EXTENDIDA
			//IR A LA POSICION DE LA PARTICION EXTENDIDA EN EL ARCHIVO BINARIO Y ELIMINAR EL CONTENIDO
			if eliminar == "full" {
				
			}
		}
		escribirArchivoBinario(mbrTemp,ruta)
	}else{
		fmt.Println("No existe la particion a eliminar")
	}
}

func agregarEspacio(nombre string, ruta string, agregar string, unidad string){
	fmt.Println("...")
	/*mbrTemp := leerArchivoBinario(ruta)
	tamanioMBR := binary.Size(mbrTemp)
	fmt.Println("Tamaño del mbr leido: ", tamanioMBR)
	//IMPRIMIR LOS DATOS LEIDOS PARA COMPROBAR QUE SE LEYO CORRECTAMENTE
	fmt.Println("Imprimiendo mbr leido: ", mbrTemp)*/
}

func main(){
	//escribirArchivoBinario()
	leerArchivoEntrada()
}