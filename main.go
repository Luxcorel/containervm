package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	Exit int = iota
	CreateStorage
	BootFromISO
	BootFromHarddrive
	VdiToQcow2
)

func main() {
	for {
		fmt.Println()
		fmt.Println("###########################################")
		fmt.Println("0: Exit")
		fmt.Println("1: Create storage")
		fmt.Println("2: Boot from ISO")
		fmt.Println("3: Boot from harddrive")
		fmt.Println("4: Convert VDI to QCOW2")
		fmt.Println("###########################################")

		fmt.Print("Enter your choice: ")

		var input int
		_, err := fmt.Scanln(&input)
		fmt.Println()
		if err != nil {
			log.Printf("could not read input: %v\n", err)
			continue
		}

		switch input {

		case Exit: // Exit
			os.Exit(0)

		case CreateStorage: // Create image
			imageSize, err := userInput("Enter image size(eg. 20G): ")
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			storageName, err := userInput("Enter image name(eg. storage): ")
			storageName += ".qcow2"
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			err = createStorage(imageSize, storageName)
			if err != nil {
				log.Printf("could not create storage: %v\n", err)
			}

		case BootFromISO: // Boot from ISO
			isoName, err := userInput("Enter ISO name(eg. ubuntu.iso): ")
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			storageName, err := userInput("Enter storage name(eg. storage): ")
			storageName += ".qcow2"
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			memorySize, err := userInput("Enter memory size(eg. 4G): ")
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			err = bootFromISO(isoName, storageName, memorySize)
			if err != nil {
				log.Printf("could not boot from ISO: %v\n", err)
			}

		case BootFromHarddrive: // Boot from harddrive
			memorySize, err := userInput("Enter memory size(eg. 4G): ")
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			storageName, err := userInput("Enter storage name(eg. storage): ")
			storageName += ".qcow2"
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			err = bootFromHarddrive(memorySize, storageName)
			if err != nil {
				log.Printf("could not boot from harddrive: %v\n", err)
			}

		case VdiToQcow2: // Convert VDI to QCOW2
			vdiName, err := userInput("Enter VDI name(eg. input): ")
			vdiName += ".vdi"
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			qcow2Name, err := userInput("Enter output name(eg. output): ")
			qcow2Name += ".qcow2"
			if err != nil {
				log.Printf("could not read input: %v\n", err)
				break
			}

			err = vdiToQcow2(vdiName, qcow2Name)
			if err != nil {
				log.Printf("could not convert vdi to qcow2 %v\n", err)
			}

		default:
			fmt.Println("Invalid option")
		}
	}
}

func Execute(cmd string, options []string, environment []string) error {
	command := exec.Command(cmd, options...)
	envs := command.Environ()
	envs = append(envs, environment...)
	command.Env = envs
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

func userInput(userText string) (string, error) {
	fmt.Println(userText)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func createStorage(imageSize, storageName string) error {
	cmd := "qemu-img"
	options := []string{"create", "-f", "qcow2", storageName, imageSize}

	err := Execute(cmd, options, nil)
	return err
}

func bootFromISO(isoName, storageName, memorySize string) error {
	cmd := "qemu-system-x86_64"
	options := []string{"-m", memorySize, "-hda", storageName, "-cdrom", isoName, "-boot", "d", "-net", "nic", "-net", "user"}
	err := Execute(cmd, options, nil)
	return err
}

func bootFromHarddrive(memorySize, storageName string) error {
	cmd := "qemu-system-x86_64"
	options := []string{"-m", memorySize, "-hda", storageName, "-net", "nic", "-net", "user"}
	err := Execute(cmd, options, nil)
	return err
}

func vdiToQcow2(vdiName, qcow2Name string) error {
	cmd := "qemu-img"
	options := []string{"convert", "-f", "vdi", "-O", "qcow2", vdiName, qcow2Name}
	err := Execute(cmd, options, nil)
	return err
}
