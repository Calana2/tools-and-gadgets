#!/usr/bin/env python3
''' Symbol-to-Stack Encoder (Null-Free) (NASM sintaxis)'''

import sys
import textwrap
import binascii

if len(sys.argv) < 3 or sys.argv[1] == "-h" or sys.argv[1] == "--help":
    print("Usage: {} <symbol_name> <arch>".format(sys.argv[0]))
    exit(0)

arch = sys.argv[2]
if arch not in ("x86","x86_64"):
    print("Arch must be 'x86' or 'x86_64'")
    exit(0)

symbol_name = sys.argv[1]

arch_reg_size = dict()
arch_reg_size['x86'] = 4
arch_reg_size['x86_64'] = 8

size_words = dict()
size_words[1] = 'byte'
size_words[2] = 'word'
size_words[4] = 'dword'
size_words[8] = 'qword'

arch_stack_register = dict()
arch_stack_register['x86'] = 'esp'
arch_stack_register['x86_64'] = 'rsp'

valid_remaining_sizes = [4,2,1]

size = arch_reg_size[arch]
reg = arch_stack_register[arch]
boxes = textwrap.wrap(symbol_name,size)
margin = " " * 18

print("\nxor ebx,ebx\npush ebx")
for box in boxes[::-1]:
    box = box[::-1]
    if len(box) < size:
        dif = size - len(box)
        box = box.rjust(size,"A")
        print("{}{}{}; {}".format("push 0x", binascii.hexlify(box.encode()).decode(), margin, box))
        if dif in(1,2,4):
            print("sub {} [{} + 0x{}], 0x{}".format(size_words[dif], reg, size-dif, "41"*dif))
        else:
            while dif > 0:
                for vsize in valid_remaining_sizes:
                    if dif >= vsize:
                        hex = "41" * vsize
                        print("sub {} [{} + 0x{}], 0x{}".format(size_words[vsize], reg, size-dif , hex))
                        dif -= vsize
                        break
        continue
    print("{}{}{}; {}".format("push 0x", binascii.hexlify(box.encode()).decode(), margin, box))

print("push " + reg)
