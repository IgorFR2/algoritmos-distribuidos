package main
/**
Relogio vetorial
	1) Processo realiza ação incrementa vetor contador na posição  do seu PID
	2) Preocesso envia mensagem para demais com vetor de contadores
	3) Processo recebe mensagem com vetor de contadores
	4) Processo varre todas as posições e atualiza os que forem maiores que o seu (o de sua posição sempre será menor igual)
*/
import (
	"fmt"
	"time"
)

// Mensagem agora envia o vetor em vez do valor
type Message struct {
	Body      string
	Timestamp [3]int
}

// 1) Evento agora incrementa o contador na posição do processo em vez de uma variavel inteira
// Cuidado que vetor é [0,1,2...] e o 1º PID é 1, logo o 3º tenta acessar 4ªposição
func event(pid int, counter [3]int) [3]int {
	counter[pid-1] += 1
	fmt.Printf("Event in process pid=%v. Counter=%v. Counters=%v\n", pid, counter[pid-1],counter)
	return counter
}

// Função auxiliar, não muda.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// 4) Atualizar os valores do vetor. Counter[pid] sempre será >= que a Msg
func calcTimestamp(recvTimestamp, counter [3]int) [3]int {
	for i:=0; i < len(counter);i++{
		counter[i] = max(counter[i], recvTimestamp[i])
	}
	// return max(recvTimestamp, counter) + 1
	return counter
}

// 2) Enviar mensagem da mesma forma, porém com vetor de contadores.
func sendMessage(ch chan Message, pid int, counter [3]int) [3]int {
	//counter[pid-1] += 1
	counter[pid-1]++
	ch <- Message{"Test msg!!!", counter}
	fmt.Printf("Message sent from pid=%v. Counter=%v.\n", pid, counter)
	return counter

}
// 3) Recebimento com vetor, irá atualizar no processo o que for necessário
func receiveMessage(ch chan Message, pid int, counter [3]int) [3]int {
	message := <-ch
	counter[pid-1]++
	counter = calcTimestamp(message.Timestamp, counter)
	fmt.Printf("Message received at pid=%v. Counter=%v\n", pid, counter)
	return counter
}


func processOne(ch12, ch21 chan Message) {
	pid := 1
	counter := [3]int{0,0,0}
	counter = event(pid, counter)
	counter = sendMessage(ch12, pid, counter)
	counter = event(pid, counter)
	counter = receiveMessage(ch21, pid, counter)
	counter = event(pid, counter)
}

func processTwo(ch12, ch21, ch23, ch32 chan Message) {
	pid := 2
	counter :=  [3]int{0,0,0}
	counter = receiveMessage(ch12, pid, counter)
	counter = sendMessage(ch21, pid, counter)
	counter = sendMessage(ch23, pid, counter)
	counter = receiveMessage(ch32, pid, counter)

}

func processThree(ch23, ch32 chan Message) {
	pid := 3
	counter :=  [3]int{0,0,0}
	counter = receiveMessage(ch23, pid, counter)
	counter = sendMessage(ch32, pid, counter)

}

func main() {
	oneTwo := make(chan Message, 100)
	twoOne := make(chan Message, 100)
	twoThree := make(chan Message, 100)
	threeTwo := make(chan Message, 100)

	go processOne(oneTwo, twoOne)
	go processTwo(oneTwo, twoOne, twoThree, threeTwo)
	go processThree(twoThree, threeTwo)

	time.Sleep(5 * time.Second)
}
