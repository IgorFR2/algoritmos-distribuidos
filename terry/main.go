// Implementação do algoritmo de travessia de Tarry
package main

import (
	"fmt"
)

type Token struct {
	Sender string
}

type Neighbour struct {
	Id   string
	From chan Token
	To   chan Token
}

func redirect(in chan Token, neigh Neighbour) {
	token := <-neigh.From
	in <- token
}

func process(id string, token Token, neighs ...Neighbour) {
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Token, 1)
	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(in, neigh)
	}

	if token.Sender == "init" {
		// Processo iniciador
		fmt.Printf("* %s é raiz.\n", id)
		// Como iniciador não tem pai, o token terá "Init"
		// Ao enviar para o próximo, o token terá o "id" do processo atual.
		token.Sender = id
		neighs[0].To <- token
		size := len(neighs)
		for i := 1; i < size; i++ {
			tk := <-in
			fmt.Printf("[%v] From %s to %s\n", id, tk.Sender, id)
			tk.Sender = id
			neighs[i].To <- tk
		}
		tk := <-in
		fmt.Printf("[%v] From %s to %s\n", id, tk.Sender, id)
		fmt.Println("Fim!")
	} else {
		// Processo não iniciador. Passar token.
		tk := <-in
		fmt.Printf("[%v] From %s to %s\n", id, tk.Sender, id)
		// Se não tiver pai (""), será quem o enviou
		for _, neigh := range neighs {
			if pai.Id == "" {
				pai = nmap[tk.Sender]
				fmt.Printf("[%v] %s é pai de %s\n", id, pai.Id, id)
			}
			// Entrega o token para o vizinho se ele não for o pai
			if pai.Id != neigh.Id {
				tk.Sender = id
				neigh.To <- tk
				tk = <-in
				fmt.Printf("[%v] Recebido de %s em %s\n", id, tk.Sender, id)
			}
		}
		// Token volta para o pai depois de ter enviado para todos os vizinhos
		tk.Sender = id
		pai.To <- tk
	}

}

func main() {

	pW := make(chan Token, 1)
	pS := make(chan Token, 1)
	pR := make(chan Token, 1)
	wP := make(chan Token, 1)
	wS := make(chan Token, 1)
	sP := make(chan Token, 1)
	sW := make(chan Token, 1)
	rQ := make(chan Token, 1)
	rP := make(chan Token, 1)
	qR := make(chan Token, 1)

	go process("W", Token{}, Neighbour{"P", pW, wP}, Neighbour{"S", sW, wS})
	go process("S", Token{}, Neighbour{"P", pS, sP}, Neighbour{"W", wS, sW})
	go process("R", Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q", Token{}, Neighbour{"R", rQ, qR})
	process("P", Token{"init"}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
}
