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
        //DriverHandle ds = driver.newDataSource(source);
//System.out.printf("ds.reference.length=%d\n", ds.reference.array().length);
        ModuleBrowser moduleBrowser = new ModuleBrowser(null);
        DriverHandle moduleBrowserHnd = loadModule(source, resource, moduleBrowser);
        return moduleBrowser.module;
    }

    native static DriverHandle loadModule(DataSource datasource, String resource, Browser yang_browser);
}
