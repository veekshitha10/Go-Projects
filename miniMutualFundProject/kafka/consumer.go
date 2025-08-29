 package kafka

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/twmb/franz-go/pkg/kgo"
// )

// // func NewMessagingC(topic string, brokers []string) *Messaging {
// // 	return &Messaging{make(chan []byte), topic, brokers}
// // }
// func (msg *Messaging) ConsumeRecords() {

// 	// flag.StringVar(&Topic, "topic", "omnenext.demo.v1", "--topic=omnenext.demo.v1")
// 	// flag.StringVar(&ConsumerGroup, "cg", " ", "--cg=demo-consumer-group")
// 	// flag.Parse()

// 	seeds := []string{"localhost:19092", "localhost:29092", "localhost:39092"}
// 	// One client can both produce and consume!
// 	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
// 	cl, err := kgo.NewClient(
// 		kgo.SeedBrokers(seeds...),
// 		kgo.ConsumerGroup(msg.ConsumerGroup),
// 		kgo.ConsumeTopics(msg.Topic),
// 		//kgo.RequiredAcks(kgo.AllISRAcks()), // or kgo.RequireOneAck(), kgo.RequireNoAck()
// 		//kgo.DisableIdempotentWrite()
// 		//kgo.RetryTimeout()
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer cl.Close()

// 	ctx := context.Background()
// 	time.Sleep(time.Second * 5)
// 	for {
// 		fetches := cl.PollFetches(ctx)
// 		if errs := fetches.Errors(); len(errs) > 0 {
// 			// All errors are retried internally when fetching, but non-retriable errors are
// 			// returned from polls so that users can notice and take action.
// 			panic(fmt.Sprint(errs))
// 		}

// 		// We can iterate through a record iterator...
// 		iter := fetches.RecordIter()
// 		for !iter.Done() {
// 			record := iter.Next()
// 			fmt.Println("Partition-->", record.Partition, "Topic-->", record.Topic, string(record.Value), "from an iterator!")
// 						msg.ChMessagingC<-record.Value

// 		}

// 		// // or a callback function.
// 		// fetches.EachPartition(func(p kgo.FetchTopicPartition) {
// 		// 	for _, record := range p.Records {
// 		// 		fmt.Println(string(record.Value), "from range inside a callback!")
// 		// 	}

// 		// 	// We can even use a second callback!
// 		// 	p.EachRecord(func(record *kgo.Record) {
// 		// 		fmt.Println(string(record.Value), "from a second callback!")
// 		// 	})
// 		// })
// 	}

// }
