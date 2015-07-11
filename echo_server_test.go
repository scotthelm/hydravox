package main

// func TestMain(m *testing.M) {
// 	es := new(EchoServer)
// 	es.Port = ":7778"
// 	es.Start()
// }

// func TestEcho(t *testing.T) {
// 	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:7778")
// 	checkError(err, t)
// 	conn, err := net.DialTCP("tcp", nil, tcpAddr)
// 	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
// 	checkError(err, t)
// 	_, err = conn.Write([]byte("Helo"))
// 	checkError(err, t)
// 	_, err = ioutil.ReadAll(conn)
// 	checkError(err, t)
// }

// func checkError(err error, t *testing.T) {
// 	if err != nil {
// 		t.Fatal("Error: ", err.Error())
// 	}
// }
