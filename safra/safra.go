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

func incremente(inteiro *int){
	*inteiro++
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
func redirect(tokens chan Token, mensagem chan Mensagem, neigh Neighbour) {
	pacote := <-neigh.From
	if(pacote.(type)=="Mensagem"){
			mensagem <- pacote
		} else {
			token <- pacote
		}
}

func process(id string, running *sync.WaitGroup, mensagem Mensagem, neighs ...Neighbour) {
	contador := 99999
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Mensagens, 1)  // Canal de Mensagens
	intk := make(chan Token) // Canal de Tokens do Safra

	nmap := make(map[string]Neighbour) // Vizinhos []
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(intk, in, neigh)
	}

	msg := mensagem
	var tk Token
	var ativo sync.WaitGroup()
/*
****************** LAÇO PARA FICAR ESPERANDO CASO RECEBA OUTRA MENSAGEM ***********
*/

	for {

		// Aqui vem o caso base: nó raiz.
		/*
			O nó raiz enviará o token quando enviar as mensagens para geral.
		*/
		if msg.Sender == "init" {
				
			fmt.Printf("* %s é raiz.\n", id)
			// Ao enviar para o próximo, o mensagem terá o "id" do processo atual.
			// mensagem.Sender = id
			// mensagem.Distancia = 0
			msg.Sender = id
			msg.Distancia = 0
			contador = 0
			pai.Id = "init"// Colocar um pai não vazio para evitar erro
			
			// neighs[0].To <- mensagem
			size := len(neighs)
			// for i := 1; i < size; i++ {
			for i := 0; i < size; i++ {
				neighs[i].To <- msg
			}
			
			tk <- intk

			running.Wait()



		} else {
			// Processo não iniciador
			// Blz, aqui que o jogo começa:
			msg := <-in // Processo aguardando receber mensagem
			incremente(&msg.Distancia)
			
			// Mensagem recebida, processo ativo
			ativo.Add(1)
			
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
						neigh.To <- msg
					}
				}

			// Enviou para todos, processo passivo
			ativo.Done()
		}	
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

	var running sync.WaitGroup()
	running.Add(1)
	go process("W", &running, Token{}, Neighbour{"P", pW, wP}, Neighbour{"S", sW, wS})
	go process("S", &running, Token{}, Neighbour{"P", pS, sP}, Neighbour{"W", wS, sW})
	go process("R", &running, Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q", &running, Token{}, Neighbour{"R", rQ, qR})
	go process("P", &running, Token{"init",0}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
	ativo.Wait()
}
