package provider

import (
	"context"
	"time"
)

// AnchorResult is the outcome of anchoring a chain to an external ledger.
type AnchorResult struct {
	Anchored bool
	Time     *time.Time
	TxHash   *string
}

// AnchorProvider anchors a hash-chain root to an external ledger. Implementations
// may be a no-op (default) or tie into Ethereum/other via a provider interface.
type AnchorProvider interface {
	// Anchor records that a chain root hash has been anchored. Returned result
	// reflects whether a real external tx was written (TxHash set) or not.
	Anchor(ctx context.Context, chainID, rootHash string) (AnchorResult, error)
}

// noopProvider is the default provider: it never writes to an external ledger,
// so existing flows are unaffected. Anchored stays false and TxHash nil.
type noopProvider struct{}

// NewNoopProvider returns the default no-op anchor provider.
func NewNoopProvider() AnchorProvider {
	return &noopProvider{}
}

func (p *noopProvider) Anchor(ctx context.Context, chainID, rootHash string) (AnchorResult, error) {
	return AnchorResult{Anchored: false}, nil
}
