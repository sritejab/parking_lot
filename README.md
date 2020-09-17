# Parking Lot server implementation
This is a basic server implementation of a parking lot system using golang and sqlite3. 
## Assumptions for the project
1) It is a multifloor parking lot in a IT park/ IT SEZ with each floor assigned uniquely to one company
2) Each floor can accommodate either 2 wheelers or 4 wheelers only.
3) Each floor is made up of slots. Slot numbers are unique throughout parking lot.
4) Tickets are generated for each parked vehicke and can be used to retrieve the slot and floor details of the vehicke if needed(if user forgets)
## Usage
Requirements: Golang 1.11 and above
1) Clone the repo
2) cd into parking_lot folder
3) In command prompt, enter "go mod init parking_lot"
4) Run the database commands manually to initialize the database with required tables and data with commands from db_initialize.txt file
5) cd to src folder and build the server.go file ("go build server.go")
6) Run the server.exe and specify port number. Ex.: server.exe "8080"
7) Make required API calls from your favourite API client.
## APIs Implemented
1) allocateSlot - A POST method to allocate a slot to either a two wheeler or a 4 wheeler based on the JSON information posted in body. It returns a ticket ID for the vehicle parked
Example Json:
```json
{
  "companyName":"Microsoft",
  "vehicleNum" : "AP20IT2020",
  "numWheels" : 2
}
```
2) findParkedSpotByTicketID - A GET method to retrieve the parking slot and floor incase the user forgot where he parked the vehicle. This method uses the Ticket ID generated to find the location of the vehicle.
We need to give a header "ticketId" and its value(type: text/plain)
3) findParkedSpotByVehicleNum - A GET method to retrieve the parking slot and floor incase the user forgot where he parked the vehicle. This method uses the Vehicle number to identify the vehicle.
We need to give a header "vehicleId" and its value(type: text/plain)
4) checkFreeSpaceOnFloor - A GET mentod to query how many free spots of 2 wheelers and 4 wheelers are available on a particular floor. The floor number is sent as a header "floorNum" and value "text/plain" type.

Example URL: localhost:8080/findParkedSpotByTicketID

## Nice to Have in future
1) Ability to accommodate vechiles having different number of wheels apart from 2 and 4
2) Ability to allocate multiple floors to a company/ share a floor brtween companies
3) More refined and well planned database schema and database functions in code
4) Increasing and decreasing the capacity of the floor
