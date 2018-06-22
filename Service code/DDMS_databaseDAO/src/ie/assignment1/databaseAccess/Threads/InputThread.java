package ie.assignment1.databaseAccess.Threads;

import ie.assignment1.databaseAccess.DAO.DAO;
import org.json.JSONObject;

public class InputThread implements Runnable {

    private Thread runner;
    private String message;
    private DAO dao;
    public InputThread(byte[] message,DAO dao) {
        runner = new Thread(this);
        this.dao = dao;
        try {
            this.message = new String(message, "UTF-8");
        } catch (Exception e) {
            System.out.println("Exception occurred for sentiment " + message);
        }
        runner.start();
    }
    public void run() {
        JSONObject jsonObj = new JSONObject(message);
        double sentiment = jsonObj.getDouble("Sentiment");
        String measurement = jsonObj.getString("Source");
        long time = jsonObj.getLong("Time");
        System.out.println("Writing sentiment "+sentiment + " to measurement " + measurement + " at " + time);
        dao.writeData(sentiment,measurement,time);
    }
}