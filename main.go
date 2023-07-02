package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type RicartAgrawala struct {
	me            int
	totalProcs    int
	timestamp     int
	responseCount int
	wantCS        bool
	replies       chan bool
}

func main() {
	// Obtener los argumentos de la línea de comandos
	args := os.Args[1:]
	totalProcs, _ := strconv.Atoi(args[0])
	rows, _ := strconv.Atoi(args[1])
	cols, _ := strconv.Atoi(args[2])
	filename := args[3]

	// Leer el archivo y construir la matriz
	matrix := readMatrixFromFile(filename, rows, cols)

	// Crear los canales para las respuestas
	replies := make(chan bool, totalProcs-1)

	// Crear los procesos
	processes := make([]*RicartAgrawala, totalProcs)
	for i := 0; i < totalProcs; i++ {
		processes[i] = &RicartAgrawala{
			me:            i + 1,
			totalProcs:    totalProcs,
			replies:       replies,
		}
	}

	// Iniciar los procesos en paralelo
	var wg sync.WaitGroup
	wg.Add(totalProcs)
	for i := 0; i < totalProcs; i++ {
		go func(i int) {
			defer wg.Done()
			processes[i].start(rows, matrix)
		}(i)
	}

	// Esperar a que todos los procesos terminen
	wg.Wait()

	// Imprimir la matriz resultante
	printMatrix(matrix)
}

func readMatrixFromFile(filename string, rows, cols int) ([][]string) {
	file, _ := os.Open(filename)
	defer file.Close()

	matrix := make([][]string, rows)
	scanner := bufio.NewScanner(file)
	for i := 0; i < rows && scanner.Scan(); i++ {
		line := scanner.Text()
		println(line)
		matrix[i] = strings.Fields(line)
	}

	return matrix
}

func (ra *RicartAgrawala) start(rows int, matrix [][]string) {
	for row := 0; row < rows; row++ {
		ra.requestCriticalSection(row, matrix)
	}
}

func (ra *RicartAgrawala) requestCriticalSection(row int, matrix [][]string) {
	ra.timestamp++
	ra.wantCS = true

	for proc := 0; proc < ra.totalProcs; proc++ {
		if proc+1 == ra.me {
			continue
		}
		go func(proc int) {
			ra.replyHandler(proc)
		}(proc)
	}

	ra.checkCriticalSection(row, matrix)
}

func (ra *RicartAgrawala) replyHandler(proc int) {
	for ra.wantCS {
		// Simular una respuesta retrasada
		// time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		// Recibir respuestas
		ra.replies <- true
	}
}

func (ra *RicartAgrawala) checkCriticalSection(row int, matrix [][]string) {
	for ra.responseCount < ra.totalProcs-1 {
		select {
		case <-ra.replies:
			ra.responseCount++
		}
	}

	ra.enterCriticalSection(row, matrix)
	ra.exitCriticalSection()
}

func (ra *RicartAgrawala) enterCriticalSection(row int, matrix [][]string) {
	palindromes := findPalindromesInRow(row, matrix)
	if len(palindromes) > 0 {
		for _, p := range palindromes {
			replacePalindrome(row, p, matrix, ra.me)
		}
		fmt.Printf("Proceso %d encuentra palíndromo(s) en la fila %d\n", ra.me, row+1)
	} else {
		fmt.Printf("Proceso %d no encontró ningún palíndromo en la fila %d\n", ra.me, row+1)
	}
}

func (ra *RicartAgrawala) exitCriticalSection() {
	ra.timestamp++
	ra.wantCS = false
	ra.responseCount = 0
}

func findPalindromesInRow(row int, matrix [][]string) []string {
	palindromes := make([]string, 0)

	rowString := strings.Join(matrix[row], "")
	println(rowString)
	n := len(rowString)

	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			substring := rowString[i:j]
			if isPalindrome(substring) {
				palindromes = append(palindromes, substring)
			}
		}
	}

	return palindromes
}

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func replacePalindrome(row int, palindrome string, matrix [][]string, proc int) {
	for col := range matrix[row] {
		if strings.HasPrefix(matrix[row][col], palindrome) {
			matrix[row][col] = strconv.Itoa(proc)
		}
	}
}

func printMatrix(matrix [][]string) {
	for _, row := range matrix {
		fmt.Println(strings.Join(row, " "))
	}
}
