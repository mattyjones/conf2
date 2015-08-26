package org.conf2.yang.driver;

import org.conf2.yang.DataSource;

import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.Map;

public class Driver {
    private static boolean init;
    static {
        System.loadLibrary("yangc2j");
    }

    public Driver() {
        if (!init) {
            initializeDriver();
            init = true;
        }
    }

    private native void initializeDriver();

    public native String echoTest(DataSource loader, String resourceId);

    public native DriverHandle newDataSource(DataSource source);

    public static String printException(Throwable t) {
        StringWriter out = new StringWriter();
        String msg = t.getMessage();
        if (msg != null) {
            out.append(t.getMessage());
            out.append('\n');
        } else {
            out.append(t.toString());
            out.append('\n');
        }
        PrintWriter pout = new PrintWriter(out);
        t.printStackTrace(new PrintWriter(out));
        pout.flush();
        return out.toString();
    }
}
