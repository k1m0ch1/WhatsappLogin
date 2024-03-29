package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/joho/godotenv"
	qrcodeTerminal "github.com/mdp/qrterminal/v3"
)

type waHandler struct {
	c         *whatsapp.Conn
	startTime uint64
}

type param struct {
	phoneNumber string
	dirOutput   string
}

var params param

func main() {

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sessionOutput := flag.String("o", mydir+"/sessions", "an output dir")
	pN := flag.String("p", "6288", "phone number")

	flag.Parse()

	params.dirOutput = *sessionOutput
	params.phoneNumber = *pN

	if params.dirOutput == mydir && params.phoneNumber == "6288" {
		fmt.Println("example usage " + os.Args[0] + " -o ~/sessions -p 62872123123")
		os.Exit(1)
	}

	wac, err := whatsapp.NewConn(1 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}
	var major, minor, patch int = 2, 2123, 7

	if clientVersion := os.Getenv("CLIENT_VERSION"); clientVersion != "" {
		cv := strings.Split(clientVersion, ".")
		// Major, minor, patch all there
		major, err = strconv.Atoi(cv[0])
		if err != nil {
			fmt.Println("Client version format wrong, major section is not int", err)
			major = 2
		}
		minor, err = strconv.Atoi(cv[1])
		if err != nil {
			fmt.Println("Client version format wrong, minor section is not int", err)
			minor = 2123
		}
		patch, err = strconv.Atoi(cv[2])
		if err != nil {
			fmt.Println("Client version format wrong, patch section is not int", err)
			patch = 7
		}
	}

	wac.SetClientVersion(major, minor, patch)
	wac.AddHandler(&waHandler{wac, uint64(time.Now().Unix())})

	if err := login(wac); err != nil {
		log.Fatalf("error logging in: %v\n", err)
	}

	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	session, err := wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
	if err := writeSession(session, params.phoneNumber); err != nil {
		log.Fatalf("error saving session: %v", err)
	}
}

func login(wac *whatsapp.Conn) error {
	qr := make(chan string)
	go func() {
		config := qrcodeTerminal.Config{
			Level:     qrcodeTerminal.L,
			Writer:    os.Stdout,
			BlackChar: qrcodeTerminal.BLACK,
			WhiteChar: qrcodeTerminal.WHITE,
			QuietZone: 1,
		}
		qrcodeTerminal.GenerateWithConfig(<-qr, config)
	}()
	session, err := wac.Login(qr)
	if err != nil {
		return fmt.Errorf("error during login: %v\n", err)
	}
	err = writeSession(session, params.phoneNumber)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func writeSession(session whatsapp.Session, phoneNumber string) error {
	file, err := os.Create(getSessionName(phoneNumber))
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func getSessionName(phoneNumber string) string {
	if _, err := os.Stat(params.dirOutput); os.IsNotExist(err) {
		os.MkdirAll(params.dirOutput, os.ModePerm)
	}
	return params.dirOutput + "/" + phoneNumber + ".gob"
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *waHandler) HandleError(err error) {
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}
