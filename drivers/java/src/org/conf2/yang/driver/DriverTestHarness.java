package org.conf2.yang.driver;

import org.conf2.yang.Module;
import org.conf2.yang.SimpleStreamSource;
import org.conf2.yang.StreamSource;
import org.conf2.yang.browse.BrowseRead;
import org.conf2.yang.browse.Browser;
import org.conf2.yang.browse.ModuleBrowser;
import org.conf2.yang.browse.Selection;

import java.io.File;

/**
 *
 */
public class DriverTestHarness implements Browser {
    private Driver driver;
    private Module module;
    private DriverHandle harnessHandle;
    private Selection root;
    private Selection[] testSelection = new Selection[1];

    public DriverTestHarness(Driver d) {
        this.driver = d;
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
        module = driver.loadModule(yangPathSource, "test/functional.yang");
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
