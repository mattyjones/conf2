package org.conf2.yang.driver;

import org.conf2.yang.Handle;
import org.conf2.yang.Module;
import org.conf2.yang.StreamSource;
import org.conf2.yang.browse.Browser;
import org.conf2.yang.browse.ModuleBrowser;

import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

public abstract class Driver {
    // TODO: Use weakreference queue and release objects automatically
    private Map<Object, DriverHandle> handles = new HashMap<Object, DriverHandle>();
    private static boolean loaded;

    protected Driver() {
        if (!loaded) {
            loadLibrary();
            initializeDriver();
            loaded = true;
        }
    }

    public abstract void loadLibrary();

    public static Driver yangDriver() {
        return new Driver() {
            @Override
            public void loadLibrary() {
                //System.loadLibrary("yangc2j");
            }
        };
    }

    public void release() {
        for (Handle handle : handles.values()) {
            handle.Release();
        }
        handles.clear();
    }

    public DriverHandle newHandle(long hndId, Object obj) {
System.out.println("Adding handle " + hndId + " for " + obj.toString());
        DriverHandle h = new DriverHandle(this, hndId);
        handles.put(obj, h);
        return h;
    }

    public DriverHandle getHandle(Object o) {
System.out.println("Looking for handle " + o.toString() + " from handles: " + handles.size());
        return handles.get(o);
    }

    public void releaseHandle(DriverHandle handle) {
        releaseHandle(handle.handle);
    }

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

    public Module loadModule(StreamSource source, String resource) {
//System.out.printf("ds.reference.length=%d\n", ds.reference.array().length);
        ModuleBrowser moduleBrowser = new ModuleBrowser(null);
        long handle = loadModule(source, resource, moduleBrowser);
        if (handle == 0) {
            throw new DriverError("Could not load module " + resource);
        }
        newHandle(handle, moduleBrowser.module);
        return moduleBrowser.module;
    }

    private native void releaseHandle(long handleId);
    private native long loadModule(StreamSource resourceSourceHand, String resource, Browser yang_browser);
    private native void initializeDriver();
}
