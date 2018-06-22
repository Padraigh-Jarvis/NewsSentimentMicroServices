package ie.assignment1.databaseAccess.DAO;

import org.influxdb.InfluxDB;
import org.influxdb.InfluxDBFactory;
import org.influxdb.dto.Point;
import org.influxdb.dto.Query;
import org.influxdb.dto.QueryResult;

import java.util.concurrent.TimeUnit;

public class DAO {

    private final String HOSTNAME = "http://influxdb:8086";
    private final String USERNAME = "root";
    private final String PASSWORD = "root";
    private final String DBNAME = "assignment2";
    private InfluxDB influxDB;



    public DAO(){
        influxDB = InfluxDBFactory.connect(HOSTNAME,USERNAME,PASSWORD);
        String setDbString = "CREATE DATABASE "+DBNAME;
        Query query = new Query(setDbString,DBNAME);
        influxDB.query(query);
        influxDB.setDatabase(DBNAME);

    }

    public void writeData(double result,String measurement,long time){
        influxDB.write(Point.measurement(measurement)
                .time(time, TimeUnit.MILLISECONDS)
                .addField("sentiment", result)
                .build());
    }
    public QueryResult readData(String command){
        Query query = new Query(command,DBNAME);
        QueryResult result = influxDB.query(query);
        if(result.getResults().get(0).getSeries() == null)
            return null;
        else
            return result;

    }


}
