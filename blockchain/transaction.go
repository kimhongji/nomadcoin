package blockchain

import (
	"errors"
	"time"

	"github.com/kimhongji/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool = &mempool{} //memory 에 올라와 있음 , mempool을 위한 데이터 베이스도 필요 없음

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	TxId  string `json:"txId"`  // input 으로 사용할 output 의 tx id
	Index int    `json:"index"` // 위 tx 의 output 들 중에서 어떤 txOut 인지 index
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// UTxOut 어떤 output 이 쓰였는지 , 안쓰였는지 확인할 수 있게 도와주는 구조체
type UTxOut struct {
	TxId   string
	Index  int
	Amount int
}

func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if txIn.TxId == uTxOut.TxId {
				if txIn.Index == uTxOut.Index {
					return true
				}
			}
		}
	}
	return false
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	// 원래는 인풋 들어가면 아웃풋 같은거 자동으로 뭔가 처리 되어야 함 처음이라 이런 구조
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx
}

// txInput:  [2(from), 2(from)]
// amount: 3
// txOutput: [3(to), 1(me)]
func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := Blockchain().UTxOutsByAddress(from)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{uTxOut.TxId, uTxOut.Index, from})
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	tx := &Tx{
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("hayz", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("hayz") //block 을 채굴했을 때만
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
