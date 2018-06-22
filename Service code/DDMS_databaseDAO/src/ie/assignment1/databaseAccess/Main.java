package ie.assignment1.databaseAccess;

import com.rabbitmq.client.*;
import ie.assignment1.databaseAccess.DAO.DAO;
import ie.assignment1.databaseAccess.Threads.*;
import org.json.JSONObject;

public class Main {

    private final static String DATA_INPUT_QUEUE = "toDatabase";
    private final static String QUERY_QUEUE_NAME = "DatabaseFetch";
    private final static String TEST_QUEUE_NAME = "testqueue";
    private final static String HOSTNAME = "rabbitmq";
    private static DAO dao;

    private static void connectInputQueue() {
        try {
            ConnectionFactory factory = new ConnectionFactory();
            factory.setHost(HOSTNAME);
            Connection connection = factory.newConnection();
            Channel channel = connection.createChannel();

            channel.queueDeclare(DATA_INPUT_QUEUE, false, false, false, null);
            //System.out.println(" [*] Waiting for new sentiment data");

            Consumer consumer = new DefaultConsumer(channel) {
                @Override
                public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body) {
                    new InputThread(body,dao);
                }
            };
            channel.basicConsume(DATA_INPUT_QUEUE, true, consumer);
        } catch (Exception e) {
            System.out.println("Error: An IO error occurred");
        }
    }

    private static void connectQueryQueue() {
        try {
            ConnectionFactory factory = new ConnectionFactory();
            factory.setHost(HOSTNAME);
            Connection connection = factory.newConnection();
            Channel channel = connection.createChannel();

            channel.queueDeclare(QUERY_QUEUE_NAME, false, false, false, null);

//            System.out.println(" [*] Waiting for database requests.");

            Consumer consumer = new DefaultConsumer(channel) {
                @Override
                public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body) {
                    try{
                        new RequestThread(dao,channel,properties,body);
                }catch (Exception e){e.printStackTrace();}
                }
            };
            channel.basicConsume(QUERY_QUEUE_NAME, true, consumer);
        } catch (Exception e) {
            System.out.println("Error: An IO error occurred");
        }
    }
    private static void connectTestQueue() {
        try {
            ConnectionFactory factory = new ConnectionFactory();
            factory.setHost(HOSTNAME);
            Connection connection = factory.newConnection();
            Channel channel = connection.createChannel();

            channel.queueDeclare(TEST_QUEUE_NAME, false, false, false, null);

            System.out.println(" [*] Waiting for test requests.");

            Consumer consumer = new DefaultConsumer(channel) {
                @Override
                public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body) {
                    String testString ="";
                    try {
                        testString = new String(body, "UTF-8");
                    }catch (Exception e){
                        System.out.println(e.getMessage());
                    }
                    JSONObject jsonObj = new JSONObject(testString);
                    jsonObj.put("Databaseaccessqueue",QUERY_QUEUE_NAME);
                    try {
                        channel.basicPublish("", properties.getReplyTo(), properties, jsonObj.toString().getBytes());
                    }catch (Exception e){
                        System.out.println("Error: An IO error occurred " + e);
                    }
                }
            };
            channel.basicConsume(TEST_QUEUE_NAME, true, consumer);
        } catch (Exception e) {
            System.out.println("Error: An IO error occurred");
        }
    }
    public static void main(String[] argv) {
        dao = new DAO();
        connectInputQueue();
        connectQueryQueue();
        connectTestQueue();

    }
}

