package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	//"strconv"
	"bytes"
	//"encoding/json"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal((err))
	}

}

func checkForArgs() {
	if len(os.Args) != 3 {
		log.Fatal("Usage client <hostname> <port>")
	}

	fmt.Printf("Connected to server %s on port %d \n", os.Args[1], os.Args[2])
}

var (
	//serverResponse bytes.Buffer = []byte
	selectYourOperationText = "Press i for inserting new record\nPress r for removing an existing record using ID\nPress f to search for a record \nPress l to list all records\nSelect the operation or enter 'stop' to exit >>> "
	op string
	id string
	name string
	state string
	response = make([]byte, 1024)
)

type citizen struct{
	cvID string
	cvName string
	cvState string
}

func main() {

	checkForArgs()

	sock, err := net.Dial("tcp", os.Args[1]+":"+os.Args[2])
	checkErr(err)

	handleOperation(sock)
	sock.Close()
}

func handleOperation(sock net.Conn) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(selectYourOperationText)
	op,_,_ := reader.ReadLine()

	switch string(op) {
	case "i":
		handleInsertion(sock)
	case "r":
		handleDeletion(sock)
	case "f":
		handleSearching(sock)
	case "l":
		handleListing(sock)
	case "stop":
		break
	default :
		fmt.Println("Please enter a valid input")
		handleOperation(sock)
	}
}

func handleInsertion(sock net.Conn) {
	//reader := bufio.NewReader(os.Stdin)
	op = "i"
	//op := []byte("i")
	for {
		fmt.Print("Enter ID > ")
		id = ""
		fmt.Scan(&id)
		//id,_ := reader.ReadBytes('\n')
		//byteId := []byte(string(id))
		
		fmt.Print(("Enter name for ID " + string(id) + " > "))
		name = ""
		fmt.Scan(&name)
		//name, _ := reader.ReadBytes('\n')
		//byteName := []byte(name)

		fmt.Print(("Enter state for name " + string(name) + " and ID " + string(id) + " > "))
		state = ""
		fmt.Scan(&state)
		//state, _ := reader.ReadBytes('\n')
		//byteState := []byte(state)

		//newArrayRecord := [][]byte{op,byteId,byteName,byteState}
		newArrayRecord := [][]byte{[]byte(op),[]byte(id),[]byte(name),[]byte(state)}
		newSliceRecord := bytes.Join(newArrayRecord, []byte(","))
		
		_,errW := sock.Write((newSliceRecord))
		checkErr(errW)

		fmt.Print("Record {" + string(newSliceRecord) +  "} is sent to the server\n")
		//mt.Println(string(bytes.Split(newRecord, []byte(","))[0]))
		
		resLen,errR := sock.Read(response)
		checkErr(errR)

		fmt.Println(string(response[:resLen]))
		// fmt.Println("Name Saved :: "+ string(line))
		handleOperation(sock)
	}
}

func handleDeletion(sock net.Conn) {

	op = "r"
	for {
		fmt.Print("Enter ID > ")
		id = ""
		fmt.Scan(&id)
		
		_,err := sock.Write(bytes.Join([][]byte{[]byte(op),[]byte(id)}, []byte(",")))
		checkErr(err)

		fmt.Println("Record Deleted ...")
		var response []byte
		sock.Read(response)
		fmt.Println(response)
		handleOperation(sock)
	}
}

func handleSearching(sock net.Conn) {
	op = "f"
	for {
		fmt.Print("Enter ID > ")
		id = ""
		fmt.Scan(&id)

		_,err := sock.Write(bytes.Join([][]byte{[]byte(op),[]byte(id)}, []byte(",")))
		checkErr(err)

		fmt.Println("Seraching for record ...")

		var response []byte
		sock.Read(response)
		fmt.Println(string(response))

		handleOperation(sock)
		
	}
}

func handleListing(sock net.Conn) {

	op = "l"
	_,err := sock.Write([]byte("l"))
	checkErr(err)

	var response []byte 
	sock.Read(response)
	fmt.Println(response)
}