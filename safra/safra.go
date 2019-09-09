/* 
Implementação do algoritmo de término de Safra derivado do Chandy-Misra.

Funcionamento do Algoritmo Safra:
1) Funciona paralelamente
2) Token de safra só envia quando processo não estiver enviando mensagem.
3) Processo terá contador extra para mensagens
3.1) Contador mensagens incrementa ao enviar e decrementa ao receber
4) Token acumula contadores de mensagens
4.1) Processos serão acumulados apenas uma vez por rodada, indicados por marcador
4.2) Processo é marcado ao 


Processo recebe mensagem -> msgcounter--
Raiz recebe Stk do ultimo vizinho -> Soma == 0? Termine ; Repita
*/
package main
import (
	"fmt"
	"sync"
)

func incremente(inteiro *int, distancia int){
	*inteiro += distancia
}

/*
	Apenas Safra enviará Token. 
	Chandy-Misra enviará Mensagens.
*/
type Token struct {
	Sender string
	Contador int
}

type Mensagem struct {
	Sender string
	Distancia int
}

type Neighbour struct {
	Id   string
	Distancia int
	From chan interface{} // Assim que usa a interface: agora é um canal que aceita qualquer coisa 
	To   chan interface{} // (usar mecanismos de tipo para direcionar os tokens aos lugares corretos)
}
/*
	Redirect agora irá receber os 2 canais: para token e para mensagens
*/
func redirect(tokens chan Token, mensagens chan Mensagem, neigh Neighbour) {
	pacote := <-neigh.From
	switch pacote.(type){
	case Mensagem:
			mensagens <- pacote.(Mensagem)
	case Token:
			tokens <- pacote.(Token)
	}
}

func safra(id string, contador *int, running *sync.WaitGroup, ativo *sync.WaitGroup, tokens chan Token, vizinhos []Neighbour, nmap map[string]Neighbour){
	// Aguardar entrada de token
	var pai Neighbour
	for token := range tokens{
		if token.Sender == "init" {
			for {
				fmt.Printf("\t\t\t\t[%s][Safra] Iniciado.\n", id)
				token.Sender = id
				ativo.Wait()

				vizinhos[0].To <- token
				size := len(vizinhos)
				for i := 1; i < size; i++ {
					token := <-tokens
					fmt.Printf("\t\t\t\t[%v][Safra] Enviando de %s para %s\n", id, token.Sender, id)
					token.Sender = id
					ativo.Wait()

					vizinhos[i].To <- token
				}
				token := <-tokens
				fmt.Printf("\t\t\t\t[%v][Safra] Enviando de %s para %s\n", id, token.Sender, id)
				if(token.Contador == 0){
					fmt.Println("Fim!")
					break
				}
			}
		} else {
			if pai.Id == "" {
				pai = nmap[token.Sender]
			}

			fmt.Printf("\t\t\t\t[%v] Enviando de %s para %s\n", id, token.Sender, id)
			token.Contador = token.Contador + *contador
			token.Sender = id
			vizinhos[0].To <- token
			size := len(vizinhos)
			for i := 1; i < size; i++ {
				token := <-tokens
				fmt.Printf("\t\t\t\t[%v][Safra] Enviando de %s para %s\n", id, token.Sender, id)
				token.Sender = id
				vizinhos[i].To <- token
			}
			token := <-tokens
			fmt.Printf("\t\t\t\t[%v][Safra] Enviando de %s para %s\n", id, token.Sender, id)
			pai.To <- token
		}
		running.Done()

	}
}

func process(id string, running *sync.WaitGroup, mensagem Mensagem, neighs ...Neighbour) {
	contador := 99999
	contadorMsg := 0
	var pai Neighbour
	
	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Mensagem)//, 1)  // Canal de Mensagens
	intk := make(chan Token) // Canal de Tokens do Safra
	nmap := make(map[string]Neighbour) // Vizinhos []
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(intk, in, neigh)
	}
	//msg := mensagem
	// in <- mensagem
	//var tk Token
	var ativo sync.WaitGroup
