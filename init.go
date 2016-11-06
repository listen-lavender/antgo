package antgo

import "sync"

var DefaultConfig *Config

var ExitChan chan struct{}    // notify all goroutines to shutdown
var WaitGroup *sync.WaitGroup //

func init() {
	DefaultConfig = &Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20}
	ExitChan = make(chan struct{})
	WaitGroup = &sync.WaitGroup{}
}
