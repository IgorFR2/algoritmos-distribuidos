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
func incremente(inteiro *int){
	*inteiro++
}
type Token struct {
	Sender string
	Distancia int
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
	contador := 99999
	var pai Neighbour

	// Redeirecionando todos os canais de entrada para um único canal "in" de entrada
	in := make(chan Token, 1)  // Unificados
	
	nmap := make(map[string]Neighbour) // Vizinhos []
	for _, neigh := range neighs {
		nmap[neigh.Id] = neigh
		go redirect(in, neigh)
	}
	tk := token
/*
****************** LAÇO PARA FICAR ESPERANDO CASO RECEBA OUTRA MENSAGEM ***********
*/

	for {

		// Aqui vem o caso base: nó raiz.
		// 1) Processo verifica conteúdo do token
		// 1.1) Se token == 'init' então ele é o raiz
		// if token.Sender == "init" {
		if tk.Sender == "init" {
			// Processo iniciador
			fmt.Printf("* %s é raiz.\n", id)
			// Ao enviar para o próximo, o token terá o "id" do processo atual.
			// token.Sender = id
			// token.Distancia = 0
			tk.Sender = id
			tk.Distancia = 0
			contador = 0
			pai.Id = "init"// Colocar um pai não vazio para evitar erro
			
			// neighs[0].To <- token
			size := len(neighs)
			// for i := 1; i < size; i++ {
			for i := 0; i < size; i++ {
				// tk := <-in // Aqui ele fica parado esperando
				// fmt.Printf("From %s to %s\n", tk.Sender, id)
				// tk.Sender = id
				// tk.Distancia = contador // Bom, ao que parece o raiz só fica aqui.
				neighs[i].To <- tk
			}
			// Aqui está esperando receber resposta do ultimo vizinho
			// Depois disso vai enviar para o pai.
			// tk := <-in
			// fmt.Printf("[%s] From %s to %s. (Raiz) Contador: %v / Token Dist: %v\n",id, tk.Sender, id, contador, tk.Distancia)
			// fmt.Println("Fim!")
			} else {
				// Processo não iniciador
				// Blz, aqui que o jogo começa:
				tk := <-in // Processo aguardando receber token
				incremente(&tk.Distancia)
				
				
				fmt.Printf("[%s] From %s to %s. (Fora) Contador: %v / Token Dist: %v\n",id, tk.Sender, id, contador, tk.Distancia)
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
				for _, neigh := range neighs {
						// Entrega o token para o vizinho, se ele não for o pai
						if pai.Id != neigh.Id {
							tk.Sender = id
							neigh.To <- tk
							// tk = <-in
							// fmt.Printf("[%s] From %s to %s. (Dentro) Contador: %v / Token Dist: %v\n", id, tk.Sender, id,contador, tk.Distancia)
						}
					}
					// Token volta para o pai depois de ter enviado para todos os vizinhos
					// tk.Sender = id
					// pai.To <- tk
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

	go process("W", Token{}, Neighbour{"P", pW, wP}, Neighbour{"S", sW, wS})
	go process("S", Token{}, Neighbour{"P", pS, sP}, Neighbour{"W", wS, sW})
	go process("R", Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP})
	go process("Q", Token{}, Neighbour{"R", rQ, qR})
	process("P", Token{"init",0}, Neighbour{"W", wP, pW}, Neighbour{"S", sP, pS}, Neighbour{"R", rP, pR})
}
