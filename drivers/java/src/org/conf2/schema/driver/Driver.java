package org.conf2.schema.driver;

import org.conf2.schema.Handle;
import org.conf2.schema.Module;
import org.conf2.schema.StreamSource;
import org.conf2.schema.browse.Browser;
import org.conf2.schema.browse.ModuleBrowser;

import java.io.IOException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.Map;

public class Driver {
    // TODO: Use weakreference queue and release objects automatically
    private Map<Object, DriverHandle> handles = new HashMap<Object, DriverHandle>();
    private static boolean loaded;

    public Driver() {
        if (!loaded) {
            System.loadLibrary("conf2j");
            initializeDriver();
            loaded = true;
        }
    }

    public void release() {
        for (Handle handle : handles.values()) {
            handle.Release();
        }
        handles.clear();
    }

    public DriverHandle newHandle(long hndId, Object obj) {
        DriverHandle h = new DriverHandle(this, hndId);
        handles.put(obj, h);
        return h;
    }

    public DriverHandle getHandle(Object o) {
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

    private native void releaseHandle(long handleId);
    private native void initializeDriver();
}
