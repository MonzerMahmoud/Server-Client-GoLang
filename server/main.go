package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	//"strconv"
)

const (
	HOST_NAME = "localhost"
	PORT_NUM = "5050"
	CONN_TYPE = "tcp"
)

var (
	publicCivilReg []citizen
	numberOfConnections int = 0
	id string
	name string
	state string
	serRes = make(map[string]int)
)

type citizen struct{
	cvID string
	cvName string
	cvState string
}

func main() {

	
	psock, err := net.Listen(CONN_TYPE, HOST_NAME+":"+PORT_NUM)
	if err != nil {
		log.Fatalf("Cannot open port %s at %s\n", PORT_NUM, HOST_NAME)
	}

	fmt.Println("Server is running on PORT :: "+string(PORT_NUM))

	defer psock.Close()
	
	//record := make(map[string]string)
	//var index = 0
	for {
		sock, err := psock.Accept()
		if err != nil {
			log.Fatal("Cannot accept a new connection")
		}

		go handleOperation(sock)
	}
}

func handleOperation(sock net.Conn) {

	fmt.Println("New client connected")

	numberOfConnections = numberOfConnections + 1

	fmt.Println("Number of connected user is " + fmt.Sprint(numberOfConnections))
	
	buf := make([]byte, 1024)

	for {
		recLen, err := sock.Read(buf)
		if err != nil {
		log.Println(err)
		break
		}
		if recLen <= 0 {
		break
		}
		
		//newRecord := bytes.Trim(a, ", ")
		
		switch string(bytes.Split(buf, []byte(","))[0]) {
		case "i":
			handleInsertion(buf, sock)
		case "r":
			handleDeletion(sock,string(bytes.Split(buf, []byte(","))[1]))
		case "f":
			handleSearching(sock,string(bytes.Split(buf, []byte(","))[1]))
		case "l":
			handleListing(sock)
		}
		//a := string(bytes.Split(buf, []byte(","))[0])
		//b := string(bytes.Split(buf, []byte(","))[1])
		//fmt.Println(a)
		//fmt.Println(b)

		//record[id] = name
		//p//rintln(record)
		//fmt.Printf("Server recived :: %s",rdLen)
		//fmt.Println(record)
		//sock.Write(buf[:rdLen])
	}
	sock.Close()
	fmt.Println("One client disconnected.")
	numberOfConnections = numberOfConnections - 1
	fmt.Println("Number of connected user is " + fmt.Sprint(numberOfConnections))

}

func handleInsertion(buf []byte, sock net.Conn) {

	id = ""
	id := string(bytes.Split(buf, []byte(","))[1])

	name = ""
	name := string(bytes.Split(buf, []byte(","))[2])

	state = ""
	state := string(bytes.Split(buf, []byte(","))[3])

	isAvail,_ := isIDAvailable(id)

	if isAvail {
		newCitizen := citizen{cvID: id, cvName: name, cvState: state}

		fmt.Printf("New record has been added ... %s\n", newCitizen)
		publicCivilReg = append(publicCivilReg, newCitizen)

		fmt.Printf("The Current Reg :: %s\n",publicCivilReg)
		//serverResponse := []byte("Record Added Successfully")

		serRes["cvResp"] = 1
		serRes["cvCode"] = 0
				
		sock.Write(bytes.Join([][]byte{[]byte("1"), []byte("0")} , []byte(",")))
	} else {
		fmt.Println("ID is already taken")

		sock.Write(bytes.Join([][]byte{[]byte("2"), []byte("1")} , []byte(",")))
	}
}

func handleDeletion(sock net.Conn, index string) {

	var tempReg []citizen
	for _, i := range publicCivilReg {
		if i.cvID != index {
			tempReg = append(tempReg, i)
		}
	}

	publicCivilReg = tempReg
	fmt.Println("Record has been deleted ...")
	fmt.Println(publicCivilReg)
	sock.Write([]byte("Success"))
}

func handleSearching(sock net.Conn, index string) {

	isAvail,record := isIDAvailable(index)
	//var record citizen
	
	if isAvail {
		fmt.Println("There is no record with ID of " + index)
		sock.Write([]byte("There is no record with ID of " + index))
	} else {
		fmt.Println(string(bytes.Join([][]byte{[]byte(record.cvID),[]byte(record.cvName),[]byte(record.cvState)}, []byte(","))))
		//sock.Write(bytes.Join([][]byte{[]byte(record.cvID),[]byte(record.cvName),[]byte(record.cvState)}, []byte(",")))
		newRecord := citizenStructToSlice(*record)
		sock.Write(newRecord)
	}	
}

func handleListing(sock net.Conn) {

	fmt.Println(publicCivilReg)
	sock.Write([]byte("Listing ..."))
}

func isIDAvailable(index string) (bool, *citizen) {

	var record citizen
	for _, i := range publicCivilReg {
		if i.cvID == index {
			record.cvID = i.cvID
			record.cvName = i.cvName
			record.cvState = i.cvState

			return false, &record
		}
	}
	return true, nil
}

func citizenStructToSlice(record citizen) []byte {
	return bytes.Join([][]byte{[]byte(record.cvID),[]byte(record.cvName),[]byte(record.cvState)}, []byte(","))
}

func sliceToCitizenStruct(slice []byte) citizen {
	return citizen{cvID: string(bytes.Split(slice, []byte(","))[1]), cvName: string(bytes.Split(slice, []byte(","))[2]), cvState: string(bytes.Split(slice, []byte(","))[3])}
}