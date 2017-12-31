package main

import (
	"fmt"
)

func main() {
	S, err := CreateSession("krake.one", "C9F6251FA50892B3877ACECA523ACFD925CAA7D9FA245D9C50DD00083A39F199")
	if err != nil {
		fmt.Println(err)
		return
	}

	S.AddHandler(userListener)
	S.AddHandler(channelListener)
	S.AddHandler(userListenerTwo)

	token, created, err := S.Login(false, "Chris", "lol", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token, "\nCreated:", created)
	for true {
		S.liveRunner()
	}
}

func userListener(session Session, received UserReceive) {
	fmt.Println("I heard about a user called", received.inner.name)
}

func userListenerTwo(session Session, received UserReceive) {
	fmt.Println("Can confirm that!")
}

func channelListener(session Session, received ChannelReceive) {
	fmt.Println("I heard about a channel called", received.inner.name)
}
