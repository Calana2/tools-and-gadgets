# Small DNS client

import random
import sys, socket

qtypes = {
        "A": 1,
        "NS": 2,
        "CNAME": 5,
        "SOA": 6,
        "PTR": 12,
        "MX": 15,
        "TXT": 16,
        "AAAA": 28,
        "SRV": 33,
        "ANY": 255
    }

def encode_name(name):
    parts = name.split('.')
    res = b""
    for part in parts:
        res += len(part).to_bytes(1) + part.encode()
    res += b"\00"
    return res

def decode_name(data,len):
    name = ""
    i = 0
    while i < len:
        length = data[i]
        if length == 0:
            break
        i += 1
        name += data[i:i+length].decode() + "."
        i += length
    return name

def encode_qtype(qtype):
    return qtypes.get(qtype.upper(), 1).to_bytes(2, byteorder='big')

def decode_type(code):
    for name,v in qtypes.items():
        if v == code: return name

def decode_class(clas):
    return "IN"

def parse_IP(ip):
    octet1 = str(ip & 0xff)
    octet2 = str((ip >> 8)  & 0xff)
    octet3 = str((ip >> 16) & 0xff)
    octet4 = str((ip >> 24) & 0xff)
    return f"{octet4}.{octet3}.{octet2}.{octet1}"
    
def parse_response(res):
    id = int.from_bytes(res[:2])
    flags = int.from_bytes(res[2:4])
    qdcount = int.from_bytes(res[4:6])
    ancount = int.from_bytes(res[6:8])
    nscount = int.from_bytes(res[8:10])
    arcount = int.from_bytes(res[10:12])

    print("ID        :", id)
    print("FLAGS     :", hex(flags))
    print("QDCOUNT   :", qdcount)
    print("ANCOUNT   :", ancount)
    print("NSCOUNT   :", nscount)
    print("ARCOUNT   :", arcount)

    offset = 12
    for _ in range(qdcount):
        while res[offset] != 0:
            offset += res[offset] + 1 # LENGTH BYTE + TEXT
        offset += 1                   # NULL BYTE
        offset += 4                   # QCLASS, QTYPR
    rtype, rclass, ttl, rdlength = 0,0,0,0
    rdata = b""
    for _ in range(ancount):
        offset +=2;              # NAME (0xC0 0x0C)
        rtype = (int.from_bytes(res[offset:offset+2]))
        offset += 2
        rclass = (int.from_bytes(res[offset:offset+2]))
        offset += 2
        ttl = (int.from_bytes(res[offset:offset+4]))
        offset += 4
        rdlength = int.from_bytes(res[offset:offset+2])
        offset += 2
        rdata = res[offset:offset+rdlength]
        offset += rdlength
        print("\n** ANSWER **")
        print("RTYPE     :", decode_type(rtype))
        print("RCLASS    :", decode_class(rclass))
        print("TTL       :", ttl, "seconds")
        print("RDLENGTH  :", rdlength, "bytes")
        if rtype == 1 and rdlength == 4:  # A IPV4
            rdata = parse_IP(int.from_bytes(rdata,byteorder="little"))
        elif rtype == 5: # CNAME
            rdata = decode_name(rdata,rdlength)
        # TODO the rest of them
        print("RDATA     :", rdata)

if __name__ == '__main__':
    if len(sys.argv) < 5:
        print(f'Usage {sys.argv[0]} <ip> <port> <query_type> <domain1> [domain2]...')
        sys.exit(0)

    ip = sys.argv[1]
    port = sys.argv[2]
    qtype = sys.argv[3]
    for domain in sys.argv[4:]:
        # ** Headers **
        # ID
        DNS_packet = random.getrandbits(16).to_bytes(2)
        # FLAGS
        DNS_packet += b"\x01\x00" # QR=0, OPCODE=QUERY, RD=1
        # QDCOUNT
        DNS_packet += b"\x00\x01"
        # ANCOUNT
        DNS_packet += b"\x00\x00"
        # NSCOUNT
        DNS_packet += b"\x00\x00"
        # ARCOUNT
        DNS_packet += b"\x00\x00"
        # ** Question Section **
        # QNAME, QTYPE, QCLASS
        DNS_packet += encode_name(domain)
        DNS_packet += encode_qtype(qtype)
        DNS_packet += b"\x00\x01"  # IN         
        try: 
            assert len(DNS_packet) <= 512, "Query exceeded the 512 byte length"
        except AssertionError as e:
            print(e); sys.exit(1)

        s = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
        s.sendto(DNS_packet, (ip, int(port)))
        response, _ = s.recvfrom(1024)
        parse_response(response)
