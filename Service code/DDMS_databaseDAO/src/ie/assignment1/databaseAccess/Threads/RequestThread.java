package ie.assignment1.databaseAccess.Threads;

import ie.assignment1.databaseAccess.DAO.DAO;
import com.rabbitmq.client.AMQP;

import com.rabbitmq.client.Channel;
import org.influxdb.dto.QueryResult;

import java.util.ArrayList;

public class RequestThread implements Runnable {
    private DAO dao;
    private Thread runner;
    private Channel channel;
    private AMQP.BasicProperties properties;
    private byte[] sources;
    public RequestThread(DAO dao, Channel channel, AMQP.BasicProperties properties, byte[] sources) {
        this.dao=dao;
        this.channel=channel;
        this.properties=properties;
        this.sources=sources;
        runner = new Thread(this);
        runner.start();
    }

    @Override
    public void run() {
        try{
            System.out.println("Average request received");
            AMQP.BasicProperties replyProps = new AMQP.BasicProperties
                    .Builder()
                    .correlationId(properties.getCorrelationId())
                    .build();

            String sourcesAsString = new String(sources);
            String[] listOfSources = sourcesAsString.split(",");
            String returnValue="";

            for (String source : listOfSources){
                System.out.println(source);
                QueryResult result = dao.readData("SELECT sentiment FROM "+ source + " WHERE time > now() - 1d" );
                if(result == null)
                    returnValue+="No data for " + source + " in the last day";
                else
                    returnValue+=parseData(result);
                returnValue+=":";
            }

            channel.basicPublish("",properties.getReplyTo(),replyProps,returnValue.getBytes("UTF-8"));

        }catch (Exception e){e.printStackTrace();}
    }
    private String parseData(QueryResult queryResult){
        double average =0;
        QueryResult.Series values = queryResult.getResults().get(0).getSeries().get(0);
        for(Object obj : values.getValues()){
            ArrayList al = (ArrayList)obj ;
            average+=Double.parseDouble(al.get(1).toString());
        }

        average=average/values.getValues().size();
        return ""+Math.floor((average*100))/100;
    }

}
