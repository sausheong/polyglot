import org.zeromq.ZMQ;
import java.util.UUID;

 
public class Hello {

  public static void main(String[] args) {
    ZMQ.Context context = ZMQ.context(1);
    String routeid = "GET/_/hello/java";
    String identity = UUID.randomUUID().toString();
    
    ZMQ.Socket socket = context.socket(ZMQ.REQ);
    socket.setIdentity(identity.getBytes());
    socket.connect ("tcp://localhost:4321");
 
    System.out.printf("%s - %s responder ready\n", routeid, identity);
    
    socket.send(routeid, 0);
    try {
      while (true) {
        String request = socket.recvStr();        
        
        socket.send(routeid, ZMQ.SNDMORE);
        socket.send("200", ZMQ.SNDMORE);
        socket.send("{\"Content-Type\": \"text/html\"}", ZMQ.SNDMORE);
        socket.send(request);
      }      
    } catch (Exception e) {
      socket.close();
      context.term();      
    } 
  }
}