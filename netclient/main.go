package main

import (
	"fmt"
        "github.com/gravitl/netmaker/netclient/functions"
        "golang.zx2c4.com/wireguard/wgctrl"
        nodepb "github.com/gravitl/netmaker/grpc"
	"flag"
	"os"
        "os/exec"
        "strconv"
	"log"
)

const (
	// name of the service
	name        = "wcdaemon"
	description = "Wirecat Daemon Service"
)

var password string
var group string
var server string
var accesskey string

var (
        wgclient *wgctrl.Client
)

var (
        wcclient nodepb.NodeServiceClient
)

func main() {
	tpassword := flag.String("p", "changeme", "This node's password for accessing the server regularly")
	taccesskey := flag.String("k", "badkey", "an access key generated by the server and used for one-time access (install only)")
	tserver := flag.String("s", "localhost:50051", "The location (including port) of the remote gRPC server.")
	tgroup := flag.String("g", "badgroup", "The node group you are attempting to join.")
	tnoauto := flag.Bool("na", false, "No auto mode. If true, netmaker will not be installed as a system service and you will have to retrieve updates manually via checkin command.")
	command := flag.String("c", "required", "The command to run")


        flag.Parse()



         getID := exec.Command("id", "-u")
         out, err := getID.Output()

         if err != nil {
                 log.Fatal(err)
         }
         id, err := strconv.Atoi(string(out[:len(out)-1]))

         if err != nil {
                 log.Fatal(err)
         }

         if id != 0 {
                 log.Fatal("This program must be run with elevated privileges (sudo). This program installs a SystemD service and configures WireGuard and networking rules. Please re-run with sudo/root.")
         }


        switch *command {
		case "required":
                        fmt.Println("command flag 'c' is required. Pick one of |install|checkin|update|remove|")
                        os.Exit(1)
			log.Fatal("Exiting")
                case "install":
			fmt.Println("Beginning agent installation.")
			err := functions.Install(*taccesskey, *tpassword, *tserver, *tgroup, *tnoauto)
			if err != nil {
				fmt.Println("Error installing: ", err)
				fmt.Println("Cleaning up (uninstall)")
				err = functions.Remove()
				if err != nil {
                                        fmt.Println("Error uninstalling: ", err)
                                        fmt.Println("Wiping local.")
					err = functions.WipeLocal()
					if err != nil {
						fmt.Println("Error removing artifacts: ", err)
					}
                                        err = functions.RemoveSystemDServices()
                                        if err != nil {
                                                fmt.Println("Error removing services: ", err)
                                        }
				}
				os.Exit(1)
			}
		case "service-install":
                        fmt.Println("Beginning service installation.")
                        err := functions.ConfigureSystemD()
                        if err != nil {
                                fmt.Println("Error installing service: ", err)
                                os.Exit(1)
                        }
                case "service-uninstall":
                        fmt.Println("Beginning service uninstall.")
                        err := functions.RemoveSystemDServices()
                        if err != nil {
                                fmt.Println("Error installing service: ", err)
                                os.Exit(1)
                        }
		case "checkin":
			fmt.Println("Beginning node check in.")
			err := functions.CheckIn()
			if err != nil {
				fmt.Println("Error checking in: ", err)
				os.Exit(1)
			}
		case "remove":
                        fmt.Println("Beginning node cleanup.")
			err := functions.Remove()
                        if err != nil {
					/*
                                        fmt.Println("Error uninstalling: ", err)
                                        fmt.Println("Wiping local.")
                                        err = functions.WipeLocal()
                                        if err != nil {
                                                fmt.Println("Error removing artifacts: ", err)
                                        }
                                        err = functions.RemoveSystemDServices()
                                        if err != nil {
                                                fmt.Println("Error removing services: ", err)
                                        }
					*/
                                fmt.Println("Error deleting node: ", err)
                                os.Exit(1)
                        }
	}
	fmt.Println("Command " + *command + " Executed Successfully")
}
