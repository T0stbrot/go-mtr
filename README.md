# Traceroute tool written in go
Features:
- Faster than the Windows integrated tracert tool, by a lot.
- Shows ASN and RDNS information via my external API Service, so it can show RDNS even if your normal DNS blocks it.
- Very light and easy to install, just drag mtr.exe inside your System32 folder, then you can call it from any CMD.

# It doesn't work for me, what do i do?
Sometimes you need to allow ICMP(v6) trough Windows Firewall rules:

`Open a CMD as Administrator in Windows`

`netsh advfirewall firewall add rule name=AllowICMP protocol=ICMPv4 dir=in action=allow`

`netsh advfirewall firewall add rule name=AllowICMPv6 protocol=ICMPv6 dir=in action=allow`

Im looking to "fixing" this so this maybe works without those rules.