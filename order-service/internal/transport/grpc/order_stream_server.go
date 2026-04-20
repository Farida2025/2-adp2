package grpc

import (
	"log"
	"sync"

	pb "github.com/Farida2025/assignment2-generated/order"
)

type OrderStreamServer struct {
	pb.UnimplementedOrderServiceServer
	subscribers map[string][]chan string
	mu          sync.RWMutex
}

func NewOrderStreamServer() *OrderStreamServer {
	return &OrderStreamServer{
		subscribers: make(map[string][]chan string),
	}
}

func (s *OrderStreamServer) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderService_SubscribeToOrderUpdatesServer) error {
	orderID := req.GetOrderId()
	updateChan := make(chan string, 1)

	s.mu.Lock()
	s.subscribers[orderID] = append(s.subscribers[orderID], updateChan)
	s.mu.Unlock()

	log.Printf("Client subscribed to notification: %s", orderID)

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Cancelled: %s", orderID)
			return nil
		case status := <-updateChan:
			err := stream.Send(&pb.OrderStatusUpdate{
				OrderId: orderID,
				Status:  status,
			})
			if err != nil {
				return err
			}
		}
	}
}

func (s *OrderStreamServer) NotifyStatusChange(orderID, newStatus string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if chans, ok := s.subscribers[orderID]; ok {
		for _, ch := range chans {
			ch <- newStatus
		}
	}
}
