# MTR - Traceroute
![{003288F6-E30C-4BD5-A2C6-0A1F9A290291}](https://github.com/user-attachments/assets/c6e46c3c-5443-418c-a9aa-658a6a6c89ba)

# Features:
- Faster than the integrated `tracert` Tool. It achieves this by only sending 1 probe for each TTL and reduced delays.
- Shows ASN and RDNS information via my external API Service, this allows for seeing rDNS information even when being blocked by your default DNS Server.
- Very light and its easy to install, just drag mtr.exe inside your System32 folder, add the Firewall Rules, then you can call it from any CMD.

# Usage
- `go-mtr <ipv4/ipv6>`

# It doesn't work for me, what do i do?
You need to allow ICMP(v6) trough Windows Firewall rules:

`Open a CMD as Administrator in Windows`

`netsh advfirewall firewall add rule name=AllowICMP protocol=ICMPv4 dir=in action=allow`

`netsh advfirewall firewall add rule name=AllowICMPv6 protocol=ICMPv6 dir=in action=allow`

# It gets detected by my Antivirus, is this malicious?
No i ensure you it isn't malicious, many Antivirus vendors falsely flag any program written in Go as malicious.
~~This Program only depends on 1 non-internal go library: [github.com/t0stbrot/go-ping](https://github.com/t0stbrot/go-ping), you can see all the source code of the library and this tool itself on here.~~ This library is now included in this Repo directly.
