package main

import (
    "os"
    "fmt"
    "net"
    "log"
    "time"
    "os/user"
    "os/signal"
    "syscall"
	"strings"
	"bytes"
)

func main () {
    /*init*/
    //Get the user who is logging in
    started := time.Now()
    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    //Put a logfile in place for any errors or interesting events
    logfile:=usr.HomeDir+"/notsh.log"
    var f *os.File
    if _, er := os.Stat(logfile); os.IsNotExist(er) {
        f, err = os.Create(logfile)   
    } else {
        f, err = os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    }
    if err!=nil{
        panic(err)
    }
    //Stop logging from going to stdout/stderr
    log.SetOutput(f)
    log.Printf("Set output for logging. Initialised at %s",started)
    //Listen on a socket file for messages from admin
    masterSock := usr.HomeDir+"/"+ grid()+".session"
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        os.Remove(masterSock)
        os.Exit(1)
    }()
    /*Gather details*/
    //Collect any information that the admin might want about the session.
    //TODO: Detect ssh -L port forwarding. User should still be able to escape with <enter>,~,C
    /*Setup*/
    //make sure no existing sock exists
	err = os.Remove(masterSock)
	if err == nil {
		fmt.Println("Overwrote existing master sock")
	}
	session, err := net.Listen("unix", masterSock)
	if err != nil {

		log.Fatal("Write: ", err)
	}
    defer os.Remove(masterSock)
    println("You are now connected.")
    /*Wait for message to inspect or kill*/
    //Set up a channel to suspend until a signal is accepted
    var quit = make(chan int)
    //listen to any incoming connections on master sock
    go func(){
        for {
            conn, _ := session.Accept()
            buffer := make([]byte,1024)//max msg size
            conn.Write([]byte(
                fmt.Sprintf(`
This connection was started at %s
The user appears to be: %s (%s)
Enter KILL to close this connection, or MSG to send:
`,started,usr.Name,usr.Username)            ))
            
             _,err := conn.Read(buffer[:])
            for err==nil {
                execute(buffer,&conn,quit)
                buffer = make([]byte,1024)
                 _,err = conn.Read(buffer[:])
            }
        }
    }()
    <-quit
}

func grid() string {
	//out, err := exec.Command("uuidgen").Output()
    return fmt.Sprintf("%d",time.Now().UnixNano())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// return string(out[:len(out)-1])
}


func execute (buffer []byte, conn *net.Conn, quit chan int) {
    tbuff:=bytes.Trim(buffer, "\x00")
    cmd:=strings.Replace(string(tbuff),"\n","",-1)
    switch (cmd){
        case "MSG", "Msg", "msg":
            (*conn).Write([]byte(fmt.Sprintf("Enter message to display:\n")))
            _,err := (*conn).Read(buffer[:])
            if err != nil {
                log.Fatal("Write: ", err)
            }
            msg:=string(bytes.Trim(buffer, "\x00"))
            print(msg,"\n")
            (*conn).Write([]byte(fmt.Sprintf("Displayed: %s\n",msg)))
        case "KILL", "Kill", "kill":
            quit<-1
            (*conn).Write([]byte(fmt.Sprintf("Connection killed.\n")))
        case "":
        default:
            (*conn).Write([]byte(fmt.Sprintf("Unknown instruction: %s",cmd)))
    }
}