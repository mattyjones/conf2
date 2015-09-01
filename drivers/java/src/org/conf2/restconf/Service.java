package org.conf2.restconf;

import org.conf2.yang.browse.Browser;
import org.conf2.yang.StreamSource;
import org.conf2.yang.driver.Driver;
import org.conf2.yang.driver.DriverHandle;

/**
 * Restconf Service
 */
public class Service {
    private DriverHandle service;
    private Driver driver;
    private static boolean loaded;

    public Service(Driver driver) {
        this.driver = driver;
        service = driver.newHandle(newService(), this);
    }

    public void start() {
        startService(service.getId());
    }

    public void setDocRoot(StreamSource docroot) {
        sendSetDocRoot(service.getId(), docroot);
    }

    public void registerBrowser(Browser browser) {
        long moduleHndId = driver.getHandle(browser.getModule()).getId();
        registerBrowserWithService(service.getId(), moduleHndId, browser);
    }

    native void sendSetDocRoot(long serviceHndId, StreamSource loader);

    native long newService();

    native void startService(long serviceHndId);

    native void registerBrowserWithService(long serviceHndId, long moduleHndId, Browser browser);
}
