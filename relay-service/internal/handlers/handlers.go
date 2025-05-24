package handlers

import (
	pb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/proto"
	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/repository"
)

// var allMessages = []*pb.Message{ //testing purposes
// 	{SenderId: "66", ReceiverId: "55", Timestamp: 100, Data: "Poruka 1"},
// 	{SenderId: "66", ReceiverId: "55", Timestamp: 105, Data: "Poruka 2"},
// 	{SenderId: "55", ReceiverId: "66", Timestamp: 120, Data: "Poruka 3"},
// 	{SenderId: "55", ReceiverId: "66", Timestamp: 130, Data: "Poruka 4"},
// }

func StoreMessage(message *pb.Message) (status *pb.Status) {
	// allMessages = append(allMessages, message)
	repository.SaveMessage(message)
	return &pb.Status{Status: 0}
}

func GetMessagesForUserId(id *pb.UserId) []*pb.Message {
	var result []*pb.Message
	// for i := range allMessages {
	// 	if allMessages[i].ReceiverId == id.Id {
	// 		result = append(result, allMessages[i])
	// 	}
	// }
	// return result
	result, _ = repository.GetMessages(id)
	return result
}
