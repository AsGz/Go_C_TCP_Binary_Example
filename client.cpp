/*
   client.cpp

   Test client for the tcpsockets classes. 

   ------------------------------------------

   Copyright © 2013 [Vic Hargrave - http://vichargrave.com]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

#include <stdio.h>
#include <stdlib.h>
#include <sstream>  
#include <string>  
#include <iostream>  
#include "tcpconnector.h"

using namespace std;

void hexDump (char *desc, void *addr, int len) {
    int i;
    unsigned char buff[17];
    unsigned char *pc = (unsigned char*)addr;

    // Output description if given.
    if (desc != NULL)
        printf ("%s:\n", desc);

    if (len == 0) {
        printf("  ZERO LENGTH\n");
        return;
    }
    if (len < 0) {
        printf("  NEGATIVE LENGTH: %i\n",len);
        return;
    }

    // Process every byte in the data.
    for (i = 0; i < len; i++) {
        // Multiple of 16 means new line (with line offset).

        if ((i % 16) == 0) {
            // Just don't print ASCII for the zeroth line.
            if (i != 0)
                printf ("  %s\n", buff);

            // Output the offset.
            printf ("  %04x ", i);
        }

        // Now the hex code for the specific character.
        printf (" %02x", pc[i]);

        // And store a printable ASCII character for later.
        if ((pc[i] < 0x20) || (pc[i] > 0x7e))
            buff[i % 16] = '.';
        else
            buff[i % 16] = pc[i];
        buff[(i % 16) + 1] = '\0';
    }

    // Pad out last line if not exactly 16 characters.
    while ((i % 16) != 0) {
        printf ("   ");
        i++;
    }
    // And print the final ASCII bit.
    printf ("  %s\n", buff);
}

int main(int argc, char** argv)
{
    if (argc != 3) {
        printf("usage: %s <port> <ip>\n", argv[0]);
        exit(1);
    }

    int len;
    string message;
    char line[256];
    TCPConnector* connector = new TCPConnector();
    TCPStream* stream = connector->connect(argv[2], atoi(argv[1]));
    if (stream) {
        uint8_t head = 0x88;
        uint8_t t = 0x7;
        string s = "I am a string data";
        uint16_t l = uint16_t(1 + 4 + 8 + 4 + 10 + s.size());
        uint8_t cmdId = 0x1;
        uint32_t nTime = uint32_t(1234567788);
        double float64 = double(11111111.11);
        float float32 = float(3.333333333);
        uint8_t byte[10] = {1,2,3,4,5,6,7,8,9,10};
        
        ostringstream buf;
        buf << head 
            << t 
            << l
            << cmdId
            << nTime
            << float64
            << float32
            << byte
            << s;
        
        printf("send len:%d\n", l+4);
        hexDump("send buff", (void*)buf.str().c_str(), (int)buf.str().size());
        //stream->send(buf.str().c_str(), buf.str().size());
        printf("head:%lu, t:%lu, len:%lu, cmdId:%lu, nTime:%lu, float64:%lu, float32:%lu, byte:%lu, sring:%lu\n", 
            sizeof(head),
            sizeof(t),
            sizeof(l),
            sizeof(cmdId),
            sizeof(nTime),
            sizeof(float64),
            sizeof(float32),
            sizeof(byte),
            s.size());

        //采用stringstream序列化的二进制数据,有问题,只好采用单个发送方式
        stream->send((char*)&head, sizeof(head));
        stream->send((char*)&t, sizeof(t));
        stream->send((char*)&l, sizeof(l));
        stream->send((char*)&cmdId, sizeof(cmdId));
        stream->send((char*)&nTime, sizeof(nTime));
        stream->send((char*)&float64, sizeof(float64));
        stream->send((char*)&float32, sizeof(float32));
        stream->send((char*)byte, sizeof(byte));
        stream->send((char*)s.c_str(), s.size());
         
        len = stream->receive(line, sizeof(line));

        line[len] = 0;
        printf("receive all:%d\n", len);
        hexDump("receive buff", (void*)line, len);
        delete stream;
    }

    exit(0);
}
