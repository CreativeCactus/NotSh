#NotSh

Give your users ssh, but not shell!

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
 
#Demo

...
