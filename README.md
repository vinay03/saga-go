## Saga
A framework to implement Choreography based SAGA pattern in Golang. You can use this framework to implement SAGA pattern in local as well as distributed systems.


## Structure
- Entities
	- Stage (Sub-Transaction)
	- `Saga`
	- Transaction
	<!-- - Template -->
	- Coordinator
- Storage
	- In-Memory
	- DB
		- SQL
		- No-SQL
		- Key-Value
	- AMQP Protocol
		- Redis
		- RabbitMQ


## Saga
## Saga Stage
<!-- ## Saga Template -->
## Saga Transaction + Coordinator
<!-- ## Saga Coordinator -->
	-> Start
	-> StartAsync
	-> Load
	-> Abort


Actions
> Start a SAGA transaction
	- Provide saga id, context and data
	- Trigger an event `OrderCreated::start`
	- this should trigger 