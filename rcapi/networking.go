package rcapi

import (    
    "golang.org/x/net/proxy"
    "net/http"
    "net/url" 
    tls "github.com/refraction-networking/utls"
    fhttp "github.com/useflyent/fhttp"  
    "github.com/gorilla/websocket" 
    "time"   
    "log"
    "io"
    "net"
    "context"
    "errors"
    "strings"
    "io/ioutil"
    "crypto/sha256"
    "rollercoin-bot/constants"
)

func (user *RCUser) WriteWsMessage(data []byte) {
    for {
        err := user.WSConn.WriteMessage(websocket.TextMessage, data)
        if err != nil {
            log.Printf("WriteWsMessage: failed to write ws message: %v\n", err)
            user.DialWsRetry()
        } else {
            return
        }
    }
}

func (user *RCUser) ReadWsMessage() []byte {
    for {
        _, message, err := user.WSConn.ReadMessage()
        if err != nil {
            log.Printf("ReadWsMessage: failed to receive ws message: %v\n", err)
            user.DialWsRetry()
        } else {
            return message
        }
    }
}

func (user *RCUser) DialWs() error {
    wsUrl := constants.BaseWebsocketURL + "/cmd?token=" + user.token

    dialer := NewTorProxyWsDialer(user.Proxy)

    header := make(http.Header)

    header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
    header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

    var err error

    user.WSConn, _, err = dialer.Dial(wsUrl, header)
    if err != nil {
        return err
    }

    return nil
}

func (user *RCUser) DialWsRetry() {
    for {
        err := user.DialWs()

        if err != nil {
            log.Printf("DialWsRetry: %s\n", err.Error())
            time.Sleep(constants.RetryDelay)
        } else {
            return
        }
    }
}

var respBody []byte

