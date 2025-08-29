package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"mutualfundminiproject/models"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Messaging struct {
	ChMessaging chan []byte
	Topic       string
	Brokers     []string
	ConsumerGroup string
	ChMessagingC   chan []byte
}

func NewMessaging(topic string, brokers []string,consumergroup string,) *Messaging {
	return &Messaging{make(chan []byte,2), topic, brokers,consumergroup,make(chan []byte,2)}
	}


func (msg *Messaging) ProduceRecords() {
	if msg.Topic == "" {
		panic("invalid topc")
	}

	if len(msg.Brokers) < 1 {
		panic("invalid brokers")
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),
		// kgo.ConsumerGroup(ConsumerGroup),
		// kgo.ConsumeTopics(Topic),
		kgo.RequiredAcks(kgo.AllISRAcks()), // or kgo.RequireOneAck(), kgo.RequireNoAck()
		//kgo.DisableIdempotentWrite()
	)
	if err != nil {
		panic(err)
	}

	defer cl.Close()
	ctx := context.Background()
	for message := range msg.ChMessaging {

		record := &kgo.Record{Topic: msg.Topic, Value: message, Key: nil}

		//
		cl.Produce(ctx, record, func(r *kgo.Record, err error) {
			//defer wg.Done()
			if err != nil {
				fmt.Printf("record had a produce error: %v\n", err)
			}
			user := new(models.UserTable)
			json.Unmarshal(r.Value, user)
			//fmt.Println(user)
			fmt.Println("Producer-->", r.ProducerID, "Topid-->", r.Topic, "Partition:", r.Partition, "Offset:", r.Offset, "Value:", user)
		})

		// // Alternatively, ProduceSync exists to synchronously produce a batch of records.
		// if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
		// 	fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
		// }
	}
	cl.Flush(ctx)
	log.Print("Closed publishing data")

}

func (msg *Messaging) ConsumeRecords() {
   print("coming here in consumer1")

	// flag.StringVar(&Topic, "topic", "omnenext.demo.v1", "--topic=omnenext.demo.v1")
	// flag.StringVar(&ConsumerGroup, "cg", " ", "--cg=demo-consumer-group")
	// flag.Parse()

	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),
		kgo.ConsumerGroup(msg.ConsumerGroup),
		kgo.ConsumeTopics(msg.Topic),
		//kgo.RequiredAcks(kgo.AllISRAcks()), // or kgo.RequireOneAck(), kgo.RequireNoAck()
		//kgo.DisableIdempotentWrite()
		//kgo.RetryTimeout()
	)
	   print("coming here in consumer2")

	if err != nil {
		panic(err)
	}
	defer cl.Close()
   print("coming here in consumer3")
	ctx := context.Background()
	time.Sleep(time.Second * 5)
	for {
		fetches := cl.PollFetches(ctx)
		   print("coming here in consumer4")

		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		for !iter.Done() {
			   print("coming here in consumer5")

			record := iter.Next()
			fmt.Println("Partition-->", record.Partition, "Topic-->", record.Topic, string(record.Value), "from an iterator!")
						msg.ChMessagingC<-record.Value

		}

		// // or a callback function.
		// fetches.EachPartition(func(p kgo.FetchTopicPartition) {
		// 	for _, record := range p.Records {
		// 		fmt.Println(string(record.Value), "from range inside a callback!")
		// 	}

		// 	// We can even use a second callback!
		// 	p.EachRecord(func(record *kgo.Record) {
		// 		fmt.Println(string(record.Value), "from a second callback!")
		// 	})
		// })
	}

}
