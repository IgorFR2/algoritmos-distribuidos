/* 
Implementação do algoritmo de travessia de Tarry
	1) Raiz recebe mensagem enviado por "init"
	2) Processo assina o mensagem e envia para um filho e espera mensagem para enviar ao próximo
	2.1) Processo registra como "pai" o primeiro a enviar mensagem
	2.2) Todos vizinhos não "pai" são "filhos"
	3) Processo envia mensagem para "pai"
	4) Raiz encerra programa
*/
package main

import (
	"fmt"
)

type Mensagem struct {
	Sender string
}

type Vizinho struct {
	Id   string
	From chan Mensagem
	To   chan Mensagem
}

// Utilizar um unico canal de leitura para mensagens
func redirect(in chan Mensagem, vizinho Vizinho) {
	mensagem := <-vizinho.From
	in <- mensagem
}

func process(id string, mensagem Mensagem, vizinhos ...Vizinho) {
	var pai Vizinho

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Mensagem, 1)
	nmap := make(map[string]Vizinho)
	for _, vizinho := range vizinhos {
		nmap[vizinho.Id] = vizinho
		go redirect(in, vizinho)
	}

	if mensagem.Sender == "init" {
		// Processo iniciador
		fmt.Printf("%s é o processo raiz.\n", id)
		// Como iniciador não tem pai, o mensagem terá "Init"
		// Ao enviar para o próximo, o mensagem terá o "id" do processo atual.
		mensagem.Sender = id
		vizinhos[0].To <- mensagem
		size := len(vizinhos)
		for i := 1; i < size; i++ {
			tk := <-in
			fmt.Printf("[%v] Enviando de %s para %s\n", id, tk.Sender, id)
			tk.Sender = id
			vizinhos[i].To <- tk
		}
		tk := <-in
		fmt.Printf("[%v] Enviando de %s para %s\n", id, tk.Sender, id)
		fmt.Println("Fim!")
	} else {
		// Processo não iniciador. Passar mensagem adiante.
		tk := <-in
		fmt.Printf("[%v] Enviando de %s para %s\n", id, tk.Sender, id)
		// Se não tiver pai (""), será quem o enviou
		for _, vizinho := range vizinhos {
			if pai.Id == "" {
				pai = nmap[tk.Sender]
				fmt.Printf("[%v] %s é pai de %s\n", id, pai.Id, id)
			}
			// Entrega o mensagem para o vizinho se ele não for o pai
			if pai.Id != vizinho.Id {
				tk.Sender = id
				vizinho.To <- tk
				tk = <-in
				fmt.Printf("[%v] Enviando de %s para %s\n", id, tk.Sender, id)
			}
		}
		// Mensagem volta para o pai depois de ter enviado para todos os vizinhos
		tk.Sender = id
		pai.To <- tk
	}

}

func main() {

	pW := make(chan Mensagem, 1)
	pS := make(chan Mensagem, 1)
	pR := make(chan Mensagem, 1)
	wP := make(chan Mensagem, 1)
	wS := make(chan Mensagem, 1)
	sP := make(chan Mensagem, 1)
	sW := make(chan Mensagem, 1)
	rQ := make(chan Mensagem, 1)
	rP := make(chan Mensagem, 1)
	qR := make(chan Mensagem, 1)

	go process("W", Mensagem{}, Vizinho{"P", pW, wP}, Vizinho{"S", sW, wS})
	go process("S", Mensagem{}, Vizinho{"P", pS, sP}, Vizinho{"W", wS, sW})
	go process("R", Mensagem{}, Vizinho{"Q", qR, rQ}, Vizinho{"P", pR, rP})
	go process("Q", Mensagem{}, Vizinho{"R", rQ, qR})
	process("P", Mensagem{"init"}, Vizinho{"W", wP, pW}, Vizinho{"S", sP, pS}, Vizinho{"R", rP, pR})
}
