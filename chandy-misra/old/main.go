/* 
Implementação do algoritmo de Chandy-Misra derivado da travessia de Tarry 

Funcionamento do Algoritmo Chandy Misra:
1) Processo verifica conteúdo do token
1.1) Se token == 'init' então ele é o raiz
1.2.1) Se não ele adiciona a distancia do antecessor somado a distancia no token
1.2.2) Se a soma for menor que a do pai, o antecessor vira o novo pai
2) Escreve o nome do processo no token e a nova distancia percorrida
3) Enviar token para todos vizinhos
4.1) Enviar token para o pai 
4.2) Encerrar se for nó raiz (pai=="É raiz")

*/
package main

import ("fmt")

// Função auxiliar 
func incremente(inteiro *int){
	//inteiro = inteiro +1
	*inteiro++
}

// Token do algoritmo travis. Contém só a string de quem envia.
type Token struct {
	Sender string
	// Token terá distancia percorrida 
	Distancia int
}

type Neighbour struct {
	Id   string
	From chan Token
	To   chan Token
}

// Ao que entendi aqui é um funil.
// Nesta função ele recebe de um vizinho e direciona para a saida unica.
func redirect(in chan Token, neigh Neighbour) {
	token := <-neigh.From
	in <- token
}

// Blz, aqui tem o principal: processos.
// O processo tem vizinhos, dos quais recebe mensagens e envia.
// Vizinho tem um pai, só envia mensagem para o pai quando mandar pra todo mundo e receber retorno. 
func process(id string, running *sync.WaitGroup, token Token, neighs ...Neighbour) {
	contador := 4294967295 // Esse cara vai armazenar a distancia até a raiz. MaxUint32 = 4294967295
	// Lembrando que vizinho tem 2 canais (entrada,saida : From/To) e um nome (Id)
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Token, 1) // Canal genérico de entrada
	nmap := make(map[string]Neighbour) // Ao que parece essa função faz um VETOR com os nomes dos vizinhos
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh // É, realmente a estrutura lembra o dict do python em que vc pode chamar dict["chave"]
		go redirect(in, neigh)
	}

	// Aqui vem o caso base: nó raiz.
	// 1) Processo verifica conteúdo do token
	// 1.1) Se token == 'init' então ele é o raiz
	if token.Sender == "init" {
		// Processo iniciador
		fmt.Printf("* %s é raiz.\n", id)
		// Como iniciador não tem pai, o token terá "Init"
		// Ao enviar para o próximo, o token terá o "id" do processo atual.
		token.Sender = id
		token.Distancia = 0 // Distancia raiz->raiz é zero, quem receber o token incremente
		contador = 0
		neighs[0].To <- token
		size := len(neighs)
		for i := 1; i < size; i++ {
			// Blz, já enviei para o vizinho 0, agora mandar para o resto.
			// tk eu não sei o que é, parece mais uma variavel temporária não declarada
			// * Esse tk está esperando chegar mensagem! Ele bloqueia o processo até receber algo.
			// De qualquer forma ta recebendo mensagem do canal unificado (in),
			// assinando o toke e repassando
			tk := <-in
			fmt.Printf("From %s to %s\n", tk.Sender, id)
			tk.Sender = id
			tk.Distancia = contador // Bom, ao que parece o raiz só fica aqui.
			neighs[i].To <- tk
		}
		// Aqui está esperando receber resposta do ultimo vizinho
		// Depois disso vai enviar para o pai.
		tk := <-in
		fmt.Printf("[%s] From %s to %s. (Raiz) Contador: %v / Token Dist: %v\n",id, tk.Sender, id, contador, tk.Distancia)
		fmt.Println("Fim!")
	} else {
		// Processo não iniciador
		// Blz, aqui que o jogo começa:
		tk := <-in // Processo aguardando receber token
		incremente(&tk.Distancia)
		
		
		fmt.Printf("[%s] From %s to %s. (Fora) Contador: %v / Token Dist: %v\n",id, tk.Sender, id, contador, tk.Distancia)
		for _, neigh := range neighs {
			// Blz, verificar se o cara tem pai.
			if pai.Id == "" {
				pai = nmap[tk.Sender]
				fmt.Printf("[%s] %s é pai de %s (Orfão). Contador: %v / Token Dist: %v\n", id, pai.Id, id, contador, tk.Distancia)
				contador = tk.Distancia // Aqui a distancia do token já ta incrementada, não tendo pai o contador ta max_int.
			} else if contador > tk.Distancia{
				// Se o novo cara for melhor que o pai, ele será o novo pai e atualiza contador
				pai = nmap[tk.Sender]  // Novo pai
				fmt.Printf("[%s] %s é o novo pai de %s (Orfão). Contador: %v / Token Dist: %v\n", id, pai.Id, id, contador, tk.Distancia)
				contador = tk.Distancia // Contador fica com a distancia acumulada pelo token
			}
			// Entrega o token para o vizinho, se ele não for o pai
			if pai.Id != neigh.Id {
				tk.Sender = id
				neigh.To <- tk
				tk = <-in
				fmt.Printf("[%s] From %s to %s. (Dentro) Contador: %v / Token Dist: %v\n", id, tk.Sender, id,contador, tk.Distancia)
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

	var running sync.WaitGroup()
	running.Add(1)
	go process("W", &running, Token{}, Neighbour{"P", pW, wP}, Neighbour{"S", sW, wS})
	go process("S", &running, Token{}, Neighbour{"P", pS, sP}, Neighbour{"W", wS, sW})
	go process("R", &running, Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q", &running, Token{}, Neighbour{"R", rQ, qR})
	go process("P", &running, Token{"init",0}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
	running.Wait()
}
