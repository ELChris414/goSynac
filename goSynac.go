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
		return strings.EqualFold(digestStr, securityString)
	})

	conn, err := openssl.Dial("tcp", ip, ctx, openssl.InsecureSkipHostVerification)
	if err != nil {
		return
	}

	session.stream = conn
	session.Handlers.status = 0
	session.Users = make(map[uintptr]User)
	session.Channels = make(map[uintptr]Channel)
	return
}

func (session Session) listen() (t int, output interface{}, err error) {
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
	fmt.Println("Received:", inputBytes)
	t, output, err = processMsgpack(inputBytes)
	if err != nil {
		err = errors.New("Unmarshaling content failed, " + err.Error())
		return
	}
	return
}

func (session *Session) liveRunner() {
	if session.Handlers.status > 0 {
		t, data, err := session.listen()
		if err != nil {
			panic(err)
		}
		switch findPacket(t) {
		case "error":
			err = errors.New(findError(data.(int)))
			panic(err)
		case "userReceive":
			data := data.(UserReceive)
			session.Users[data.Inner.ID] = data.Inner
			session.runHandler("UR", data)
		}
		fmt.Println(t, data)
	}
}

// runHandler runs the appropriate handlers for an event
func (session *Session) runHandler(t string, handler interface{}) {
	switch t {
	case "UR":
		for _, i := range session.Handlers.UR {
			i(session, handler.(UserReceive))
		}
	}

}

// AddHandler adds a handler
func (session *Session) AddHandler(handler interface{}) error {
	var err error
	switch handler.(type) {
	case func(*Session, MessageReceive):
		handler := handler.(func(*Session, MessageReceive))
		session.Handlers.MR = append(session.Handlers.MR, handler)
	case func(*Session, MessageDeleteReceive):
		handler := handler.(func(*Session, MessageDeleteReceive))
		session.Handlers.MDR = append(session.Handlers.MDR, handler)
	case func(*Session, UserReceive):
		handler := handler.(func(*Session, UserReceive))
		session.Handlers.UR = append(session.Handlers.UR, handler)
	case func(*Session, ChannelReceive):
		handler := handler.(func(*Session, ChannelReceive))
		session.Handlers.CR = append(session.Handlers.CR, handler)
	default:
		panic("Invalid handler function")
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
	return Wrapper{typey, struct{ Content interface{} }{structy}}
}

// Login logs you in Synac
func (session *Session) Login(bot bool, name string, password string, token string) (tokenO string, created bool, err error) {
	lg := Login{
		Bot:      bot,
		Name:     name,
		Password: password,
		Token:    token,
	}
	packet := packIt(lg, findRPacket("login"))
	fmt.Println("Sent:", packet)
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
	t, logins, err := session.listen()
	if err != nil {
		return
	}
	fmt.Println(findPacket(t))
	switch findPacket(t) {
	case "loginSuccess":
		data := logins.(LoginSuccess)
		created = data.Created
		tokenO = data.Token
		session.ID = data.ID
		session.Handlers.status = 1
	case "error":
		err = errors.New(findError(logins.(int)))
	default:
		err = errors.New("Something unkown happened.")
	}

	return
}
