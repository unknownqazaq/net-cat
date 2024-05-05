# net-cat
## About Project
This project consists on recreating the `NetCat` in a `Server-Client` Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open `TCP` connections, send UDP packages, listen on arbitrary `TCP` and UDP ports and many more.

To see more information about NetCat inspect the manual `man nc`.
## Usage
```
$ go run . $port 
```
### Example
- `1st` Terminal
```bash
$ go run . 8989
Listening on the port :8989
Client connected: [::1]:63747
Client disconnected: [::1]:63747

$ 
```
- `2nd` Terminal
```bash
$ go run . 8989
$ nc localhost 8989
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Yenlik
[2020-01-20 16:03:43][Yenlik]:hello
[2020-01-20 16:03:46][Yenlik]:How are you?
[2020-01-20 16:04:10][Yenlik]:
Lee has joined our chat...
[2020-01-20 16:04:15][Yenlik]:
[2020-01-20 16:04:32][Lee]:Hi everyone!
[2020-01-20 16:04:32][Yenlik]:
[2020-01-20 16:04:35][Lee]:How are you?
[2020-01-20 16:04:35][Yenlik]:great, and you?
[2020-01-20 16:04:41][Yenlik]:
[2020-01-20 16:04:44][Lee]:good!
[2020-01-20 16:04:44][Yenlik]:
[2020-01-20 16:04:50][Lee]:alright, see ya!
[2020-01-20 16:04:50][Yenlik]:bye-bye!
[2020-01-20 16:04:57][Yenlik]:
Lee has left our chat...
[2020-01-20 16:04:59][Yenlik]:
[2020-01-20 16:05:02][Yenlik]:^C
$ 
```
- `3rd` Terminal
```bash
$go run . 8989
Yenliks-MacBook-Air:simpleTCPChat ybokina$ nc localhost 8989
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Lee
[2020-01-20 16:03:43][Yenlik]:hello
[2020-01-20 16:03:46][Yenlik]:How are you?
[2020-01-20 16:04:15][Lee]:Hi everyone!
[2020-01-20 16:04:32][Lee]:How are you?
[2020-01-20 16:04:35][Lee]:
[2020-01-20 16:04:41][Yenlik]:great, and you?
[2020-01-20 16:04:41][Lee]:good!
[2020-01-20 16:04:44][Lee]:alright, see ya!
[2020-01-20 16:04:50][Lee]:
[2020-01-20 16:04:57][Yenlik]:bye-bye!
[2020-01-20 16:04:57][Lee]:^C
$ 
```
## Author
- [unknownqazaq]
