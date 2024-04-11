## Saga Pattern in Go


## Structure
- Entities
	- `Saga`
		-> AddStage
		-> VerifyStageConfig
	- Stage (Sub-Transaction)
		-> Exec
	- SagaOperation
	- SagaTemplate
		-> 
	- Operator
		-> Start
		-> StartAsync
		-> Load
		-> Abort
- Storage
	- In-Memory
	- DB
		- SQL
		- No-SQL
		- Key-Value
	- AMQP Protocol
		- Redis
		- RabbitMQ
- 