package org.conf2.yang.driver;

import org.conf2.yang.Handle;
import org.conf2.yang.Module;
import org.conf2.yang.StreamSource;
import org.conf2.yang.browse.Browser;
import org.conf2.yang.browse.ModuleBrowser;

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
            System.loadLibrary("yangc2j");
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

    public Module loadModule(StreamSource source, String resource) {
        ModuleBrowser moduleBrowser = new ModuleBrowser(null);
        long handle = loadModule(source, resource, moduleBrowser);
        if (handle == 0) {
            throw new DriverError("Could not load module " + resource);
        }
        newHandle(handle, moduleBrowser.module);
        return moduleBrowser.module;
    }

    public static int utf8StrLen(CharSequence s) {
        int count = 0;
        for (int i = 0, len = s.length(); i < len; i++) {
            char ch = s.charAt(i);
            if (ch <= 0x7F) {
                count++;
            } else if (ch <= 0x7FF) {
                count += 2;
            } else if (Character.isHighSurrogate(ch)) {
                count += 4;
                ++i;
            } else {
                count += 3;
            }
        }
        return count;
    }

    public static ByteBuffer encodeCStrArray(String[] strlist) throws IOException {
        byte bNULL = 0;
        int datalen = 0;
        for (String s : strlist) {
            datalen += utf8StrLen(s) + 1;
        }
        ByteBuffer out = ByteBuffer.allocateDirect(datalen);
        for (String s : strlist) {
            out.put(s.getBytes());
            out.put(bNULL);
        }
        return out;
    }

    public static ByteBuffer encodeCIntArray(int[] intlist) throws IOException {
        ByteBuffer out = ByteBuffer.allocateDirect(4 * intlist.length);
        for (int i = 0; i < intlist.length; i++) {
            out.putInt(intlist[i]);
        }
        return out;
    }

    public static ByteBuffer encodeCBoolArray(boolean[] boollist) throws IOException {
        short sTrue = 1;
        short sFalse = 0;
        ByteBuffer out = ByteBuffer.allocateDirect(2 * boollist.length);
        for (int i = 0; i < boollist.length; i++) {
            out.putShort(boollist[i] ? sTrue : sFalse);
        }
        return out;
    }

    private native void releaseHandle(long handleId);
    private native long loadModule(StreamSource resourceSourceHand, String resource, Browser yang_browser);
    private native void initializeDriver();
}
