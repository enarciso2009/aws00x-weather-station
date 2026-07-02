package main

import (
    "fmt"
    "net"
)

func main() {
    ln, err := net.Listen("tcp", ":8899")
    if err != nil {
        panic(err)
    }

    fmt.Println("Escutando na porta 8899")

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Erro ao aceitar conexão:", err)
            continue
        }

        fmt.Println("Conectado:", conn.RemoteAddr())

        go handleConnection(conn)
    }
}

func handleConnection(c net.Conn) {
    defer c.Close()

    buf := make([]byte, 4096)

    for {
        n, err := c.Read(buf)
        if err != nil {
            fmt.Println("Desconectado")
            return
        }

        data := buf[:n]

        // Ignora o pacote de registro do G781
        if n == 10 && string(data) == "www.usr.cn" {
            fmt.Println("Pacote de registro recebido.")
            continue
        }

        // Garante tamanho mínimo esperado
        if len(data) < 22 {
            continue
        }

        // Temperatura
        temp := float64(data[17]) / 10.0

        // Umidade
        humidity := float64(
            uint16(data[18])<<8|uint16(data[19]),
        ) / 10.0

        // Pressão
        pressure := float64(
            uint16(data[20])<<8|uint16(data[21]),
        ) / 10.0

        fmt.Println()
        fmt.Println("=====================================")
        fmt.Printf("Temperatura : %.1f °C\n", temp)
        fmt.Printf("Umidade     : %.1f %%\n", humidity)
        fmt.Printf("Pressão     : %.1f hPa\n", pressure)
        fmt.Println("=====================================")
    }
}
