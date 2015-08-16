package org.conf2.yang.driver;

import org.conf2.yang.DataSource;

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

    public native static String echoTest(DataSource loader, String resourceId);

    public native DriverHandle newDataSource(DataSource source);
}
