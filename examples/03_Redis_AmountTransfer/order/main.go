package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	saga "github.com/vinay03/saga-go"

	sagas "github.com/vinay03/saga-go/examples/03_Redis_AmountTransfer/sagas"
)

var SagaCoordinator *saga.Coordinator

type transferPayload struct {
	From   int     `json:"from"`
	To     int     `json:"to"`
	Amount float32 `json:"amount"`
}

func transfer(payload transferPayload) error {
	fmt.Println(payload)
	SagaCoordinator.Start("Transfer", payload)
	return nil
}

func main() {
	SagaCoordinator = saga.GetCoordinatorInstance()

	sagas.SetupSagas()

	server := gin.Default()
	server.POST("/transfer", func(c *gin.Context) {

		var payload transferPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := transfer(payload)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	})
	server.Run(":8080")

}
