package errorhandler

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"net"
	"os"
	"syscall"
)

func HandleError(err error) {
	if err != nil {
		if a, ok := err.(net.Error);ok {
			fmt.Println(a.Error())
			fmt.Println(a.Timeout())
			fmt.Println(a.Temporary())

			if c, d :=a.(*net.OpError);d {
				switch t := c.Err.(type) {
				case *net.DNSError:
					log.Println("net.DNSError:", t)
				case *net.InvalidAddrError:
					log.Println("net.InvalidAddrError:", t)
				case *net.UnknownNetworkError:
					log.Println("net.UnknownNetworkError:", t)
				case *net.AddrError:
					log.Println("net.AddrError:", t)
				case *net.DNSConfigError:
					log.Println("net.DNSConfigError:", t)
				case *os.SyscallError:
					log.Printf("os.SyscallError:%+v", t)
					if errno, ok := t.Err.(syscall.Errno); ok {
						log.Printf("errno:%d\n", errno)
						switch errno {
						case syscall.ECONNREFUSED:
							log.Println("connect refused")
						case syscall.ETIMEDOUT:
							log.Println("timeout")
						case syscall.ECONNABORTED:
							log.Println("connect aborted")

						case syscall.EHOSTDOWN:
							log.Println("EHOSTDOWN")

						case windows.WSAECONNREFUSED:
							log.Println("WSAECONNREFUSED")

						case syscall.EHOSTUNREACH:
							log.Println("EHOSTUNREACH")

						}
					} else {
						log.Println("t.Err.(syscall.Errno)", "assertion fail")
					}
				}
			}
		}
	}
}