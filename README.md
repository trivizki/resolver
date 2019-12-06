# resolver
University project

# Installation

**For Developers**

Install git

0. `sudo apt install git`

Creating your GOPATH directory.

1. `mkdir -p $HOME/go/src`

2. `cd $HOME/go/src`

Clone this repo into your GOPATH.

3. `git clone https://github.com/trivizki/resolver.git`

4. `cd $HOME/go/src/resolver`

5. `sudo chmod +x ./install.sh`

Run the installation script.

6. `sudo ./install.sh`

Return user permissions

7. `sudo chown -R USER:USER ~`

8. Start Develop.

# Runnig

Compile
1. `make build`

2. configure `build/conf.yml` according to your device (pay attention to network interface's names).

Run the binary.

3. `make run`
