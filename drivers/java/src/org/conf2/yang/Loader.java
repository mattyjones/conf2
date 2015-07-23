package org.conf2.yang;


import org.conf2.yang.comm.DataSource;
import org.conf2.yang.comm.Driver;
import org.conf2.yang.comm.DriverHandle;

/**
 *
 */
public class Loader {
    private Driver driver;

    public Loader(Driver driver) {
        this.driver = driver;
    }

    public Module loadModule(DataSource source, String resource) {
        DriverHandle ds = driver.newDataSource(source);
        loadModule(ds, resource);
        // TODO Get stream of module meta
        return null;
    }

    native static DriverHandle loadModule(DriverHandle hnd, String resource);
}