func (user *RCUser) PostRequest(URL string, data io.Reader) ([]byte, error) {
    ja3_transport := NewTorProxyTransportConfig(user.Proxy)

    client := &fhttp.Client{
        Transport: ja3_transport,
    }

    req, err := fhttp.NewRequest("POST", URL, data)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("content-type", "application/json")
    req.Header.Set("authorization", "Bearer " + user.token)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

func (user *RCUser) PostRequestRetry(url string, data io.Reader) []byte {
    for {
        response, err := user.PostRequest(url, data)

        if err != nil {
            log.Printf("PostRequestRetry: %v\n", err.Error())
            time.Sleep(constants.RetryDelay)
        } else {
            return response
        }
    }
}


/*
Tor 12.5.2 version

TLSv1.3 Record Layer: Handshake Protocol: Client Hello
    Content Type: Handshake (22)
    Version: TLS 1.0 (0x0301)
    Length: 512
    Handshake Protocol: Client Hello
        Handshake Type: Client Hello (1)
        Length: 508
        Version: TLS 1.2 (0x0303)
        Random: 25acc7ccbed526c29b35a6ae6bdb94fa4ae7b91f3e7d1d4cde90c8ef8873992d
        Session ID Length: 32
        Session ID: 7f2bf2be7c517739d409af5dc8a337317b621b3a78f5125e3ccc94830b7e1744
        Cipher Suites Length: 22
        Cipher Suites (11 suites)
        Compression Methods Length: 1
        Compression Methods (1 method)
        Extensions Length: 413
        Extension: server_name (len=19)
        Extension: extended_master_secret (len=0)
        Extension: renegotiation_info (len=1)
        Extension: supported_groups (len=14)
        Extension: ec_point_formats (len=2)
        Extension: application_layer_protocol_negotiation (len=14)
        Extension: status_request (len=5)
        Extension: delegated_credentials (len=10)
        Extension: key_share (len=107)
        Extension: supported_versions (len=5)
        Extension: signature_algorithms (len=24)
        Extension: record_size_limit (len=2)
        Extension: padding (len=158)
        [JA3 Fullstring: 771,4865-4867-4866-49195-49199-52393-52392-49196-49200-156-157,0-23-65281-10-11-16-5-34-51-43-13-28-21,29-23-24-25-256-257,0]
        [JA3: c79653a3a53172c2304e6da72cd7aa2a]
*/

func newTorSpec() *tls.ClientHelloSpec {
    torSpec := &tls.ClientHelloSpec{
        CipherSuites: []uint16{
            4865, 4867, 4866,
            49195, 49199, 52393, 52392, 49196, 49200,
            156, 157,
        },
        Extensions: []tls.TLSExtension{
            &tls.SNIExtension{}, // 0
            &tls.ExtendedMasterSecretExtension{}, // 23
            &tls.RenegotiationInfoExtension{ // 65281
                Renegotiation: tls.RenegotiateOnceAsClient,
            },
            &tls.SupportedCurvesExtension{Curves: []tls.CurveID{23, 24, 25, 29, 256, 257}}, // 10
            &tls.SupportedPointsExtension{SupportedPoints: []uint8{}}, // 11
            &tls.ALPNExtension{ // 16
                AlpnProtocols: []string{"http/1.1"},
            },
            &tls.StatusRequestExtension{}, // 5
            &tls.KeyShareExtension{ // 51
                KeyShares: []tls.KeyShare{
                    {Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
                    {Group: tls.X25519},
                },
            },  
            &tls.SupportedVersionsExtension{Versions: []uint16{ // 43
                tls.GREASE_PLACEHOLDER,
                tls.VersionTLS13,
                tls.VersionTLS12,
                tls.VersionTLS11,
                tls.VersionTLS10}},
            &tls.SignatureAlgorithmsExtension{ // 13
                SupportedSignatureAlgorithms: []tls.SignatureScheme{
                    tls.ECDSAWithP256AndSHA256,
                    tls.ECDSAWithP384AndSHA384,
                    tls.ECDSAWithP521AndSHA512,
                    tls.PSSWithSHA256,
                    tls.PSSWithSHA384,
                    tls.PSSWithSHA512,
                    tls.PKCS1WithSHA256,
                    tls.PKCS1WithSHA384,
                    tls.PKCS1WithSHA512,
                    tls.ECDSAWithSHA1,
                    tls.PKCS1WithSHA1,
                },
            },
            &tls.FakeRecordSizeLimitExtension{}, // 28
            &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle}, // 21
        },
        CompressionMethods: []byte{0},
        GetSessionID: sha256.Sum256,
    }

    return torSpec
}

func NewTorProxyWsDialer(s5_url string) (dialer *websocket.Dialer) {
    spec := newTorSpec()

    config := &tls.Config{
            InsecureSkipVerify: true,
            MaxVersion:         tls.VersionTLS13,
    }

    dialer = websocket.DefaultDialer
    dialer.NetDialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
        proxyURI, e := url.Parse(s5_url)
        if e != nil {
            return nil, e
        }

        if proxyURI.Scheme != "socks5" {
            return nil, errors.New("unsupported proxy scheme: " + proxyURI.Scheme)
        }

        d := &net.Dialer{}
        proxyDialer, e := proxy.FromURL(proxyURI, d)
        if e != nil {
            return nil, e
        }

        dialConn, e := proxyDialer.Dial("tcp", addr)
        if e != nil {
            return nil, e
        }

        config.ServerName = strings.Split(addr, ":")[0]

        uTLSConn := tls.UClient(dialConn, config, tls.HelloCustom)
        if err := uTLSConn.ApplyPreset(spec); err != nil {
            return nil, err
        }
        if err := uTLSConn.Handshake(); err != nil {
            return nil, err
        }
        return uTLSConn, nil
    }

    return dialer
}

// https://github.com/aj3423/whatsapp-go/blob/master/net/ja3transport.go
func NewTorProxyTransportConfig(s5_url string) *fhttp.Transport {
    config := &tls.Config{
            InsecureSkipVerify: true,
            MaxVersion:         tls.VersionTLS13,
    }

    spec := newTorSpec()

    dialtls := func(network, addr string) (net.Conn, error) {
        proxyURI, e := url.Parse(s5_url)
        if e != nil {
            return nil, e
        }

        if proxyURI.Scheme != "socks5" {
            return nil, errors.New("unsupported proxy scheme: " + proxyURI.Scheme)
        }

        d := &net.Dialer{}
        proxyDialer, e := proxy.FromURL(proxyURI, d)
        if e != nil {
            return nil, e
        }

        dialConn, e := proxyDialer.Dial("tcp", addr)
        if e != nil {
            return nil, e
        }

        config.ServerName = strings.Split(addr, ":")[0]

        uTLSConn := tls.UClient(dialConn, config, tls.HelloCustom)
        if err := uTLSConn.ApplyPreset(spec); err != nil {
            return nil, err
        }
        if err := uTLSConn.Handshake(); err != nil {
            return nil, err
        }
        return uTLSConn, nil
    }

    return &fhttp.Transport{DialTLS: dialtls}
}