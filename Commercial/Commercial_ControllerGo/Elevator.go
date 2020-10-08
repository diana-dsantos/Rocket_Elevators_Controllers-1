package main

import (
	"fmt"
	"math"
)

// STATIC PROPERTIES (in Golang, it's treated as a global)
var OriginFloor int = 1

type Elevator struct {

	// PROPERTIES
	ID            int
	Status        string // online|offline
	Movement      string // up|down|idle
	MaxWeightKG   int
	CurrentFloor  int
	NextFloor     int
	Column        Column
	FloorDisplay  FloorDisplay
	RequestsQueue []int
	Door          ElevatorDoor
}

// METHODS
// Change properties of an ACTIVE elevator in one line - USE ONLY FOR TESTING
func (elevator *Elevator) ChangePropertiesActive(newCurrentFloor int, newNextFloor int?) {

	elevator.CurrentFloor = newCurrentFloor
	elevator.NextFloor = newNextFloor

	if elevator.CurrentFloor > elevator.NextFloor
	{
		elevator.Movement = "down"

	} else {

		elevator.Movement = "up"
	}

	r := Request{elevator.NextFloor, elevator.Movement}
	elevator.RequestsQueue = append(elevator.RequestsQueue, r)
}

// Change properties of an IDLE elevator in one line - USE ONLY FOR TESTING
func (elevator *Elevator) ChangePropertiesIdle(newCurrentFloor int) {

	elevator.CurrentFloor = newCurrentFloor
	elevator.Movement = "idle"
}

// Make elevator go to its scheduled next floor
func (elevator *Elevator) GoToNextFloor() {

	if (elevator.CurrentFloor != elevator.NextFloor) {
		
		if elevator.CurrentFloor > 0 {
			
			fmt.Printf("Elevator %d of Column %d, currently at floor %d, is about to go to floor %d...", 
						elevator.ID, elevator.Column.ID, elevator.CurrentFloor, elevator.NextFloor)

		} else if elevator.NextFloor < 0 {

			fmt.Printf("Elevator %d of Column %d, currently at floor B%d, is about to go to floor B%d...", 
						elevator.ID, elevator.Column.ID, math.Abs(elevator.CurrentFloor), math.Abs(elevator.NextFloor))
		
		} else if elevator.nextFloor > 0 {

			fmt.Printf("Elevator %d of Column %d, currently at floor B%d, is about to go to floor %d...", 
						elevator.ID, elevator.Column.ID, math.Abs(elevator.CurrentFloor), elevator.NextFloor)
		}
		fmt.Println("=================================================================")
		elevator.FloorDisplay.DisplayFloor()
		
		// Traverse through the floors
		for elevator.CurrentFloor != elevator.NextFloor {
			
			// Do not display floors that are not part of the column's range
			if elevator.Movement == "up" {

				if elevator.CurrentFloor + 1 < elevator.Column.LowestFloor {
					
					fmt.Printf("\n... Quickly traversing through the floors not in column %d's usual elevator range ...\n",
							    elevator.Column.ID);
					elevator.CurrentFloor = elevator.Column.LowestFloor
				} else {
					elevator.CurrentFloor++
				}

			} else {

				if elevator.CurrentFloor - 1 < elevator.Column.LowestFloor {

					fmt.Printf("\n... Quickly traversing through the floors not in column %d's usual elevator range ...\n",
								 elevator.Column.ID);
					elevator.CurrentFloor = OriginFloor
				} else {
					elevatorCurrentFloor++
				}
			}

			elevator.FloorDisplay.DisplayFloor()
		}
		fmt.Println("=================================================================")
		
		if elevator.CurrentFloor > 0 {

			fmt.Printf("Elevator %d of Column %d has reached its requested floor! It is now at floor %d...", 
						elevator.ID, elevator.Column.ID, elevator.CurrentFloor)
		} else {

			fmt.Printf("Elevator %d of Column %d has reached its requested floor! It is now at floor B%d...", 
						elevator.ID, elevator.Column.ID, math.Abs(elevator.CurrentFloor))
		}
	}
}

