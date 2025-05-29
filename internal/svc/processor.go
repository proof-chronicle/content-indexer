package svc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/proofchronicle/content-indexer/config"
	pb "github.com/proofchronicle/content-indexer/internal/client/chain_gateway"
	"github.com/proofchronicle/content-indexer/internal/consumer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Processor struct {
	cfg config.Config
}

// NewProcessor creates a new Processor with the given configuration.
func NewProcessor(cfg config.Config) *Processor {
	return &Processor{cfg: cfg}
}

// Handle processes an incoming message.
func (p *Processor) Handle(msg consumer.Message) error {
	return p.sendToGateway(msg)
}

// sendToGateway sends the message to the ChainGateway gRPC service using the non-deprecated NewClient method.
func (p *Processor) sendToGateway(msg consumer.Message) error {
	addr := p.cfg.GatewayAddr
	if addr == "" {
		return fmt.Errorf("gateway address is not configured")
	}

	log.Printf("Connecting to Chain Gateway at %s", addr)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to create gRPC client for %q: %w", addr, err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("warning: failed to close connection: %v", cerr)
		}
	}()

	client := pb.NewChainGatewayClient(conn)

	// Prepare and call Store
	storeReq := &pb.StoreRequest{
		Record: &pb.ContentRecord{
			Uid:       msg.Uid,
			CreatedAt: time.Now().Format(time.RFC3339),
			Hash:      msg.Hash,
			Url:       msg.Url,
		},
	}

	storeCtx, storeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer storeCancel()

	storeResp, err := client.Store(storeCtx, storeReq)
	if err != nil {
		return fmt.Errorf("store RPC failed: %w", err)
	}
	log.Printf("Stored: success=%v txid=%s", storeResp.Success, storeResp.TransactionId)

	// Prepare and call Retrieve
	retrieveReq := &pb.RetrieveRequest{TransactionId: storeResp.TransactionId}
	retrieveCtx, retrieveCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer retrieveCancel()

	retrieveResp, err := client.Retrieve(retrieveCtx, retrieveReq)
	if err != nil {
		return fmt.Errorf("retrieve RPC failed: %w", err)
	}
	log.Printf("Retrieved record: %+v", retrieveResp.Record)

	return nil
}
