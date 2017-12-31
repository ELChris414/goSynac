 
file:///home/elchris414/go/src/github.com/elchris414/goSchmic/generals.go
package main

import "github.com/spacemonkeygo/openssl"

// Constants
const defaultPort = 8439
const limitUserName = 128
const limitChannelName = 128
const limitAttrName = 128
const limitAttrAmount = 2048
const limitMessage = 16384

const limitMessageList = 64

const errAttrInvalidPos = 1
const errAttrLockedName = 2
const errLimitReached = 4
const errLoginBanned = 5
const errLoginBot = 6
const errLoginInvalid = 7
const errMaxConnPerIP = 8
const errMissingField = 9
const errMissingPermission = 10
const errNameTaken = 11
const errUnknownChannel = 12
const errUnknownGroup = 13
const errUnknownMessage = 14
const errUnknownUser = 15

const close = 0
const err = 1
const rateLimit = 2
const statusLogin = 3
const attributeCreate = 4
const attributeUpdate = 5
const attributeDelete = 6
const channelCreate = 7
const channelUpdate = 8
const channelDelete = 9
const messageList = 10
const messageCreate = 11
const messageUpdate = 12
const messageDelete = 13
const command = 14
const userUpdate = 15
const loginUpdate = 16
const loginSuccess = 17
const userReceive = 18
const attributeReceive = 19
const attributeDeleteReceive = 20
const channelReceive = 21
const channelDeleteReceive = 22
const messageReceive = 23
const messageDeleteReceive = 24
const commandReceive = 25

// Attribute attributes in Schmic
type Attribute struct {
	Allow uint8
	Deny  uint8
	ID    uintptr
	Name  string
	Pos   uintptr
}

// User users in Schmic
type User struct {
	Attributes []uintptr
	Bot        bool
	ID         uintptr
	Name       string
	Nick       string
}

// Channel channels in Schmic
type Channel struct {
	allow []uintptr
	deny  []uintptr
	id    uintptr
	name  string
}

// Message messages in Schmic
type Message struct {
	author        uintptr
	channel       uintptr
	id            uintptr
	text          []uint8
	timestamp     int64
	timestampEdit int64
}

// Packets

// Login packet in Schmic
type Login struct {
	Bot      bool
	Name     string
	Password string
	Token    string
}

// AttributeCreate AttributeCreate
type AttributeCreate struct {
	allow uint8
	deny  uint8
	name  string
	pos   uintptr
}

// AttributeUpdate AttributeUpdate
type AttributeUpdate struct {
	inner Attribute
}

// AttributeDelete AttributeDelete
type AttributeDelete struct {
	id uintptr
}

// ChannelCreate ChannelCreate
type ChannelCreate struct {
	allow []uintptr
	deny  []uintptr
	name  string
}

// ChannelUpdate ChannelUpdate
type ChannelUpdate struct {
	inner Channel
}

// ChannelDelete ChannelDelete
type ChannelDelete struct {
	channel uintptr
}

// MessageList MessageList
type MessageList struct {
	after   uintptr
	before  uintptr
	channel uintptr
	limit   uintptr
}

// MessageCreate MessageCreate
type MessageCreate struct {
	channel uintptr
	text    []uint8
}

// MessageUpdate MessageUpdate
type MessageUpdate struct {
	id   uintptr
	text []uint8
}

// MessageDelete MessageDelete
type MessageDelete struct {
	id uintptr
}

// Command Command
type Command struct {
	author    uintptr
	parts     []string
	recipient uintptr
}

// Close Close
type Close struct{}

// MessageReceive is an event handler
type MessageReceive struct {
	*Message
	new bool
}

// MessageDeleteReceive is an event handler
type MessageDeleteReceive struct {
	id uintptr
}

// Handlers is a struct that stores handlers
type Handlers struct {
	status int
	MR     []func(SchmicSession, MessageReceive)
	MDR    []func(SchmicSession, MessageDeleteReceive)
}

// SchmicSession A session for the Schmic chat
type SchmicSession struct {
	Attributes map[uintptr]Attribute
	Channel    uintptr
	Channels   map[uintptr]Channel
	ID         uintptr
	Users      map[uintptr]User
	Stream     *openssl.Conn
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

func findError(err int8) (errS string) {
	switch err {
	case errAttrInvalidPos:
		errS = "ERR_ATTR_INVALID_POS"
	case errAttrLockedName:
		errS = "ERR_ATTR_LOCKED_NAME"
	case errLimitReached:
		errS = "ERR_LIMIT_REACHED"
	case errLoginBanned:
		errS = "ERR_LOGIN_BANNED"
	case errLoginBot:
		errS = "ERR_LOGIN_BOT"
	case errLoginInvalid:
		errS = "ERR_LOGIN_INVALID"
	case errMaxConnPerIP:
		errS = "ERR_MAX_CONN_PER_IP"
	case errMissingField:
		errS = "ERR_MISSING_FIELD"
	case errMissingPermission:
		errS = "ERR_MISSING_PERMISSION"
	case errUnknownGroup:
		errS = "ERR_UNKNOWN_GROUP"
	case errUnknownChannel:
		errS = "ERR_UNKNOWN_CHANNEL"
	case errUnknownMessage:
		errS = "ERR_UNKNOWN_MESSAGE"
	case errUnknownUser:
		errS = "ERR_UNKNOWN_USER"
	case errNameTaken:
		errS = "ERR_NAME_TAKEN"
	case errInvalidChannel:
		errS = "ERR_INVALID_CHANNEL"
	default:
		errS = "ERR_UKNOWN_ERR"
	}
	return
}

func findPacket(thing interface{}) (errS string) {
	switch thing.(int8) {
	case close:
		errS = "close"
	case err:
		errS = "error"
	case rateLimit:
		errS = "rateLimit"
	case attributeCreate:
		errS = "attributeCreate"
	case attributeDelete:
		errS = "attributeDelete"
	case attributeUpdate:
		errS = "attributeUpdate"
	case channelCreate:
		errS = "channelCreate"
	case channelUpdate:
		errS = "channelUpdate"
	case channelDelete:
		errS = "channelDelete"
	case messageList:
		errS = "messageList"
	case messageCreate:
		errS = "messageCreate"
	case messageUpdate:
		errS = "messageUpdate"
	case messageDelete:
		errS = "messageDelete"
	case command:
		errS = "command"
	case userUpdate:
		errS = "userUpdate"
	case loginUpdate:
		errS = "loginUpdate"
	case loginSuccess:
		errS = "loginSuccess"
	case userReceive:
		errS = "userReceive"
	case attributeReceive:
		errS = "attributeReceive"
	case attributeDeleteReceive:
		errS = "attributeDeleteReceive"
	case channelReceive:
		errS = "channelReceive"
	case channelDeleteReceive:
		errS = "channelDeleteReceive"
	case messageReceive:
		errS = "messageReceive"
	case messageDeleteReceive:
		errS = "messageDeleteReceive"
	case commandReceive:
		errS = "commandReceive"
	}
	return
}
