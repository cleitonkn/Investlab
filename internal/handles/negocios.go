package handles

import (
	"errors"

	"../querys"
	"../structs"
)

func insertDinheiro(user string, valor float64) {
	trans := new(structs.Transacao)
	trans.ID = querys.QueryGetUseID(user)
	trans.Descricao = "deposito"
	trans.Tipo = 1
	trans.Valor = valor
	querys.InsertTransacao(*trans)
}

func insertCompra(userCompra, codigo string, valor float64, qtd int) error {

	lances := querys.QueryGetCarteiraUser(codigo, valor)
	for i := 0; i < qtd; i++ {

		transVenda := new(structs.Transacao)
		transVenda.ID = lances[i].IDUser
		transVenda.Tipo = 2
		transVenda.Valor = valor
		transVenda.Descricao = "venda acao:" + codigo

		if transVenda.ID == -1 {
			err := errors.New("usuario de venda não encontrado")
			return err
		}

		transCompra := new(structs.Transacao)
		transCompra.ID = querys.QueryGetUseID(userCompra)
		transCompra.Descricao = "compra acao" + codigo
		transCompra.Tipo = 3
		transCompra.Valor = -valor

		if transCompra.ID == -1 {
			err := errors.New("usuario de compra não encontrado")
			return err
		}
		saldo := querys.QueryGetSaldo(transCompra.ID).Valor

		if saldo < valor {
			err := errors.New("sem saldo")
			return err
		}

		carteira := new(structs.Carteira)
		carteira.Codigo = codigo
		carteira.ID = transCompra.ID
		carteira.Lance = 0
		carteira.Venda = false
		querys.InsertCarteira(*carteira)
		querys.RemoveCarteira(lances[i].IDCarteira)
		querys.InsertTransacao(*transVenda)
		querys.InsertTransacao(*transCompra)
	}

	return nil
}
