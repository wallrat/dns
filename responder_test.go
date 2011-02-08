package responder

import (
        "os"
	"testing"
        "fmt"
	"dns"
	"net"
	"time"
)

type myserv Server

func createpkg(id uint16, tcp bool, remove net.Addr) []byte {
	m := new(dns.Msg)
	m.MsgHdr.Id = id
	m.MsgHdr.Authoritative = true
	m.MsgHdr.AuthenticatedData = false
	m.MsgHdr.RecursionAvailable = true
	m.MsgHdr.Response = true
	m.MsgHdr.Opcode = dns.OpcodeQuery
	m.MsgHdr.Rcode = dns.RcodeSuccess
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{"miek.nl.", dns.TypeTXT, dns.ClassINET}
	m.Answer = make([]dns.RR, 1)
	t := new(dns.RR_TXT)
	t.Hdr = dns.RR_Header{Name: "miek.nl.", Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}
	if tcp {
		t.Txt = "Dit is iets anders TCP"
	} else {
		t.Txt = "Dit is iets anders UDP"
	}
	m.Answer[0] = t
	out, _ := m.Pack()
	return out
}

func (s *myserv) ResponderUDP(c *net.UDPConn, a net.Addr, in []byte) {
	inmsg := new(dns.Msg)
	inmsg.Unpack(in)
	if inmsg.MsgHdr.Response == true {
		// Uh... answering to an response??
		// dont think so
		return
	}
	out := createpkg(inmsg.MsgHdr.Id, false, a)
	SendUDP(out, c, a)
	// Meta.QLen/RLen/QueryStart/QueryEnd can be filled in at
	// this point for logging purposses or anything else
}

func (s *myserv) ResponderTCP(c *net.TCPConn, in []byte) {
	inmsg := new(dns.Msg)
	inmsg.Unpack(in)
	if inmsg.MsgHdr.Response == true {
		// Uh... answering to an response??
		// dont think so
		return
	}
	out := createpkg(inmsg.MsgHdr.Id, true, c.RemoteAddr())
	SendTCP(out, c)
}

func TestResponder(t *testing.T) {
	/* udp servertje */
	su := new(Server)
	su.Address = "127.0.0.1"
	su.Port = "8053"
	var us *myserv
	uch := make(chan os.Error)
	go su.NewResponder(us, uch)

	/* tcp servertje */
	st := new(Server)
	st.Address = "127.0.0.1"
	st.Port = "8053"
	st.Tcp = true
	var ts *myserv
	tch := make(chan os.Error)
	go st.NewResponder(ts, tch)
	time.Sleep(1 * 1e9)
	uch <- nil
	tch <- nil
}

/*
func TestReflectorResponder(t *testing.T) {
	stop := make(chan os.Error)
	s := new(Server)
	s.Port = "8053"
	s.Address = "127.0.0.1"

	stoptcp := make(chan os.Error)
	stcp := new(Server)
	stcp.Port = "8053"
	stcp.Address = "127.0.0.1"
	stcp.Tcp = true

	go stcp.NewResponder(Reflector, stoptcp)
	go s.NewResponder(Reflector, stop)

	time.Sleep(1 * 1e9)
	stop <- nil
	stoptcp <- nil
}
*/

type servtsig Server

func createpkgtsig(id uint16, tcp bool, remove net.Addr) []byte {
	m := new(dns.Msg)
	m.MsgHdr.Id = id
	m.MsgHdr.Authoritative = true
	m.MsgHdr.AuthenticatedData = false
	m.MsgHdr.RecursionAvailable = true
	m.MsgHdr.Response = true
	m.MsgHdr.Opcode = dns.OpcodeQuery
	m.MsgHdr.Rcode = dns.RcodeSuccess
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{"miek.nl.", dns.TypeTXT, dns.ClassINET}
	m.Answer = make([]dns.RR, 1)
	t := new(dns.RR_TXT)
	t.Hdr = dns.RR_Header{Name: "miek.nl.", Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}
	if tcp {
		t.Txt = "Dit is iets anders TCP"
	} else {
		t.Txt = "Dit is iets anders UDP"
	}
	m.Answer[0] = t
	out, _ := m.Pack()
	return out
}

func (s *servtsig) ResponderUDP(c *net.UDPConn, a net.Addr, in []byte) {
	inmsg := new(dns.Msg)
	inmsg.Unpack(in)
        fmt.Printf("%v\n", inmsg)
	if inmsg.MsgHdr.Response == true {
		// Uh... answering to an response??
		// dont think so
		return
	}
        rr := inmsg.Extra[len(inmsg.Extra)-1]
        switch t := rr.(type) {
        case *dns.RR_TSIG:
                v := t.Verify(inmsg, "awwLOtRfpGE+rRKF2+DEiw==")
                println(v)
        }


	out := createpkgtsig(inmsg.MsgHdr.Id, false, a)
	SendUDP(out, c, a)
	// Meta.QLen/RLen/QueryStart/QueryEnd can be filled in at
	// this point for logging purposses or anything else
}

func (s *servtsig) ResponderTCP(c *net.TCPConn, in []byte) {
	inmsg := new(dns.Msg)
	inmsg.Unpack(in)
	if inmsg.MsgHdr.Response == true {
		// Uh... answering to an response??
		// dont think so
		return
	}
	out := createpkgtsig(inmsg.MsgHdr.Id, true, c.RemoteAddr())
	SendTCP(out, c)
}

func TestResponderTsig(t *testing.T) {
	/* udp servertje */
	su := new(Server)
	su.Address = "127.0.0.1"
	su.Port = "8053"
	var us *servtsig
	uch := make(chan os.Error)
	go su.NewResponder(us, uch)

	/* tcp servertje */
	st := new(Server)
	st.Address = "127.0.0.1"
	st.Port = "8053"
	st.Tcp = true
	var ts *servtsig
	tch := make(chan os.Error)
	go st.NewResponder(ts, tch)
	time.Sleep(1 * 1e9)
	uch <- nil
	tch <- nil
}