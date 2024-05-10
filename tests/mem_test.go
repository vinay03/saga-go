package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "net/http/pprof"

	saga "github.com/vinay03/saga-go"
)

var _ = Describe("In-Memory Message Carrier", func() {
	// BeforeEach(func() {
	// })

	// AfterEach(func() {
	// })

	BeforeAll(func() {
		InitInMemSaga()
	})

	AfterAll(func() {

	})

	It("Algorithm with Multiple Targets", func() {
		coord := saga.GetCoordinatorInstance()

		sampleSaga := coord.Saga
		// data := struct{}{}
		// start := time.Now()
	})

})
