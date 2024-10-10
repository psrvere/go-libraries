package syscall

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// AF = Address Family
// AF_INET - IPV4 - used for communicating over IPV4
// AF_INET6 - IPV6 - used for communicating over IPV6
// AF_UNIX OR AF_LOCAL - used for communication between processes on the same machine
// AF_PACKET - raw access to network interfaces
// AF_BLUETOOTH - for bluetooth communication

// These are socket types
// SOCK_STREAM - TCP socket
// TCP - Transmission Control Protocol
// SOCK_DGRAM - UDP socket

// protocol
// 0 - default protocol for the selected domain and typ
// AF_INET + SOCK_STREAM = TCP
// AF_INET + SOCK_DGRAM = UDP
func CreateNewSocket() int {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Printf("Error creating socket: %v", err)
	}
	return fd
}

func CloseSocket(fd int) {
	if err := syscall.Close(fd); err != nil {
		log.Printf("Error closing socket. fd: %v err:%v", fd, err)
	}
}

// Levels
// SOL_SOCKET - option applies to socket not to TCP or IP
// SOL = socket level

// Options
// SO_REUSEADDR allows socket to forcibly bind to a port that is in
// TIME_WAIT state. This is useful when we want to restart a server
// and bind it to a port that was recently used. That port may still
// be processing packets from the previous connections on that port.
func SetSocketOption(fd int) {
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		log.Printf("Error setting socket options. fd: %v err:%v", fd, err)
	}
}

// When a file descriptor is set to non blocking mode, system calls (read, write, accept, etc.)
// made on that socket will return immediately without waiting for the operation to complete
// Why use non blocking mode?
// It is used in event driven or asynchronous programs
// Example - a server handling multiple clients, use non blocking sockets with an event loop
// This way server is free to pick up other tasks
func SetNonBlock(fd int) {
	if err := syscall.SetNonblock(fd, true); err != nil {
		log.Printf("Error setting non block. fd: %v err: %v", fd, err)
	}
}

func BindSocket(fd int, ip string, port int) {
	// net.ParseIP return net.IP which can be v4 or v6
	// It checks if IP is a valid IP
	ip4 := net.ParseIP(ip)
	if ip4 == nil {
		log.Printf("Invalid ip: %v", ip)
	}

	// Bind socket to the port
	if err := syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		log.Printf("Error binding socket. fd: %v err: %v", fd, err)
	}
}

// The second arg in Listen is the backlog which defines the maximum number
// of pending connections that the system will queue for the socket before
// it starts rejecting new ones
func ListenOnSocket(fd int) {
	if err := syscall.Listen(fd, 5); err != nil {
		log.Printf("Error listening on socket. fd: %v err: %v", fd, err)
	}
}

// Signals
// SIGTERM - Termination Signal - sent to gracefully shutdown a program
// SIGINT - Interrupt Signal - Cmd + C
// We can add as many signala in Notify function
// If no sinals are added, then all signals will be sent on the channel
func TerminateProgram() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// fire a go routine to handle shutdown
	go func(sigs chan os.Signal) {
		<-sigs
		// do graceful shutdown
	}(sigs)
}

// various syscall errors
func CheckError(err error) {
	// connection has been reset locally
	if errors.Is(err, syscall.EBADF) || errors.Is(err, os.ErrClosed) {
		log.Println("connection closed locally")
	}

	// connection has been reset by remote peer
	if errors.Is(err, syscall.ECONNRESET) {
		log.Println("connection reset")
	}

	// This may happen while during a non blocking I/O operation
	// Like reading from a non blocking socket where there is not data available to read
	// typically wait and retry after some time
	if errors.Is(err, syscall.EAGAIN) {
		log.Printf("Resource temporary available")
	}

	// It happens when non blocking operation can not proceed wihtout blocking
	// Example: buffer is full
	if errors.Is(err, syscall.EWOULDBLOCK) {
		log.Printf("Operation would block")
	}

	// When we attempt to write to a socket or pipe that has been closed
	// by the other end
	if errors.Is(err, syscall.EPIPE) {
		log.Printf("Broken pipe")
	}

	// network operation exceeds the defined time limit for completion
	if errors.Is(err, syscall.ETIMEDOUT) {
		log.Printf("Operation timed out")
	}
}

func ReadFromSocket(fd int) []byte {
	// The size of buffer defines maximum amount of operation to read in one operation
	// If the data can be longer than say 1024 bytes than program should run in a infinite for loop
	// reading data into a temporary buffer and appending to a larger buffer
	buf := make([]byte, 1024)

	n, err := syscall.Read(fd, buf)
	if err != nil {
		log.Printf("Error reading data. fd: %v err: %v", fd, err)
	}
	log.Printf("Read a total of %v bytes", n)
	return buf
}

func WriteToSocket(fd int, buf []byte) {
	n, err := syscall.Write(fd, buf)
	if err != nil {
		log.Printf("Error writing data. fd: %v err: %v", fd, err)
	}
	log.Printf("A total of %v bytes were written", n)
}

func GetProcessIDRaw() {
	// Syscall function is used to invoke a call directly for which there might not be a wrapper
	// Args:
	// trap uintptr - system call number
	// a1, a2, a3 uintptr - arguments to the system call. They differ based on system call
	// Returns:
	// r1, r2 - uintptr - return values from system call
	// err - 0 if the call was successful
	// syscall.Syscall is same as syscall.RawSyscall
	pid, _, err := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	if err != 0 {
		log.Printf("Error getting process id: %v\n", err)
	}
	log.Printf("Process id is: %v\n", pid)
}

func GetProcessId() int {
	return syscall.Getpid()
}

func GerParentProcessID() int {
	return syscall.Getppid()
}

func ForkChildProcess() uintptr {
	// SYS_FORK will fork the existing process i.e. it will create a new process by duplicating
	// the calling process
	child, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		log.Printf("Error forking: %v\n", err)
	}
	return child
}

func Exit(n int) {
	syscall.Exit(n)
}