/*
****************** LAÇO PARA FICAR ESPERANDO CASO RECEBA OUTRA MENSAGEM ***********
*/
	// for {
		// Aqui vem o caso base: nó raiz.
		/*
			O nó raiz enviará o token quando enviar as mensagens para geral.
		*/
		if mensagem.Sender == "init" {

			// Token Ativo
			ativo.Add(1)
			fmt.Printf("* %s é raiz.\n", id)
			// Ao enviar para o próximo, o mensagem terá o "id" do processo atual.
			// mensagem.Sender = id
			// mensagem.Distancia = 0
			mensagem.Sender = id
			mensagem.Distancia = 0
			contador = 0
			pai.Id = "init"// Colocar um pai não vazio para evitar erro
			
			// neighs[0].To <- mensagem
			size := len(neighs)
			// for i := 1; i < size; i++ {
			for i := 0; i < size; i++ {
				neighs[i].To <- mensagem
			}
			ativo.Done()
			go safra(id, &contadorMsg, running, &ativo , intk, neighs, nmap)
			
			for range in{
				contadorMsg--
			}
			// tk <- int

		} else {
			go safra(id, &contadorMsg, running, &ativo , intk, neighs, nmap)
			for mgs := range in{
				// Mensagem recebida, processo ativo
				ativo.Add(1)
	
				// Processo não iniciador
				// Blz, aqui que o jogo começa:
				// msg := <-in // Processo aguardando receber mensagem
				incremente(&msg.Distancia, nmap[msg.Sender].Distancia)
				contadorMsg--
				
				fmt.Printf("[%s] From %s to %s. (Fora) Contador: %v / Token Dist: %v\n",id, msg.Sender, id, contador, msg.Distancia)
				// Blz, verificar se o cara tem pai.
				if pai.Id == "" {
					pai = nmap[msg.Sender]
					fmt.Printf("[%s] %s é pai de %s (Orfão). Contador: %v / Token Dist: %v\n", id, pai.Id, id, contador, msg.Distancia)
					contador = msg.Distancia // Aqui a distancia do mensagem já ta incrementada, não tendo pai o contador ta max_int.
					} else if contador > msg.Distancia{
						// Se o novo cara for melhor que o pai, ele será o novo pai e atualiza contador
						pai = nmap[msg.Sender]  // Novo pai
						fmt.Printf("[%s] %s é o novo pai de %s (Orfão). Contador: %v / Token Dist: %v\n", id, pai.Id, id, contador, msg.Distancia)
						contador = msg.Distancia // Contador fica com a distancia acumulada pelo mensagem
					}
				for _, neigh := range neighs {
						// Entrega o mensagem para o vizinho, se ele não for o pai
						if pai.Id != neigh.Id {
							msg.Sender = id
							// contadorMsg++ ??
							neigh.To <- msg
							contadorMsg++
						}
					}
	
				// Enviou para todos, processo passivo
				ativo.Done()
			}
		}	
	}
}

		
func main() {

	pW := make(chan interface{}, 1)
	pS := make(chan interface{}, 1)
	pR := make(chan interface{}, 1)
	wP := make(chan interface{}, 1)
	wS := make(chan interface{}, 1)
	sP := make(chan interface{}, 1)
	sW := make(chan interface{}, 1)
	rQ := make(chan interface{}, 1)
	rP := make(chan interface{}, 1)
	qR := make(chan interface{}, 1)

	var running sync.WaitGroup

	running.Add(1)
	go process("W", &running, Mensagem{}, Neighbour{"P", 1, pW, wP}, Neighbour{"S", 1, sW, wS})
	go process("S", &running, Mensagem{}, Neighbour{"P", 1, pS, sP}, Neighbour{"W", 1, wS, sW})
	go process("R", &running, Mensagem{}, Neighbour{"Q", 1, qR, rQ}, Neighbour{"P", 1, pR, rP})
	go process("Q", &running, Mensagem{}, Neighbour{"R", 1, rQ, qR})
	go process("P", &running, Mensagem{"init",0}, Neighbour{"W", 1, wP, pW}, Neighbour{"S", 1, sP, pS}, Neighbour{"R", 1, rP, pR})
	running.Wait()
}
