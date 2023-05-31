package main

import "os"

func createEmptyFile(fileName string) {
    if _, fileStatErr := os.Stat(fileName); fileStatErr != nil {
        file, fileCreateErr := os.Create(fileName)
        if fileCreateErr != nil {
            panic("Failed to create empty file '" + fileName + "': " + fileCreateErr.Error())
        }
        file.Close()
    }
}

func main() {
    createEmptyFile("auth_secret")
    createEmptyFile("tunnel_config.json")
}
