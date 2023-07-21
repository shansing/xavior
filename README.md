# xavior

a tcp forwarder with xor

**Experimental. Use at your own risk!**

CLI params:

> -l listened address  
> -r remote address forwarding to  
> -sp password when sending data  
> -rp password when receving data  

Due to the XOR characteristics, there is no need to distinguish the server and client, so the following command can be executed on both sides:

```sh
 ./xavior -l "0.0.0.0:2222" -r "127.0.0.1:22" -sp "123456" -rp "123456"
```

Or you may want to use different passwords for downstream or uplink, then execute the following command on both sides :

```sh
 ./xavior -l "0.0.0.0:2222" -r "127.0.0.1:22" -sp "123456" -rp "654321"
```

That adds a little security to active detection; but note that it is still insecure!

