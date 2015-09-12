package org.conf2.schema.driver;

import org.conf2.schema.Module;
import org.conf2.schema.SimpleStreamSource;
import org.conf2.schema.StreamSource;
import org.conf2.schema.browse.BrowseRead;
import org.conf2.schema.browse.Browser;
import org.conf2.schema.browse.ModuleBrowser;
import org.conf2.schema.browse.Selection;
import org.conf2.schema.yang.ModuleLoader;

import java.io.File;

/**
 *
 */
public class DriverTestHarness implements Browser {
    private Driver driver;
    private Module module;
    private ModuleLoader loader;
    private DriverHandle harnessHandle;
    private Selection root;
    private Selection[] testSelection = new Selection[1];

    public DriverTestHarness(Driver d) {
        this.driver = d;
        this.loader = new ModuleLoader(this.driver);
        loadModule();
        long handle = newTestHarness(d.getHandle(module).getId(), this);
        harnessHandle = d.newHandle(handle, this);
        root = new Selection();
        root.meta = module;
        root.Enter = () -> {
            root.found = true;
            return testSelection[0];
        };
    }

    public void release() {
        harnessHandle.Release();
    }

    public String getReport() {
        return report(harnessHandle.getId());
    }

    void loadModule() {
        String yangPath = System.getenv().get("YANGPATH");
        if (yangPath == null) {
            throw new RuntimeException("Missing YANGPATH environment variable");
        }
        StreamSource yangPathSource = new SimpleStreamSource(new File(yangPath));
        module = loader.loadModule(yangPathSource, "test/functional.yang");
    }

    public boolean runTest(String testname, Selection s) {
        this.testSelection[0] = s;
        return runTest(harnessHandle.getId(), testname);
    }

    @Override
    public Selection getRootSelector() {
        return root;
    }

    @Override
    public Module getModule() {
        return module;
    }

    private native long newTestHarness(long module_hnd, Browser harnessTester);
    private native boolean runTest(long harnessHandle, String testname);
    private native String report(long harnessHandle);
}
