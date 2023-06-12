package mempool

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

type MempoolClient struct {
	baseURL string
}

func NewClient(netParams *chaincfg.Params) *MempoolClient {
	baseURL := ""
	if netParams.Net == wire.MainNet {
		baseURL = "https://mempool.space/api"
	} else if netParams.Net == wire.TestNet3 {
		baseURL = "https://mempool.space/testnet/api"
	} else if netParams.Net == chaincfg.SigNetParams.Net {
		baseURL = "https://mempool.space/signet/api"
	} else {
		log.Fatal("mempool don't support other netParams")
	}
	return &MempoolClient{
		baseURL: baseURL,
	}
}

func (c *MempoolClient) request(method, subPath string, requestBody io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, subPath)
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	return body, nil
}

type UnspentOutput struct {
	Outpoint *wire.OutPoint
	Output   *wire.TxOut
}

type BTCAPIClient interface {
	GetRawTransaction(txHash *chainhash.Hash) (*wire.MsgTx, error)
	BroadcastTx(tx *wire.MsgTx) (*chainhash.Hash, error)
	ListUnspent(address btcutil.Address) ([]*UnspentOutput, error)
}

var _ BTCAPIClient = (*MempoolClient)(nil)
