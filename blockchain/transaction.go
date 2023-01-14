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

var Mempool *mempool = &mempool{} //memory 에 올라와 있음 , mempool을 위한 데이터 베이스도 필요 없음

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
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
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

	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{Owner: txOut.Owner, Amount: txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}

	if total > amount {
		txOut := &TxOut{Owner: from, Amount: total - amount}
		txOuts = append(txOuts, txOut)
	}
	txOuts = append(txOuts, &TxOut{Owner: to, Amount: amount})

	tx := &Tx{
		Id:        "",
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
