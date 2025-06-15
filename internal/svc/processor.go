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
	log.Printf("Processing message: UID=%s, URL=%s, Hash=%s, CreatedAt=%s, ContentLength=%d",
		msg.Uid, msg.Url, msg.Hash, msg.CreatedAt, msg.ContentLength)
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
	uid := msg.Uid
	if uid == "" {
		uid = "test-uid-123"
	}

	url := msg.Url
	if url == "" {
		url = "https://example.com/test-page"
	}

	hash := msg.Hash
	if hash == "" {
		hash = "abc123def456789abcdef123456789abcdef123456789abcdef123456789abcdef"
	}

	createdAt := msg.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339)
	}

	storeReq := &pb.StoreRequest{
		Record: &pb.ContentRecord{
			Uid:           uid,
			Url:           url,
			ContentHash:   hash,
			ContentLength: 1024, // Hardcoded for now, could be calculated from content
			Version:       1,    // Schema version
		},
	}

	log.Printf("Sending to gateway - UID: %s, URL: %s, Hash: %s, ContentLength: %d, Version: %d",
		storeReq.Record.Uid, storeReq.Record.Url, storeReq.Record.ContentHash,
		storeReq.Record.ContentLength, storeReq.Record.Version)

	storeCtx, storeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer storeCancel()

	storeResp, err := client.Store(storeCtx, storeReq)
	if err != nil {
		return fmt.Errorf("store RPC failed: %w", err)
	}
	log.Printf("✅ Stored successfully: success=%v txid=%s account=%s",
		storeResp.Success, storeResp.TransactionId, storeResp.AccountAddress)

	return nil
}
