# xavior

a tcp forwarder/tunnel with xor

*Experimental and insecure. NO WARRANTY. Use at your own risk!*

## Usage

```
xavior
  -l    <listenHost>         listened address  
  -r    <remoteHost>         remote address forwarding to  
  -sp   <sendPassword>       password when sending data  
  -rp   <receivePassword>    password when receving data  
```

Due to the nature of XOR, there is no need to distinguish the server and client, so the following command can be executed on both sides:

```sh
 ./xavior -l "0.0.0.0:2222" -r "127.0.0.1:22" -sp "123456" -rp "123456"
```

You can also specify different passwords for downstream and uplink. That adds a little security to active detection; but note that it is still insecure!

## Example

Let's say you want to encrypt traffic between your computer and the SSH server (10.10.10.10:22).

Execute this on your SSH server (10.10.10.10):

```sh
 ./xavior -l "0.0.0.0:2222" -r "127.0.0.1:22" -sp "ChangeTh1sToARandomLonger0ne123456" -rp "ChangeTh1sToARandomLonger0ne654321"
```

On your computer, execute:

```sh
 ./xavior -l "127.0.0.1:5555" -r "10.10.10.10:2222" -sp "ChangeTh1sToARandomLonger0ne123456" -rp "ChangeTh1sToARandomLonger0ne654321"
```

Right, no need to swap the `-sp` and `-rp`. Then you can log in to your SSH server though the XOR tunnel:

```sh
ssh username@127.0.0.1 -p5555
```

which is the replacement of

```sh
ssh username@10.10.10.10 -p22
```

Windows is also supported.

## Contribute

Well, it is just like a learning project. Please play your precious time on other ones.
