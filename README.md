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
		- Kafka


# Coordinator
# CoordinatorSaga
# Saga
# Stage


<!-- ## Saga Template -->
## Saga Transaction + Coordinator
<!-- ## Saga Coordinator -->
	-> Start
	-> StartAsync
	-> Load
	-> Abort

## Tasks
- `DONE` Make SAGA execution asynchronous
- Add support for logger
	- Add support for logger agents and/or centralized logging
- 
<!-- 
Inspiration
EchoVault:
- [LinkedIn Post](https://www.linkedin.com/feed/update/urn:li:activity:7186276898723799040?utm_source=share&utm_medium=member_desktop)
- https://github.com/EchoVault/EchoVault
 -->
