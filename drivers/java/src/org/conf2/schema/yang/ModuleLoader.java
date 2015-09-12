package org.conf2.schema.yang;

import org.conf2.schema.Module;
import org.conf2.schema.StreamSource;
import org.conf2.schema.browse.Browser;
import org.conf2.schema.browse.ModuleBrowser;
import org.conf2.schema.driver.Driver;
import org.conf2.schema.driver.DriverError;

/**
 *
 */
public class ModuleLoader {
    private Driver driver;
    public ModuleLoader(Driver d) {
        this.driver = d;
    }

    public Module loadModule(StreamSource source, String resource) {
        ModuleBrowser moduleBrowser = new ModuleBrowser(null);
        long handle = loadModule(source, resource, moduleBrowser);
        if (handle == 0) {
            throw new DriverError("Could not load module " + resource);
        }
        driver.newHandle(handle, moduleBrowser.module);
        return moduleBrowser.module;
    }

    private native long loadModule(StreamSource resourceSourceHand, String resource, Browser yang_browser);
}
