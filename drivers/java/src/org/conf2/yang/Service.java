package org.conf2.yang;

/**
 * Restconf Service
 */
public class Service {
    private DataSource docroot;
    private DriverHandle service;

    static {
        System.loadLibrary("yangj");
    }

    public Service(Driver driver) {
    }

    public void start() {
        service = startNewService();
    }

    public void setDocRoot(DataSource docroot) {
        this.docroot = docroot;
        sendSetDocRoot(service, docroot);
    }

    native void sendSetDocRoot(DriverHandle serviceHnd, DataSource loader);

    native DriverHandle startNewService();

    public native static String echoTest(DataSource loader, String resourceId);
}
