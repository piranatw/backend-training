package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
)

// Returns a random integer from 0 to cap (inclusive)
func randInt(cap int) int {
    seed := time.Now().UnixNano()
    source := rand.NewSource(seed)
    randInstance := rand.New(source)
    return randInstance.Intn(cap + 1)
}

func main() {
    cap := 0
    fmt.Print("Hi, this is a number guessing game. Please enter the upper bound: ")
    fmt.Scan(&cap)

    goal := randInt(cap)
    guess := 0

    maximumAttempts := int(math.Ceil(math.Log2(float64(cap))))
    for try := 0; try < maximumAttempts; try++ {
        fmt.Print("Please enter your guess: ")
        fmt.Scan(&guess)
        switch {
        case guess == goal:
            fmt.Println("Correct! You got it right.")
            return
        case guess > goal:
            fmt.Println("Too high. Try again.")
        default:
            fmt.Println("Too low. Try again.")
        }
    }
    fmt.Println("You ran out of attempts. Game over.")
}