package main

import (
	"fmt"
)

func main() {
	S, err := CreateSession("localhost", "ECACDC7A85CB6B6C31F87535B97D95FFA8FFABAC751752CD01EC2193B6393AE")
	if err != nil {
		fmt.Println(err)
		return
	}

	S.AddHandler(userListener)
	S.AddHandler(channelListener)
	S.AddHandler(groupListener)
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

func userListener(session SchmicSession, received UserReceive) {
	fmt.Println("I heard about a user called", received.inner.Name)
}

func userListenerTwo(session SchmicSession, received UserReceive) {
	fmt.Println("Can confirm that!")
}

func channelListener(session SchmicSession, received ChannelReceive) {
	fmt.Println("I heard about a channel called", received.inner.Name)
}

func groupListener(session SchmicSession, received GroupReceive) {
	fmt.Println("I heard about a group called", received.inner.Name)
	fmt.Println("It happens to be", received.new)
}
