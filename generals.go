package main

import "github.com/spacemonkeygo/openssl"

// Constants
const defaultPort = 8439
const typingTimeout = 10
const limitChannelName = 128
const limitUserName = 128
const limitMessage = 16384

const limitBulk = 64

var synacErrors = make(map[int8]string)
var packets = make(map[int8]string)
var rpackets = make(map[string]int)

// TODO PERMISSIONS

// Channel stores a channel
type Channel struct {
	DefaultModeBot  uint8
	DefaultModeUser uint8
	ID              uintptr
	Name            string
}

// Message stores a message
type Message struct {
	Author        uintptr
	Channel       uintptr
	I             uintptr
	Text          []uint8
	Timestamp     int64
	TimestampEdit int64
}

// User stores a user
type User struct {
	Admin bool
	Ban   bool
	Bot   bool
	ID    uintptr
	Nodes map[uintptr]uint8
	Name  string
}

type ChannelCreate struct {
	DefaultModeBot  uint8
	DefaultModeUser uint8
	Name            string
}

type ChannelDelete struct {
	ID uintptr
}

type ChannelUpdate struct {
	Inner Channel
}

type Command struct {
	Args      []string
	Recipient uintptr
}

type Login struct {
	Bot      bool
	Name     string
	Password string
	Token    string
}

type LoginUpdate struct {
	Name             string
	Password_current string
	Password_new     string
	Reset_token      bool
}

type MessageCreate struct {
	Channel uintptr
	Text    []uint8
}

type MessageDelete struct {
	ID uintptr
}

type MessageDeleteBulk struct {
	Channel uintptr
	IDs     []uintptr
}

type MessageList struct {
	After   uintptr
	Before  uintptr
	Channel uintptr
	Limit   uintptr
}

type MessageUpdate struct {
	ID   uintptr
	Text []uint8
}

type PrivateMessage struct {
	Text      []uint8
	Recipient uintptr
}

type Typing struct {
	Channel uintptr
}

type UserUpdate struct {
	Admin       bool
	Ban         bool
	ChannelMode map[uintptr]uint8 // may be wrong
	ID          uintptr
}

type ChannelDeleteReceive struct {
	Inner Channel
}

type ChannelReceive struct {
	Inner Channel
}

type CommmandReceive struct {
	Args   []string
	Author uintptr
}

type LoginSuccess struct {
	Created bool
	ID      uintptr
	Token   string
}

type MessageDeleteReceive struct {
	ID uintptr
}

type MessageReceive struct {
	Inner Message
	New   bool
}

type PMreceive struct {
	Author uintptr
	Text   []uint8
}

type TypingReceive struct {
	Author  uintptr
	Channel uintptr
}

type UserReceive struct {
	Inner User
}

// Handlers stores handlers
type Handlers struct {
	status int
	MR     []func(*Session, MessageReceive)
	MDR    []func(*Session, MessageDeleteReceive)
	UR     []func(*Session, UserReceive)
	CR     []func(*Session, ChannelReceive)
}

// Session A session for the Synac chat
type Session struct {
	Channel  uintptr
	Channels map[uintptr]Channel
	ID       uintptr
	Users    map[uintptr]User
	stream   *openssl.Conn
	Handlers
}

// Wrapping because Go really sucks
type Wrapping struct {
	Content interface{}
}

// Wrapper because Go sucks
type Wrapper struct {
	Type    int
	Content Wrapping
}

func initialize() {
	// ERRORS
	synacErrors[1] = "ERR_LIMIT_REACHED"
	synacErrors[2] = "ERR_LOGIN_BANNED"
	synacErrors[3] = "ERR_LOGIN_BOT"
	synacErrors[4] = "ERR_LOGIN_INVALID"
	synacErrors[5] = "ERR_MAX_CONN_PER_IP"
	synacErrors[6] = "ERR_MISSING_FIELD"
	synacErrors[7] = "ERR_MISSING_PERMISSION"
	synacErrors[8] = "ERR_NAME_TAKEN"
	synacErrors[9] = "ERR_UNKNOWN_BOT"
	synacErrors[10] = "ERR_UNKNOWN_CHANNEL"
	synacErrors[11] = "ERR_UNKNOWN_MESSAGE"
	synacErrors[12] = "ERR_UNKNOWN_USER"

	// PACKETS
	rpackets["close"] = 0
	rpackets["err"] = 1
	rpackets["rateLimit"] = 2
	rpackets["channelCreate"] = 3
	rpackets["channelDelete"] = 4
	rpackets["channelUpdate"] = 5
	rpackets["command"] = 6
	rpackets["login"] = 7
	rpackets["loginUpdate"] = 8
	rpackets["messageCreate"] = 9
	rpackets["messageDelete"] = 10
	rpackets["messageDeleteBulk"] = 11
	rpackets["messageList"] = 12
	rpackets["messageUpdate"] = 13
	rpackets["privateMessage"] = 14
	rpackets["typing"] = 15
	rpackets["userUpdate"] = 16

	rpackets["channelDeleteReceive"] = 17
	rpackets["channelReceive"] = 18
	rpackets["commandReceive"] = 19
	rpackets["loginSuccess"] = 20
	rpackets["messageDeleteReceive"] = 21
	rpackets["messageListReceived"] = 22
	rpackets["pmReceive"] = 23
	rpackets["userReceive"] = 24

	for k, v := range rpackets {
		packets[int8(v)] = k
	}
}

func findError(err int8) string {
	return synacErrors[err]
}

func findPacket(thing interface{}) string {
	return packets[thing.(int8)]
}

func findRPacket(packet string) int {
	return rpackets[packet]
}
