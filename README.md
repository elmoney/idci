# idci
Interbank Delegated Client Identification

### Prerequisites
* [Git client](https://git-scm.com/downloads)
* [Go](https://golang.org/) - 1.6 or later
* 

### Steps

#### Install GO!
*https://golang.org/doc/install

#### Set your GOPATH
Make sure you have properly setup your Host's [GOPATH environment variable](https://github.com/golang/go/wiki/GOPATH). This allows for both building within the Host and the VM.

#### Clone the Peer project
Create a fork of the fabric repository using the GitHub web interface. Next, clone your fork in the appropriate location.

cd $GOPATH/src
mkdir -p github.com/hyperledger-bankid
cd github.com/hyperledger-bankid
git clone https://github.com/elmoney/idci

#### Install RocksDB
  - [RocksDB](https://github.com/facebook/rocksdb/blob/master/INSTALL.md) version 4.1 and its dependencies
```
apt-get install -y libsnappy-dev zlib1g-dev libbz2-dev
cd /tmp
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
git checkout v4.1
PORTABLE=1 make shared_lib
INSTALL_PATH=/usr/local make install-shared
```
- Execute the following commands:
```
cd $GOPATH/src/github.com/hyperledger-bankid
CGO_CFLAGS=" " CGO_LDFLAGS="-lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy" go install
```

#### Build project
cd $GOPATH/src/github.com/hyperledger-bankid
go build

#### Build/Run project
go build -o peer && reset && CORE_PEER_ID=vp1 ./peer peer

