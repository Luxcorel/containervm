package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	CreateImage int = iota
	BootFromISO
	BootFromHarddrive
	VdiToQcow2
	Exit
)

func main() {
	//TODO: Make these customizable
	createImage := "qemu-img"
	createImageOptions := []string{"create", "-f", "qcow2", "ubuntu_vm.qcow2", "20G"}

	createVm := "qemu-system-x86_64"
	createVmOptions := []string{"-m", "4G", "-hda", "ubuntu_vm.qcow2", "-cdrom", "test.iso", "-boot", "d", "-net", "nic", "-net", "user"}

	startVm := "qemu-system-x86_64"
	startVmOptions := []string{"-m", "4G", "-hda", "ubuntu_vm.qcow2", "-net", "nic", "-net", "user"}

	vdiToQcow2 := "qemu-img"
	vdiToQcow2Options := []string{"convert", "-f", "vdi", "-O", "qcow2", "ubuntu_vm.vdi", "ubuntu_vm.qcow2"}

	for {
		fmt.Println("What do you want to do?")
		fmt.Println("0: Create image")
		fmt.Println("1: Boot from ISO")
		fmt.Println("2: Boot from harddrive")
		fmt.Println("3: Convert VDI to QCOW2")
		fmt.Println("4: Exit")
		var input int
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatalf("could not read input: %v", input)
		}

		// TODO: Move exit to case 0
		switch input {
		case CreateImage: // Create image
			err := Execute(createImage, createImageOptions, nil)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
		case BootFromISO: // Boot from ISO
			err := Execute(createVm, createVmOptions, nil)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
		case BootFromHarddrive: // Boot from harddrive
			err := Execute(startVm, startVmOptions, nil)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
		case VdiToQcow2: // Convert VDI to QCOW2
			err := Execute(vdiToQcow2, vdiToQcow2Options, nil)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
		case Exit: // Exit
			os.Exit(0)
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
