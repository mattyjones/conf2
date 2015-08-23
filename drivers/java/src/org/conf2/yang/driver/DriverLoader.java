package org.conf2.yang.driver;

import org.conf2.yang.browse.Browser;
import org.conf2.yang.browse.ModuleBrowser;
import org.conf2.yang.DataSource;
import org.conf2.yang.Loader;
import org.conf2.yang.Module;

/**
 *
 */
public class DriverLoader implements Loader {
    private Driver driver;

    public DriverLoader(Driver driver) {
        this.driver = driver;
    }

    public Module loadModule(DataSource source, String resource) {
        DriverHandle ds = driver.newDataSource(source);
System.out.println("ds.hnd", ds.reference);
        ModuleBrowser moduleBrowser = new ModuleBrowser(null);
        DriverHandle moduleBrowserHnd = loadModule(ds, resource, moduleBrowser);
        return moduleBrowser.module;
    }

    native static DriverHandle loadModule(DriverHandle datasource_hnd, String resource, Browser yang_browser);
}
