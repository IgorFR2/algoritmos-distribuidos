	/* 
	Implementação do algoritmo de Eleição baseado no de Terry
		
	*/
	package main

	import (
		"fmt"
	)

	type Mensagem struct {
		Sender string
		Eleito string
		Valor int
	}

	type Vizinho struct {
		Id   string
		From chan Mensagem
		To   chan Mensagem
	}

	// Utilizar um unico canal de leitura para mensagens
	func redirect(in chan Mensagem, vizinho Vizinho) {
		for{
			mensagem := <-vizinho.From
			in <- mensagem
		}
	}

	func process(valor int, id string, mensagem Mensagem, vizinhos ...Vizinho) {
		var pai Vizinho
		var eleito string;

		// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
		in := make(chan Mensagem, 1)
		nmap := make(map[string]Vizinho)
		for _, vizinho := range vizinhos {
			nmap[vizinho.Id] = vizinho
			go redirect(in, vizinho)
		}
		var msg Mensagem
		
		if mensagem.Sender == "init" {
			// Processo iniciador
			fmt.Printf("%s é o processo raiz.\n", id)
			var eleito string
			eleito = id	
			size := len(vizinhos)
			msg.Eleito = eleito
			if msg.Valor < valor { 
				msg.Valor = valor
				msg.Eleito = id
			}
			for j := 0; j < 2; j++ {
			   fmt.Printf("\n[DEBUGGER MANUAL] \n[%v] p.eleito:'%v' msg.Eleito:'%v'\n    Processo:'%v' Msg:'%v'\n    Sender:'%v'\n\n", id, eleito, msg.Eleito, valor, msg.Valor, msg.Sender)
				msg.Sender = id
				vizinhos[0].To <- msg // Envia para 1o vizinho
				for i := 1; i < size; i++ {
					msg := <-in // Receber retorno do (i-1)-esimo
					fmt.Printf("[%v] Enviando de %s para %s\n", id, msg.Sender, id)
					if msg.Valor > valor {
						eleito = msg.Eleito // Atualiza o ganhador atual
					} 
					msg.Sender = id
					vizinhos[i].To <- msg  // Enviar para o (i+1)-esimo
				}
				msg = <-in // Esperar retorno do ultimo
				fmt.Printf("[Ultimo] Enviando de %s para %s (%v,%v)\n", msg.Sender, id, msg.Eleito,msg.Valor)
				eleito = msg.Eleito // Atualizar vencedor
				fmt.Printf("[Ultimo] %s é o novo processo eleito.\n", msg.Eleito) // Proximo laço usando o eleito: (Q,4)
			}
			fmt.Println("Fim!")
		} else {
			for msg := range in{
				fmt.Printf("[%v] Enviando de %s para %s (%v,%v)\n", id, msg.Sender, id, msg.Eleito,msg.Valor)
				for _, vizinho := range vizinhos { // Se não tiver pai (""), será quem o enviou
					if pai.Id == "" {
						pai = nmap[msg.Sender]
						fmt.Printf("[%v] %s é pai de %s\n", id, pai.Id, id)
					}
					if msg.Valor > valor {
						eleito = msg.Eleito
						fmt.Printf("[%s] %s é o novo eleito de %s.\n", id, eleito, id)
					} else {
						eleito = id
						msg.Eleito = id
						msg.Valor = valor
						fmt.Printf("[%s] %s é o novo processo liderando a eleição.\n", id, eleito)
					}
					if pai.Id != vizinho.Id {// Entrega o mensagem para o vizinho se ele não for o pai
						msg.Sender = id
						vizinho.To <- msg
						msg = <-in
						fmt.Printf("[%v] Enviando de %s para %s (%v,%v)\n", id, msg.Sender, id, msg.Eleito,msg.Valor)
					}
				}
				msg.Sender = id// Mensagem volta para o pai depois de ter enviado para todos os vizinhos
				pai.To <- msg
			}
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
		
		go process(1, "W", Mensagem{}, Vizinho{"P", pW, wP}, Vizinho{"S", sW, wS})
		go process(2,"S", Mensagem{}, Vizinho{"P", pS, sP}, Vizinho{"W", wS, sW})
		go process(3,"R", Mensagem{}, Vizinho{"Q", qR, rQ}, Vizinho{"P", pR, rP})
		go process(4,"Q", Mensagem{}, Vizinho{"R", rQ, qR})
		process(2,"P", Mensagem{"init","",-1}, Vizinho{"W", wP, pW}, Vizinho{"S", sP, pS}, Vizinho{"R", rP, pR})
	}
