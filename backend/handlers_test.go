package main

import "testing"

func TestIsValidStatus(t *testing.T) {
	// Testa um cenário de sucesso
	if !isValidStatus("A Fazer") {
		t.Errorf("Esperava true para 'A Fazer', recebeu false")
	}

	// Testa um cenário de falha (defensivo)
	if isValidStatus("Finalizado") {
		t.Errorf("Esperava false para um status inválido como 'Finalizado'")
	}
}