// Make elevator go to origin floor
func (elevator *Elevator) GoToOrigin() {

	elevator.NextFloor = OriginFloor
	fmt.Printf("Elevator %d of Column %d going back to RC (floor %d)...".
				elevator.ID, elevator.Column.ID, OriginFloor)
	GoToNextFloor()
}

// Set what should be the movement direction of the elevator for its upcoming request
func (elevator *Elevator) SetMovement() {

	floorDifference := elevator.CurrentFloor - elevator.RequestsQueue[0].Floor

	if floorDifference > 0
	{
		elevator.Movement = "down"
	} else if (floorDifference < 0) {

		elevator.Movement = "up"
	} else {

		elevator.Movement = "idle"
	}
}

// Sort requests, for added efficiency
func (elevator *Elevator) SortRequestsQueue() {
	
	request := elevator.RequestsQueue[0]

	// Remove any requests which are useless i.e. requests that are already on their desired floor
	for i, req := range elevator.RequestsQueue {

		if req.Floor == elevator.CurrentFloor {
			
			elevator.RequestsQueue = append(elevator.RequestsQueue[:i-1], elevator.RequestsQueue[i+1:])
		}
	}

	SetMovement()

	if len(elevator.RequestsQueue) > 1 {

		if elevator.Movement == "up" {

			// Sort the queue in ascending order
			sort.SliceStable(elevator.RequestsQueue, func(i int, j int) bool {
				return elevator.RequestsQueue[i].Floor < elevator.RequestsQueue[j].Floor
			})

			//  Push any request to the end of the queue that would require a direction change
			for _, req := range elevator.RequestsQueue {

				if req.Direction != elevator.Movement || req.Floor < elevator.CurrentFloor {

					elevator.RequestsQueue = append(elevator.RequestsQueue[:i-1], elevator.RequestsQueue[i+1:])
					elevator.RequestsQueue = append(req);
				}
			}

		} else {

			// Reverse the sorted queue (will now be in descending order)
			sort.SliceStable(elevator.RequestsQueue, func(i int, j int) bool {
				return elevator.RequestsQueue[j].Floor < elevator.RequestsQueue[i].Floor
			})

			//  Push any request to the end of the queue that would require a direction change
			for _, req := range elevator.RequestsQueue {

				if req.Direction != elevator.Movement || req.Floor > elevator.CurrentFloor {

					elevator.RequestsQueue = append(elevator.RequestsQueue[:i-1], elevator.RequestsQueue[i+1:])
					elevator.RequestsQueue = append(req);
				}
			}




		}
	}
}

// Complete the elevator requests
func (elevator *Elevator) DoRequests() {

	if len(elevator.RequestsQueue) > 0 {

		// Make sure queue is sorted before any request is completed
		SortRequestsQueue()
		requestToComplete := elevator.RequestsQueue[0]

		// Go to requested floor
		if (elevator.Door.Status != "closed") {
			CloseDoor()
		}
		elevator.NextFloor = requestToComplete.Floor
		GoToNextFloor()

		// Remove request after it is complete
		OpenDoor()
		elevator.RequestsQueue = append(elevator.RequestsQueue[1:])

		// Automatically close door
		CloseDoor()
	}
	// Automatically go idle temporarily if 0 requests or at the end of request
	elevator.Movement = "idle"
}

// Check if elevator is at full capacity
func (elevator *Elevator) CheckWeight(currentWeightKG int){
	
	// currentWeightKG calculated thanks to weight sensors
	if currentWeightKG > elevator.MaxWeightKG {

		// Display 10 warnings
		for _ := 0; _ < 10; _++ {
			fmt.Printf("\nALERT: Maximum weight capacity reached on Elevator %d of Column %d",
						elevator.ID, elevator.Column.ID)
						
		}

		// Freeze elevator until weight goes back to normal
		elevator.Movement = "idle"
		OpenDoor()
	}
}