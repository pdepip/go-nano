package node

import (
    "log"
    "net"
    "bytes"
)

const (
    DefaultPort  = ":7075"      // Default port that nodes on network use
    PacketSize   = 512          // Default packet size

    VersionMax   = 0x05         // Highest accepted version
    VersionMin   = 0x04         // Lowest accepted version
    VersionUsing = 0x05         // Version we are currently using
)

// Non-idiomatic constant names to mirror reference implementation (cpp node)
const (
    Message_invalid byte = iota
    Message_not_a_type
    Message_keepalive
    Message_publish
    Message_confirm_req
    Message_configm_ack
    Message_bulk_pull
    Message_bulk_push
    Message_frontier_req
)

var (
    MagicNumber = [2]byte{'R', 'C'} // Wtf is this

    // Default node to connect to
    DefaultPeer  = Peer{
        //IP:   net.ParseIP("::ffff:138.68.2.234"),
        IP:   net.ParseIP("::ffff:35.196.55.211"),

        Port: 7075,
    }
)


// Placeholder to handle all incoming packets
func handlePacket(buf []byte) {
    log.Println(string(buf))
}


// Listens on default port for all incoming udp messages
// Passes all received messges to handlePacket
func ListenUDP() {

    log.Println("Listening for UDP Packets on: ", DefaultPort)

    ln, err := net.ListenPacket("udp", DefaultPort)
    if err != nil {
        log.Fatalln("[ERROR] net.ListenPacket:", err)
    }

    buf := make([]byte, PacketSize)

    for {
        n, _, err := ln.ReadFrom(buf)
        if err != nil {
            log.Println("hi:", err)
            continue
        }
        log.Println("got one")

        if n > 0 {
            handlePacket(buf[:n])
        }
    }
}


// Connects to default node
func FindDefaultNode() (err error) {

    addr, err := DefaultPeer.Addr()
    if err != nil {
        return
    }

    outConn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        return 
    }

    msg := NewConnectMessage()
    
    // Convert message to byte array (look up fastest method)
    buf := bytes.NewBuffer(nil)
    msg.Write(buf)

    // Send
    outConn.Write(buf.Bytes())
    log.Println("Sent:", string(buf.Bytes()))
    return

}


// Messages sent when connecting to a new node
type MessageConnect struct {
    Header MessageHeader
}


// Creates a new connect message
func NewConnectMessage() *MessageConnect {
    var msg MessageConnect
    msg.Header.MagicNumber  = MagicNumber
    msg.Header.VersionMax   = VersionMax
    msg.Header.VersionMin   = VersionMin
    msg.Header.VersionUsing = VersionUsing
    msg.Header.MessageType  = Message_keepalive
    return &msg
}

// Writes message connect struct to byte array
func (m *MessageConnect) Write(buf *bytes.Buffer) (err error) {
    
    var errs []error

    errs = append(errs,
        buf.WriteByte(m.Header.MagicNumber[0]),
        buf.WriteByte(m.Header.MagicNumber[1]),
        buf.WriteByte(m.Header.VersionMax),
        buf.WriteByte(m.Header.VersionUsing),
        buf.WriteByte(m.Header.VersionMin),
        buf.WriteByte(m.Header.MessageType),
        buf.WriteByte(m.Header.Extensions),
        buf.WriteByte(m.Header.BlockType),
    )

    for _, err := range errs {
        if err != nil {
            return err
        }
    }
    return

}


// Header for all messages sent on the network
type MessageHeader struct {
    MagicNumber  [2]byte
    VersionMax   byte
    VersionUsing byte
    VersionMin   byte
    MessageType  byte
    Extensions   byte
    BlockType    byte
}
