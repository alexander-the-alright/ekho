# ekho
##### _Needlessly Networked MOTD generator_
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

#### NOTE
In previous versions, the suggested method of installation involved placing bash scripts in ```/etc/update-motd.d/```, however, there have been issues getting that installation method working (namely, having MOTD run user-defined commands), and the current quick and extremely dirty suggested installation method is to just append the user's shell config files (```.bashrc```, ```.bash_profile```, ```.zshrc```, etc) with the paths to the executables.

For example, if the server used zsh, the installation could look something like this:
``` sh
echo "~/abs/path/to/server.o &" >> ~/.zshrc
```
Similarly for a client using Bash:
``` sh
echo "~/abs/path/to/client.o" >> ~/.bashrc
```

Any help figuring out the issue with MOTD would be appreciated.

Finally done!
### Usage
The client uses the following flags
``` sh
-h      - display help message
-a <s>  - adds required argument <s> as a quote to server
-l      - list all quotes
-r [i]  - remove; may be run with or without argument
-ip     - specify alternate server destination
```

### How It Works 
The sequence diagram can be viewed in ```diag``` folder, which shows the sequence of networked messages for each usage of ekho.
The source code is in the same folder, ```ekho-seq.txt```, and can be uploaded to [sequencediagram.org][seqdiag] to rerender the diagram and view any updates.


### To Be Implemented
- Multiple transactions for large quote file
- Log all uses/transaction data each day
- Request, send, and receive the logfile for any day
- User-specific quotes/quotefiles
- Installation (client-side, at least) via MOTD


[golang]: <https://go.dev/doc/install>
[seqdiag]: <sequencediagram.org>
