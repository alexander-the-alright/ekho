# ekho
##### _Needlessly Networked MOTD generator_ [WIP]
## Why?
To test my abilities to develop networked programs, and to enhance my skills in writing such.
## Requirements
### Hardware
- Client machine
- Server machine (not required to be distinct from Client)
### Software
- [Golang]
- Linux distro (highly recommened)
- ```make``` is also used in this guide, but the makefile is incredibly simple, and can be skipped entirely by running ```go build``` instead
### Installation
Golang can be installed on Mac with
``` sh
brew install golang
```
Clone the repo
``` sh
git clone https://github.com/alexander-the-alright/ekho
cd ekho
```
Using your favorite text editor, change the destination IP and Port to whatever works for you. This is defined on line 247 in client.go as 127.0.0.1:1300.

Now the server file needs to be run on the preferred port, as specified in the client file. This is defined on line 257 in server.go.


The binary(s) can be obtained using the makefile
If the server is being run on a separate machine, only run make with client flag
``` sh
make client
```
If both binaries are being run on the same machine, ```make all``` will suffice. Although, both of these commands can be bypassed with ```go build```, as this is all ```make``` does anyway.
``` sh
make all
```
At this point, the user should feel free to change the names of the Bash scripts and binaries as they see fit, making sure to update the calls inside the Bash scripts as necessary.

The client binary needs to run on login, the directory on the Raspberry Pi where these scripts go is ```/etc/update-motd.d/```.
``` sh
sudo chmod u+x quote-run.sh
sudo mv quote-run.sh /etc/update-motd.d/01-quote-run
```
The server binary and complementary script need to be running at all times. Should they be run on a distinct machine, this is also the time to move them, as well as the list of quotes they pull from.
``` sh
scp server.go user@ip:~/path/to/file/
scp server-run.sh user@ip:~/path/to/file/
scp list.q user@ip:~/path/to/file/
```
On that machine, the server will need to be compiled, so the makefile can be copied with ```scp``` just the same, or compiled using ```go build```. It is important that list.q is in the same directory as the go binary.

From there, on the server machine, to avoid inconvenient server crashes, move the server script, just like the client script.
``` sh
sudo chmod u+x server-run.sh
sudo mv server-run.sh /etc/update-motd.d/01-server-run
```
This installation method requires the user to be logged in at all times, and will run multiple instances should the same user log in multiple times, so it is recommended to have a either have a main user that never logs out, or a dummy user that auto-logs in on reboot. (If someone knows a better way to have one instance of the software always run, please let me know)

Finally done!
### Usage
The client uses the following flags
``` sh
-h      - display help message
-a <s>  - adds required argument <s> as a quote to server
-l      - list all quotes
-r [i]  - remove; may be run with or without argument
-d      - specify alternate server destination
-log    - [NOT YET IMPLEMENTED] gets logfile from server (contains use and error histories)
```

[golang]: <https://go.dev/doc/install>
