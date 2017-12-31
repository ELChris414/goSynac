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
	defaultModeBot  uint8
	defaultModeUser uint8
	id              uintptr
	name            string
}

// Message stores a message
type Message struct {
	author        uintptr
	channel       uintptr
	id            uintptr
	text          []uint8
	timestamp     int64
	timestampEdit int64
}

// User stores a user
type User struct {
	admin bool
	ban   bool
	bot   bool
	id    uintptr
	nodes map[uintptr]uint8
	name  string
}

type ChannelCreate struct {
	defaultModeBot  uint8
	defaultModeUser uint8
	name            string
}

type ChannelDelete struct {
	id uintptr
}

type ChannelUpdate struct {
	inner Channel
}

type Command struct {
	args      []string
	recipient uintptr
}

type Login struct {
	bot      bool
	name     string
	password string
	token    string
}

type LoginUpdate struct {
	name             string
	password_current string
	password_new     string
	reset_token      bool
}

type MessageCreate struct {
	channel uintptr
	text    []uint8
}

type MessageDelete struct {
	id uintptr
}

type MessageDeleteBulk struct {
	channel uintptr
	ids     []uintptr
}

type MessageList struct {
	after   uintptr
	before  uintptr
	channel uintptr
	limit   uintptr
}

type MessageUpdate struct {
	id   uintptr
	text []uint8
}

type PrivateMessage struct {
	text      []uint8
	recipient uintptr
}

type Typing struct {
	channel uintptr
}

type UserUpdate struct {
	admin       bool
	ban         bool
	channelMode map[uintptr]uint8 // may be wrong
	id          uintptr
}

type ChannelDeleteReceive struct {
	inner Channel
}

type ChannelReceive struct {
	inner Channel
}

type CommmandReceive struct {
	args   []string
	author uintptr
}

type LoginSuccess struct {
	created bool
	id      uintptr
	token   string
}

type MessageDeleteReceive struct {
	id uintptr
}

type MessageReceive struct {
	inner Message
	new   bool
}

type PMreceive struct {
	author uintptr
	text   []uint8
}

type TypingReceive struct {
	author  uintptr
	channel uintptr
}

type UserReceive struct {
	inner User
}

// Handlers stores handlers
type Handlers struct {
	status int
	MR     []func(Session, MessageReceive)
	MDR    []func(Session, MessageDeleteReceive)
	UR     []func(Session, UserReceive)
	CR     []func(Session, ChannelReceive)
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

// Wrapping because Go REALLY sucks
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
