package main

import (
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

func transfer(from int, to int, amount float32) error {
	payload := transferPayload{
		From:   from,
		To:     to,
		Amount: amount,
	}
	SagaCoordinator.Start("Transfer", payload)
	return nil
}

func main() {
	SagaCoordinator = saga.GetCoordinatorInstance()

	sagas.SetupSagas()

	server := gin.Default()
	server.POST("/transfer", func(c *gin.Context) {
		// Parse the JSON object from the request body
		// Get the "from", "to", and "amount" fields from the JSON object
		// Call the transfer function from the account package
		// If the transfer is successful, return a status code of 200
		// If the transfer is unsuccessful, return a status code of 400
		var payload transferPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := transfer(payload.From, payload.To, payload.Amount)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	})
	server.Run(":8080")

}
