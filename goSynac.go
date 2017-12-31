file:///home/elchris414/go/src/github.com/elchris414/goSchmic/goSchmic.go
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

// CreateSession Creates an SchmicSession
func CreateSession(ip string, securityString string) (session SchmicSession, err error) {
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

	session.Stream = conn
	session.Handlers.status = 0
	return
}

func (session SchmicSession) listen() (output []interface{}, err error) {
	twoBytes := make([]byte, 2)
	_, err = session.Stream.Read(twoBytes)
	if err != nil {
		err = errors.New("Reading from stream failed, " + err.Error())
		return
	}
	size := binary.BigEndian.Uint16(twoBytes)
	inputBytes := make([]byte, size)
	_, err = session.Stream.Read(inputBytes)
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

func (session *SchmicSession) liveRunner() {
	if session.Handlers.status > 0 {
		data, err := session.listen()
		if err != nil {
			panic(err)
		}
		switch findPacket(data[0]) {
		case "error":
			err = errors.New(findError(data[1].([]interface{})[0].(int8)))
			panic(err)
		case "attributeReceive":
			deeper := data[1].([]interface{})[0].([]interface{})
			allow, _ := deeper[0].(uint8)
			deny, _ := deeper[1].(uint8)
			id, _ := deeper[2].(uintptr)
			name := deeper[3].(string)
			pos, _ := deeper[4].(uintptr)
			session.Attributes[id] = Attribute{
				Allow: allow,
				Deny:  deny,
				ID:    id,
				Name:  name,
				Pos:   pos,
			}
		case "userReceive":
			deeper := data[1].([]interface{})[0].([]interface{})
			evenDeeper := deeper
			attributes, _ := deeper[0].([]uintptr)
			bot := deeper[1].(bool)
			id, _ := deeper[2].(uintptr)
			name := deeper[3].(string)
			nick, _ := deeper[4].(string)
			session.Users[id] = User{
				Attributes: attributes,
				Bot:        bot,
				ID:         id,
				Name:       name,
				Nick:       nick,
			}
		}
		fmt.Println(data)
	}
}

// Close closes an SchmicSession
func (session SchmicSession) Close() {
	session.Stream.Close()
}

// AddHandler adds a handler
func (session *SchmicSession) AddHandler(handler interface{}) error {
	var err error
	switch handler.(type) {
	case func(SchmicSession, MessageReceive):
		handler := handler.(func(SchmicSession, MessageReceive))
		session.Handlers.MR = append(session.Handlers.MR, handler)
	case func(SchmicSession, MessageDeleteReceive):
		handler := handler.(func(SchmicSession, MessageDeleteReceive))
		session.Handlers.MDR = append(session.Handlers.MDR, handler)
	}
	return err
}

// Write writes bytes in an SchmicSession
func (session SchmicSession) Write(input []byte) (written int, err error) {
	size := len(input)
	if size > math.MaxUint16 {
		err = errors.New("packet too large")
		return
	}
	writable := make([]byte, 2)
	binary.BigEndian.PutUint16(writable[0:], uint16(size))
	written, err = session.Stream.Write(writable)
	if err != nil {
		err = errors.New("writing error" + err.Error())
		return
	}
	written, err = session.Stream.Write(input)
	if err != nil {
		err = errors.New("writing error" + err.Error())
		return
	}
	return
}

func packIt(structy interface{}, typey int) Wrapper {
	return Wrapper{typey, Wrapping{structy}}
}

// Login logs you in Schmic
func (session *SchmicSession) Login(bot bool, name string, password string, token string) (tokenO string, created bool, err error) {
	lg := Login{
		Bot:      bot,
		Name:     name,
		Password: password,
		Token:    token,
	}
	packet := packIt(lg, login)
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
 
