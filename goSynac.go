package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/spacemonkeygo/openssl"
	"github.com/vmihailenco/msgpack"
)

// CreateSession Creates an goSynac.Session
func CreateSession(ip string, securityString string) (session Session, err error) {
	initialize()
	IP := net.ParseIP(ip)
	if IP != nil {
		err = errors.New("Invalid IP")
		return
	}

	_, port, _ := net.SplitHostPort(ip)
	if port == "" {
		ip += ":8439"
	}

	ctx, err := openssl.NewCtx()
	if err != nil {
		return
	}

	ctx.SetVerify(openssl.VerifyPeer, func(ok bool, store *openssl.CertificateStoreCtx) bool {
		pkey, err := store.GetCurrentCert().PublicKey()
		if err != nil {
			return false
		}
		pem, err := pkey.MarshalPKIXPublicKeyPEM()
		if err != nil {
			return false
		}
		result, err := openssl.SHA256(pem)
		if err != nil {
			return false
		}
		digestStr := hex.EncodeToString(result[:])
		if strings.EqualFold(strings.ToUpper(digestStr), securityString) {
			return false
		}
		return true
	})

	conn, err := openssl.Dial("tcp", ip, ctx, openssl.InsecureSkipHostVerification)
	if err != nil {
		return
	}

	session.stream = conn
	session.Handlers.status = 0
	return
}

func (session Session) listen() (output []interface{}, err error) {
	twoBytes := make([]byte, 2)
	_, err = session.stream.Read(twoBytes)
	if err != nil {
		err = errors.New("Reading from stream failed, " + err.Error())
		return
	}
	size := binary.BigEndian.Uint16(twoBytes)
	inputBytes := make([]byte, size)
	_, err = session.stream.Read(inputBytes)
	if err != nil {
		err = errors.New("Reading from stream failed, " + err.Error())
		return
	}
	err = msgpack.Unmarshal(inputBytes, &output)
	if err != nil {
		err = errors.New("Unmarshaling content failed, " + err.Error())
		return
	}
	return
}

func (session *Session) liveRunner() {
	if session.Handlers.status > 0 {
		data, err := session.listen()
		if err != nil {
			panic(err)
		}
		switch findPacket(data[0]) {
		case "error":
			err = errors.New(findError(data[1].([]interface{})[0].(int8)))
			panic(err)
		case "userReceive":
			deeper := data[1].([]interface{})[0].([]interface{})
			admin, _ := deeper[0].(bool)
			ban, _ := deeper[1].(bool)
			bot, _ := deeper[2].(bool)
			id, _ := deeper[3].(uintptr)
			nodes, _ := deeper[4].(map[uintptr]uint8)
			name, _ := deeper[5].(string)
			session.Users[id] = User{
				admin: admin,
				ban:   ban,
				bot:   bot,
				id:    id,
				nodes: nodes,
				name:  name,
			}
		}
		fmt.Println(data)
	}
}

// Close closes an Session
func (session Session) Close() {
	session.stream.Close()
}

// AddHandler adds a handler
func (session *Session) AddHandler(handler interface{}) error {
	var err error
	switch handler.(type) {
	case func(Session, MessageReceive):
		handler := handler.(func(Session, MessageReceive))
		session.Handlers.MR = append(session.Handlers.MR, handler)
	case func(Session, MessageDeleteReceive):
		handler := handler.(func(Session, MessageDeleteReceive))
		session.Handlers.MDR = append(session.Handlers.MDR, handler)
	case func(Session, UserReceive):
		handler := handler.(func(Session, UserReceive))
		session.Handlers.UR = append(session.Handlers.UR, handler)
	case func(Session, ChannelReceive):
		handler := handler.(func(Session, ChannelReceive))
		session.Handlers.CR = append(session.Handlers.CR, handler)
	}
	return err
}

// Write writes bytes in an Session
func (session Session) Write(input []byte) (written int, err error) {
	size := len(input)
	if size > math.MaxUint16 {
		err = errors.New("packet too large")
		return
	}
	writable := make([]byte, 2)
	binary.BigEndian.PutUint16(writable[0:], uint16(size))
	written, err = session.stream.Write(writable)
	if err != nil {
		err = errors.New("writing error" + err.Error())
		return
	}
	written, err = session.stream.Write(input)
	if err != nil {
		err = errors.New("writing error" + err.Error())
		return
	}
	return
}

func packIt(structy interface{}, typey int) Wrapper {
	return Wrapper{typey, Wrapping{structy}}
}

// Login logs you in Synac
func (session *Session) Login(bot bool, name string, password string, token string) (tokenO string, created bool, err error) {
	lg := Login{
		bot:      bot,
		name:     name,
		password: password,
		token:    token,
	}
	packet := packIt(lg, findRPacket("login"))
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).StructAsArray(true)
	err = enc.Encode(&packet)
	if err != nil {
		return
	}
	_, err = session.Write(buf.Bytes())
	if err != nil {
		return
	}
	logins, err := session.listen()
	if err != nil {
		return
	}
	switch findPacket(logins[0]) {
	case "loginSuccess":
		deeper := logins[1].([]interface{})[0].([]interface{})
		created = deeper[0].(bool)
		id, _ := deeper[1].(uintptr)
		tokenO = deeper[2].(string)
		session.ID = id
		session.Handlers.status = 1
	case "error":
		err = errors.New(findError(logins[1].([]interface{})[0].(int8)))
	default:
		err = errors.New("Something unkown happened.")
	}

	return
}
