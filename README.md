#NotSh

Give your users ssh, but not shell!
This will prevent users from running commands but allow port tunneling. It probably isn't 100% foolproof, but it certainly works for basic cases. 

#Features

In addition to preventing users from executing commands, there are some basic features which will be expanded upon as needed (PRs welcome!).

- Send a message to a session
- Kill a session

#Install

Something like this should work


```
wget ...../notsh -O /bin/notsh
echo "/bin/notsh" >>  /etc/shells
chmod +x /bin/notsh
chsh USERNAME
	/bin/notsh
```

Now try to log in!

Each login will create a /home/user/1478759535685599744.session where the number is the unix nano time of the login. There is also a log file at /home/user/notsh.log. You can connect via

```
socat UNIX-CONNECT:/home/user/1478759535685599744.session -
```

#Warning! 

Do not do this to your root account, or if the account is supposed to be for accessing the system. Likewise, do not replace /bin/bash, as nobody will be able to do anything!
 
#Demo

...Coming soon.
