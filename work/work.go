package work

import (
	"fmt"

	"github.com/topritchett/game-server/proxmox"
)

var workVMNames = []string{"instadocker", "if-prod-db", "es-01", "es-02", "es-03"}

func StartWorkVMs(auth, qemuUrl string) string {
	fmt.Println("Attempting to start IF VM")
	for _, vmName := range workVMNames {
		vmID := proxmox.GetVMID(auth, qemuUrl, vmName)
		if vmID != "0" {
			fmt.Println("Attempting to start VM")
			response, err := proxmox.StartVM(auth, qemuUrl, vmID)
			if err != nil {
				fmt.Println("Error", err)
			}
			fmt.Println(response)
		}
	}
	return "Started work VMs"
}

func PauseWorkVMs(auth, qemuUrl string) string {
	fmt.Println("Attempting to pause IF VMs")
	for _, vmName := range workVMNames {
		vmID := proxmox.GetVMID(auth, qemuUrl, vmName)
		if vmID != "0" {
			fmt.Println("Attempting to pause VM")
			response, err := proxmox.PauseVM(auth, qemuUrl, vmID)
			if err != nil {
				fmt.Println("Error", err)
			}
			fmt.Println(response)
		}
	}
	return "Paused work VMs"
}
