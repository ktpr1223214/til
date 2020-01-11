import java.io.*;
import java.net.*;
import java.util.*;


class TCPServer implements Runnable {
    private ServerSocket serverSocket;

    public TCPServer(int port) throws IOException {
        this.serverSocket = new ServerSocket(port);
    }

    public void run() {
        for (;;) {
            try {
                Socket socket = serverSocket.accept();
                new ConnectionHandler(socket).run();
            } catch (IOException e) {
                exception.printStackTrace();
            }
        }
    }
}
