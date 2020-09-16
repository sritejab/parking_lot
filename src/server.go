package main

import (
	"fmt"
	"net/http"
	"parking_lot/src/db"
	"strconv"
	"encoding/json"
	"github.com/golang/gddo/httputil/header"
	"bytes"
	"os/signal"
	"os"
	"syscall"
)

func main() {
	// code to close db connection when exiting the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<- sigs
		db.Close()
		fmt.Println("Exiting")
		os.Exit(0)
	}()
	// server functionality
	var port string
	ip := "127.0.0.1"
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		fmt.Println("Please enter port")
		fmt.Println("Usage: server.exe 'port'")
		db.Close()
		fmt.Println("Exiting")
		os.Exit(0)
	}

	url := ip + ":" + port
	fmt.Println("Parking Lot server running on ", url)
	http.HandleFunc("/allocateSlot", AllocateSlot) // POST call
	http.HandleFunc("/findParkedSpotByTicketID", FindParkedSpotByTicketID ) //GET call
	http.HandleFunc("/findParkedSpotByVehicleNum", FindParkedSpotByVehicleNum )// GET call
	http.HandleFunc("/checkFreeSpaceOnFloor", CheckFreeSpaceOnFloor )// GET call
	http.ListenAndServe(url, nil)
	
}

type floorJsonStruct struct{
	CompanyName string `json:"companyName"`
	VehicleNum string `json:"vehicleNum"`
	NumWheels int `json:"numWheels"`
}

func AllocateSlot(w http.ResponseWriter, r *http.Request) {
	var companyName string
	var vehicleNum string
	var numOfWheels int

	//get the sent body for required information
	if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, 400)
            return
        }
    } else{
		msg:= "Content type missing"
		http.Error(w, msg, 400)
	}

	//decoding the json body from POST call
	buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
	var slotInFloor floorJsonStruct
	err := json.NewDecoder(buf).Decode(&slotInFloor)
	
    if err != nil {
		errMsg := "Json not valid"
		http.Error(w, errMsg, 400)
	} else{
		companyName = slotInFloor.CompanyName
		vehicleNum = slotInFloor.VehicleNum
		numOfWheels = slotInFloor.NumWheels
		
	}
	
	//query with companyName and get floor num
	query := `Select floorNum from lot where company = "` + companyName + `";`
	floorNum := db.Query(query) // nice-to-have: Handling a company having multiple floors of parking assigned to it

	//Check if there are empty slots for that vehicle type in that floor
	if numOfWheels == 2 {
		query = `select numTwoWheelerSlots from lot where floorNum = ` + strconv.Itoa(floorNum) + `;`
	} else if numOfWheels == 4 {
		query = `select numTwoWheelerSlots from lot where floorNum = ` + strconv.Itoa(floorNum)
	} else{
		msg := "Cannot accomodate vehicles other than 2 and 4 wheelers" //nice-to-have: to be able to park vechicles other than 2 & 4 wheelers
		http.Error(w, msg, 400)
		return
	}
	openSlots := db.Query(query)
	
	//Accomodate the vehicle based on space availability
	if openSlots > 0 {
		//allocate a slot on that floor and save the vehicle data
		queryStmt:= `INSERT INTO floor (floorNum, vehicleNo, numWheels) Values`
		valueStr := `(` + strconv.Itoa(floorNum) + `, "` + vehicleNum + `", ` + strconv.Itoa(numOfWheels) + `);` 
		queryStmt = queryStmt + valueStr
		db.ExecuteStatement(queryStmt)
		//Generate a ticket for that vehicle and add to ticket db
		ticketQueryStmt := `INSERT INTO ticket (floorNum, vehicleNo) VALUES ` + `(` + strconv.Itoa(floorNum) + `, "` + vehicleNum +  `");` 
		db.ExecuteStatement(ticketQueryStmt)
		query := `Select ticketNum from ticket where vehicleNo = "` + vehicleNum + `";`
		ticketNum := db.Query(query)
		outputTicketString := []byte("your ticket no. is: " + strconv.Itoa(ticketNum))
		w.Write(outputTicketString)
		//update open slots
		if numOfWheels == 2 {
			query = `UPDATE lot set numTwoWheelerSlots = ` + strconv.Itoa((openSlots -1)) + `;`
		} else{// add case for vehicles other than 2,4 wheelers
			query = `UPDATE lot set numFourWheelerSlots = ` + strconv.Itoa((openSlots -1)) + `;`
		}
		db.ExecuteStatement(query)
	} else{
		responseData:= "Cannot allocate slot, floor capacity full"
		byteResponseData := []byte(responseData)
		w.Write(byteResponseData)
	}
}

func FindParkedSpotByTicketID(w http.ResponseWriter, r *http.Request) {
	var vehicleNum, responseData string
	var slotNum, floorNum int

	ticketId, err := strconv.Atoi(r.Header.Get("ticketId"))
	if err != nil {
		fmt.Println("error converting string to int", err)
	}

	query:= `select vehicleNo from ticket where ticketNum = ` + strconv.Itoa(ticketId) + `;`
	vehicleNum = db.QueryForString(query)
	if vehicleNum != "" {
		query = `select slotNum from floor where vehicleNo = "` + vehicleNum + `";`
		slotNum = db.Query(query)
		if slotNum != 0 {
			query = `select floorNum from floor where slotNum = ` + strconv.Itoa(slotNum) + `;`
			floorNum = db.Query(query)
			if floorNum != 0 {
				responseData = "Vehicle is parked at floor: " + strconv.Itoa(floorNum) + ", slot no: " + strconv.Itoa(slotNum)
			} else {
				responseData = "No such floor"
			}
		} else {
			responseData = "No such slot"
		}
	} else {
		responseData = "No such vehicle"
	}
	byteResponseData := []byte(responseData)
	w.Write(byteResponseData)
}

func FindParkedSpotByVehicleNum(w http.ResponseWriter, r *http.Request) {
	var vehicleNum, responseData string
	var slotNum, floorNum int

	vehicleNum = r.Header.Get("vehicleId")

	if vehicleNum != "" {
		query := `select slotNum from floor where vehicleNo = "` + vehicleNum + `";`
		slotNum = db.Query(query)
		if slotNum != 0 {
			query = `select floorNum from floor where slotNum = ` + strconv.Itoa(slotNum) + `;`
			floorNum = db.Query(query)
			if floorNum != 0 {
				responseData = "Vehicle is parked at floor: " + strconv.Itoa(floorNum) + ", slot no: " + strconv.Itoa(slotNum)
			} else {
				responseData = "No such floor"
			}
		} else {
			responseData = "No such slot"
		}
	} else {
		responseData = "No such vehicle"
	}
	byteResponseData := []byte(responseData)
	w.Write(byteResponseData)
}

func CheckFreeSpaceOnFloor(w http.ResponseWriter, r *http.Request) {
	//get floor num
	floorNum := r.Header.Get("floorNum")
	//get 2 wheeler and 4 wheeler slot capacity of the floor
	query := `Select numTwoWheelerSlots from lot where floorNum = ` + floorNum + `;`
	openTwoWheelerSlots := db.Query(query)
	query = `Select numFourWheelerSlots from lot where floorNum = ` + floorNum + `;`
	openFourWheelerSlots := db.Query(query)
	//get count of occupied spots and subtract
	responseData := strconv.Itoa(openTwoWheelerSlots) + ` Two Wheeler slots and ` + strconv.Itoa(openFourWheelerSlots) + ` four wheeler slots are available on floor no. ` + floorNum
	byteResponseData := []byte(responseData)
	w.Write(byteResponseData)
}